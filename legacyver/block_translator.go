package legacyver

import (
	"bytes"
	_ "embed"

	"github.com/cespare/xxhash/v2"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/ethaniccc/legacy-version/internal/chunk"
	"github.com/ethaniccc/legacy-version/mapping"
	"github.com/hashicorp/go-version"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	BlockVersionLatest int32 = (1 << 24) | (21 << 16) | (50 << 8) | 29
)

var (
	// LatestNetworkPersistentEncoding is the Encoding used for sending a Chunk over network. It uses NBT, unlike NetworkEncoding.
	LatestNetworkPersistentEncoding = chunk.NewNetworkPersistentEncoding(blockMappingLatest, BlockVersionLatest)
	// LatestBlockPaletteEncoding is the paletteEncoding used for encoding a palette of block states encoded as NBT.
	LatestBlockPaletteEncoding = chunk.NewBlockPaletteEncoding(blockMappingLatest, BlockVersionLatest)
	cavesAndCliffsVersion, _   = version.NewVersion("1.18.0")
)

type BlockTranslator interface {
	// DowngradeBlockPackets downgrades the input block packets to legacy block packets.
	DowngradeBlockPackets([]packet.Packet, *minecraft.Conn) (result []packet.Packet)
	// UpgradeBlockPackets upgrades the input block packets to the latest block packets.
	UpgradeBlockPackets([]packet.Packet, *minecraft.Conn) (result []packet.Packet)
	// DowngradeLevelChunk downgrades the given LevelChunk packet to a legacy format.
	DowngradeLevelChunk(*packet.LevelChunk) error
	// BlockMapping returns the block mapping used by this translator.
	BlockMapping() mapping.Block
}
type DefaultBlockTranslator struct {
	mapping                 mapping.Block
	latest                  mapping.Block
	pse                     chunk.Encoding
	pe                      chunk.PaletteEncoding
	oldFormat               bool
	dimensionDefinitions    []protocol.DimensionDefinition
	currentDimension        int32
	UseBlockNetworkIDHashes bool
	BaseGameVersion         string
}

func NewBlockTranslator(mapping mapping.Block, latestMapping mapping.Block, pse chunk.Encoding, pe chunk.PaletteEncoding, oldFormat bool) *DefaultBlockTranslator {
	return &DefaultBlockTranslator{mapping: mapping, latest: latestMapping, pse: pse, pe: pe, oldFormat: oldFormat}
}

func (t *DefaultBlockTranslator) BlockMapping() mapping.Block {
	return t.mapping
}

func (t *DefaultBlockTranslator) DowngradeLevelChunk(pk *packet.LevelChunk) error {
	count := int(pk.SubChunkCount)
	if count == protocol.SubChunkRequestModeLimitless || count == protocol.SubChunkRequestModeLimited {
		return nil
	}

	buf := bytes.NewBuffer(pk.RawPayload)
	writeBuf := bytes.NewBuffer(nil)
	if !pk.CacheEnabled {
		c, err := chunk.NetworkDecode(t.latest.Air(), buf, count, false, t.getRange(), LatestNetworkPersistentEncoding, LatestBlockPaletteEncoding, t.latest, t.UseBlockNetworkIDHashes)
		if err != nil {
			return err
		}
		c = t.DowngradeChunk(c)

		payload, err := chunk.NetworkEncode(t.mapping.Air(), c, t.oldFormat, t.pe, t.mapping, false)
		if err != nil {
			return err
		}
		writeBuf.Write(payload)
		pk.SubChunkCount = uint32(len(c.Sub()))
	}
	safeBytes := buf.Bytes()

	countBorder, err := buf.ReadByte()
	if err != nil {
		pk.RawPayload = append(writeBuf.Bytes(), safeBytes...)
		return err
	}
	borderBytes := make([]byte, countBorder)
	if _, err = buf.Read(borderBytes); err != nil {
		pk.RawPayload = append(writeBuf.Bytes(), safeBytes...)
		return err
	}
	writeBuf.WriteByte(countBorder)
	writeBuf.Write(borderBytes)

	enc := nbt.NewEncoderWithEncoding(writeBuf, nbt.NetworkLittleEndian)
	dec := nbt.NewDecoderWithEncoding(buf, nbt.NetworkLittleEndian)
	for {
		var decNbt map[string]any
		if err = dec.Decode(&decNbt); err != nil {
			break
		}
		t.mapping.DowngradeBlockActorData(decNbt)

		if err = enc.Encode(decNbt); err != nil {
			break
		}
	}
	pk.RawPayload = append(writeBuf.Bytes(), buf.Bytes()...)
	return nil
}

