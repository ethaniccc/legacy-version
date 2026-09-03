package proto

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	recipeShapeless int32 = iota
	recipeShaped
	recipeFurnace
	recipeFurnaceData
	recipeMulti
	recipeShulkerBox
	recipeShapelessChemistry
	recipeShapedChemistry
	recipeSmithingTransform
	recipeSmithingTrim
)

// MarshalCraftingData keeps the pre-1.26.40 tagged recipe vector and the
// 1.26.40 typed vectors in one version-aware implementation.
func MarshalCraftingData(r protocol.IO, pk *packet.CraftingData) {
	if IsProtoGTE(r, ID2168) {
		protocol.FuncIOSlice(r, &pk.ShapedRecipes, MarshalShapedRecipe)
		protocol.FuncIOSlice(r, &pk.ShapelessRecipes, MarshalShapelessRecipe)
		protocol.FuncIOSlice(r, &pk.MultiRecipes, marshalMultiRecipe)
		protocol.FuncIOSlice(r, &pk.UserDataShapelessRecipes, func(r protocol.IO, x *protocol.UserDataShapelessRecipe) {
			MarshalShapelessRecipe(r, &x.ShapelessRecipe)
		})
		protocol.FuncIOSlice(r, &pk.ShapelessChemistryRecipes, func(r protocol.IO, x *protocol.ShapelessChemistryRecipe) {
			MarshalShapelessRecipe(r, &x.ShapelessRecipe)
		})
		protocol.FuncIOSlice(r, &pk.ShapedChemistryRecipes, func(r protocol.IO, x *protocol.ShapedChemistryRecipe) { MarshalShapedRecipe(r, &x.ShapedRecipe) })
		protocol.FuncIOSlice(r, &pk.SmithingTransformRecipes, marshalSmithingTransformRecipe)
		protocol.FuncIOSlice(r, &pk.SmithingTrimRecipes, marshalSmithingTrimRecipe)
		return
	}
	if IsReader(r) {
		var count uint32
		r.Varuint32(&count)
		for range count {
			var recipeType int32
			r.Varint32(&recipeType)
			switch recipeType {
			case recipeShapeless:
				var x protocol.ShapelessRecipe
				MarshalShapelessRecipe(r, &x)
				pk.ShapelessRecipes = append(pk.ShapelessRecipes, x)
			case recipeShaped:
				var x protocol.ShapedRecipe
				MarshalShapedRecipe(r, &x)
				pk.ShapedRecipes = append(pk.ShapedRecipes, x)
			case recipeMulti:
				var x protocol.MultiRecipe
				marshalMultiRecipe(r, &x)
				pk.MultiRecipes = append(pk.MultiRecipes, x)
			case recipeShulkerBox:
				// Shulker box recipes are the UserDataShapelessRecipes
				// of 1.26.40 and later, with the shapeless layout.
				var x protocol.UserDataShapelessRecipe
				MarshalShapelessRecipe(r, &x.ShapelessRecipe)
				pk.UserDataShapelessRecipes = append(pk.UserDataShapelessRecipes, x)
			case recipeShapelessChemistry:
				var x protocol.ShapelessChemistryRecipe
				MarshalShapelessRecipe(r, &x.ShapelessRecipe)
				pk.ShapelessChemistryRecipes = append(pk.ShapelessChemistryRecipes, x)
			case recipeShapedChemistry:
				var x protocol.ShapedChemistryRecipe
				MarshalShapedRecipe(r, &x.ShapedRecipe)
				pk.ShapedChemistryRecipes = append(pk.ShapedChemistryRecipes, x)
			case recipeSmithingTransform:
				var x protocol.SmithingTransformRecipe
				marshalSmithingTransformRecipe(r, &x)
				pk.SmithingTransformRecipes = append(pk.SmithingTransformRecipes, x)
			case recipeSmithingTrim:
				var x protocol.SmithingTrimRecipe
				marshalSmithingTrimRecipe(r, &x)
				pk.SmithingTrimRecipes = append(pk.SmithingTrimRecipes, x)
			case recipeFurnace, recipeFurnaceData:
				marshalDiscardedFurnaceRecipe(r, recipeType == recipeFurnaceData)
			default:
				r.UnknownEnumOption(recipeType, "crafting data recipe type")
				return
			}
		}
		return
	}
	count := len(pk.ShapelessRecipes) + len(pk.ShapedRecipes) + len(pk.MultiRecipes) + len(pk.UserDataShapelessRecipes) +
		len(pk.ShapelessChemistryRecipes) + len(pk.ShapedChemistryRecipes) + len(pk.SmithingTransformRecipes) + len(pk.SmithingTrimRecipes)
	count32 := uint32(count)
	r.Varuint32(&count32)
	for i := range pk.ShapelessRecipes {
		writeRecipeType(r, recipeShapeless)
		MarshalShapelessRecipe(r, &pk.ShapelessRecipes[i])
	}
	for i := range pk.ShapedRecipes {
		writeRecipeType(r, recipeShaped)
		MarshalShapedRecipe(r, &pk.ShapedRecipes[i])
	}
	for i := range pk.MultiRecipes {
		writeRecipeType(r, recipeMulti)
		marshalMultiRecipe(r, &pk.MultiRecipes[i])
	}
	for i := range pk.UserDataShapelessRecipes {
		writeRecipeType(r, recipeShulkerBox)
		MarshalShapelessRecipe(r, &pk.UserDataShapelessRecipes[i].ShapelessRecipe)
	}
	for i := range pk.ShapelessChemistryRecipes {
		writeRecipeType(r, recipeShapelessChemistry)
		MarshalShapelessRecipe(r, &pk.ShapelessChemistryRecipes[i].ShapelessRecipe)
	}
	for i := range pk.ShapedChemistryRecipes {
		writeRecipeType(r, recipeShapedChemistry)
		MarshalShapedRecipe(r, &pk.ShapedChemistryRecipes[i].ShapedRecipe)
	}
	for i := range pk.SmithingTransformRecipes {
		writeRecipeType(r, recipeSmithingTransform)
		marshalSmithingTransformRecipe(r, &pk.SmithingTransformRecipes[i])
	}
	for i := range pk.SmithingTrimRecipes {
		writeRecipeType(r, recipeSmithingTrim)
		marshalSmithingTrimRecipe(r, &pk.SmithingTrimRecipes[i])
	}
}

