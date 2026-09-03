package legacyver

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/akmalfairuz/legacy-version/legacyver/legacypacket"
	legacyproto "github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestProtocol2168MarshalsLikeLatest(t *testing.T) {
	tests := []packet.Packet{
		&packet.AnvilDamage{},
		&packet.ClientBoundMapItemData{},
		&packet.ClientboundUpdateSoundData{},
		&packet.CraftingData{},
		&packet.CreativeContent{},
		&packet.DimensionData{},
		&packet.GameRulesChanged{},
		&packet.InventoryContent{},
		&packet.InventorySlot{},
		&packet.InventoryTransaction{TransactionData: &protocol.NormalTransactionData{}},
		&packet.ItemStackRequest{},
		&packet.ItemStackResponse{},
		&packet.LevelChunk{},
		&packet.MobArmourEquipment{},
		&packet.MobEquipment{},
		&packet.MoveActorDelta{},
		&packet.MovePlayer{},
		&packet.PlaySound{},
		&packet.PlayerAuthInput{},
		&packet.PlayerList{},
		&packet.PlayerLocation{},
		&packet.PlayerSkin{},
		&packet.PlayerUpdateEntityOverrides{},
		&packet.PrimitiveShapes{},
		&packet.ResourcePackClientResponse{},
		&packet.ResourcePacksInfo{},
		&packet.ServerBoundDiagnostics{},
		&packet.SetScore{},
		&packet.SetScoreboardIdentity{},
		&packet.StructureBlockUpdate{},
		&packet.SubChunk{},
		&packet.SubChunkRequest{},
		&packet.Transfer{},
	}
	for _, pk := range tests {
		t.Run(packetName(pk), func(t *testing.T) {
			marshalFn, ok := packets[pk.ID()]
			if !ok {
				t.Fatalf("packet %T is not registered for compatibility marshaling", pk)
			}
			var native, compatible bytes.Buffer
			pk.Marshal(protocol.NewWriter(&native, 0))
			marshalFn(legacyproto.NewWriter(protocol.NewWriter(&compatible, 0), legacyproto.ID2168, 0), pk)
			if !bytes.Equal(native.Bytes(), compatible.Bytes()) {
				t.Fatalf("1.26.40 marshal differs from latest\nnative:     %x\ncompatible: %x", native.Bytes(), compatible.Bytes())
			}
		})
	}
}

func TestProtocol1001EntityMetadataByteLayout(t *testing.T) {
	metadata := protocol.EntityMetadata{1: byte(1)}
	var encoded bytes.Buffer
	legacyproto.MarshalEntityMetadata(New1001().NewWriter(&encoded, 0), &metadata)
	expected := []byte{1, 1, byte(protocol.EntityDataTypeByte), 1}
	if !bytes.Equal(encoded.Bytes(), expected) {
		t.Fatalf("protocol 1001 metadata differs\nexpected: %x\nactual:   %x", expected, encoded.Bytes())
	}
}

