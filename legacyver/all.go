package legacyver

import (
	"fmt"
	"sync"

	"github.com/ethaniccc/legacy-version/internal/chunk"
	"github.com/ethaniccc/legacy-version/mapping"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var (
	EnableChunkTranslation = true

	blockTranslatorsMu sync.RWMutex
	blockTranslators   = make(map[int32]BlockTranslator)
)

// All returns a slice of all legacy protocol versions that are supported. dragonflyMapping
// must be set to true if you're using Dragonfly.
func All(dragonflyMapping bool) []minecraft.Protocol {
	return []minecraft.Protocol{
		New860(dragonflyMapping),
		New859(dragonflyMapping),
		New844(dragonflyMapping),
		New827(dragonflyMapping),
		New819(dragonflyMapping),
		New818(dragonflyMapping),
		New800(dragonflyMapping),
		New786(dragonflyMapping),
		New776(dragonflyMapping),
		New766(dragonflyMapping),
		New748(dragonflyMapping),
		New729(dragonflyMapping),
		New712(dragonflyMapping),
		New686(dragonflyMapping),
		New685(dragonflyMapping),
		New671(dragonflyMapping),
		New662(dragonflyMapping),
		New649(dragonflyMapping),
	}
}

// lookupOrCreateBlockTranslator looks up a block translator for the given protocol version.
func lookupOrCreateBlockTranslator(protocolVersion, blockVersion int32, blockStateData []byte) BlockTranslator {
	blockTranslatorsMu.RLock()
	if translator, ok := blockTranslators[protocolVersion]; ok {
		blockTranslatorsMu.RUnlock()
		return translator
	}
	blockTranslatorsMu.RUnlock()

	blockTranslatorsMu.Lock()
	defer blockTranslatorsMu.Unlock()

	if translator, ok := blockTranslators[protocolVersion]; ok {
		return translator
	}

	blockMapping := mapping.NewBlockMapping(blockStateData)
	ret := NewBlockTranslator(blockMapping, blockMappingLatest, chunk.NewNetworkPersistentEncoding(blockMapping, blockVersion), chunk.NewBlockPaletteEncoding(blockMapping, blockVersion), false)
	blockTranslators[protocolVersion] = ret
	return ret
}

// DowngradeLevelChunkPacket downgrades the given level chunk packet to the specified protocol version.
func DowngradeLevelChunkPacket(protocolVersion int32, pk *packet.LevelChunk) error {
	blockTranslatorsMu.RLock()
	translator, ok := blockTranslators[protocolVersion]
	blockTranslatorsMu.RUnlock()
	if !ok {
		return fmt.Errorf("no block translator found for protocol version %d", protocolVersion)
	}

	return translator.DowngradeLevelChunk(pk)
}
