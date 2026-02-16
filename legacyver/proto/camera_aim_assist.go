package proto

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// CameraAimAssistCategory is an aim assist category that defines priorities for specific blocks and entities.
type CameraAimAssistCategory struct {
	// Name is the name of the category which can be used by a CameraAimAssistPreset.
	Name string
	// Priorities represents the block and entity specific priorities as well as the default priorities for
	// this category.
	Priorities CameraAimAssistPriorities
}

// Marshal encodes/decodes a CameraAimAssistCategory.
func (x *CameraAimAssistCategory) Marshal(r protocol.IO) {
	r.String(&x.Name)
	protocol.Single(r, &x.Priorities)
}

func (x *CameraAimAssistCategory) FromLatest(v protocol.CameraAimAssistCategory) CameraAimAssistCategory {
	x.Name = v.Name
	x.Priorities = (&CameraAimAssistPriorities{}).FromLatest(v.Priorities)
	return *x
}

func (x *CameraAimAssistCategory) ToLatest() protocol.CameraAimAssistCategory {
	return protocol.CameraAimAssistCategory{
		Name:       x.Name,
		Priorities: x.Priorities.ToLatest(),
	}
}

// CameraAimAssistPriorities represents block and entity priorities for targeting.
type CameraAimAssistPriorities struct {
	// Entities is a list of priorities for specific entity identifiers.
	Entities []protocol.CameraAimAssistPriority
	// Blocks is a list of priorities for specific block identifiers.
	Blocks []protocol.CameraAimAssistPriority
	// BlockTags is a list of priorities for specific block tags.
	BlockTags []protocol.CameraAimAssistPriority
	// EntityDefault is the default priority for entities.
	EntityDefault protocol.Optional[int32]
	// BlockDefault is the default priority for blocks.
	BlockDefault protocol.Optional[int32]
}

// Marshal encodes/decodes a CameraAimAssistPriorities.
func (x *CameraAimAssistPriorities) Marshal(r protocol.IO) {
	protocol.Slice(r, &x.Entities)
	protocol.Slice(r, &x.Blocks)
	protocol.Slice(r, &x.BlockTags)
	protocol.OptionalFunc(r, &x.EntityDefault, r.Int32)
	protocol.OptionalFunc(r, &x.BlockDefault, r.Int32)
}

func (x *CameraAimAssistPriorities) FromLatest(v protocol.CameraAimAssistPriorities) CameraAimAssistPriorities {
	x.Entities = v.Entities
	x.Blocks = v.Blocks
	x.BlockTags = v.BlockTags
	x.EntityDefault = v.EntityDefault
	x.BlockDefault = v.BlockDefault
	return *x
}

func (x *CameraAimAssistPriorities) ToLatest() protocol.CameraAimAssistPriorities {
	return protocol.CameraAimAssistPriorities{
		Entities:      x.Entities,
		Blocks:        x.Blocks,
		BlockTags:     x.BlockTags,
		EntityDefault: x.EntityDefault,
		BlockDefault:  x.BlockDefault,
	}
}

// CameraAimAssistPreset defines a base preset that can be extended upon when sending an aim assist.
type CameraAimAssistPreset struct {
	// Identifier represents the identifier of this preset.
	Identifier string
	// BlockExclusions is a list of block identifiers that should be ignored by the aim assist.
	BlockExclusions []string
	// EntityExclusions is a list of entity identifiers that should be ignored by the aim assist.
	EntityExclusions []string
	// BlockTagExclusions is a list of block tags that should be ignored by the aim assist.
	BlockTagExclusions []string
	// LiquidTargets is a list of entity identifiers that should be targetted when inside of a liquid.
	LiquidTargets []string
	// ItemSettings is a list of settings for specific item identifiers.
	ItemSettings []protocol.CameraAimAssistItemSettings
	// DefaultItemSettings is the identifier of a category to use when a mapped item is not held.
	DefaultItemSettings protocol.Optional[string]
	// HandSettings is the identifier of a category to use when no item is held.
	HandSettings protocol.Optional[string]
}

// Marshal encodes/decodes a CameraAimAssistPreset.
func (x *CameraAimAssistPreset) Marshal(r protocol.IO) {
	r.String(&x.Identifier)
	protocol.FuncSlice(r, &x.BlockExclusions, r.String)
	protocol.FuncSlice(r, &x.EntityExclusions, r.String)
	protocol.FuncSlice(r, &x.BlockTagExclusions, r.String)
	protocol.FuncSlice(r, &x.LiquidTargets, r.String)
	protocol.Slice(r, &x.ItemSettings)
	protocol.OptionalFunc(r, &x.DefaultItemSettings, r.String)
	protocol.OptionalFunc(r, &x.HandSettings, r.String)
}