func TestProtocol2168TranslationMatchesGophertunnelProtocol(t *testing.T) {
	flags := protocol.NewInputFlags(packet.InputFlagCount)
	flags.Set(packet.InputFlagPerformItemInteraction)
	emptyFlags := protocol.NewInputFlags(packet.InputFlagCount)
	fixtures := map[string]packet.Packet{
		"empty/play_sound":        &packet.PlaySound{},
		"full/play_sound":         &packet.PlaySound{SoundName: "note.pling", Position: [3]float32{1.25, -2.5, 3.75}, Volume: .75, Pitch: 1.5, LoopCount: 3, Handle: protocol.Option(uint64(77))},
		"empty/player_auth_input": &packet.PlayerAuthInput{InputData: emptyFlags},
		"full/player_auth_input": &packet.PlayerAuthInput{
			InputData: flags, ItemInteractionData: protocol.Option(protocol.UseItemTransactionData{
				LegacySetItemSlots: protocol.Option([]protocol.LegacySetItemSlot{{ContainerID: 2, Slots: []byte{1, 4}}}),
				Actions:            protocol.Option([]protocol.InventoryAction{}), ActionType: protocol.UseItemActionClickBlock,
				TriggerType: protocol.TriggerTypePlayerInput, BlockPosition: protocol.BlockPos{10, 64, -10}, BlockFace: 2,
				ClientPrediction: protocol.ClientPredictionSuccess, ClientCooldownState: protocol.ClientCooldownStateOn,
			}),
		},
		"empty/sub_chunk": &packet.SubChunk{},
		"full/sub_chunk": &packet.SubChunk{CacheEnabled: true, SubChunkEntries: []protocol.SubChunkEntry{{
			Offset: protocol.SubChunkOffset{1, -2, 3}, Result: protocol.SubChunkResultSuccess,
			RawPayload: protocol.Option([]byte{1, 2, 3}), HeightMapType: protocol.HeightMapDataHasData,
			HeightMapData: protocol.Option(make([]int8, 256)), RenderHeightMapType: protocol.HeightMapDataNone,
			BlobHash: protocol.Option(uint64(99)),
		}}},
		"empty/player_list":      &packet.PlayerList{},
		"full/player_list":       &packet.PlayerList{Entries: []protocol.PlayerListEntry{{ActionType: protocol.PlayerListActionRemove}}},
		"empty/set_score":        &packet.SetScore{},
		"full/set_score":         &packet.SetScore{Entries: []protocol.ScoreboardEntry{{IdentityType: protocol.ScoreboardIdentityFakePlayer, EntryID: 4, ObjectiveName: "obj", Score: 7, DisplayName: "line"}}},
		"empty/update_abilities": &packet.UpdateAbilities{},
		"full/set_actor_data": &packet.SetActorData{
			EntityRuntimeID: 0x123456789, Tick: 0xabcdef,
			EntityMetadata: protocol.EntityMetadata{
				1: byte(0x7f), 2: int16(-1234), 3: int32(-56789), 4: float32(12.5), 5: "dummy metadata",
				6: map[string]any{"name": "dummy"}, 7: protocol.BlockPos{-12, 64, 99},
				8: int64(-9876543210), 9: mgl32.Vec3{1.25, -2.5, 3.75},
			},
			EntityProperties: protocol.EntityProperties{
				IntegerProperties: []protocol.IntegerEntityProperty{{Index: 3, Value: -77}, {Index: 8, Value: 99}},
				FloatProperties:   []protocol.FloatEntityProperty{{Index: 4, Value: 1.5}, {Index: 9, Value: -2.75}},
			},
		},
		"full/update_abilities": &packet.UpdateAbilities{AbilityData: protocol.AbilityData{
			EntityUniqueID: 123, PlayerPermissions: 2, CommandPermissions: 1,
			Layers: []protocol.AbilityLayer{{Type: protocol.AbilityLayerTypeBase, Abilities: 3, Values: 1, FlySpeed: .1, VerticalFlySpeed: .2, WalkSpeed: .3}},
		}},
		"empty/hurt_armour":                    &packet.HurtArmour{},
		"full/hurt_armour":                     &packet.HurtArmour{Cause: -3, Damage: 7, ArmourSlots: 9},
		"empty/player_location":                &packet.PlayerLocation{Type: packet.PlayerLocationTypeHide},
		"full/player_location":                 &packet.PlayerLocation{Type: packet.PlayerLocationTypeCoordinates, EntityUniqueID: -44, Position: [3]float32{1, 2, 3}},
		"empty/player_update_entity_overrides": &packet.PlayerUpdateEntityOverrides{},
		"full/player_update_entity_overrides": &packet.PlayerUpdateEntityOverrides{
			EntityUniqueID: 55, PropertyIndex: 3, Type: packet.PlayerUpdateEntityOverridesTypeInt, IntValue: -19,
		},
		"empty/primitive_shapes": &packet.PrimitiveShapes{},
		"full/primitive_shapes": &packet.PrimitiveShapes{Shapes: []protocol.PrimitiveShape{{
			NetworkID: 8, Type: protocol.Option(uint8(protocol.PrimitiveShapeText)),
			Location: protocol.Option(mgl32.Vec3{1, 2, 3}), Scale: protocol.Option(float32(2)),
			Rotation: protocol.Option(mgl32.Vec3{4, 5, 6}), TotalTimeLeft: protocol.Option(float32(7)),
			MaxRenderDistance: protocol.Option(float32(64)), Colour: protocol.Option(color.RGBA{R: 1, G: 2, B: 3, A: 4}),
			DimensionID: protocol.Option(int32(-1)), AttachedToEntityID: protocol.Option(int64(-99)),
			ExtraShapeData: &protocol.TextShape{Text: "test", UseRotation: true, BackgroundColour: protocol.Option(color.RGBA{R: 5, G: 6, B: 7, A: 8}), DepthTest: true},
		}}},
		"empty/server_presence_info": &packet.ServerPresenceInfo{},
		"full/server_presence_info": &packet.ServerPresenceInfo{PresenceInfo: protocol.Option(protocol.PresenceInfo{
			RichPresenceID: protocol.Option("presence"),
		})},
	}
	translated := &Protocol{id: legacyproto.ID2168, ver: "1.26.40"}
	for name, fixture := range fixtures {
		t.Run(name, func(t *testing.T) {
			native := marshalThroughProtocol(t, minecraft.DefaultProtocol, fixture)
			compat := marshalThroughProtocol(t, translated, fixture)
			if !reflect.DeepEqual(native, compat) {
				t.Fatalf("protocol 2168 output differs\nnative:     %x\ntranslated: %x", native, compat)
			}
			assertProtocolRoundTrip(t, translated, fixture.ID(), native)
		})
	}
	t.Run("full/command_parameter_float", func(t *testing.T) {
		parameter := protocol.CommandParameter{
			Name: "value", Type: protocol.CommandArgValid | protocol.CommandArgTypeFloat, Optional: true, Options: 3,
		}
		var native, compatible bytes.Buffer
		parameter.Marshal(minecraft.DefaultProtocol.NewWriter(&native, 0))
		legacyproto.MarshalCommandParameter(translated.NewWriter(&compatible, 0), &parameter)
		if !bytes.Equal(native.Bytes(), compatible.Bytes()) {
			t.Fatalf("protocol 2168 command parameter differs\nnative:     %x\ntranslated: %x", native.Bytes(), compatible.Bytes())
		}
		decoded := protocol.CommandParameter{}
		readerBuffer := bytes.NewBuffer(native.Bytes())
		legacyproto.MarshalCommandParameter(translated.NewReader(readerBuffer, 0, true), &decoded)
		if readerBuffer.Len() != 0 {
			t.Fatalf("protocol 2168 command parameter left %d bytes unread", readerBuffer.Len())
		}
	})
}

func TestProtocol2168AllEmptyPacketsMatchGophertunnel(t *testing.T) {
	translated := &Protocol{id: legacyproto.ID2168, ver: "1.26.40"}
	for _, listener := range []bool{false, true} {
		direction := "server"
		if listener {
			direction = "client"
		}
		for id, constructor := range minecraft.DefaultProtocol.Packets(listener) {
			pk := constructor()
			initialiseCurrentEmptyPacket(pk, legacyproto.ID2168)
			t.Run(fmt.Sprintf("%s/%d/%T", direction, id, pk), func(t *testing.T) {
				native, err := tryMarshalNative(pk)
				if err != nil {
					t.Skipf("empty packet is not a valid native payload: %v", err)
				}
				compatible := marshalCompatibilityPacket(t, translated, pk)
				if !bytes.Equal(native, compatible) {
					t.Fatalf("protocol 2168 empty packet differs\nnative:     %x\ntranslated: %x", native, compatible)
				}
				assertProtocolRoundTrip(t, translated, id, [][]byte{native})
			})
		}
	}
}

func TestProtocol2168AllPopulatedCollectionsMatchGophertunnel(t *testing.T) {
	translated := &Protocol{id: legacyproto.ID2168, ver: "1.26.40"}
	for _, listener := range []bool{false, true} {
		direction := "server"
		if listener {
			direction = "client"
		}
		for id, constructor := range minecraft.DefaultProtocol.Packets(listener) {
			pk := constructor()
			initialiseCurrentEmptyPacket(pk, legacyproto.ID2168)
			populateCollectionFields(reflect.ValueOf(pk), false, 0)
			normalisePopulatedPacket(pk, legacyproto.ID2168)
			t.Run(fmt.Sprintf("%s/%d/%T", direction, id, pk), func(t *testing.T) {
				native, err := tryMarshalNative(pk)
				if err != nil {
					t.Fatalf("populated native packet is invalid: %v", err)
				}
				compatible := marshalCompatibilityPacket(t, translated, pk)
				if !bytes.Equal(native, compatible) {
					t.Fatalf("protocol 2168 populated packet differs\nnative:     %x\ntranslated: %x", native, compatible)
				}
				assertProtocolRoundTrip(t, translated, id, [][]byte{native})
			})
		}
	}
}