func (t *DefaultBlockTranslator) DowngradeBlockPackets(pks []packet.Packet, conn *minecraft.Conn) (result []packet.Packet) {
	for _, pk := range pks {
		switch pk := pk.(type) {
		case *packet.LevelChunk:
			if !EnableChunkTranslation {
				break
			}
			if err := t.DowngradeLevelChunk(pk); err != nil {
				//fmt.Println(err)
				break
			}
		case *packet.SubChunk:
			if !EnableChunkTranslation {
				break
			}
			r := t.getRange()
			if t.oldFormat {
				r = cube.Range{0, 255}
			}

			for i, entry := range pk.SubChunkEntries {
				if entry.Result == protocol.SubChunkResultSuccess {
					buf := bytes.NewBuffer(entry.RawPayload)
					writeBuf := bytes.NewBuffer(nil)
					if !pk.CacheEnabled && !conn.ClientCacheEnabled() {
						ind := byte(i)
						subChunk, err := chunk.DecodeSubChunk(t.latest.Air(), r, buf, &ind, chunk.NetworkEncoding, LatestNetworkPersistentEncoding, LatestBlockPaletteEncoding, t.latest, t.UseBlockNetworkIDHashes)
						if err != nil {
							//fmt.Println(err)
							continue
						}
						t.DowngradeSubChunk(subChunk)
						writeBuf.Write(chunk.EncodeSubChunk(subChunk, chunk.NetworkEncoding, t.pe, chunk.SubChunkVersion9, r, int(ind), t.mapping, false))
					}

					enc := nbt.NewEncoderWithEncoding(writeBuf, nbt.NetworkLittleEndian)
					dec := nbt.NewDecoderWithEncoding(buf, nbt.NetworkLittleEndian)
					for {
						var decNbt map[string]any
						if err := dec.Decode(&decNbt); err != nil {
							break
						}
						t.mapping.DowngradeBlockActorData(decNbt)

						if err := enc.Encode(decNbt); err != nil {
							break
						}
					}

					entry.RawPayload = append(writeBuf.Bytes(), buf.Bytes()...)
					entry.BlobHash = xxhash.Sum64(entry.RawPayload)
					pk.SubChunkEntries[i] = entry
				}
			}
		case *packet.ClientCacheMissResponse:
			r := t.getRange()
			if t.oldFormat {
				r = cube.Range{0, 255}
			}

			for i, blob := range pk.Blobs {
				buf := bytes.NewBuffer(blob.Payload)
				ind := byte(0)
				subChunk, err := chunk.DecodeSubChunk(t.latest.Air(), r, buf, &ind, chunk.NetworkEncoding, LatestNetworkPersistentEncoding, LatestBlockPaletteEncoding, t.latest, t.UseBlockNetworkIDHashes)
				if err != nil {
					// Has a possibility to be a biome, ignore then
					continue
				}
				t.DowngradeSubChunk(subChunk)
				blob.Payload = append(chunk.EncodeSubChunk(subChunk, chunk.NetworkEncoding, t.pe, chunk.SubChunkVersion9, r, int(ind), t.mapping, false), buf.Bytes()...)
				blob.Hash = xxhash.Sum64(blob.Payload)
				pk.Blobs[i] = blob
			}
		case *packet.UpdateSubChunkBlocks:
			for i, block := range pk.Blocks {
				block.BlockRuntimeID = t.DowngradeBlockRuntimeID(block.BlockRuntimeID)
				pk.Blocks[i] = block
			}
			for i, block := range pk.Extra {
				block.BlockRuntimeID = t.DowngradeBlockRuntimeID(block.BlockRuntimeID)
				pk.Extra[i] = block
			}
		case *packet.UpdateBlock:
			pk.NewBlockRuntimeID = t.DowngradeBlockRuntimeID(pk.NewBlockRuntimeID)
		case *packet.UpdateBlockSynced:
			pk.NewBlockRuntimeID = t.DowngradeBlockRuntimeID(pk.NewBlockRuntimeID)
		case *packet.InventoryTransaction:
			if transactionData, ok := pk.TransactionData.(*protocol.UseItemTransactionData); ok {
				transactionData.BlockRuntimeID = t.DowngradeBlockRuntimeID(transactionData.BlockRuntimeID)
				pk.TransactionData = transactionData
			}
		case *packet.LevelEvent:
			switch pk.EventType {
			case packet.LevelEventParticleLegacyEvent | 20: // terrain
				fallthrough
			case packet.LevelEventParticlesDestroyBlock:
				fallthrough
			case packet.LevelEventParticlesDestroyBlockNoSound:
				pk.EventData = int32(t.DowngradeBlockRuntimeID(uint32(pk.EventData)))
			case packet.LevelEventParticlesCrackBlock:
				face := pk.EventData >> 24
				rid := t.DowngradeBlockRuntimeID(uint32(pk.EventData & 0xffff))
				pk.EventData = int32(rid) | (face << 24)
			}
		case *packet.LevelSoundEvent:
			switch pk.SoundType {
			case packet.SoundEventBreak:
				fallthrough
			case packet.SoundEventPlace:
				fallthrough
			case packet.SoundEventHit:
				fallthrough
			case packet.SoundEventLand:
				fallthrough
			case packet.SoundEventItemUseOn:
				pk.ExtraData = int32(t.DowngradeBlockRuntimeID(uint32(pk.ExtraData)))
			}
		case *packet.AddActor:
			if pk.EntityType == "minecraft:falling_block" {
				pk.EntityMetadata = t.downgradeEntityMetadata(pk.EntityMetadata)
			}
		case *packet.SetActorData:
			//pk.EntityMetadata = t.downgradeEntityMetadata(pk.EntityMetadata)
		case *packet.StartGame:
			t.UseBlockNetworkIDHashes = pk.UseBlockNetworkIDHashes
			t.BaseGameVersion = pk.BaseGameVersion
			t.latest.Adjust(pk.Blocks)
			t.mapping.Adjust(pk.Blocks)
		case *packet.ResourcePackStack:
			var packs []protocol.StackResourcePack
			for _, pack := range pk.TexturePacks {
				if pack.UUID == "0fba4063-dba1-4281-9b89-ff9390653530" {
					continue
				}
				packs = append(packs, pack)
			}
			pk.TexturePacks = packs
		}
		result = append(result, pk)
	}
	return result
}

