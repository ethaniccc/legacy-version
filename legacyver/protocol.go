package legacyver

import (
	"strings"

	"github.com/ethaniccc/legacy-version/legacyver/legacypacket"
	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/samber/lo"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var (
	packetPoolClient packet.Pool
	packetPoolServer packet.Pool
)

func init() {
	packetPoolClient = packet.NewClientPool()
	packetPoolServer = packet.NewServerPool()

	for pkId, cur := range packetPoolClient {
		packetPoolClient[pkId] = convertPacketFunc(pkId, cur)
	}

	for pkId, cur := range packetPoolServer {
		packetPoolServer[pkId] = convertPacketFunc(pkId, cur)
	}
}

func convertPacketFunc(pid uint32, cur func() packet.Packet) func() packet.Packet {
	switch pid {
	case packet.IDDimensionData:
		return func() packet.Packet { return &legacypacket.DimensionData{} }
	case packet.IDCameraAimAssist:
		return func() packet.Packet { return &legacypacket.CameraAimAssist{} }
	case packet.IDCameraPresets:
		return func() packet.Packet { return &legacypacket.CameraPresets{} }
	case packet.IDContainerRegistryCleanup:
		return func() packet.Packet { return &legacypacket.ContainerRegistryCleanup{} }
	case packet.IDEmote:
		return func() packet.Packet { return &legacypacket.Emote{} }
	case packet.IDInventoryContent:
		return func() packet.Packet { return &legacypacket.InventoryContent{} }
	case packet.IDInventorySlot:
		return func() packet.Packet { return &legacypacket.InventorySlot{} }
	case packet.IDItemStackResponse:
		return func() packet.Packet { return &legacypacket.ItemStackResponse{} }
	case packet.IDMobEffect:
		return func() packet.Packet { return &legacypacket.MobEffect{} }
	case packet.IDPlayerAuthInput:
		return func() packet.Packet { return &legacypacket.PlayerAuthInput{} }
	case packet.IDResourcePacksInfo:
		return func() packet.Packet { return &legacypacket.ResourcePacksInfo{} }
	case packet.IDTransfer:
		return func() packet.Packet { return &legacypacket.Transfer{} }
	case packet.IDUpdateAttributes:
		return func() packet.Packet { return &legacypacket.UpdateAttributes{} }
	case packet.IDAddPlayer:
		return func() packet.Packet { return &legacypacket.AddPlayer{} }
	case packet.IDAddActor:
		return func() packet.Packet { return &legacypacket.AddActor{} }
	case packet.IDSetActorLink:
		return func() packet.Packet { return &legacypacket.SetActorLink{} }
	case packet.IDCameraInstruction:
		return func() packet.Packet { return &legacypacket.CameraInstruction{} }
	case packet.IDChangeDimension:
		return func() packet.Packet { return &legacypacket.ChangeDimension{} }
	case packet.IDCorrectPlayerMovePrediction:
		return func() packet.Packet { return &legacypacket.CorrectPlayerMovePrediction{} }
	case packet.IDDisconnect:
		return func() packet.Packet { return &legacypacket.Disconnect{} }
	case packet.IDEditorNetwork:
		return func() packet.Packet { return &legacypacket.EditorNetwork{} }
	case packet.IDMobArmourEquipment:
		return func() packet.Packet { return &legacypacket.MobArmourEquipment{} }
	case packet.IDPlayerArmourDamage:
		return func() packet.Packet { return &legacypacket.PlayerArmourDamage{} }
	case packet.IDSetTitle:
		return func() packet.Packet { return &legacypacket.SetTitle{} }
	case packet.IDStopSound:
		return func() packet.Packet { return &legacypacket.StopSound{} }
	case packet.IDInventoryTransaction:
		return func() packet.Packet { return &legacypacket.InventoryTransaction{} }
	case packet.IDItemStackRequest:
		return func() packet.Packet { return &legacypacket.ItemStackRequest{} }
	case packet.IDCraftingData:
		return func() packet.Packet { return &legacypacket.CraftingData{} }
	case packet.IDContainerClose:
		return func() packet.Packet { return &legacypacket.ContainerClose{} }
	case packet.IDText:
		return func() packet.Packet { return &legacypacket.Text{} }
	case packet.IDBookEdit:
		return func() packet.Packet { return &legacypacket.BookEdit{} }
	case packet.IDStartGame:
		return func() packet.Packet { return &legacypacket.StartGame{} }
	case packet.IDCodeBuilderSource:
		return func() packet.Packet { return &legacypacket.CodeBuilderSource{} }
	case packet.IDItemRegistry:
		return func() packet.Packet { return &legacypacket.ItemRegistry{} }
	case packet.IDStructureBlockUpdate:
		return func() packet.Packet { return &legacypacket.StructureBlockUpdate{} }
	case packet.IDBossEvent:
		return func() packet.Packet { return &legacypacket.BossEvent{} }
	case packet.IDCameraAimAssistPresets:
		return func() packet.Packet { return &legacypacket.CameraAimAssistPresets{} }
	case packet.IDCommandBlockUpdate:
		return func() packet.Packet { return &legacypacket.CommandBlockUpdate{} }
	case packet.IDCreativeContent:
		return func() packet.Packet { return &legacypacket.CreativeContent{} }
	case packet.IDUpdateAbilities:
		return func() packet.Packet { return &legacypacket.UpdateAbilities{} }
	case packet.IDClientCheatAbility:
		return func() packet.Packet { return &legacypacket.ClientCheatAbility{} }
	case packet.IDClientMovementPredictionSync:
		return func() packet.Packet { return &legacypacket.ClientMovementPredictionSync{} }
	case packet.IDLevelSoundEvent:
		return func() packet.Packet { return &legacypacket.LevelSoundEvent{} }
	case packet.IDSetHud:
		return func() packet.Packet { return &legacypacket.SetHud{} }
	case packet.IDResourcePackStack:
		return func() packet.Packet { return &legacypacket.ResourcePackStack{} }
	case packet.IDUpdatePlayerGameType:
		return func() packet.Packet { return &legacypacket.UpdatePlayerGameType{} }
	case packet.IDSetActorMotion:
		return func() packet.Packet { return &legacypacket.SetActorMotion{} }
	case packet.IDBiomeDefinitionList:
		return func() packet.Packet { return &legacypacket.BiomeDefinitionList{} }
	case packet.IDPlayerList:
		return func() packet.Packet { return &legacypacket.PlayerList{} }
	case packet.IDSubChunk:
		return func() packet.Packet { return &legacypacket.SubChunk{} }
	case packet.IDGameRulesChanged:
		return func() packet.Packet { return &legacypacket.GameRulesChanged{} }
	case packet.IDAnimate:
		return func() packet.Packet { return &legacypacket.Animate{} }
	case packet.IDAvailableCommands:
		return func() packet.Packet { return &legacypacket.AvailableCommands{} }
	case packet.IDCommandOutput:
		return func() packet.Packet { return &legacypacket.CommandOutput{} }
	case packet.IDCommandRequest:
		return func() packet.Packet { return &legacypacket.CommandRequest{} }
	case packet.IDEvent:
		return func() packet.Packet { return &legacypacket.Event{} }
	case packet.IDInteract:
		return func() packet.Packet { return &legacypacket.Interact{} }
	default:
		return cur
	}
}

type Protocol struct {
	ver string
	id  int32

	blockTranslator BlockTranslator
	itemTranslator  ItemTranslator
	Items           []protocol.ItemEntry
}

func (p *Protocol) Ver() string {
	return p.ver
}

func (p *Protocol) ID() int32 {
	return p.id
}

func (p *Protocol) Packets(listener bool) packet.Pool {
	if listener {
		return packetPoolClient
	}
	return packetPoolServer
}

func (p *Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return proto.NewReader(protocol.NewReader(r, shieldID, enableLimits), p.id)
}

func (p *Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return proto.NewWriter(protocol.NewWriter(w, shieldID), p.id)
}

func (p *Protocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	return p.blockTranslator.UpgradeBlockPackets(
		p.itemTranslator.UpgradeItemPackets(p.upgradePackets([]packet.Packet{pk}, conn), conn),
		conn)
}

func (p *Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	return p.downgradePackets(p.blockTranslator.DowngradeBlockPackets(
		p.itemTranslator.DowngradeItemPackets([]packet.Packet{pk}, conn),
		conn), conn)
}