func writeRecipeType(r protocol.IO, recipeType int32) { r.Varint32(&recipeType) }

func MarshalShapedRecipe(r protocol.IO, recipe *protocol.ShapedRecipe) {
	r.String(&recipe.RecipeID)
	r.Varint32(&recipe.Width)
	r.Varint32(&recipe.Height)
	if IsProtoGTE(r, ID2168) {
		protocol.FuncSlice(r, &recipe.Input, r.ItemDescriptorCount)
	} else {
		protocol.FuncSliceOfLen(r, uint32(recipe.Width*recipe.Height), &recipe.Input, r.ItemDescriptorCount)
	}
	protocol.FuncSlice(r, &recipe.Output, r.Item)
	r.UUID(&recipe.UUID)
	r.String(&recipe.Block)
	r.Varint32(&recipe.Priority)
	if IsProtoGTE(r, ID671) {
		r.Bool(&recipe.AssumeSymmetry)
	}
	if IsProtoGTE(r, ID2168) {
		protocol.OptionalFunc(r, &recipe.UnlockRequirement, func(x *protocol.RecipeUnlockRequirement) { marshalCurrentRecipeUnlockRequirement(r, x) })
	} else if IsProtoGTE(r, ID685) {
		unlock, _ := recipe.UnlockRequirement.Value()
		marshalLegacyRecipeUnlockRequirement(r, &unlock)
		recipe.UnlockRequirement = protocol.Option(unlock)
	}
	r.Varuint32(&recipe.RecipeNetworkID)
}

func MarshalShapelessRecipe(r protocol.IO, recipe *protocol.ShapelessRecipe) {
	r.String(&recipe.RecipeID)
	protocol.FuncSlice(r, &recipe.Input, r.ItemDescriptorCount)
	protocol.FuncSlice(r, &recipe.Output, r.Item)
	r.UUID(&recipe.UUID)
	r.String(&recipe.Block)
	r.Varint32(&recipe.Priority)
	if IsProtoGTE(r, ID2168) {
		protocol.OptionalFunc(r, &recipe.UnlockRequirement, func(x *protocol.RecipeUnlockRequirement) { marshalCurrentRecipeUnlockRequirement(r, x) })
	} else if IsProtoGTE(r, ID685) {
		unlock, _ := recipe.UnlockRequirement.Value()
		marshalLegacyRecipeUnlockRequirement(r, &unlock)
		recipe.UnlockRequirement = protocol.Option(unlock)
	}
	r.Varuint32(&recipe.RecipeNetworkID)
}

func marshalCurrentRecipeUnlockRequirement(r protocol.IO, x *protocol.RecipeUnlockRequirement) {
	r.Varint32(&x.Context)
	present := x.Context == protocol.RecipeUnlockContextNone
	r.Bool(&present)
	if present {
		protocol.FuncSlice(r, &x.Ingredients, r.ItemDescriptorCount)
	}
}

func marshalMultiRecipe(r protocol.IO, x *protocol.MultiRecipe) {
	r.UUID(&x.UUID)
	r.Varuint32(&x.RecipeNetworkID)
}

func marshalSmithingTransformRecipe(r protocol.IO, x *protocol.SmithingTransformRecipe) {
	r.String(&x.RecipeID)
	r.ItemDescriptorCount(&x.Template)
	r.ItemDescriptorCount(&x.Base)
	r.ItemDescriptorCount(&x.Addition)
	r.Item(&x.Result)
	r.String(&x.Block)
	r.Varuint32(&x.RecipeNetworkID)
}

func marshalSmithingTrimRecipe(r protocol.IO, x *protocol.SmithingTrimRecipe) {
	r.String(&x.RecipeID)
	r.ItemDescriptorCount(&x.Template)
	r.ItemDescriptorCount(&x.Base)
	r.ItemDescriptorCount(&x.Addition)
	r.String(&x.Block)
	r.Varuint32(&x.RecipeNetworkID)
}

func marshalLegacyRecipeUnlockRequirement(r protocol.IO, x *protocol.RecipeUnlockRequirement) {
	context := uint8(x.Context)
	r.Uint8(&context)
	x.Context = int32(context)
	if x.Context == protocol.RecipeUnlockContextNone {
		protocol.FuncSlice(r, &x.Ingredients, r.ItemDescriptorCount)
	}
}

func marshalDiscardedFurnaceRecipe(r protocol.IO, withMetadata bool) {
	var networkID int32
	r.Varint32(&networkID)
	if withMetadata {
		var metadata int32
		r.Varint32(&metadata)
	}
	var output protocol.ItemStack
	r.Item(&output)
	var block string
	r.String(&block)
}
