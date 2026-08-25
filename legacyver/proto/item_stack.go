package proto

import (
	"fmt"
	_ "unsafe"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func MarshalItemStackRequest(r protocol.IO, x *protocol.ItemStackRequest) {
	r.Varint32(&x.RequestID)
	protocol.FuncSlice(r, &x.Actions, func(p *protocol.StackRequestAction) {
		IOStackRequestAction(r, p)
	})
	protocol.FuncSlice(r, &x.FilterStrings, r.String)
	r.Int32(&x.FilterCause)
}

func MarshalStackRequestAction(r protocol.IO, x protocol.StackRequestAction) {
	switch act := x.(type) {
	case *protocol.TakeStackRequestAction:
		MarshalTakeStackRequestAction(r, act)
	case *protocol.PlaceStackRequestAction:
		MarshalPlaceStackRequestAction(r, act)
	case *protocol.SwapStackRequestAction:
		MarshalSwapStackRequestAction(r, act)
	case *protocol.DropStackRequestAction:
		MarshalDropStackRequestAction(r, act)
	case *protocol.DestroyStackRequestAction:
		MarshalDestroyStackRequestAction(r, act)
	case *protocol.ConsumeStackRequestAction:
		MarshalConsumeStackRequestAction(r, act)
	case *protocol.CreateStackRequestAction:
		r.Uint8(&act.ResultsSlot)
	case *protocol.LabTableCombineStackRequestAction, *protocol.CraftNonImplementedStackRequestAction:
		// These actions carry no additional data.
	case *protocol.BeaconPaymentStackRequestAction:
		r.Varint32(&act.PrimaryEffect)
		r.Varint32(&act.SecondaryEffect)
	case *protocol.MineBlockStackRequestAction:
		r.Varint32(&act.HotbarSlot)
		r.Varint32(&act.PredictedDurability)
		if IsProtoGTE(r, ID2168) {
			r.Int32(&act.StackNetworkID)
		} else {
			r.Varint32(&act.StackNetworkID)
		}
	case *protocol.CraftRecipeStackRequestAction:
		MarshalCraftRecipeStackRequestAction(r, act)
	case *protocol.AutoCraftRecipeStackRequestAction:
		MarshalAutoCraftRecipeStackRequestAction(r, act)
	case *protocol.CraftCreativeStackRequestAction:
		MarshalCraftCreativeStackRequestAction(r, act)
	case *protocol.CraftRecipeOptionalStackRequestAction:
		r.Varuint32(&act.RecipeNetworkID)
		r.Int32(&act.FilterStringIndex)
	case *protocol.CraftGrindstoneRecipeStackRequestAction:
		MarshalCraftGrindstoneRecipeStackRequestAction(r, act)
	case *protocol.CraftLoomRecipeStackRequestAction:
		MarshalCraftLoomRecipeStackRequestAction(r, act)
	case *protocol.CraftResultsDeprecatedStackRequestAction:
		MarshalCraftResultsDeprecatedStackRequestAction(r, act)
	default:
		r.UnknownEnumOption(fmt.Sprintf("%T", x), "stack request action type")
	}
}

func MarshalItemStackResponse(r protocol.IO, x *protocol.ItemStackResponse) {
	r.Uint8(&x.Status)
	r.Varint32(&x.RequestID)
	if IsProtoGTE(r, ID2168) {
		var containerInfo protocol.Optional[[]protocol.StackResponseContainerInfo]
		if len(x.ContainerInfo) != 0 {
			containerInfo = protocol.Option(x.ContainerInfo)
		}
		protocol.DoubleOptionalFunc(r, &containerInfo, func(entries *[]protocol.StackResponseContainerInfo) {
			protocol.FuncIOSlice(r, entries, MarshalStackResponseContainerInfo)
		})
		if entries, ok := containerInfo.Value(); ok {
			x.ContainerInfo = entries
		}
	} else if x.Status == protocol.ItemStackResponseStatusOK {
		protocol.FuncIOSlice(r, &x.ContainerInfo, MarshalStackResponseContainerInfo)
	}
}

func MarshalStackResponseContainerInfo(r protocol.IO, x *protocol.StackResponseContainerInfo) {
	MarshalFullContainerName(r, &x.Container)
	protocol.FuncIOSlice(r, &x.SlotInfo, MarshalStackResponseSlotInfo)
}

func MarshalStackResponseSlotInfo(r protocol.IO, x *protocol.StackResponseSlotInfo) {
	r.Uint8(&x.Slot)
	r.Uint8(&x.HotbarSlot)
	r.Uint8(&x.Count)
	if IsProtoGTE(r, ID2168) {
		var stackNetworkID protocol.Optional[int32]
		if x.StackNetworkID > 0 {
			stackNetworkID = protocol.Option(x.StackNetworkID)
		}
		protocol.DoubleOptionalFunc(r, &stackNetworkID, r.Varint32)
		if value, ok := stackNetworkID.Value(); ok {
			x.StackNetworkID = value
		}
	} else {
		r.Varint32(&x.StackNetworkID)
	}
	if IsProtoLT(r, ID2168) && x.Slot != x.HotbarSlot {
		r.InvalidValue(x.HotbarSlot, "hotbar slot", "hot bar slot must be equal to normal slot")
	}
	r.String(&x.CustomName)
	if IsProtoGTE(r, ID2168) {
		// The oomph gophertunnel fork encodes FilteredCustomName as an
		// optional string at 2168.
		protocol.OptionalFunc(r, &x.FilteredCustomName, r.String)
	} else if IsProtoGTE(r, ID766) {
		// Older protocols always encode a plain string.
		filtered, _ := x.FilteredCustomName.Value()
		r.String(&filtered)
		if IsReader(r) {
			if filtered != "" {
				x.FilteredCustomName = protocol.Option(filtered)
			} else {
				x.FilteredCustomName = protocol.Optional[string]{}
			}
		}
	}
	r.Varint32(&x.DurabilityCorrection)
	if IsProtoGTE(r, ID2168) && (x.DurabilityCorrection < -32768 || x.DurabilityCorrection > 32767) {
		r.InvalidValue(x.DurabilityCorrection, "durability correction", "must fit in an int16")
	}
}

func MarshalTakeStackRequestAction(r protocol.IO, a *protocol.TakeStackRequestAction) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
	StackReqSlotInfo(r, &a.Destination)
}