func populateCollectionFields(value reflect.Value, inCollection bool, depth int) {
	if !value.IsValid() || depth > 12 {
		return
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return
		}
		populateCollectionFields(value.Elem(), inCollection, depth+1)
		return
	}
	if value.Kind() == reflect.Interface {
		if !value.IsNil() {
			populateCollectionFields(value.Elem(), inCollection, depth+1)
		}
		return
	}
	switch value.Kind() {
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			if value.Field(i).CanSet() {
				populateCollectionFields(value.Field(i), inCollection, depth+1)
			}
		}
	case reflect.Slice:
		if value.IsNil() || value.Len() == 0 {
			value.Set(reflect.MakeSlice(value.Type(), 1, 1))
		}
		for i := 0; i < value.Len(); i++ {
			populateCollectionFields(value.Index(i), true, depth+1)
		}
	case reflect.Array:
		for i := 0; i < value.Len(); i++ {
			populateCollectionFields(value.Index(i), true, depth+1)
		}
	case reflect.Map:
		if value.Type().Elem().Kind() == reflect.Interface {
			value.Set(reflect.MakeMap(value.Type()))
			switch value.Type().Key().Kind() {
			case reflect.String:
				value.SetMapIndex(reflect.ValueOf("dummy").Convert(value.Type().Key()), reflect.ValueOf(int32(1)))
			case reflect.Uint32:
				value.SetMapIndex(reflect.ValueOf(uint32(1)).Convert(value.Type().Key()), reflect.ValueOf(byte(1)))
			}
		}
	case reflect.String:
		if inCollection && value.Len() == 0 {
			value.SetString("dummy")
		}
	case reflect.Bool:
		// A zero boolean is valid dummy data and avoids enabling dependent union fields.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Keep enum-like integers at their usually valid zero value.
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Keep enum-like integers at their usually valid zero value.
	case reflect.Float32, reflect.Float64:
		// Zero is valid dummy floating-point data.
	}
}

func normalisePopulatedPacket(pk packet.Packet, protocolID int32) {
	initialiseCurrentEmptyPacket(pk, protocolID)
	switch pk := pk.(type) {
	case *packet.PlayerList:
		pk.Entries = []protocol.PlayerListEntry{
			{ActionType: protocol.PlayerListActionAdd, Username: "trusted", Skin: validDummySkin(true)},
			{ActionType: protocol.PlayerListActionAdd, Username: "untrusted", Skin: validDummySkin(false)},
		}
	case *packet.PlayerSkin:
		pk.Skin = validDummySkin(true)
	case *packet.ItemStackRequest:
		pk.Requests = []protocol.ItemStackRequest{{Actions: []protocol.StackRequestAction{&protocol.TakeStackRequestAction{}}, FilterStrings: []string{"dummy"}}}
	case *packet.PrimitiveShapes:
		pk.Shapes = []protocol.PrimitiveShape{{Type: protocol.Option(uint8(protocol.PrimitiveShapeText)), ExtraShapeData: &protocol.TextShape{Text: "dummy"}}}
	case *packet.SubChunk:
		pk.SubChunkEntries = []protocol.SubChunkEntry{{Result: protocol.SubChunkResultSuccess, RawPayload: protocol.Option([]byte{1})}}
	case *packet.CraftingData:
		normaliseCurrentCraftingData(pk)
	case *packet.GameRulesChanged:
		pk.GameRules = []protocol.GameRule{{Name: "dummy", Value: false}}
	case *packet.SetScore:
		pk.Entries = []protocol.ScoreboardEntry{{IdentityType: protocol.ScoreboardIdentityFakePlayer, DisplayName: "dummy"}}
	case *packet.StartGame:
		pk.GameRules = []protocol.GameRule{{Name: "dummy", Value: false}}
	case *packet.LevelSoundEvent:
		if protocolID < legacyproto.ID1001 {
			pk.SoundType = packet.SoundEventItemUseOn
		}
	case *packet.CameraSpline:
		pk.Splines = []protocol.CameraSplineDefinition{{
			Name: "dummy", ControlPoints: []mgl32.Vec3{{}}, ProgressKeyFrames: []protocol.CameraProgressOption{{}}, RotationKeyFrames: []protocol.CameraRotationOption{{}},
		}}
	}
}

func normaliseCurrentCraftingData(pk *packet.CraftingData) {
	descriptor := func() protocol.ItemDescriptorCount {
		return protocol.ItemDescriptorCount{Descriptor: &protocol.InvalidItemDescriptor{}, Count: 1}
	}
	for i := range pk.ShapedRecipes {
		pk.ShapedRecipes[i].Width, pk.ShapedRecipes[i].Height = 1, 1
		pk.ShapedRecipes[i].Input = []protocol.ItemDescriptorCount{descriptor()}
	}
	for i := range pk.ShapelessRecipes {
		pk.ShapelessRecipes[i].Input = []protocol.ItemDescriptorCount{descriptor()}
	}
	for i := range pk.UserDataShapelessRecipes {
		pk.UserDataShapelessRecipes[i].Input = []protocol.ItemDescriptorCount{descriptor()}
	}
	for i := range pk.ShapelessChemistryRecipes {
		pk.ShapelessChemistryRecipes[i].Input = []protocol.ItemDescriptorCount{descriptor()}
	}
	for i := range pk.ShapedChemistryRecipes {
		pk.ShapedChemistryRecipes[i].Width, pk.ShapedChemistryRecipes[i].Height = 1, 1
		pk.ShapedChemistryRecipes[i].Input = []protocol.ItemDescriptorCount{descriptor()}
	}
	for i := range pk.SmithingTransformRecipes {
		pk.SmithingTransformRecipes[i].Template = descriptor()
		pk.SmithingTransformRecipes[i].Base = descriptor()
		pk.SmithingTransformRecipes[i].Addition = descriptor()
	}
	for i := range pk.SmithingTrimRecipes {
		pk.SmithingTrimRecipes[i].Template = descriptor()
		pk.SmithingTrimRecipes[i].Base = descriptor()
		pk.SmithingTrimRecipes[i].Addition = descriptor()
	}
}