func (p *Protocol) downgradePackets(pks []packet.Packet, conn *minecraft.Conn) []packet.Packet {
	translator, ok := p.blockTranslator.(*DefaultBlockTranslator)
	if !ok {
		return pks
	}
	for pkIndex, pk := range pks {
		switch pk := pk.(type) {
		case *packet.DimensionData:
			translator.dimensionDefinitions = pk.Definitions
		case *packet.ClientCacheStatus:
			// pk.Enabled = false // TODO: enable when chunk translation is not broken
		case *packet.SetActorMotion:
			pks[pkIndex] = &legacypacket.SetActorMotion{
				EntityRuntimeID: pk.EntityRuntimeID,
				Velocity:        pk.Velocity,
				Tick:            pk.Tick,
			}
		case *packet.ResourcePackStack:
			pks[pkIndex] = &legacypacket.ResourcePackStack{
				TexturePackRequired:          pk.TexturePackRequired,
				TexturePacks:                 pk.TexturePacks,
				BaseGameVersion:              pk.BaseGameVersion,
				Experiments:                  pk.Experiments,
				ExperimentsPreviouslyToggled: pk.ExperimentsPreviouslyToggled,
				IncludeEditorPacks:           pk.IncludeEditorPacks,
			}
		case *packet.UpdatePlayerGameType:
			pks[pkIndex] = &legacypacket.UpdatePlayerGameType{
				GameType:       pk.GameType,
				PlayerUniqueID: pk.PlayerUniqueID,
				Tick:           pk.Tick,
			}
		case *packet.CameraPresets:
			presets := make([]proto.CameraPreset, len(pk.Presets))
			for i, p := range pk.Presets {
				presets[i] = (&proto.CameraPreset{}).FromLatest(p)
			}
			pks[pkIndex] = &legacypacket.CameraPresets{
				Presets: presets,
			}
		case *packet.PlayerAuthInput:
			inputData := pk.InputData
			if p.ID() < proto.ID766 {
				inputData = fitBitset(inputData, 64)
			}
			pks[pkIndex] = &legacypacket.PlayerAuthInput{
				Pitch:                  pk.Pitch,
				Yaw:                    pk.Yaw,
				Position:               pk.Position,
				MoveVector:             pk.MoveVector,
				HeadYaw:                pk.HeadYaw,
				InputData:              inputData,
				InputMode:              pk.InputMode,
				PlayMode:               pk.PlayMode,
				InteractionModel:       pk.InteractionModel,
				InteractPitch:          pk.InteractPitch,
				InteractYaw:            pk.InteractYaw,
				Tick:                   pk.Tick,
				Delta:                  pk.Delta,
				ItemInteractionData:    pk.ItemInteractionData,
				ItemStackRequest:       (&proto.ItemStackRequest{}).FromLatest(pk.ItemStackRequest),
				BlockActions:           pk.BlockActions,
				VehicleRotation:        pk.VehicleRotation,
				ClientPredictedVehicle: pk.ClientPredictedVehicle,
				AnalogueMoveVector:     pk.AnalogueMoveVector,
				CameraOrientation:      pk.CameraOrientation,
				RawMoveVector:          pk.RawMoveVector,
			}
		case *packet.ItemStackResponse:
			responses := make([]proto.ItemStackResponse, len(pk.Responses))
			for i, r := range pk.Responses {
				responses[i] = (&proto.ItemStackResponse{}).FromLatest(r)
			}
			pks[pkIndex] = &legacypacket.ItemStackResponse{Responses: responses}
		case *packet.ResourcePacksInfo:
			texturePacks := make([]proto.TexturePackInfo, len(pk.TexturePacks))
			packURLs := make([]protocol.PackURL, 0)
			for i, t := range pk.TexturePacks {
				texturePacks[i] = (&proto.TexturePackInfo{}).FromLatest(t)
				if t.DownloadURL != "" {
					packURLs = append(packURLs, protocol.PackURL{
						UUIDVersion: t.UUID.String() + "_" + t.Version,
						URL:         t.DownloadURL,
					})
				}
			}
			pks[pkIndex] = &legacypacket.ResourcePacksInfo{
				TexturePackRequired:        pk.TexturePackRequired,
				HasAddons:                  pk.HasAddons,
				HasScripts:                 pk.HasScripts,
				ForceDisableVibrantVisuals: pk.ForceDisableVibrantVisuals,
				WorldTemplateUUID:          pk.WorldTemplateUUID,
				WorldTemplateVersion:       pk.WorldTemplateVersion,
				TexturePacks:               texturePacks,
				PackURLs:                   packURLs,
			}
		case *packet.InventorySlot:
			pks[pkIndex] = &legacypacket.InventorySlot{
				WindowID:             pk.WindowID,
				Slot:                 pk.Slot,
				Container:            (&proto.FullContainerName{}).FromLatest(pk.Container),
				DynamicContainerSize: 0,
				StorageItem:          pk.StorageItem,
				NewItem:              pk.NewItem,
			}
		case *packet.InventoryContent:
			pks[pkIndex] = &legacypacket.InventoryContent{
				WindowID:             pk.WindowID,
				Content:              pk.Content,
				Container:            (&proto.FullContainerName{}).FromLatest(pk.Container),
				DynamicContainerSize: 0,
				StorageItem:          pk.StorageItem,
			}
		case *packet.MobEffect:
			pks[pkIndex] = &legacypacket.MobEffect{
				EntityRuntimeID: pk.EntityRuntimeID,
				Operation:       pk.Operation,
				EffectType:      pk.EffectType,
				Amplifier:       pk.Amplifier,
				Particles:       pk.Particles,
				Duration:        pk.Duration,
				Tick:            pk.Tick,
				Ambient:         pk.Ambient,
			}
		case *packet.CameraAimAssist:
			pks[pkIndex] = &legacypacket.CameraAimAssist{
				Preset:          pk.Preset,
				Angle:           pk.Angle,
				Distance:        pk.Distance,
				TargetMode:      pk.TargetMode,
				Action:          pk.Action,
				ShowDebugRender: pk.ShowDebugRender,
			}
		case *packet.UpdateAttributes:
			attributes := make([]proto.Attribute, len(pk.Attributes))
			for i, a := range pk.Attributes {
				attributes[i] = (&proto.Attribute{}).FromLatest(a)
			}
			pks[pkIndex] = &legacypacket.UpdateAttributes{
				EntityRuntimeID: pk.EntityRuntimeID,
				Attributes:      attributes,
				Tick:            pk.Tick,
			}
		case *packet.ContainerRegistryCleanup:
			removedContainers := make([]proto.FullContainerName, len(pk.RemovedContainers))
			for i, c := range pk.RemovedContainers {
				removedContainers[i] = (&proto.FullContainerName{}).FromLatest(c)
			}
			pks[pkIndex] = &legacypacket.ContainerRegistryCleanup{
				RemovedContainers: removedContainers,
			}
		case *packet.Emote:
			pks[pkIndex] = &legacypacket.Emote{
				EntityRuntimeID: pk.EntityRuntimeID,
				EmoteLength:     pk.EmoteLength,
				EmoteID:         pk.EmoteID,
				XUID:            pk.XUID,
				PlatformID:      pk.PlatformID,
				Flags:           pk.Flags,
			}
		case *packet.Transfer:
			pks[pkIndex] = &legacypacket.Transfer{
				Address:     pk.Address,
				Port:        pk.Port,
				ReloadWorld: pk.ReloadWorld,
			}
		case *packet.AddActor:
			links := make([]proto.EntityLink, len(pk.EntityLinks))
			for i, l := range pk.EntityLinks {
				links[i] = (&proto.EntityLink{}).FromLatest(l)
			}
			pks[pkIndex] = &legacypacket.AddActor{
				EntityUniqueID:   pk.EntityUniqueID,
				EntityRuntimeID:  pk.EntityRuntimeID,
				EntityType:       pk.EntityType,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				BodyYaw:          pk.BodyYaw,
				Attributes:       pk.Attributes,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				EntityLinks:      links,
			}
		case *packet.AddPlayer:
			links := make([]proto.EntityLink, len(pk.EntityLinks))
			for i, l := range pk.EntityLinks {
				links[i] = (&proto.EntityLink{}).FromLatest(l)
			}
			pks[pkIndex] = &legacypacket.AddPlayer{
				UUID:             pk.UUID,
				Username:         pk.Username,
				EntityRuntimeID:  pk.EntityRuntimeID,
				PlatformChatID:   pk.PlatformChatID,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				HeldItem:         pk.HeldItem,
				GameType:         pk.GameType,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				AbilityData:      (&proto.AbilityData{}).FromLatest(pk.AbilityData),
				EntityLinks:      links,
				DeviceID:         pk.DeviceID,
				BuildPlatform:    pk.BuildPlatform,
			}
		case *packet.SetActorLink:
			pks[pkIndex] = &legacypacket.SetActorLink{
				EntityLink: (&proto.EntityLink{}).FromLatest(pk.EntityLink),
			}
		case *packet.CameraInstruction:
			var iSet protocol.Optional[proto.CameraInstructionSet]
			if v, ok := pk.Set.Value(); ok {
				iSet = protocol.Option((&proto.CameraInstructionSet{}).FromLatest(v))
			}
			var spline protocol.Optional[proto.CameraSplineInstruction]
			if v, ok := pk.Spline.Value(); ok {
				spline = protocol.Option((&proto.CameraSplineInstruction{}).FromLatest(v))
			}
			pks[pkIndex] = &legacypacket.CameraInstruction{
				Set:              iSet,
				Clear:            pk.Clear,
				Fade:             pk.Fade,
				Target:           pk.Target,
				RemoveTarget:     pk.RemoveTarget,
				FieldOfView:      pk.FieldOfView,
				Spline:           spline,
				AttachToEntity:   pk.AttachToEntity,
				DetachFromEntity: pk.DetachFromEntity,
			}
		case *packet.ChangeDimension:
			translator.currentDimension = pk.Dimension
			pks[pkIndex] = &legacypacket.ChangeDimension{
				Dimension:       pk.Dimension,
				Position:        pk.Position,
				Respawn:         pk.Respawn,
				LoadingScreenID: pk.LoadingScreenID,
			}
		case *packet.CorrectPlayerMovePrediction:
			pks[pkIndex] = &legacypacket.CorrectPlayerMovePrediction{
				PredictionType:         pk.PredictionType,
				Position:               pk.Position,
				Delta:                  pk.Delta,
				Rotation:               pk.Rotation,
				VehicleAngularVelocity: pk.VehicleAngularVelocity,
				OnGround:               pk.OnGround,
				Tick:                   pk.Tick,
			}
		case *packet.Disconnect:
			pks[pkIndex] = &legacypacket.Disconnect{
				Reason:                  pk.Reason,
				HideDisconnectionScreen: pk.HideDisconnectionScreen,
				Message:                 pk.Message,
				FilteredMessage:         pk.FilteredMessage,
			}
		case *packet.EditorNetwork:
			pks[pkIndex] = &legacypacket.EditorNetwork{
				RouteToManager: pk.RouteToManager,
				Payload:        pk.Payload,
			}
		case *packet.MobArmourEquipment:
			pks[pkIndex] = &legacypacket.MobArmourEquipment{
				EntityRuntimeID: pk.EntityRuntimeID,
				Helmet:          pk.Helmet,
				Chestplate:      pk.Chestplate,
				Leggings:        pk.Leggings,
				Boots:           pk.Boots,
				Body:            pk.Body,
			}
		case *packet.PlayerArmourDamage:
			pks[pkIndex] = &legacypacket.PlayerArmourDamage{
				List: pk.List,
			}
		case *packet.SetTitle:
			pks[pkIndex] = &legacypacket.SetTitle{
				ActionType:       pk.ActionType,
				Text:             pk.Text,
				FadeInDuration:   pk.FadeInDuration,
				RemainDuration:   pk.RemainDuration,
				FadeOutDuration:  pk.FadeOutDuration,
				XUID:             pk.XUID,
				PlatformOnlineID: pk.PlatformOnlineID,
				FilteredMessage:  pk.FilteredMessage,
			}
		case *packet.StopSound:
			pks[pkIndex] = &legacypacket.StopSound{
				SoundName:       pk.SoundName,
				StopAll:         pk.StopAll,
				StopMusicLegacy: pk.StopMusicLegacy,
			}
		case *packet.InventoryTransaction:
			trData := pk.TransactionData
			if x, ok := trData.(*protocol.UseItemTransactionData); ok {
				trData = (&proto.UseItemTransactionData{}).FromLatest(x)
			}
			pks[pkIndex] = &legacypacket.InventoryTransaction{
				LegacyRequestID:    pk.LegacyRequestID,
				LegacySetItemSlots: pk.LegacySetItemSlots,
				Actions:            pk.Actions,
				TransactionData:    trData,
			}
		case *packet.ItemStackRequest:
			requests := make([]proto.ItemStackRequest, len(pk.Requests))
			for i, r := range pk.Requests {
				requests[i] = (&proto.ItemStackRequest{}).FromLatest(r)
			}
			pks[pkIndex] = &legacypacket.ItemStackRequest{Requests: requests}
		case *packet.Text:
			pks[pkIndex] = &legacypacket.Text{
				TextType:         pk.TextType,
				NeedsTranslation: pk.NeedsTranslation,
				SourceName:       pk.SourceName,
				Message:          pk.Message,
				Parameters:       pk.Parameters,
				XUID:             pk.XUID,
				PlatformChatID:   pk.PlatformChatID,
				FilteredMessage:  pk.FilteredMessage,
			}
		case *packet.BookEdit:
			pks[pkIndex] = &legacypacket.BookEdit{
				InventorySlot:       pk.InventorySlot,
				ActionType:          pk.ActionType,
				PageNumber:          pk.PageNumber,
				SecondaryPageNumber: pk.SecondaryPageNumber,
				Text:                pk.Text,
				PhotoName:           pk.PhotoName,
				Title:               pk.Title,
				Author:              pk.Author,
				XUID:                pk.XUID,
			}
		case *packet.ContainerClose:
			pks[pkIndex] = &legacypacket.ContainerClose{
				WindowID:      pk.WindowID,
				ContainerType: pk.ContainerType,
				ServerSide:    pk.ServerSide,
			}
		case *packet.CraftingData:
			recipes := make([]proto.Recipe, len(pk.Recipes))
			for i, r := range pk.Recipes {
				recipes[i] = proto.RecipeFromLatest(r)
			}
			pks[pkIndex] = &legacypacket.CraftingData{
				Recipes:                      recipes,
				PotionRecipes:                pk.PotionRecipes,
				PotionContainerChangeRecipes: pk.PotionContainerChangeRecipes,
				MaterialReducers:             pk.MaterialReducers,
				ClearRecipes:                 pk.ClearRecipes,
			}
		case *packet.StartGame:
			translator.currentDimension = pk.Dimension
			// Adjust game version
			pk.GameVersion = p.ver
			pk.BaseGameVersion = p.ver

			items := make([]proto.LegacyItemRegistryEntry, len(p.Items))
			for i, it := range p.Items {
				items[i] = (&proto.LegacyItemRegistryEntry{}).FromLatest(it)
			}
			items = p.itemTranslator.DowngradeLegacyItemRegistry(items)

			pks[pkIndex] = &legacypacket.StartGame{
				Items:                          items,
				EntityUniqueID:                 pk.EntityUniqueID,
				EntityRuntimeID:                pk.EntityRuntimeID,
				PlayerGameMode:                 pk.PlayerGameMode,
				PlayerPosition:                 pk.PlayerPosition,
				Pitch:                          pk.Pitch,
				Yaw:                            pk.Yaw,
				WorldSeed:                      pk.WorldSeed,
				SpawnBiomeType:                 pk.SpawnBiomeType,
				UserDefinedBiomeName:           pk.UserDefinedBiomeName,
				Dimension:                      pk.Dimension,
				Generator:                      pk.Generator,
				WorldGameMode:                  pk.WorldGameMode,
				Hardcore:                       pk.Hardcore,
				Difficulty:                     pk.Difficulty,
				WorldSpawn:                     pk.WorldSpawn,
				AchievementsDisabled:           pk.AchievementsDisabled,
				EditorWorldType:                pk.EditorWorldType,
				CreatedInEditor:                pk.CreatedInEditor,
				ExportedFromEditor:             pk.ExportedFromEditor,
				DayCycleLockTime:               pk.DayCycleLockTime,
				EducationEditionOffer:          pk.EducationEditionOffer,
				EducationFeaturesEnabled:       pk.EducationFeaturesEnabled,
				EducationProductID:             pk.EducationProductID,
				RainLevel:                      pk.RainLevel,
				LightningLevel:                 pk.LightningLevel,
				ConfirmedPlatformLockedContent: pk.ConfirmedPlatformLockedContent,
				MultiPlayerGame:                pk.MultiPlayerGame,
				LANBroadcastEnabled:            pk.LANBroadcastEnabled,
				XBLBroadcastMode:               pk.XBLBroadcastMode,
				PlatformBroadcastMode:          pk.PlatformBroadcastMode,
				CommandsEnabled:                pk.CommandsEnabled,
				TexturePackRequired:            pk.TexturePackRequired,
				GameRules:                      pk.GameRules,
				Experiments:                    pk.Experiments,
				ExperimentsPreviouslyToggled:   pk.ExperimentsPreviouslyToggled,
				BonusChestEnabled:              pk.BonusChestEnabled,
				StartWithMapEnabled:            pk.StartWithMapEnabled,
				PlayerPermissions:              pk.PlayerPermissions,
				ServerChunkTickRadius:          pk.ServerChunkTickRadius,
				HasLockedBehaviourPack:         pk.HasLockedBehaviourPack,
				HasLockedTexturePack:           pk.HasLockedTexturePack,
				FromLockedWorldTemplate:        pk.FromLockedWorldTemplate,
				MSAGamerTagsOnly:               pk.MSAGamerTagsOnly,
				FromWorldTemplate:              pk.FromWorldTemplate,
				WorldTemplateSettingsLocked:    pk.WorldTemplateSettingsLocked,
				OnlySpawnV1Villagers:           pk.OnlySpawnV1Villagers,
				PersonaDisabled:                pk.PersonaDisabled,
				CustomSkinsDisabled:            pk.CustomSkinsDisabled,
				EmoteChatMuted:                 pk.EmoteChatMuted,
				BaseGameVersion:                pk.BaseGameVersion,
				LimitedWorldWidth:              pk.LimitedWorldWidth,
				LimitedWorldDepth:              pk.LimitedWorldDepth,
				NewNether:                      pk.NewNether,
				EducationSharedResourceURI:     pk.EducationSharedResourceURI,
				ForceExperimentalGameplay:      pk.ForceExperimentalGameplay,
				LevelID:                        pk.LevelID,
				WorldName:                      pk.WorldName,
				TemplateContentIdentity:        pk.TemplateContentIdentity,
				Trial:                          pk.Trial,
				PlayerMovementSettings:         (&proto.PlayerMovementSettings{}).FromLatest(pk.PlayerMovementSettings),
				Time:                           pk.Time,
				EnchantmentSeed:                pk.EnchantmentSeed,
				Blocks:                         pk.Blocks,
				MultiPlayerCorrelationID:       pk.MultiPlayerCorrelationID,
				ServerAuthoritativeInventory:   pk.ServerAuthoritativeInventory,
				GameVersion:                    pk.GameVersion,
				PropertyData:                   pk.PropertyData,
				ServerBlockStateChecksum:       pk.ServerBlockStateChecksum,
				ClientSideGeneration:           pk.ClientSideGeneration,
				WorldTemplateID:                pk.WorldTemplateID,
				ChatRestrictionLevel:           pk.ChatRestrictionLevel,
				DisablePlayerInteractions:      pk.DisablePlayerInteractions,
				ServerID:                       pk.ServerID,
				WorldID:                        pk.WorldID,
				ScenarioID:                     pk.ScenarioID,
				OwnerID:                        pk.OwnerID,
				UseBlockNetworkIDHashes:        pk.UseBlockNetworkIDHashes,
				ServerAuthoritativeSound:       pk.ServerAuthoritativeSound,
			}
		case *packet.CodeBuilderSource:
			pks[pkIndex] = &legacypacket.CodeBuilderSource{
				Operation:  pk.Operation,
				Category:   pk.Category,
				CodeStatus: pk.CodeStatus,
			}
		case *packet.ItemRegistry:
			p.Items = pk.Items
			items := make([]proto.ItemEntry, len(pk.Items))
			for i, it := range pk.Items {
				items[i] = (&proto.ItemEntry{}).FromLatest(it)
			}
			items = p.itemTranslator.DowngradeItemEntries(items)

			if p.ID() < proto.ID776 {
				items = lo.Filter(items, func(item proto.ItemEntry, index int) bool {
					return !strings.HasPrefix(item.Name, "minecraft:") && item.ComponentBased
				}) // only custom items here

				if len(items) < 1 {
					return []packet.Packet{}
				}
			}
			pks[pkIndex] = &legacypacket.ItemRegistry{Items: items}
		case *packet.StructureBlockUpdate:
			pks[pkIndex] = &legacypacket.StructureBlockUpdate{
				Position:              pk.Position,
				StructureName:         pk.StructureName,
				FilteredStructureName: pk.FilteredStructureName,
				DataField:             pk.DataField,
				IncludePlayers:        pk.IncludePlayers,
				ShowBoundingBox:       pk.ShowBoundingBox,
				StructureBlockType:    pk.StructureBlockType,
				Settings:              pk.Settings,
				RedstoneSaveMode:      pk.RedstoneSaveMode,
				ShouldTrigger:         pk.ShouldTrigger,
				Waterlogged:           pk.Waterlogged,
			}
		case *packet.BossEvent:
			pks[pkIndex] = &legacypacket.BossEvent{
				BossEntityUniqueID:   pk.BossEntityUniqueID,
				EventType:            pk.EventType,
				PlayerUniqueID:       pk.PlayerUniqueID,
				BossBarTitle:         pk.BossBarTitle,
				FilteredBossBarTitle: pk.FilteredBossBarTitle,
				HealthPercentage:     pk.HealthPercentage,
				ScreenDarkening:      pk.ScreenDarkening,
				Colour:               pk.Colour,
				Overlay:              pk.Overlay,
			}
		case *packet.CameraAimAssistPresets:
			categories := make([]proto.CameraAimAssistCategory, len(pk.Categories))
			for i, c := range pk.Categories {
				categories[i] = (&proto.CameraAimAssistCategory{}).FromLatest(c)
			}
			presets := make([]proto.CameraAimAssistPreset, len(pk.Presets))
			for i, p := range pk.Presets {
				presets[i] = (&proto.CameraAimAssistPreset{}).FromLatest(p)
			}
			pks[pkIndex] = &legacypacket.CameraAimAssistPresets{
				Categories: categories,
				Presets:    presets,
				Operation:  pk.Operation,
			}
		case *packet.CommandBlockUpdate:
			pks[pkIndex] = &legacypacket.CommandBlockUpdate{
				Block:                   pk.Block,
				Position:                pk.Position,
				Mode:                    pk.Mode,
				NeedsRedstone:           pk.NeedsRedstone,
				Conditional:             pk.Conditional,
				MinecartEntityRuntimeID: pk.MinecartEntityRuntimeID,
				Command:                 pk.Command,
				LastOutput:              pk.LastOutput,
				Name:                    pk.Name,
				FilteredName:            pk.FilteredName,
				ShouldTrackOutput:       pk.ShouldTrackOutput,
				TickDelay:               pk.TickDelay,
				ExecuteOnFirstTick:      pk.ExecuteOnFirstTick,
			}
		case *packet.CreativeContent:
			items := make([]proto.CreativeItem, len(pk.Items))
			for i, it := range pk.Items {
				items[i] = (&proto.CreativeItem{}).FromLatest(it)
			}
			pks[pkIndex] = &legacypacket.CreativeContent{
				Groups: pk.Groups,
				Items:  items,
			}
		case *packet.UpdateAbilities:
			pks[pkIndex] = &legacypacket.UpdateAbilities{
				AbilityData: (&proto.AbilityData{}).FromLatest(pk.AbilityData),
			}
		case *packet.PlayerSkin:
			pk.Skin.GeometryDataEngineVersion = []byte(p.Ver())
		case *packet.ClientMovementPredictionSync:
			actorFlags := fitBitset(pk.ActorFlags, proto.EntityDataFlagsLength(p.ID()))
			pks[pkIndex] = &legacypacket.ClientMovementPredictionSync{
				ActorFlags:              actorFlags,
				BoundingBoxScale:        pk.BoundingBoxScale,
				BoundingBoxWidth:        pk.BoundingBoxWidth,
				BoundingBoxHeight:       pk.BoundingBoxHeight,
				MovementSpeed:           pk.MovementSpeed,
				UnderwaterMovementSpeed: pk.UnderwaterMovementSpeed,
				LavaMovementSpeed:       pk.LavaMovementSpeed,
				JumpStrength:            pk.JumpStrength,
				Health:                  pk.Health,
				Hunger:                  pk.Hunger,
				EntityUniqueID:          pk.EntityUniqueID,
				Flying:                  pk.Flying,
			}
		case *packet.LevelSoundEvent:
			pks[pkIndex] = &legacypacket.LevelSoundEvent{
				SoundType:             pk.SoundType,
				Position:              pk.Position,
				ExtraData:             pk.ExtraData,
				EntityType:            pk.EntityType,
				BabyMob:               pk.BabyMob,
				DisableRelativeVolume: pk.DisableRelativeVolume,
				EntityUniqueID:        pk.EntityUniqueID,
			}
		case *packet.SetHud:
			pks[pkIndex] = &legacypacket.SetHud{
				Elements:   pk.Elements,
				Visibility: pk.Visibility,
			}
		case *packet.BiomeDefinitionList:
			biomeDefinitions := make([]proto.BiomeDefinition, len(pk.BiomeDefinitions))
			for i, bd := range pk.BiomeDefinitions {
				biomeDefinitions[i] = (&proto.BiomeDefinition{}).FromLatest(bd)
			}
			pks[pkIndex] = &legacypacket.BiomeDefinitionList{
				BiomeDefinitions: biomeDefinitions,
				StringList:       pk.StringList,
			}
		case *packet.PlayerList:
			entries := make([]proto.PlayerListEntry, len(pk.Entries))
			for i, e := range pk.Entries {
				entries[i] = (&proto.PlayerListEntry{}).FromLatest(e)
			}
			pks[pkIndex] = &legacypacket.PlayerList{
				ActionType: pk.ActionType,
				Entries:    entries,
			}
		case *packet.SubChunk:
			entries := make([]proto.SubChunkEntry, len(pk.SubChunkEntries))
			for i, e := range pk.SubChunkEntries {
				entries[i] = (&proto.SubChunkEntry{}).FromLatest(e)
			}
			pks[pkIndex] = &legacypacket.SubChunk{
				CacheEnabled:    pk.CacheEnabled,
				Dimension:       pk.Dimension,
				Position:        pk.Position,
				SubChunkEntries: entries,
			}
		case *packet.GameRulesChanged:
			pks[pkIndex] = &legacypacket.GameRulesChanged{
				GameRules: pk.GameRules,
			}
		case *packet.Animate:
			pks[pkIndex] = &legacypacket.Animate{
				ActionType:      pk.ActionType,
				EntityRuntimeID: pk.EntityRuntimeID,
				Data:            pk.Data,
				SwingSource:     protocol.Option(swingSourceToString(pk.SwingSource)),
			}
		case *packet.AvailableCommands:
			commands := make([]proto.Command, len(pk.Commands))
			for i, c := range pk.Commands {
				commands[i] = (&proto.Command{}).FromLatest(c)
			}
			enums := make([]proto.CommandEnum, len(pk.Enums))
			for i, e := range pk.Enums {
				enums[i] = (&proto.CommandEnum{}).FromLatest(e)
			}
			chainedSubcommands := make([]proto.ChainedSubcommand, len(pk.ChainedSubcommands))
			for i, c := range pk.ChainedSubcommands {
				chainedSubcommands[i] = (&proto.ChainedSubcommand{}).FromLatest(c)
			}
			pks[pkIndex] = &legacypacket.AvailableCommands{
				EnumValues:              pk.EnumValues,
				ChainedSubcommandValues: pk.ChainedSubcommandValues,
				Suffixes:                pk.Suffixes,
				Enums:                   enums,
				ChainedSubcommands:      chainedSubcommands,
				Commands:                commands,
				DynamicEnums:            pk.DynamicEnums,
				Constraints:             pk.Constraints,
			}
		case *packet.CommandOutput:
			outputMessages := make([]proto.CommandOutputMessage, len(pk.OutputMessages))
			for i, outputMessage := range pk.OutputMessages {
				outputMessages[i] = (&proto.CommandOutputMessage{}).FromLatest(outputMessage)
			}
			pks[pkIndex] = &legacypacket.CommandOutput{
				CommandOrigin:  pk.CommandOrigin,
				OutputType:     legacypacket.LatestCommandOutputType(pk.OutputType),
				SuccessCount:   pk.SuccessCount,
				OutputMessages: outputMessages,
				DataSet:        pk.DataSet,
			}
		case *packet.CommandRequest:
			pks[pkIndex] = &legacypacket.CommandRequest{
				CommandLine:   pk.CommandLine,
				CommandOrigin: pk.CommandOrigin,
				Internal:      pk.Internal,
				Version:       pk.Version,
			}
		case *packet.Event:
			pks[pkIndex] = &legacypacket.Event{
				EntityRuntimeID: pk.EntityRuntimeID,
				UsePlayerID:     pk.UsePlayerID,
				Event:           pk.Event,
			}
		case *packet.Interact:
			pks[pkIndex] = &legacypacket.Interact{
				ActionType:            pk.ActionType,
				TargetEntityRuntimeID: pk.TargetEntityRuntimeID,
				Position:              pk.Position,
			}
		}
	}

	return pks
}

