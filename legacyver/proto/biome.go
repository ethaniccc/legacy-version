package proto

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// BiomeDefinition represents a biome definition in the game. This can be a vanilla biome or a completely
// custom biome.
type BiomeDefinition struct {
	// NameIndex represents the index of the biome name in the string list.
	NameIndex int16
	// BiomeID is the biome ID. This is optional and can be empty.
	BiomeID int16
	// Temperature is the temperature of the biome, used for weather, biome behaviours and sky colour.
	Temperature float32
	// Downfall is the amount that precipitation affects colours and block changes.
	Downfall float32
	// RedSporeDensity is the density of red spore precipitation visuals.
	RedSporeDensity float32
	// BlueSporeDensity is the density of blue spore precipitation visuals.
	BlueSporeDensity float32
	// AshDensity is the density of ash precipitation visuals.
	AshDensity float32
	// WhiteAshDensity is the density of white ash precipitation visuals.
	WhiteAshDensity float32
	// FoliageSnow ...
	FoliageSnow float32
	// Depth ...
	Depth float32
	// Scale ...
	Scale float32
	// MapWaterColour is an ARGB value for the water colour on maps in the biome.
	MapWaterColour int32
	// Rain is true if the biome has rain, false if it is a dry biome.
	Rain bool
	// Tags are a list of indices of tags in the string list. These are used to group biomes together for
	// biome generation and other purposes.
	Tags protocol.Optional[[]uint16]
	// ChunkGeneration is optional information to assist in client-side chunk generation. Almost all servers
	// can and should leave this empty to greatly reduce the size of this packet. Only BDS and servers which
	// *exactly* match the vanilla chunk generation can benefit from this.
	ChunkGeneration protocol.Optional[BiomeChunkGeneration]
}

func (x *BiomeDefinition) Marshal(r protocol.IO) {
	r.Int16(&x.NameIndex)
	if IsProtoLT(r, ID827) {
		var opt protocol.Optional[int16]
		if x.BiomeID != -1 {
			opt = protocol.Option(x.BiomeID)
		}
		protocol.OptionalFunc(r, &opt, r.Int16)
		x.BiomeID, _ = opt.Value()
	} else {
		r.Int16(&x.BiomeID)
	}
	r.Float32(&x.Temperature)
	r.Float32(&x.Downfall)
	if IsProtoLT(r, ID844) {
		r.Float32(&x.RedSporeDensity)
		r.Float32(&x.BlueSporeDensity)
		r.Float32(&x.AshDensity)
		r.Float32(&x.WhiteAshDensity)
	}
	if IsProtoGTE(r, ID844) {
		r.Float32(&x.FoliageSnow)
	}
	r.Float32(&x.Depth)
	r.Float32(&x.Scale)
	r.Int32(&x.MapWaterColour)
	r.Bool(&x.Rain)
	protocol.OptionalFunc(r, &x.Tags, func(s *[]uint16) {
		protocol.FuncSlice(r, s, r.Uint16)
	})
	protocol.OptionalMarshaler(r, &x.ChunkGeneration)
}

func (x *BiomeDefinition) FromLatest(bd protocol.BiomeDefinition) BiomeDefinition {
	x.NameIndex = bd.NameIndex
	x.BiomeID = bd.BiomeID
	x.Temperature = bd.Temperature
	x.Downfall = bd.Downfall
	//x.RedSporeDensity = bd.RedSporeDensity
	//x.BlueSporeDensity = bd.BlueSporeDensity
	//x.AshDensity = bd.AshDensity
	//x.WhiteAshDensity = bd.WhiteAshDensity
	x.FoliageSnow = bd.FoliageSnow
	x.Depth = bd.Depth
	x.Scale = bd.Scale
	x.MapWaterColour = bd.MapWaterColour
	x.Rain = bd.Rain
	x.Tags = bd.Tags
	if v, ok := bd.ChunkGeneration.Value(); ok {
		x.ChunkGeneration = protocol.Option((&BiomeChunkGeneration{}).FromLatest(v))
	}
	return *x
}

func (x *BiomeDefinition) ToLatest() protocol.BiomeDefinition {
	ret := protocol.BiomeDefinition{
		NameIndex:   x.NameIndex,
		BiomeID:     x.BiomeID,
		Temperature: x.Temperature,
		Downfall:    x.Downfall,
		//RedSporeDensity:  x.RedSporeDensity,
		//BlueSporeDensity: x.BlueSporeDensity,
		//AshDensity:       x.AshDensity,
		//WhiteAshDensity:  x.WhiteAshDensity,
		FoliageSnow:    x.FoliageSnow,
		Depth:          x.Depth,
		Scale:          x.Scale,
		MapWaterColour: x.MapWaterColour,
		Rain:           x.Rain,
		Tags:           x.Tags,
	}
	if v, ok := x.ChunkGeneration.Value(); ok {
		ret.ChunkGeneration = protocol.Option(v.ToLatest())
	}
	return ret
}