func (x *CameraAimAssistPreset) FromLatest(v protocol.CameraAimAssistPreset) CameraAimAssistPreset {
	x.Identifier = v.Identifier
	x.BlockExclusions = v.BlockExclusions
	x.EntityExclusions = v.EntityExclusions
	x.BlockTagExclusions = v.BlockTagExclusions
	x.LiquidTargets = v.LiquidTargets
	x.ItemSettings = v.ItemSettings
	x.DefaultItemSettings = v.DefaultItemSettings
	x.HandSettings = v.HandSettings
	return *x
}

func (x *CameraAimAssistPreset) ToLatest() protocol.CameraAimAssistPreset {
	return protocol.CameraAimAssistPreset{
		Identifier:          x.Identifier,
		BlockExclusions:     x.BlockExclusions,
		EntityExclusions:    x.EntityExclusions,
		BlockTagExclusions:  x.BlockTagExclusions,
		LiquidTargets:       x.LiquidTargets,
		ItemSettings:        x.ItemSettings,
		DefaultItemSettings: x.DefaultItemSettings,
		HandSettings:        x.HandSettings,
	}
}

// CameraRotationOption represents a rotation key frame for camera spline instructions.
type CameraRotationOption struct {
	// Value is the rotation value.
	Value mgl32.Vec3
	// Time is the time for this rotation option.
	Time float32
}

// Marshal encodes/decodes a CameraRotationOption.
func (x *CameraRotationOption) Marshal(r protocol.IO) {
	r.Vec3(&x.Value)
	r.Float32(&x.Time)
}

func (x *CameraRotationOption) FromLatest(v protocol.CameraRotationOption) CameraRotationOption {
	x.Value = v.Value
	x.Time = v.Time
	return *x
}

func (x *CameraRotationOption) ToLatest() protocol.CameraRotationOption {
	return protocol.CameraRotationOption{
		Value: x.Value,
		Time:  x.Time,
	}
}

// CameraSplineInstruction represents a camera instruction that creates a spline path.
type CameraSplineInstruction struct {
	// TotalTime is the total time for the spline animation.
	TotalTime float32
	// EaseType is the easing function type used for interpolation.
	EaseType uint8
	// Curve is a list of points that define the spline curve.
	Curve []mgl32.Vec3
	// ProgressKeyFrames is a list of key frames for the spline progress.
	ProgressKeyFrames []mgl32.Vec2
	// RotationOptions is a list of rotation options for the spline.
	RotationOptions []CameraRotationOption
}

// Marshal encodes/decodes a CameraSplineInstruction.
func (x *CameraSplineInstruction) Marshal(r protocol.IO) {
	r.Float32(&x.TotalTime)
	r.Uint8(&x.EaseType)
	protocol.FuncSlice(r, &x.Curve, r.Vec3)
	protocol.FuncSlice(r, &x.ProgressKeyFrames, r.Vec2)
	protocol.Slice(r, &x.RotationOptions)
}

func (x *CameraSplineInstruction) FromLatest(v protocol.CameraSplineInstruction) CameraSplineInstruction {
	x.TotalTime = v.TotalTime
	if ease, ok := v.SplineType.Value(); ok {
		x.EaseType = ease
	}
	x.Curve = v.Curve
	x.ProgressKeyFrames = make([]mgl32.Vec2, len(v.ProgressKeyFrames))
	for i, kf := range v.ProgressKeyFrames {
		x.ProgressKeyFrames[i] = mgl32.Vec2{kf.Value, kf.Time}
	}
	x.RotationOptions = make([]CameraRotationOption, len(v.RotationOptions))
	for i, ro := range v.RotationOptions {
		x.RotationOptions[i] = (&CameraRotationOption{}).FromLatest(ro)
	}
	return *x
}

func (x *CameraSplineInstruction) ToLatest() protocol.CameraSplineInstruction {
	progressKeyFrames := make([]protocol.CameraProgressOption, len(x.ProgressKeyFrames))
	for i, kf := range x.ProgressKeyFrames {
		progressKeyFrames[i] = protocol.CameraProgressOption{
			Value: kf[0],
			Time:  kf[1],
		}
	}
	rotationOptions := make([]protocol.CameraRotationOption, len(x.RotationOptions))
	for i, ro := range x.RotationOptions {
		rotationOptions[i] = ro.ToLatest()
	}
	return protocol.CameraSplineInstruction{
		TotalTime:         x.TotalTime,
		SplineType:        protocol.Option(x.EaseType),
		Curve:             x.Curve,
		ProgressKeyFrames: progressKeyFrames,
		RotationOptions:   rotationOptions,
	}
}