func (p *Protocol) upgradePackets(pks []packet.Packet, conn *minecraft.Conn) []packet.Packet {
	for pkIndex, pk := range pks {
		switch pk := pk.(type) {
		case *packet.ClientCacheStatus:
			// pk.Enabled = false // TODO: enable when chunk translation is not broken
		case *legacypacket.SetActorMotion:
			pks[pkIndex] = &packet.SetActorMotion{
				EntityRuntimeID: pk.EntityRuntimeID,
				Velocity:        pk.Velocity,
				Tick:            pk.Tick,
			}
		case *legacypacket.ResourcePackStack:
			pks[pkIndex] = &packet.ResourcePackStack{
				TexturePackRequired:          pk.TexturePackRequired,
				TexturePacks:                 pk.TexturePacks,
				BaseGameVersion:              pk.BaseGameVersion,
				Experiments:                  pk.Experiments,
				ExperimentsPreviouslyToggled: pk.ExperimentsPreviouslyToggled,
				IncludeEditorPacks:           pk.IncludeEditorPacks,
			}
		case *legacypacket.UpdatePlayerGameType:
			pks[pkIndex] = &packet.UpdatePlayerGameType{
				GameType:       pk.GameType,
				PlayerUniqueID: pk.PlayerUniqueID,
				Tick:           pk.Tick,
			}
		case *legacypacket.CameraPresets:
			presets := make([]protocol.CameraPreset, len(pk.Presets))
			for i, p := range pk.Presets {
				presets[i] = p.ToLatest()
			}
			pks[pkIndex] = &packet.CameraPresets{
				Presets: presets,
			}
		case *packet.StartGame:
			pk.GameVersion = p.ver
			pk.BaseGameVersion = p.ver
		case *legacypacket.PlayerAuthInput:
			pks[pkIndex] = &packet.PlayerAuthInput{
				Pitch:                  pk.Pitch,
				Yaw:                    pk.Yaw,
				Position:               pk.Position,
				MoveVector:             pk.MoveVector,
				HeadYaw:                pk.HeadYaw,
				InputData:              fitBitset(pk.InputData, packet.PlayerAuthInputBitsetSize),
				InputMode:              pk.InputMode,
				PlayMode:               pk.PlayMode,
				InteractionModel:       pk.InteractionModel,
				InteractPitch:          pk.InteractPitch,
				InteractYaw:            pk.InteractYaw,
				Tick:                   pk.Tick,
				Delta:                  pk.Delta,
				ItemInteractionData:    pk.ItemInteractionData,
				ItemStackRequest:       pk.ItemStackRequest.ToLatest(),
				BlockActions:           pk.BlockActions,
				VehicleRotation:        pk.VehicleRotation,
				ClientPredictedVehicle: pk.ClientPredictedVehicle,
				AnalogueMoveVector:     pk.AnalogueMoveVector,
				CameraOrientation:      pk.CameraOrientation,
				RawMoveVector:          pk.RawMoveVector,
			}
		case *legacypacket.ItemStackResponse:
			responses := make([]protocol.ItemStackResponse, len(pk.Responses))
			for i, r := range pk.Responses {
				responses[i] = r.ToLatest()
			}
			pks[pkIndex] = &packet.ItemStackResponse{Responses: responses}
		case *legacypacket.ResourcePacksInfo:
			texturePacks := make([]protocol.TexturePackInfo, len(pk.TexturePacks))
			for i, t := range pk.TexturePacks {
				texturePacks[i] = t.ToLatest()
				if texturePacks[i].DownloadURL == "" {
					for _, u := range pk.PackURLs {
						if u.UUIDVersion == t.UUID.String()+"_"+t.Version {
							texturePacks[i].DownloadURL = u.URL
							break
						}
					}
				}
			}
			pks[pkIndex] = &packet.ResourcePacksInfo{
				TexturePackRequired:        pk.TexturePackRequired,
				HasAddons:                  pk.HasAddons,
				HasScripts:                 pk.HasScripts,
				ForceDisableVibrantVisuals: pk.ForceDisableVibrantVisuals,
				WorldTemplateUUID:          pk.WorldTemplateUUID,
				WorldTemplateVersion:       pk.WorldTemplateVersion,
				TexturePacks:               texturePacks,
			}
		case *legacypacket.InventorySlot:
			pks[pkIndex] = &packet.InventorySlot{
				WindowID:    pk.WindowID,
				Slot:        pk.Slot,
				Container:   pk.Container.ToLatest(),
				StorageItem: pk.StorageItem,
				NewItem:     pk.NewItem,
			}
		case *legacypacket.InventoryContent:
			pks[pkIndex] = &packet.InventoryContent{
				WindowID:    pk.WindowID,
				Content:     pk.Content,
				Container:   pk.Container.ToLatest(),
				StorageItem: pk.StorageItem,
			}
		case *legacypacket.MobEffect:
			pks[pkIndex] = &packet.MobEffect{
				EntityRuntimeID: pk.EntityRuntimeID,
				Operation:       pk.Operation,
				EffectType:      pk.EffectType,
				Amplifier:       pk.Amplifier,
				Particles:       pk.Particles,
				Duration:        pk.Duration,
				Tick:            pk.Tick,
				Ambient:         pk.Ambient,
			}
		case *legacypacket.CameraAimAssist:
			pks[pkIndex] = &packet.CameraAimAssist{
				Preset:          pk.Preset,
				Angle:           pk.Angle,
				Distance:        pk.Distance,
				TargetMode:      pk.TargetMode,
				Action:          pk.Action,
				ShowDebugRender: pk.ShowDebugRender,
			}
		case *legacypacket.UpdateAttributes:
			attributes := make([]protocol.Attribute, len(pk.Attributes))
			for i, a := range pk.Attributes {
				attributes[i] = a.ToLatest()
			}
			pks[pkIndex] = &packet.UpdateAttributes{
				EntityRuntimeID: pk.EntityRuntimeID,
				Attributes:      attributes,
				Tick:            pk.Tick,
			}
		case *legacypacket.ContainerRegistryCleanup:
			removedContainers := make([]protocol.FullContainerName, len(pk.RemovedContainers))
			for i, c := range pk.RemovedContainers {
				removedContainers[i] = c.ToLatest()
			}
			pks[pkIndex] = &packet.ContainerRegistryCleanup{
				RemovedContainers: removedContainers,
			}
		case *legacypacket.Emote:
			pks[pkIndex] = &packet.Emote{
				EntityRuntimeID: pk.EntityRuntimeID,
				EmoteLength:     pk.EmoteLength,
				EmoteID:         pk.EmoteID,
				XUID:            pk.XUID,
				PlatformID:      pk.PlatformID,
				Flags:           pk.Flags,
			}
		case *legacypacket.Transfer:
			pks[pkIndex] = &packet.Transfer{
				Address:     pk.Address,
				Port:        pk.Port,
				ReloadWorld: pk.ReloadWorld,
			}
		case *legacypacket.AddActor:
			links := make([]protocol.EntityLink, len(pk.EntityLinks))
			for i, l := range pk.EntityLinks {
				links[i] = l.ToLatest()
			}
			pks[pkIndex] = &packet.AddActor{
				EntityUniqueID:   pk.EntityUniqueID,
				EntityRuntimeID:  pk.EntityRuntimeID,
				EntityType:       pk.EntityType,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				BodyYaw:          pk.BodyYaw,
				Attributes:       pk.Attributes,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				EntityLinks:      links,
			}
		case *legacypacket.AddPlayer:
			links := make([]protocol.EntityLink, len(pk.EntityLinks))
			for i, l := range pk.EntityLinks {
				links[i] = l.ToLatest()
			}
			pks[pkIndex] = &packet.AddPlayer{
				UUID:             pk.UUID,
				Username:         pk.Username,
				EntityRuntimeID:  pk.EntityRuntimeID,
				PlatformChatID:   pk.PlatformChatID,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				HeldItem:         pk.HeldItem,
				GameType:         pk.GameType,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				AbilityData:      pk.AbilityData.ToLatest(),
				EntityLinks:      links,
				DeviceID:         pk.DeviceID,
				BuildPlatform:    pk.BuildPlatform,
			}
		case *legacypacket.SetActorLink:
			pks[pkIndex] = &packet.SetActorLink{
				EntityLink: pk.EntityLink.ToLatest(),
			}
		case *legacypacket.CameraInstruction:
			var iSet protocol.Optional[protocol.CameraInstructionSet]
			if v, ok := pk.Set.Value(); ok {
				iSet = protocol.Option(v.ToLatest())
			}
			var spline protocol.Optional[protocol.CameraSplineInstruction]
			if v, ok := pk.Spline.Value(); ok {
				spline = protocol.Option(v.ToLatest())
			}
			pks[pkIndex] = &packet.CameraInstruction{
				Set:              iSet,
				Clear:            pk.Clear,
				Fade:             pk.Fade,
				Target:           pk.Target,
				RemoveTarget:     pk.RemoveTarget,
				FieldOfView:      pk.FieldOfView,
				Spline:           spline,
				AttachToEntity:   pk.AttachToEntity,
				DetachFromEntity: pk.DetachFromEntity,
			}
		case *legacypacket.ChangeDimension:
			pks[pkIndex] = &packet.ChangeDimension{
				Dimension:       pk.Dimension,
				Position:        pk.Position,
				Respawn:         pk.Respawn,
				LoadingScreenID: pk.LoadingScreenID,
			}
		case *legacypacket.CorrectPlayerMovePrediction:
			pks[pkIndex] = &packet.CorrectPlayerMovePrediction{
				PredictionType:         pk.PredictionType,
				Position:               pk.Position,
				Delta:                  pk.Delta,
				Rotation:               pk.Rotation,
				VehicleAngularVelocity: pk.VehicleAngularVelocity,
				OnGround:               pk.OnGround,
				Tick:                   pk.Tick,
			}
		case *legacypacket.Disconnect:
			pks[pkIndex] = &packet.Disconnect{
				Reason:                  pk.Reason,
				HideDisconnectionScreen: pk.HideDisconnectionScreen,
				Message:                 pk.Message,
				FilteredMessage:         pk.FilteredMessage,
			}
		case *legacypacket.EditorNetwork:
			pks[pkIndex] = &packet.EditorNetwork{
				RouteToManager: pk.RouteToManager,
				Payload:        pk.Payload,
			}
		case *legacypacket.MobArmourEquipment:
			pks[pkIndex] = &packet.MobArmourEquipment{
				EntityRuntimeID: pk.EntityRuntimeID,
				Helmet:          pk.Helmet,
				Chestplate:      pk.Chestplate,
				Leggings:        pk.Leggings,
				Boots:           pk.Boots,
				Body:            pk.Body,
			}
		case *legacypacket.PlayerArmourDamage:
			pks[pkIndex] = &packet.PlayerArmourDamage{
				List: pk.List,
			}
		case *legacypacket.SetTitle:
			pks[pkIndex] = &packet.SetTitle{
				ActionType:       pk.ActionType,
				Text:             pk.Text,
				FadeInDuration:   pk.FadeInDuration,
				RemainDuration:   pk.RemainDuration,
				FadeOutDuration:  pk.FadeOutDuration,
				XUID:             pk.XUID,
				PlatformOnlineID: pk.PlatformOnlineID,
				FilteredMessage:  pk.FilteredMessage,
			}
		case *legacypacket.StopSound:
			pks[pkIndex] = &packet.StopSound{
				SoundName:       pk.SoundName,
				StopAll:         pk.StopAll,
				StopMusicLegacy: pk.StopMusicLegacy,
			}
		case *legacypacket.InventoryTransaction:
			trData := pk.TransactionData
			if x, ok := trData.(*proto.UseItemTransactionData); ok {
				trData = x.ToLatest()
			}
			pks[pkIndex] = &packet.InventoryTransaction{
				LegacyRequestID:    pk.LegacyRequestID,
				LegacySetItemSlots: pk.LegacySetItemSlots,
				Actions:            pk.Actions,
				TransactionData:    trData,
			}
		case *legacypacket.ItemStackRequest:
			requests := make([]protocol.ItemStackRequest, len(pk.Requests))
			for i, r := range pk.Requests {
				requests[i] = r.ToLatest()
			}
			pks[pkIndex] = &packet.ItemStackRequest{Requests: requests}
		case *legacypacket.CraftingData:
			recipes := make([]protocol.Recipe, len(pk.Recipes))
			for i, r := range pk.Recipes {
				recipes[i] = proto.RecipeToLatest(r)
			}
			pks[pkIndex] = &packet.CraftingData{
				Recipes:                      recipes,
				PotionRecipes:                pk.PotionRecipes,
				PotionContainerChangeRecipes: pk.PotionContainerChangeRecipes,
				MaterialReducers:             pk.MaterialReducers,
				ClearRecipes:                 pk.ClearRecipes,
			}
		case *legacypacket.ContainerClose:
			pks[pkIndex] = &packet.ContainerClose{
				WindowID:      pk.WindowID,
				ContainerType: pk.ContainerType,
				ServerSide:    pk.ServerSide,
			}
		case *legacypacket.Text:
			pks[pkIndex] = &packet.Text{
				TextType:         pk.TextType,
				NeedsTranslation: pk.NeedsTranslation,
				SourceName:       pk.SourceName,
				Message:          pk.Message,
				Parameters:       pk.Parameters,
				XUID:             pk.XUID,
				PlatformChatID:   pk.PlatformChatID,
				FilteredMessage:  pk.FilteredMessage,
			}
		case *legacypacket.BookEdit:
			pks[pkIndex] = &packet.BookEdit{
				InventorySlot:       pk.InventorySlot,
				ActionType:          pk.ActionType,
				PageNumber:          pk.PageNumber,
				SecondaryPageNumber: pk.SecondaryPageNumber,
				Text:                pk.Text,
				PhotoName:           pk.PhotoName,
				Title:               pk.Title,
				Author:              pk.Author,
				XUID:                pk.XUID,
			}
		case *legacypacket.StartGame:
			pks[pkIndex] = &packet.StartGame{
				EntityUniqueID:                 pk.EntityUniqueID,
				EntityRuntimeID:                pk.EntityRuntimeID,
				PlayerGameMode:                 pk.PlayerGameMode,
				PlayerPosition:                 pk.PlayerPosition,
				Pitch:                          pk.Pitch,
				Yaw:                            pk.Yaw,
				WorldSeed:                      pk.WorldSeed,
				SpawnBiomeType:                 pk.SpawnBiomeType,
				UserDefinedBiomeName:           pk.UserDefinedBiomeName,
				Dimension:                      pk.Dimension,
				Generator:                      pk.Generator,
				WorldGameMode:                  pk.WorldGameMode,
				Hardcore:                       pk.Hardcore,
				Difficulty:                     pk.Difficulty,
				WorldSpawn:                     pk.WorldSpawn,
				AchievementsDisabled:           pk.AchievementsDisabled,
				EditorWorldType:                pk.EditorWorldType,
				CreatedInEditor:                pk.CreatedInEditor,
				ExportedFromEditor:             pk.ExportedFromEditor,
				DayCycleLockTime:               pk.DayCycleLockTime,
				EducationEditionOffer:          pk.EducationEditionOffer,
				EducationFeaturesEnabled:       pk.EducationFeaturesEnabled,
				EducationProductID:             pk.EducationProductID,
				RainLevel:                      pk.RainLevel,
				LightningLevel:                 pk.LightningLevel,
				ConfirmedPlatformLockedContent: pk.ConfirmedPlatformLockedContent,
				MultiPlayerGame:                pk.MultiPlayerGame,
				LANBroadcastEnabled:            pk.LANBroadcastEnabled,
				XBLBroadcastMode:               pk.XBLBroadcastMode,
				PlatformBroadcastMode:          pk.PlatformBroadcastMode,
				CommandsEnabled:                pk.CommandsEnabled,
				TexturePackRequired:            pk.TexturePackRequired,
				GameRules:                      pk.GameRules,
				Experiments:                    pk.Experiments,
				ExperimentsPreviouslyToggled:   pk.ExperimentsPreviouslyToggled,
				BonusChestEnabled:              pk.BonusChestEnabled,
				StartWithMapEnabled:            pk.StartWithMapEnabled,
				PlayerPermissions:              pk.PlayerPermissions,
				ServerChunkTickRadius:          pk.ServerChunkTickRadius,
				HasLockedBehaviourPack:         pk.HasLockedBehaviourPack,
				HasLockedTexturePack:           pk.HasLockedTexturePack,
				FromLockedWorldTemplate:        pk.FromLockedWorldTemplate,
				MSAGamerTagsOnly:               pk.MSAGamerTagsOnly,
				FromWorldTemplate:              pk.FromWorldTemplate,
				WorldTemplateSettingsLocked:    pk.WorldTemplateSettingsLocked,
				OnlySpawnV1Villagers:           pk.OnlySpawnV1Villagers,
				PersonaDisabled:                pk.PersonaDisabled,
				CustomSkinsDisabled:            pk.CustomSkinsDisabled,
				EmoteChatMuted:                 pk.EmoteChatMuted,
				BaseGameVersion:                pk.BaseGameVersion,
				LimitedWorldWidth:              pk.LimitedWorldWidth,
				LimitedWorldDepth:              pk.LimitedWorldDepth,
				NewNether:                      pk.NewNether,
				EducationSharedResourceURI:     pk.EducationSharedResourceURI,
				ForceExperimentalGameplay:      pk.ForceExperimentalGameplay,
				LevelID:                        pk.LevelID,
				WorldName:                      pk.WorldName,
				TemplateContentIdentity:        pk.TemplateContentIdentity,
				Trial:                          pk.Trial,
				PlayerMovementSettings:         pk.PlayerMovementSettings.ToLatest(),
				Time:                           pk.Time,
				EnchantmentSeed:                pk.EnchantmentSeed,
				Blocks:                         pk.Blocks,
				MultiPlayerCorrelationID:       pk.MultiPlayerCorrelationID,
				ServerAuthoritativeInventory:   pk.ServerAuthoritativeInventory,
				GameVersion:                    pk.GameVersion,
				PropertyData:                   pk.PropertyData,
				ServerBlockStateChecksum:       pk.ServerBlockStateChecksum,
				ClientSideGeneration:           pk.ClientSideGeneration,
				WorldTemplateID:                pk.WorldTemplateID,
				ChatRestrictionLevel:           pk.ChatRestrictionLevel,
				DisablePlayerInteractions:      pk.DisablePlayerInteractions,
				ServerID:                       pk.ServerID,
				WorldID:                        pk.WorldID,
				ScenarioID:                     pk.ScenarioID,
				OwnerID:                        pk.OwnerID,
				UseBlockNetworkIDHashes:        pk.UseBlockNetworkIDHashes,
				ServerAuthoritativeSound:       pk.ServerAuthoritativeSound,
			}
		case *legacypacket.CodeBuilderSource:
			pks[pkIndex] = &packet.CodeBuilderSource{
				Operation:  pk.Operation,
				Category:   pk.Category,
				CodeStatus: pk.CodeStatus,
			}
		case *legacypacket.BossEvent:
			pks[pkIndex] = &packet.BossEvent{
				BossEntityUniqueID:   pk.BossEntityUniqueID,
				EventType:            pk.EventType,
				PlayerUniqueID:       pk.PlayerUniqueID,
				BossBarTitle:         pk.BossBarTitle,
				FilteredBossBarTitle: pk.FilteredBossBarTitle,
				HealthPercentage:     pk.HealthPercentage,
				ScreenDarkening:      pk.ScreenDarkening,
				Colour:               pk.Colour,
				Overlay:              pk.Overlay,
			}
		case *legacypacket.CameraAimAssistPresets:
			categories := make([]protocol.CameraAimAssistCategory, len(pk.Categories))
			for i, c := range pk.Categories {
				categories[i] = c.ToLatest()
			}
			presets := make([]protocol.CameraAimAssistPreset, len(pk.Presets))
			for i, p := range pk.Presets {
				presets[i] = p.ToLatest()
			}
			pks[pkIndex] = &packet.CameraAimAssistPresets{
				Categories: categories,
				Presets:    presets,
				Operation:  pk.Operation,
			}
		case *legacypacket.CommandBlockUpdate:
			pks[pkIndex] = &packet.CommandBlockUpdate{
				Block:                   pk.Block,
				Position:                pk.Position,
				Mode:                    pk.Mode,
				NeedsRedstone:           pk.NeedsRedstone,
				Conditional:             pk.Conditional,
				MinecartEntityRuntimeID: pk.MinecartEntityRuntimeID,
				Command:                 pk.Command,
				LastOutput:              pk.LastOutput,
				Name:                    pk.Name,
				FilteredName:            pk.FilteredName,
				ShouldTrackOutput:       pk.ShouldTrackOutput,
				TickDelay:               pk.TickDelay,
				ExecuteOnFirstTick:      pk.ExecuteOnFirstTick,
			}
		case *legacypacket.CreativeContent:
			items := make([]protocol.CreativeItem, len(pk.Items))
			for i, it := range pk.Items {
				items[i] = it.ToLatest()
			}
			pks[pkIndex] = &packet.CreativeContent{
				Groups: pk.Groups,
				Items:  items,
			}
		case *legacypacket.ItemRegistry:
			pk.Items = p.itemTranslator.UpgradeItemEntries(pk.Items)

			items := make([]protocol.ItemEntry, len(pk.Items))
			for i, it := range pk.Items {
				items[i] = it.ToLatest()
			}

			pks[pkIndex] = &packet.ItemRegistry{Items: items}
		case *legacypacket.StructureBlockUpdate:
			pks[pkIndex] = &packet.StructureBlockUpdate{
				Position:              pk.Position,
				StructureName:         pk.StructureName,
				FilteredStructureName: pk.FilteredStructureName,
				DataField:             pk.DataField,
				IncludePlayers:        pk.IncludePlayers,
				ShowBoundingBox:       pk.ShowBoundingBox,
				StructureBlockType:    pk.StructureBlockType,
				Settings:              pk.Settings,
				RedstoneSaveMode:      pk.RedstoneSaveMode,
				ShouldTrigger:         pk.ShouldTrigger,
				Waterlogged:           pk.Waterlogged,
			}
		case *legacypacket.UpdateAbilities:
			pks[pkIndex] = &packet.UpdateAbilities{AbilityData: pk.AbilityData.ToLatest()}
		case *legacypacket.ClientMovementPredictionSync:
			pks[pkIndex] = &packet.ClientMovementPredictionSync{
				ActorFlags:              fitBitset(pk.ActorFlags, protocol.EntityDataFlagCount),
				BoundingBoxScale:        pk.BoundingBoxScale,
				BoundingBoxWidth:        pk.BoundingBoxWidth,
				BoundingBoxHeight:       pk.BoundingBoxHeight,
				MovementSpeed:           pk.MovementSpeed,
				UnderwaterMovementSpeed: pk.UnderwaterMovementSpeed,
				LavaMovementSpeed:       pk.LavaMovementSpeed,
				JumpStrength:            pk.JumpStrength,
				Health:                  pk.Health,
				Hunger:                  pk.Hunger,
				EntityUniqueID:          pk.EntityUniqueID,
				Flying:                  pk.Flying,
			}
		case *legacypacket.LevelSoundEvent:
			pks[pkIndex] = &packet.LevelSoundEvent{
				SoundType:             pk.SoundType,
				Position:              pk.Position,
				ExtraData:             pk.ExtraData,
				EntityType:            pk.EntityType,
				BabyMob:               pk.BabyMob,
				DisableRelativeVolume: pk.DisableRelativeVolume,
				EntityUniqueID:        pk.EntityUniqueID,
			}
		case *legacypacket.SetHud:
			pks[pkIndex] = &packet.SetHud{
				Elements:   pk.Elements,
				Visibility: pk.Visibility,
			}
		case *legacypacket.BiomeDefinitionList:
			biomeDefinitions := make([]protocol.BiomeDefinition, len(pk.BiomeDefinitions))
			for i, bd := range pk.BiomeDefinitions {
				biomeDefinitions[i] = bd.ToLatest()
			}
			pks[pkIndex] = &packet.BiomeDefinitionList{
				BiomeDefinitions: biomeDefinitions,
				StringList:       pk.StringList,
			}
		case *legacypacket.PlayerList:
			entries := make([]protocol.PlayerListEntry, len(pk.Entries))
			for i, e := range pk.Entries {
				entries[i] = e.ToLatest()
			}
			pks[pkIndex] = &packet.PlayerList{
				ActionType: pk.ActionType,
				Entries:    entries,
			}
		case *legacypacket.SubChunk:
			entries := make([]protocol.SubChunkEntry, len(pk.SubChunkEntries))
			for i, e := range pk.SubChunkEntries {
				entries[i] = e.ToLatest()
			}
			pks[pkIndex] = &packet.SubChunk{
				CacheEnabled:    pk.CacheEnabled,
				Dimension:       pk.Dimension,
				Position:        pk.Position,
				SubChunkEntries: entries,
			}
		case *legacypacket.GameRulesChanged:
			pks[pkIndex] = &packet.GameRulesChanged{
				GameRules: pk.GameRules,
			}
		case *legacypacket.Animate:
			newSwingSource := packet.AnimateSwingSourceNone
			source, _ := pk.SwingSource.Value()
			switch source {
			case "build":
				newSwingSource = packet.AnimateSwingSourceBuild
			case "mine":
				newSwingSource = packet.AnimateSwingSourceMine
			case "interact":
				newSwingSource = packet.AnimateSwingSourceInteract
			case "attack":
				newSwingSource = packet.AnimateSwingSourceAttack
			case "useitem":
				newSwingSource = packet.AnimateSwingSourceUseItem
			case "throwitem":
				newSwingSource = packet.AnimateSwingSourceThrowItem
			case "dropitem":
				newSwingSource = packet.AnimateSwingSourceDropItem
			case "event":
				newSwingSource = packet.AnimateSwingSourceEvent
			}
			pks[pkIndex] = &packet.Animate{
				ActionType:      pk.ActionType,
				EntityRuntimeID: pk.EntityRuntimeID,
				Data:            pk.Data,
				SwingSource:     uint8(newSwingSource),
			}
		case *legacypacket.AvailableCommands:
			enums := make([]protocol.CommandEnum, len(pk.Enums))
			for i, e := range pk.Enums {
				enums[i] = e.ToLatest()
			}
			chainedSubcommands := make([]protocol.ChainedSubcommand, len(pk.ChainedSubcommands))
			for i, c := range pk.ChainedSubcommands {
				chainedSubcommands[i] = c.ToLatest()
			}
			commands := make([]protocol.Command, len(pk.Commands))
			for i, c := range pk.Commands {
				commands[i] = c.ToLatest()
			}
			pks[pkIndex] = &packet.AvailableCommands{
				EnumValues:              pk.EnumValues,
				ChainedSubcommandValues: pk.ChainedSubcommandValues,
				Suffixes:                pk.Suffixes,
				Enums:                   enums,
				ChainedSubcommands:      chainedSubcommands,
				Commands:                commands,
				DynamicEnums:            pk.DynamicEnums,
				Constraints:             pk.Constraints,
			}
		case *legacypacket.CommandOutput:
			outputMessages := make([]protocol.CommandOutputMessage, len(pk.OutputMessages))
			for i, outputMessage := range pk.OutputMessages {
				outputMessages[i] = outputMessage.ToLatest()
			}
			pks[pkIndex] = &packet.CommandOutput{
				CommandOrigin:  pk.CommandOrigin,
				OutputType:     legacypacket.LegacyCommandOutputType(pk.OutputType),
				SuccessCount:   pk.SuccessCount,
				OutputMessages: outputMessages,
				DataSet:        pk.DataSet,
			}
		case *legacypacket.CommandRequest:
			pks[pkIndex] = &packet.CommandRequest{
				CommandLine:   pk.CommandLine,
				CommandOrigin: pk.CommandOrigin,
				Internal:      pk.Internal,
				Version:       pk.Version,
			}
		case *legacypacket.Event:
			pks[pkIndex] = &packet.Event{
				EntityRuntimeID: pk.EntityRuntimeID,
				UsePlayerID:     pk.UsePlayerID,
				Event:           pk.Event,
			}
		case *legacypacket.Interact:
			pks[pkIndex] = &packet.Interact{
				ActionType:            pk.ActionType,
				TargetEntityRuntimeID: pk.TargetEntityRuntimeID,
				Position:              pk.Position,
			}
		}
	}
	return pks
}

func swingSourceToString(x uint8) string {
	switch x {
	case packet.AnimateSwingSourceNone:
		return "none"
	case packet.AnimateSwingSourceBuild:
		return "build"
	case packet.AnimateSwingSourceMine:
		return "mine"
	case packet.AnimateSwingSourceInteract:
		return "interact"
	case packet.AnimateSwingSourceAttack:
		return "attack"
	case packet.AnimateSwingSourceUseItem:
		return "useitem"
	case packet.AnimateSwingSourceThrowItem:
		return "throwitem"
	case packet.AnimateSwingSourceDropItem:
		return "dropitem"
	case packet.AnimateSwingSourceEvent:
		return "event"
	default:
		return "unknown"
	}
}