func (t *DefaultBlockTranslator) UpgradeBlockPackets(pks []packet.Packet, conn *minecraft.Conn) (result []packet.Packet) {
	for _, pk := range pks {
		switch pk := pk.(type) {
		case *packet.InventoryTransaction:
			if transactionData, ok := pk.TransactionData.(*protocol.UseItemTransactionData); ok {
				transactionData.BlockRuntimeID = t.UpgradeBlockRuntimeID(transactionData.BlockRuntimeID)
				pk.TransactionData = transactionData
			}
		case *packet.SetActorData:
			//pk.EntityMetadata = t.upgradeEntityMetadata(pk.EntityMetadata)
		}
		result = append(result, pk)
	}
	return result
}

func (t *DefaultBlockTranslator) DowngradeBlockRuntimeID(input uint32) uint32 {
	if t.latest == t.mapping {
		return input
	}
	if t.UseBlockNetworkIDHashes {
		var ok bool
		input, ok = t.latest.HashToRuntimeID(input)
		if !ok {
			return t.mapping.InfoUpdate()
		}
	}
	state, ok := t.latest.RuntimeIDToState(input)
	if !ok {
		return t.mapping.InfoUpdate()
	}
	runtimeID, ok := t.mapping.StateToRuntimeID(state)
	if !ok {
		return t.mapping.InfoUpdate()
	}
	if t.UseBlockNetworkIDHashes {
		runtimeID, ok = t.mapping.RuntimeIDToHash(runtimeID)
		if !ok {
			return t.mapping.InfoUpdate()
		}
	}
	return runtimeID
}