func validDummySkin(trusted bool) protocol.Skin {
	return protocol.Skin{
		SkinID: "dummy", SkinResourcePatch: []byte{9}, SkinImageWidth: 1, SkinImageHeight: 1, SkinData: []byte{1, 2, 3, 4},
		Animations:     []protocol.SkinAnimation{{ImageWidth: 1, ImageHeight: 1, ImageData: []byte{10, 11, 12, 13}, AnimationType: protocol.SkinAnimationHead, FrameCount: 1}},
		CapeImageWidth: 1, CapeImageHeight: 1, CapeData: []byte{5, 6, 7, 8}, SkinGeometry: []byte{14}, GeometryDataEngineVersion: []byte{15}, AnimationData: []byte{16},
		ArmSize: protocol.ArmSizeSlim, SkinColour: color.RGBA{A: 0xff}, Trusted: trusted,
		PersonaPieces:    []protocol.PersonaPiece{{PieceID: "piece", PieceType: protocol.PieceTypeBody}},
		PieceTintColours: []protocol.PersonaPieceTintColour{{PieceType: "persona_eyes", Colours: [4]color.RGBA{{A: 0xff, R: 1}, {A: 0xff, G: 2}, {A: 0xff, B: 3}, {A: 0xff, R: 4}}}},
	}
}

func tryMarshalNative(pk packet.Packet) (data []byte, err any) {
	defer func() { err = recover() }()
	var buf bytes.Buffer
	pk.Marshal(minecraft.DefaultProtocol.NewWriter(&buf, 0))
	return buf.Bytes(), nil
}

func initialiseCurrentEmptyPacket(pk packet.Packet, protocolID int32) {
	switch pk := pk.(type) {
	case *packet.PlayerAuthInput:
		size := packet.InputFlagCount
		if protocolID == legacyproto.ID1001 {
			size = 65
		}
		pk.InputData = protocol.NewInputFlags(size)
	case *packet.ResourcePackClientResponse:
		pk.Response = packet.PackResponseRefused
	case *packet.InventoryTransaction:
		pk.TransactionData = &protocol.NormalTransactionData{}
	case *packet.Event:
		pk.Event = &protocol.AchievementAwardedEvent{}
	case *packet.RequestAbility:
		pk.Value = false
	case *packet.ServerBoundPackSettingChange:
		pk.PackSetting.Value = false
	case *packet.ClientMovementPredictionSync:
		size := protocol.EntityDataFlagCount
		if protocolID < legacyproto.ID2168 {
			size = legacyproto.ClientMovementPredictionFlagsLength(protocolID)
		}
		pk.ActorFlags = protocol.NewBitset(size)
	case *packet.Text:
		pk.Message = "empty"
	case *packet.ServerBoundDataDrivenScreenClosed:
		if protocolID < legacyproto.ID1001 {
			pk.CloseReason = packet.DataDrivenScreenCloseReasonProgrammaticClose
		}
	case *packet.LevelSoundEvent:
		if protocolID < legacyproto.ID1001 {
			pk.SoundType = packet.SoundEventItemUseOn
		}
	}
}

func marshalCompatibilityPacket(t *testing.T, p *Protocol, pk packet.Packet) []byte {
	t.Helper()
	var buf bytes.Buffer
	writer := p.NewWriter(&buf, 0)
	if pk.ID() == packet.IDStartGame {
		legacypacket.StartGame(writer, pk.(*packet.StartGame), nil)
	} else if marshalFn, ok := packets[pk.ID()]; ok {
		marshalFn(writer, pk)
	} else {
		pk.Marshal(writer)
	}
	return buf.Bytes()
}