// BiomeChunkGeneration represents the information required for the client to generate chunks itself
// to create the illusion of a larger render distance.
type BiomeChunkGeneration struct {
	// Climate is optional information to specify the biome's climate.
	Climate protocol.Optional[protocol.BiomeClimate]
	// ConsolidatedFeatures is a list of features that are consolidated into a single feature.
	ConsolidatedFeatures protocol.Optional[[]protocol.BiomeConsolidatedFeature]
	// MountainParameters is optional information to specify the biome's mountain parameters.
	MountainParameters protocol.Optional[protocol.BiomeMountainParameters]
	// SurfaceMaterialAdjustments is a list of surface material adjustments.
	SurfaceMaterialAdjustments protocol.Optional[[]protocol.BiomeElementData]
	// SurfaceMaterials is a set of materials to use for the surface layers of the biome.
	SurfaceMaterials protocol.Optional[protocol.BiomeSurfaceMaterial]
	// HasDefaultOverworldSurface is true if the biome has a default overworld surface.
	HasDefaultOverworldSurface bool
	// HasSwampSurface is true if the biome has a swamp surface.
	HasSwampSurface bool
	// HasFrozenOceanSurface is true if the biome has a frozen ocean surface.
	HasFrozenOceanSurface bool
	// HasEndSurface is true if the biome has an end surface.
	HasEndSurface bool
	// MesaSurface is optional information to specify the biome's mesa surface.
	MesaSurface protocol.Optional[protocol.BiomeMesaSurface]
	// CappedSurface is optional information to specify the biome's capped surface, i.e. in the Nether.
	CappedSurface protocol.Optional[protocol.BiomeCappedSurface]
	// OverworldRules is optional information to specify the biome's overworld rules, such as rivers and hills.
	OverworldRules protocol.Optional[protocol.BiomeOverworldRules]
	// MultiNoiseRules is optional information to specify the biome's multi-noise rules.
	MultiNoiseRules protocol.Optional[protocol.BiomeMultiNoiseRules]
	// LegacyRules is a list of legacy rules for the biomes using an older format, which is just a list of
	// weighted biomes.
	LegacyRules protocol.Optional[[]protocol.BiomeConditionalTransformation]
	// ReplacementsData is a list of biome replacement data.
	ReplacementsData protocol.Optional[[]protocol.BiomeReplacementData]
	// VillageType is an optional village type for biome chunk generation.
	VillageType protocol.Optional[uint8]
}

func (x *BiomeChunkGeneration) Marshal(r protocol.IO) {
	protocol.OptionalMarshaler(r, &x.Climate)
	protocol.OptionalFunc(r, &x.ConsolidatedFeatures, func(s *[]protocol.BiomeConsolidatedFeature) {
		protocol.Slice(r, s)
	})
	protocol.OptionalMarshaler(r, &x.MountainParameters)
	protocol.OptionalFunc(r, &x.SurfaceMaterialAdjustments, func(s *[]protocol.BiomeElementData) {
		protocol.Slice(r, s)
	})
	protocol.OptionalMarshaler(r, &x.SurfaceMaterials)
	r.Bool(&x.HasDefaultOverworldSurface)
	r.Bool(&x.HasSwampSurface)
	r.Bool(&x.HasFrozenOceanSurface)
	r.Bool(&x.HasEndSurface)
	protocol.OptionalMarshaler(r, &x.MesaSurface)
	protocol.OptionalMarshaler(r, &x.CappedSurface)
	protocol.OptionalMarshaler(r, &x.OverworldRules)
	protocol.OptionalMarshaler(r, &x.MultiNoiseRules)
	protocol.OptionalFunc(r, &x.LegacyRules, func(s *[]protocol.BiomeConditionalTransformation) {
		protocol.Slice(r, s)
	})
	if IsProtoGTE(r, ID859) {
		protocol.OptionalFunc(r, &x.ReplacementsData, func(s *[]protocol.BiomeReplacementData) {
			protocol.Slice(r, s)
		})
	}
	if IsProtoGTE(r, ID924) {
		protocol.OptionalFunc(r, &x.VillageType, r.Uint8)
	}
}

func (x *BiomeChunkGeneration) FromLatest(bcg protocol.BiomeChunkGeneration) BiomeChunkGeneration {
	x.Climate = bcg.Climate
	x.ConsolidatedFeatures = bcg.ConsolidatedFeatures
	x.MountainParameters = bcg.MountainParameters
	x.SurfaceMaterialAdjustments = bcg.SurfaceMaterialAdjustments
	x.SurfaceMaterials = bcg.SurfaceMaterials
	x.HasDefaultOverworldSurface = bcg.HasDefaultOverworldSurface
	x.HasSwampSurface = bcg.HasSwampSurface
	x.HasFrozenOceanSurface = bcg.HasFrozenOceanSurface
	x.HasEndSurface = bcg.HasEndSurface
	x.MesaSurface = bcg.MesaSurface
	x.CappedSurface = bcg.CappedSurface
	x.OverworldRules = bcg.OverworldRules
	x.MultiNoiseRules = bcg.MultiNoiseRules
	x.LegacyRules = bcg.LegacyRules
	x.ReplacementsData = bcg.ReplacementsData
	x.VillageType = bcg.VillageType
	return *x
}

func (x *BiomeChunkGeneration) ToLatest() protocol.BiomeChunkGeneration {
	return protocol.BiomeChunkGeneration{
		Climate:                    x.Climate,
		ConsolidatedFeatures:       x.ConsolidatedFeatures,
		MountainParameters:         x.MountainParameters,
		SurfaceMaterialAdjustments: x.SurfaceMaterialAdjustments,
		SurfaceMaterials:           x.SurfaceMaterials,
		HasDefaultOverworldSurface: x.HasDefaultOverworldSurface,
		HasSwampSurface:            x.HasSwampSurface,
		HasFrozenOceanSurface:      x.HasFrozenOceanSurface,
		HasEndSurface:              x.HasEndSurface,
		MesaSurface:                x.MesaSurface,
		CappedSurface:              x.CappedSurface,
		OverworldRules:             x.OverworldRules,
		MultiNoiseRules:            x.MultiNoiseRules,
		LegacyRules:                x.LegacyRules,
		ReplacementsData:           x.ReplacementsData,
		VillageType:                x.VillageType,
	}
}