func (t *DefaultBlockTranslator) DowngradeChunk(input *chunk.Chunk) *chunk.Chunk {
	if t.latest == t.mapping {
		return input
	}

	start := 0
	r := t.getRange()
	if t.oldFormat {
		start = 4
		r = cube.Range{0, 255}
	}
	downgraded := chunk.New(t.mapping.Air(), r)

	i := 0
	// First downgrade the blocks.
	for _, sub := range input.Sub()[start : len(input.Sub())-start] {
		t.DowngradeSubChunk(sub)
		downgraded.Sub()[i] = sub
		i += 1
	}
	i = 0
	// Then downgrade the biome ids.
	for _, sub := range input.BiomeSub()[start : len(input.BiomeSub())-start] {
		// todo
		sub.Palette().Replace(func(v uint32) uint32 {
			return 0 // at least the client doesn't crash now
		})
		downgraded.BiomeSub()[i] = sub
		i += 1
	}

	return downgraded
}

func (t *DefaultBlockTranslator) DowngradeSubChunk(input *chunk.SubChunk) {
	if t.latest == t.mapping {
		return
	}
	for _, storage := range input.Layers() {
		storage.Palette().Replace(t.DowngradeBlockRuntimeID)
	}
}

func (t *DefaultBlockTranslator) downgradeEntityMetadata(metadata map[uint32]any) map[uint32]any {
	if t.latest == t.mapping {
		return metadata
	}
	if latestRID, ok := metadata[protocol.EntityDataKeyVariant]; ok {
		metadata[protocol.EntityDataKeyVariant] = int32(t.DowngradeBlockRuntimeID(uint32(latestRID.(int32))))
	}
	return metadata
}

func (t *DefaultBlockTranslator) UpgradeBlockRuntimeID(input uint32) uint32 {
	if t.latest == t.mapping {
		return input
	}
	if t.UseBlockNetworkIDHashes {
		var ok bool
		input, ok = t.mapping.HashToRuntimeID(input)
		if !ok {
			return t.latest.InfoUpdate()
		}
	}
	state, ok := t.mapping.RuntimeIDToState(input)
	if !ok {
		return t.latest.InfoUpdate()
	}
	runtimeID, ok := t.latest.StateToRuntimeID(state)
	if !ok {
		return t.latest.InfoUpdate()
	}
	if t.UseBlockNetworkIDHashes {
		runtimeID, ok = t.latest.RuntimeIDToHash(runtimeID)
		if !ok {
			return t.latest.InfoUpdate()
		}
	}
	return runtimeID
}

func (t *DefaultBlockTranslator) upgradeEntityMetadata(metadata map[uint32]any) map[uint32]any {
	if t.latest == t.mapping {
		return metadata
	}
	if latestRID, ok := metadata[protocol.EntityDataKeyVariant]; ok {
		metadata[protocol.EntityDataKeyVariant] = int32(t.UpgradeBlockRuntimeID(uint32(latestRID.(int32))))
	}
	return metadata
}

func (t *DefaultBlockTranslator) getRange() (r cube.Range) {
	var dimName string
	switch t.currentDimension {
	case 0:
		dimName = "minecraft:overworld"
		r = world.Overworld.Range()
		v, err := version.NewVersion(t.BaseGameVersion)
		if err == nil && v.LessThan(cavesAndCliffsVersion) {
			r = cube.Range{0, 255}
		}
	case 1:
		dimName = "minecraft:nether"
		r = world.Nether.Range()
	case 2:
		dimName = "minecraft:the_end"
		r = world.End.Range()
	}
	for _, def := range t.dimensionDefinitions {
		if def.Name == dimName {
			r = cube.Range{int(def.Range[0]), int(def.Range[1])}
			break
		}
	}
	return
}