func TestProtocol1001TranslationMatchesNativeGophertunnel(t *testing.T) {
	reference := native1001Reference(t)
	flags := protocol.NewInputFlags(packet.InputFlagCount)
	flags.Set(packet.InputFlagPerformItemInteraction)
	emptyFlags := protocol.NewInputFlags(65)
	fixtures := map[string]packet.Packet{
		"empty/play_sound": &packet.PlaySound{},
		"full/play_sound": &packet.PlaySound{
			SoundName: "note.pling", Position: [3]float32{1.25, -2.5, 3.75}, Volume: .75, Pitch: 1.5,
			Handle: protocol.Option(uint64(77)),
		},
		"empty/sub_chunk_request": &packet.SubChunkRequest{},
		"full/sub_chunk_request": &packet.SubChunkRequest{
			Dimension: -1, Position: protocol.SubChunkPos{123, -7, 456},
			Offsets: []protocol.SubChunkOffset{{1, -2, 3}, {-4, 5, -6}},
		},
		"empty/resource_pack_response": &packet.ResourcePackClientResponse{Response: packet.PackResponseRefused},
		"full/resource_pack_response": &packet.ResourcePackClientResponse{
			Response: packet.PackResponseSendPacks, PacksToDownload: []string{"pack-a_1.0.0", "pack-b_2.0.0"},
		},
		"empty/player_auth_input": &packet.PlayerAuthInput{InputData: emptyFlags},
		"full/player_auth_input": &packet.PlayerAuthInput{
			Pitch: 12.5, Yaw: -27.25, Position: [3]float32{1, 2, 3}, MoveVector: [2]float32{.25, -.5},
			HeadYaw: 8, InputData: flags, InputMode: packet.InputModeMouse, PlayMode: packet.PlayModeNormal,
			InteractionModel: packet.InteractionModelCrosshair, InteractPitch: 4, InteractYaw: 5, Tick: 99,
			Delta: [3]float32{.1, .2, .3}, AnalogueMoveVector: [2]float32{.4, .5},
			CameraOrientation: [3]float32{0, 1, 0}, RawMoveVector: [2]float32{-.25, .75},
			ItemInteractionData: protocol.Option(protocol.UseItemTransactionData{
				Actions: protocol.Option([]protocol.InventoryAction{}), ActionType: protocol.UseItemActionClickBlock,
				TriggerType: protocol.TriggerTypePlayerInput, BlockPosition: protocol.BlockPos{10, 64, -10},
				BlockFace: 2, HotBarSlot: 3, Position: [3]float32{10.5, 65, -9.5},
				ClickedPosition: [3]float32{.5, 1, .5}, BlockRuntimeID: 42,
				ClientPrediction: protocol.ClientPredictionSuccess, ClientCooldownState: protocol.ClientCooldownStateOn,
			}),
		},
		"empty/update_abilities": &packet.UpdateAbilities{},
		"full/set_actor_data": &packet.SetActorData{
			EntityRuntimeID: 0x123456789, Tick: 0xabcdef,
			EntityMetadata: protocol.EntityMetadata{
				1: byte(0x7f), 2: int16(-1234), 3: int32(-56789), 4: float32(12.5), 5: "dummy metadata",
				6: map[string]any{"name": "dummy"}, 7: protocol.BlockPos{-12, 64, 99},
				8: int64(-9876543210), 9: mgl32.Vec3{1.25, -2.5, 3.75},
			},
			EntityProperties: protocol.EntityProperties{
				IntegerProperties: []protocol.IntegerEntityProperty{{Index: 3, Value: -77}, {Index: 8, Value: 99}},
				FloatProperties:   []protocol.FloatEntityProperty{{Index: 4, Value: 1.5}, {Index: 9, Value: -2.75}},
			},
		},
		"full/level_chunk_cached": &packet.LevelChunk{
			Position: protocol.ChunkPos{12, -34}, Dimension: -1, SubChunkCount: 2, CacheEnabled: true,
			BlobHashes: []uint64{0x1122334455667788, 0x8877665544332211, 0x0102030405060708}, RawPayload: []byte{1, 2, 3, 4, 5},
		},
		"full/level_chunk_limited": &packet.LevelChunk{
			Position: protocol.ChunkPos{-8, 19}, Dimension: 2, SubChunkLimit: protocol.Option(int32(37)), RawPayload: []byte{9, 8, 7},
		},
		"full/level_chunk_limitless": &packet.LevelChunk{
			Position: protocol.ChunkPos{4, 5}, Dimension: 1, SubChunkLimit: protocol.Option(int32(-1)), RawPayload: []byte{6, 5, 4},
		},
		"full/player_list_trusted_skins": &packet.PlayerList{Entries: []protocol.PlayerListEntry{
			{ActionType: protocol.PlayerListActionAdd, UUID: uuid.MustParse("11111111-2222-3333-4444-555555555555"), EntityUniqueID: -12345, Username: "trusted", XUID: "123456789", PlatformChatID: "platform-a", BuildPlatform: 7, Skin: validDummySkin(true), Teacher: true, Host: true, SubClient: true, PlayerColour: color.RGBA{R: 1, G: 2, B: 3, A: 4}},
			{ActionType: protocol.PlayerListActionAdd, UUID: uuid.MustParse("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"), EntityUniqueID: 67890, Username: "untrusted", XUID: "987654321", PlatformChatID: "platform-b", BuildPlatform: 13, Skin: validDummySkin(false), PlayerColour: color.RGBA{R: 5, G: 6, B: 7, A: 8}},
		}},
		"full/update_abilities": &packet.UpdateAbilities{AbilityData: protocol.AbilityData{
			EntityUniqueID: 123, PlayerPermissions: 2, CommandPermissions: 1,
			Layers: []protocol.AbilityLayer{{Type: protocol.AbilityLayerTypeBase, Abilities: 3, Values: 1, FlySpeed: .1, VerticalFlySpeed: .2, WalkSpeed: .3}},
		}},
		"empty/hurt_armour":                    &packet.HurtArmour{},
		"full/hurt_armour":                     &packet.HurtArmour{Cause: -3, Damage: 7, ArmourSlots: 9},
		"empty/player_location":                &packet.PlayerLocation{Type: packet.PlayerLocationTypeHide},
		"full/player_location":                 &packet.PlayerLocation{Type: packet.PlayerLocationTypeCoordinates, EntityUniqueID: -44, Position: [3]float32{1, 2, 3}},
		"empty/player_update_entity_overrides": &packet.PlayerUpdateEntityOverrides{},
		"full/player_update_entity_overrides": &packet.PlayerUpdateEntityOverrides{
			EntityUniqueID: 55, PropertyIndex: 3, Type: packet.PlayerUpdateEntityOverridesTypeInt, IntValue: -19,
		},
		"empty/primitive_shapes": &packet.PrimitiveShapes{},
		"full/primitive_shapes": &packet.PrimitiveShapes{Shapes: []protocol.PrimitiveShape{{
			NetworkID: 8, Type: protocol.Option(uint8(protocol.PrimitiveShapeText)), Location: protocol.Option(mgl32.Vec3{1, 2, 3}),
			Scale: protocol.Option(float32(2)), Rotation: protocol.Option(mgl32.Vec3{4, 5, 6}), TotalTimeLeft: protocol.Option(float32(7)),
			MaxRenderDistance: protocol.Option(float32(64)), Colour: protocol.Option(color.RGBA{R: 1, G: 2, B: 3, A: 4}),
			DimensionID: protocol.Option(int32(-1)), AttachedToEntityID: protocol.Option(int64(-99)),
			ExtraShapeData: &protocol.TextShape{Text: "test", UseRotation: true, BackgroundColour: protocol.Option(color.RGBA{R: 5, G: 6, B: 7, A: 8}), DepthTest: true},
		}}},
		"empty/server_presence_info": &packet.ServerPresenceInfo{},
		"full/server_presence_info": &packet.ServerPresenceInfo{PresenceInfo: protocol.Option(protocol.PresenceInfo{
			RichPresenceID: protocol.Option("presence"),
		})},
	}
	translated := New1001()
	for name, fixture := range fixtures {
		t.Run(name, func(t *testing.T) {
			got := marshalThroughProtocol(t, translated, fixture)
			if len(got) != 1 || !bytes.Equal(reference[name], got[0]) {
				t.Fatalf("protocol 1001 output differs\nnative:     %x\ntranslated: %x", reference[name], got)
			}
			assertProtocolRoundTrip(t, translated, fixture.ID(), [][]byte{reference[name]})
		})
	}
	t.Run("full/player_list_trusted_skins_decode", func(t *testing.T) {
		expected := reference["full/player_list_trusted_skins"]
		constructor := translated.Packets(false)[packet.IDPlayerList]
		wirePacket := constructor()
		readerBuffer := bytes.NewBuffer(expected)
		wirePacket.Marshal(translated.NewReader(readerBuffer, 0, true))
		if readerBuffer.Len() != 0 {
			t.Fatalf("trusted-skin PlayerList left %d bytes unread", readerBuffer.Len())
		}
		latest := translated.ConvertToLatest(wirePacket, nil)
		if len(latest) != 1 {
			t.Fatalf("trusted-skin PlayerList decoded to %d packets", len(latest))
		}
		entries := latest[0].(*packet.PlayerList).Entries
		expectedEntries := fixtures["full/player_list_trusted_skins"].(*packet.PlayerList).Entries
		if !reflect.DeepEqual(entries, expectedEntries) {
			t.Fatalf("fully populated PlayerList decoded incorrectly\nexpected: %#v\nactual:   %#v", expectedEntries, entries)
		}
	})
	for _, test := range []struct {
		name          string
		subChunkCount uint32
		limit         int32
		hasLimit      bool
	}{
		{"full/level_chunk_cached", 2, 0, false},
		{"full/level_chunk_limited", ^uint32(0) - 1, 37, true},
		{"full/level_chunk_limitless", ^uint32(0), -1, true},
	} {
		t.Run(test.name+"_decode", func(t *testing.T) {
			wirePacket := translated.Packets(false)[packet.IDLevelChunk]()
			readerBuffer := bytes.NewBuffer(reference[test.name])
			wirePacket.Marshal(translated.NewReader(readerBuffer, 0, true))
			if readerBuffer.Len() != 0 {
				t.Fatalf("LevelChunk left %d bytes unread", readerBuffer.Len())
			}
			latest := translated.ConvertToLatest(wirePacket, nil)
			if len(latest) != 1 {
				t.Fatalf("LevelChunk decoded to %d packets", len(latest))
			}
			decoded := latest[0].(*packet.LevelChunk)
			if decoded.SubChunkCount != test.subChunkCount {
				t.Fatalf("SubChunkCount decoded as %d, expected %d", decoded.SubChunkCount, test.subChunkCount)
			}
			limit, hasLimit := decoded.SubChunkLimit.Value()
			if hasLimit != test.hasLimit || hasLimit && limit != test.limit {
				t.Fatalf("SubChunkLimit decoded as (%d, %t), expected (%d, %t)", limit, hasLimit, test.limit, test.hasLimit)
			}
			expected := fixtures[test.name].(*packet.LevelChunk)
			if decoded.Position != expected.Position || decoded.Dimension != expected.Dimension || decoded.CacheEnabled != expected.CacheEnabled ||
				!reflect.DeepEqual(decoded.BlobHashes, expected.BlobHashes) || !bytes.Equal(decoded.RawPayload, expected.RawPayload) {
				t.Fatalf("fully populated LevelChunk decoded incorrectly\nexpected: %#v\nactual:   %#v", expected, decoded)
			}
		})
	}
	t.Run("full/set_actor_data_decode", func(t *testing.T) {
		constructor, ok := translated.Packets(false)[packet.IDSetActorData]
		if !ok {
			constructor, ok = translated.Packets(true)[packet.IDSetActorData]
		}
		if !ok {
			t.Fatal("protocol 1001 has no SetActorData constructor")
		}
		wirePacket := constructor()
		readerBuffer := bytes.NewBuffer(reference["full/set_actor_data"])
		wirePacket.Marshal(translated.NewReader(readerBuffer, 0, true))
		if readerBuffer.Len() != 0 {
			t.Fatalf("SetActorData left %d bytes unread", readerBuffer.Len())
		}
		latest := translated.ConvertToLatest(wirePacket, nil)
		if len(latest) != 1 {
			t.Fatalf("SetActorData decoded to %d packets", len(latest))
		}
		decoded := latest[0].(*packet.SetActorData)
		expected := fixtures["full/set_actor_data"].(*packet.SetActorData)
		if !reflect.DeepEqual(decoded, expected) {
			t.Fatalf("fully populated SetActorData decoded incorrectly\nexpected: %#v\nactual:   %#v", expected, decoded)
		}
	})
	t.Run("full/command_parameter_float", func(t *testing.T) {
		expected := reference["type/command_parameter_float"]
		parameter := protocol.CommandParameter{
			Name: "value", Type: protocol.CommandArgValid | protocol.CommandArgTypeFloat, Optional: true, Options: 3,
		}
		var encoded bytes.Buffer
		legacyproto.MarshalCommandParameter(translated.NewWriter(&encoded, 0), &parameter)
		if !bytes.Equal(expected, encoded.Bytes()) {
			t.Fatalf("protocol 1001 command parameter differs\nnative:     %x\ntranslated: %x", expected, encoded.Bytes())
		}
		var decoded protocol.CommandParameter
		readerBuffer := bytes.NewBuffer(expected)
		legacyproto.MarshalCommandParameter(translated.NewReader(readerBuffer, 0, true), &decoded)
		if readerBuffer.Len() != 0 {
			t.Fatalf("protocol 1001 command parameter left %d bytes unread", readerBuffer.Len())
		}
		var roundTrip bytes.Buffer
		legacyproto.MarshalCommandParameter(translated.NewWriter(&roundTrip, 0), &decoded)
		if !bytes.Equal(expected, roundTrip.Bytes()) {
			t.Fatalf("protocol 1001 command parameter round trip changed bytes")
		}
	})
	for name, expected := range reference {
		populated := strings.HasPrefix(name, "populated/")
		if !strings.HasPrefix(name, "all/") && !populated {
			continue
		}
		parts := strings.Split(name, "/")
		listener := parts[1] == "client"
		id, err := strconv.ParseUint(parts[2], 10, 32)
		if err != nil {
			t.Fatalf("parse native packet reference %q: %v", name, err)
		}
		constructor, ok := minecraft.DefaultProtocol.Packets(listener)[uint32(id)]
		if !ok {
			t.Errorf("%s: current gophertunnel has no packet ID %d", name, id)
			continue
		}
		t.Run(name, func(t *testing.T) {
			pk := constructor()
			initialiseCurrentEmptyPacket(pk, legacyproto.ID1001)
			if populated {
				populateCollectionFields(reflect.ValueOf(pk), false, 0)
				normalisePopulatedPacket(pk, legacyproto.ID1001)
			}
			got := marshalCompatibilityPacket(t, translated, pk)
			if !bytes.Equal(expected, got) {
				t.Fatalf("protocol 1001 packet differs\nnative:     %x\ntranslated: %x", expected, got)
			}
			assertProtocolRoundTrip(t, translated, uint32(id), [][]byte{expected})
		})
	}
}