func MarshalPlaceStackRequestAction(r protocol.IO, a *protocol.PlaceStackRequestAction) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
	StackReqSlotInfo(r, &a.Destination)
}

func MarshalSwapStackRequestAction(r protocol.IO, a *protocol.SwapStackRequestAction) {
	StackReqSlotInfo(r, &a.Source)
	StackReqSlotInfo(r, &a.Destination)
}

func MarshalDropStackRequestAction(r protocol.IO, a *protocol.DropStackRequestAction) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
	r.Bool(&a.Randomly)
}

func MarshalDestroyStackRequestAction(r protocol.IO, a *protocol.DestroyStackRequestAction) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
}

func MarshalConsumeStackRequestAction(r protocol.IO, a *protocol.ConsumeStackRequestAction) {
	MarshalDestroyStackRequestAction(r, &a.DestroyStackRequestAction)
}

func MarshalCraftRecipeStackRequestAction(r protocol.IO, a *protocol.CraftRecipeStackRequestAction) {
	r.Varuint32(&a.RecipeNetworkID)
	if IsProtoGTE(r, ID712) {
		r.Uint8(&a.NumberOfCrafts)
	}
}

func MarshalAutoCraftRecipeStackRequestAction(r protocol.IO, a *protocol.AutoCraftRecipeStackRequestAction) {
	r.Varuint32(&a.RecipeNetworkID)
	r.Uint8(&a.NumberOfCrafts)
	if IsProtoGTE(r, ID712) && IsProtoLT(r, ID2168) {
		// Removed in 1.26.40. It duplicated NumberOfCrafts.
		timesCrafted := a.NumberOfCrafts
		r.Uint8(&timesCrafted)
	}
	if IsProtoGTE(r, ID2168) {
		protocol.FuncIOSlice(r, &a.Ingredients, protocol.StackRequestItemDescriptorCount)
	} else {
		protocol.FuncSlice(r, &a.Ingredients, r.ItemDescriptorCount)
	}
}

func MarshalCraftCreativeStackRequestAction(r protocol.IO, a *protocol.CraftCreativeStackRequestAction) {
	r.Varuint32(&a.CreativeItemNetworkID)
	if IsProtoGTE(r, ID712) {
		r.Uint8(&a.NumberOfCrafts)
	}
}

func MarshalCraftGrindstoneRecipeStackRequestAction(r protocol.IO, c *protocol.CraftGrindstoneRecipeStackRequestAction) {
	if IsProtoGTE(r, ID2168) {
		protocol.IntegerFunc(&c.RecipeNetworkID, r.Int32)
	} else {
		r.Varuint32(&c.RecipeNetworkID)
	}
	if IsProtoGTE(r, ID712) {
		r.Uint8(&c.NumberOfCrafts)
	}
	r.Varint32(&c.Cost)
}

func MarshalCraftLoomRecipeStackRequestAction(r protocol.IO, c *protocol.CraftLoomRecipeStackRequestAction) {
	r.String(&c.Pattern)
	if IsProtoGTE(r, ID712) {
		r.Uint8(&c.TimesCrafted)
	}
}

func MarshalCraftResultsDeprecatedStackRequestAction(r protocol.IO, a *protocol.CraftResultsDeprecatedStackRequestAction) {
	if IsProtoGTE(r, ID2168) {
		protocol.FuncSlice(r, &a.ResultItems, r.StackRequestItem)
	} else {
		// This action is deprecated and is not expected from modern clients. Its
		// old item form cannot retain the name-only 1.26.40 representation.
		protocol.FuncIOSlice(r, &a.ResultItems, legacyStackRequestItem)
	}
	r.Uint8(&a.TimesCrafted)
}

// StackReqSlotInfo reads/writes a StackRequestSlotInfo x using IO r.
func StackReqSlotInfo(r protocol.IO, x *protocol.StackRequestSlotInfo) {
	MarshalFullContainerName(r, &x.Container)
	r.Uint8(&x.Slot)
	if IsProtoGTE(r, ID2168) {
		r.Int32(&x.StackNetworkID)
	} else {
		r.Varint32(&x.StackNetworkID)
	}
}

func legacyStackRequestItem(r protocol.IO, x *protocol.StackRequestItem) {
	var item protocol.ItemStack
	item.MetadataValue = x.MetadataValue
	item.BlockRuntimeID = x.BlockRuntimeID
	item.Count = x.Count
	item.NBTData = x.NBTData
	item.CanBePlacedOn = x.CanBePlacedOn
	item.CanBreak = x.CanBreak
	item.BlockingTick = x.BlockingTick
	r.Item(&item)
}

//go:linkname lookupStackRequestAction github.com/sandertv/gophertunnel/minecraft/protocol.lookupStackRequestAction
func lookupStackRequestAction(id uint8, x *protocol.StackRequestAction) bool

//go:linkname lookupStackRequestActionType github.com/sandertv/gophertunnel/minecraft/protocol.lookupStackRequestActionType
func lookupStackRequestActionType(x protocol.StackRequestAction, id *uint8) bool