func TestHistoricalProtocolMatrixMatchesNativeGophertunnel(t *testing.T) {
	tests := []struct {
		id       int32
		protocol *Protocol
		tags     string
	}{
		{800, New800(), "historical_optional_biome_id"},
		{818, New818(), "historical_entity_flags,historical_optional_biome_id,flat_server_shapes"},
		{819, New819(), "historical_entity_flags,historical_optional_biome_id,flat_server_shapes"},
		{827, New827(), "historical_entity_flags,flat_server_shapes"},
		{844, New844(), "historical_entity_flags,flat_server_shapes"},
		{859, New859(), "historical_entity_flags,flat_debug_shapes"},
		{860, New860(), "historical_entity_flags,historical_shapes"},
		{898, New898(), "historical_entity_flags,historical_shapes"},
		{924, New924(), "historical_entity_flags,historical_shapes,legacy_camera_spline"},
		{944, New944(), "historical_entity_flags,historical_shapes,current_camera_spline"},
		{975, New975(), "historical_entity_flags,modern_shapes,current_camera_spline"},
	}
	for _, test := range tests {
		t.Run(strconv.Itoa(int(test.id)), func(t *testing.T) {
			reference := nativeHistoricalReference(t, test.id, test.tags)
			for name, expected := range reference {
				parts := strings.Split(name, "/")
				if len(parts) != 4 {
					t.Fatalf("invalid native matrix key %q", name)
				}
				populated := parts[0] == "populated"
				listener := parts[1] == "client"
				id, err := strconv.ParseUint(parts[2], 10, 32)
				if err != nil {
					t.Fatalf("parse native packet reference %q: %v", name, err)
				}
				constructor, ok := minecraft.DefaultProtocol.Packets(listener)[uint32(id)]
				if !ok {
					constructor, ok = test.protocol.Packets(listener)[uint32(id)]
				}
				if !ok {
					t.Fatalf("compatibility pool has no packet matching native protocol %d %s", test.id, name)
				}
				t.Run(name, func(t *testing.T) {
					pk := constructor()
					initialiseCurrentEmptyPacket(pk, test.id)
					if populated {
						populateCollectionFields(reflect.ValueOf(pk), false, 0)
						normalisePopulatedPacket(pk, test.id)
					}
					got := marshalCompatibilityPacket(t, test.protocol, pk)
					if !bytes.Equal(expected, got) {
						t.Fatalf("protocol %d packet differs\nnative:     %x\ntranslated: %x", test.id, expected, got)
					}
					assertProtocolRoundTrip(t, test.protocol, uint32(id), [][]byte{expected})
				})
			}
		})
	}
}

func nativeHistoricalReference(t *testing.T, protocolID int32, tags string) map[string][]byte {
	t.Helper()
	dir := filepath.Join("testdata", "protocolmatrix")
	args := []string{"run"}
	if tags != "" {
		args = append(args, "-tags", tags)
	}
	args = append(args, "-modfile=go."+strconv.Itoa(int(protocolID))+".mod", "-mod=readonly", ".", strconv.Itoa(int(protocolID)))
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOCACHE="+filepath.Join(os.TempDir(), "legacy-version-matrix-cache"))
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			t.Fatalf("run native gophertunnel protocol %d reference: %v\n%s", protocolID, err, exitErr.Stderr)
		}
		t.Fatalf("run native gophertunnel protocol %d reference: %v", protocolID, err)
	}
	var values map[string]string
	if err := json.Unmarshal(output, &values); err != nil {
		t.Fatalf("decode native gophertunnel protocol %d reference: %v", protocolID, err)
	}
	decoded := make(map[string][]byte, len(values))
	for name, value := range values {
		data, err := hex.DecodeString(value)
		if err != nil {
			t.Fatalf("decode native protocol %d %s bytes: %v", protocolID, name, err)
		}
		decoded[name] = data
	}
	return decoded
}

func assertProtocolRoundTrip(t *testing.T, p minecraft.Protocol, packetID uint32, encoded [][]byte) {
	t.Helper()
	if len(encoded) != 1 {
		t.Fatalf("round-trip helper requires one packet, got %d", len(encoded))
	}
	constructor, ok := p.Packets(false)[packetID]
	if !ok {
		constructor, ok = p.Packets(true)[packetID]
	}
	if !ok {
		t.Fatalf("packet ID %d is absent from both protocol pools", packetID)
	}
	wirePacket := constructor()
	readerBuffer := bytes.NewBuffer(encoded[0])
	wirePacket.Marshal(p.NewReader(readerBuffer, 0, true))
	if readerBuffer.Len() != 0 {
		t.Fatalf("translation decoder left %d bytes unread", readerBuffer.Len())
	}
	var roundTrip bytes.Buffer
	wirePacket.Marshal(p.NewWriter(&roundTrip, 0))
	if !bytes.Equal(encoded[0], roundTrip.Bytes()) {
		t.Fatalf("translation decode/encode changed bytes\nfirst:  %x\nsecond: %x", encoded[0], roundTrip.Bytes())
	}
}

func marshalThroughProtocol(t *testing.T, p minecraft.Protocol, pk packet.Packet) [][]byte {
	t.Helper()
	converted := p.ConvertFromLatest(pk, nil)
	encoded := make([][]byte, len(converted))
	for i, convertedPacket := range converted {
		var buf bytes.Buffer
		convertedPacket.Marshal(p.NewWriter(&buf, 0))
		encoded[i] = append([]byte(nil), buf.Bytes()...)
	}
	return encoded
}

func native1001Reference(t *testing.T) map[string][]byte {
	t.Helper()
	dir := filepath.Join("testdata", "protocol1001")
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOCACHE="+filepath.Join(os.TempDir(), "legacy-version-reference-go-cache"))
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("run native gophertunnel 1001 reference: %v", err)
	}
	var values map[string]string
	if err := json.Unmarshal(output, &values); err != nil {
		t.Fatalf("decode native gophertunnel 1001 reference: %v", err)
	}
	decoded := make(map[string][]byte, len(values))
	for name, value := range values {
		b, err := hex.DecodeString(value)
		if err != nil {
			t.Fatalf("decode %s native bytes: %v", name, err)
		}
		decoded[name] = b
	}
	return decoded
}

func TestStartGame2168MarshalsLikeLatest(t *testing.T) {
	var native, compatible bytes.Buffer
	latest := &packet.StartGame{}
	translated := &packet.StartGame{}
	latest.Marshal(protocol.NewWriter(&native, 0))
	legacypacket.StartGame(
		legacyproto.NewWriter(protocol.NewWriter(&compatible, 0), legacyproto.ID2168, 0),
		translated,
		nil,
	)
	if !bytes.Equal(native.Bytes(), compatible.Bytes()) {
		t.Fatalf("1.26.40 StartGame marshal differs from latest\nnative:     %x\ncompatible: %x", native.Bytes(), compatible.Bytes())
	}
}

func TestLegacyUnionPacketsAreSplitAndTranslated(t *testing.T) {
	p := New1001()
	playerPackets := p.downgradePackets([]packet.Packet{&packet.PlayerList{Entries: []protocol.PlayerListEntry{
		{ActionType: protocol.PlayerListActionAdd},
		{ActionType: protocol.PlayerListActionRemove},
	}}}, nil)
	if len(playerPackets) != 2 {
		t.Fatalf("expected mixed PlayerList to split into 2 packets, got %d", len(playerPackets))
	}
	for _, pk := range playerPackets {
		if _, ok := pk.(*translatedPacket); !ok {
			t.Fatalf("split PlayerList packet was not compatibility wrapped: %T", pk)
		}
	}

	scorePackets := p.downgradePackets([]packet.Packet{&packet.SetScore{Entries: []protocol.ScoreboardEntry{
		{IdentityType: protocol.ScoreboardIdentityFakePlayer},
		{IdentityType: protocol.ScoreboardIdentityRemove},
	}}}, nil)
	if len(scorePackets) != 2 {
		t.Fatalf("expected mixed SetScore to split into 2 packets, got %d", len(scorePackets))
	}
	for _, pk := range scorePackets {
		if _, ok := pk.(*translatedPacket); !ok {
			t.Fatalf("split SetScore packet was not compatibility wrapped: %T", pk)
		}
	}
}

func TestServerPlayerPostMovePositionVersionBoundary(t *testing.T) {
	pk := &packet.ServerPlayerPostMovePosition{Position: mgl32.Vec3{1, 2, 3}}
	if got := New1001().ConvertFromLatest(pk, nil); len(got) != 0 {
		t.Fatalf("protocol 1001 retained 2168-only packet: %T", got[0])
	}
	current := &Protocol{id: legacyproto.ID2168, ver: "1.26.40"}
	if got := current.ConvertFromLatest(pk, nil); len(got) != 1 || got[0] != pk {
		t.Fatalf("protocol 2168 did not retain ServerPlayerPostMovePosition")
	}
}

func TestProtocol1001RepresentativeRoundTrips(t *testing.T) {
	inputFlags := protocol.NewInputFlags(packet.InputFlagCount)
	inputFlags.Set(packet.InputFlagPerformItemStackRequest)
	tests := []packet.Packet{
		&packet.ClientBoundMapItemData{
			Scale:          protocol.Option(uint8(2)),
			MapsIncludedIn: protocol.Option([]int64{4, 8}),
		},
		&packet.CraftingData{ShapelessRecipes: []protocol.ShapelessRecipe{{
			RecipeID: "test", Input: []protocol.ItemDescriptorCount{{
				Descriptor: &protocol.DefaultItemDescriptor{Name: "minecraft:stone"}, Count: 1,
			}}, UnlockRequirement: protocol.Option(protocol.RecipeUnlockRequirement{Context: protocol.RecipeUnlockContextAlwaysUnlocked}),
		}}},
		&packet.InventoryContent{Content: []protocol.ItemInstance{{Stack: protocol.ItemStack{ItemType: protocol.ItemType{NetworkID: 1}, Count: 1}}}},
		&packet.LevelChunk{SubChunkLimit: protocol.Option(int32(12))},
		&packet.MovePlayer{Mode: packet.MoveModeTeleport, TeleportData: protocol.Option(protocol.TeleportData{})},
		&packet.PlayerAuthInput{InputData: inputFlags, ItemStackRequest: protocol.Option(protocol.ItemStackRequest{})},
		&packet.PlayerList{Entries: []protocol.PlayerListEntry{{ActionType: protocol.PlayerListActionAdd}}},
		&packet.ResourcePackClientResponse{Response: packet.PackResponseSendPacks, PacksToDownload: []string{"pack_1.0.0"}},
		&packet.SetScore{Entries: []protocol.ScoreboardEntry{{IdentityType: protocol.ScoreboardIdentityFakePlayer, DisplayName: "line"}}},
		&packet.SubChunk{CacheEnabled: true, SubChunkEntries: []protocol.SubChunkEntry{{
			Result: protocol.SubChunkResultSuccess, RawPayload: protocol.Option([]byte{1, 2}), BlobHash: protocol.Option(uint64(3)),
		}}},
	}
	for _, original := range tests {
		t.Run(packetName(original), func(t *testing.T) {
			marshalFn := packets[original.ID()]
			var first bytes.Buffer
			marshalFn(legacyproto.NewWriter(protocol.NewWriter(&first, 0), legacyproto.ID1001, 0), original)

			decoded := reflect.New(reflect.TypeOf(original).Elem()).Interface().(packet.Packet)
			readerBuffer := bytes.NewBuffer(first.Bytes())
			reader := legacyproto.NewReader(protocol.NewReader(readerBuffer, 0, true), legacyproto.ID1001, 0, true)
			marshalFn(reader, decoded)
			if readerBuffer.Len() != 0 {
				t.Fatalf("decoder left %d bytes unread", readerBuffer.Len())
			}

			var second bytes.Buffer
			marshalFn(legacyproto.NewWriter(protocol.NewWriter(&second, 0), legacyproto.ID1001, 0), decoded)
			if !bytes.Equal(first.Bytes(), second.Bytes()) {
				t.Fatalf("1.26.30 round trip changed bytes\nfirst:  %x\nsecond: %x", first.Bytes(), second.Bytes())
			}
		})
	}
}

func packetName(pk packet.Packet) string {
	return fmt.Sprintf("%T", pk)
}
