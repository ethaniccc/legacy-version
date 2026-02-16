package proto

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// CameraPreset represents a basic preset that can be extended upon by more complex instructions.
type CameraPreset struct {
	// Name is the name of the preset. Each preset must have their own unique name.
	Name string
	// Parent is the name of the preset that this preset extends upon. This can be left empty.
	Parent string
	// PosX is the default X position of the camera.
	PosX protocol.Optional[float32]
	// PosY is the default Y position of the camera.
	PosY protocol.Optional[float32]
	// PosZ is the default Z position of the camera.
	PosZ protocol.Optional[float32]
	// RotX is the default pitch of the camera.
	RotX protocol.Optional[float32]
	// RotY is the default yaw of the camera.
	RotY protocol.Optional[float32]
	// RotationSpeed is the speed at which the camera should rotate.
	RotationSpeed protocol.Optional[float32]
	// SnapToTarget determines whether the camera should snap to the target entity or not.
	SnapToTarget protocol.Optional[bool]
	// HorizontalRotationLimit is the horizontal rotation limit of the camera.
	HorizontalRotationLimit protocol.Optional[mgl32.Vec2]
	// VerticalRotationLimit is the vertical rotation limit of the camera.
	VerticalRotationLimit protocol.Optional[mgl32.Vec2]
	// ContinueTargeting determines whether the camera should continue targeting when using aim assist.
	ContinueTargeting protocol.Optional[bool]
	// TrackingRadius is the radius around the camera that the aim assist should track targets.
	TrackingRadius protocol.Optional[float32]
	// ViewOffset is only used in a follow_orbit camera and controls an offset based on a pivot point to the
	// player, causing it to be shifted in a certain direction.
	ViewOffset protocol.Optional[mgl32.Vec2]
	// EntityOffset controls the offset from the entity that the camera should be rendered at.
	EntityOffset protocol.Optional[mgl32.Vec3]
	// Radius is only used in a follow_orbit camera and controls how far away from the player the camera should
	// be rendered.
	Radius protocol.Optional[float32]
	// MinYawLimit is the minimum yaw limit of the camera.
	MinYawLimit protocol.Optional[float32]
	// MaxYawLimit is the maximum yaw limit of the camera.
	MaxYawLimit protocol.Optional[float32]
	// AudioListener defines where the audio should be played from when using this preset. This is one of the
	// constants above.
	AudioListener protocol.Optional[byte]
	// PlayerEffects is currently unknown.
	PlayerEffects protocol.Optional[bool]
	// AlignTargetAndCameraForward determines whether the camera should align the target and the camera forward
	// or not.
	AlignTargetAndCameraForward protocol.Optional[bool]
	// AimAssist defines the aim assist to use when using this preset.
	AimAssist protocol.Optional[protocol.CameraPresetAimAssist]
	// ControlScheme is the control scheme that the client should use in this camera. It is one of the following:
	//  - ControlSchemeLockedPlayerRelativeStrafe is the default behaviour, this cannot be set when the client
	//    is in a custom camera.
	//  - ControlSchemeCameraRelative makes movement relative to the camera's transform, with the client's
	//    rotation being relative to the client's movement.
	//  - ControlSchemeCameraRelativeStrafe makes movement relative to the camera's transform, with the
	//    client's rotation being locked.
	//  - ControlSchemePlayerRelative makes movement relative to the player's transform, meaning holding
	//    left/right will make the player turn in a circle.
	//  - ControlSchemePlayerRelativeStrafe makes movement the same as the default behaviour, but can be
	//    used in a custom camera.
	ControlScheme protocol.Optional[byte]
}

func (x *CameraPreset) FromLatest(cp protocol.CameraPreset) CameraPreset {
	x.Name = cp.Name
	x.Parent = cp.Parent
	x.PosX = cp.PosX
	x.PosY = cp.PosY
	x.PosZ = cp.PosZ
	x.RotX = cp.RotX
	x.RotY = cp.RotY
	x.RotationSpeed = cp.RotationSpeed
	x.SnapToTarget = cp.SnapToTarget
	x.HorizontalRotationLimit = cp.HorizontalRotationLimit
	x.VerticalRotationLimit = cp.VerticalRotationLimit
	x.ContinueTargeting = cp.ContinueTargeting
	x.TrackingRadius = cp.TrackingRadius
	x.ViewOffset = cp.ViewOffset
	x.EntityOffset = cp.EntityOffset
	x.Radius = cp.Radius
	x.MinYawLimit = cp.MinYawLimit
	x.MaxYawLimit = cp.MaxYawLimit
	x.AudioListener = cp.AudioListener
	x.PlayerEffects = cp.PlayerEffects
	x.AimAssist = cp.AimAssist
	x.ControlScheme = cp.ControlScheme
	return *x
}

func (x *CameraPreset) ToLatest() protocol.CameraPreset {
	return protocol.CameraPreset{
		Name:                    x.Name,
		Parent:                  x.Parent,
		PosX:                    x.PosX,
		PosY:                    x.PosY,
		PosZ:                    x.PosZ,
		RotX:                    x.RotX,
		RotY:                    x.RotY,
		RotationSpeed:           x.RotationSpeed,
		SnapToTarget:            x.SnapToTarget,
		HorizontalRotationLimit: x.HorizontalRotationLimit,
		VerticalRotationLimit:   x.VerticalRotationLimit,
		ContinueTargeting:       x.ContinueTargeting,
		TrackingRadius:          x.TrackingRadius,
		ViewOffset:              x.ViewOffset,
		EntityOffset:            x.EntityOffset,
		Radius:                  x.Radius,
		MinYawLimit:             x.MinYawLimit,
		MaxYawLimit:             x.MaxYawLimit,
		AudioListener:           x.AudioListener,
		PlayerEffects:           x.PlayerEffects,
		AimAssist:               x.AimAssist,
		ControlScheme:           x.ControlScheme,
	}
}

// Marshal encodes/decodes a CameraPreset.
func (x *CameraPreset) Marshal(r protocol.IO) {
	r.String(&x.Name)
	r.String(&x.Parent)
	protocol.OptionalFunc(r, &x.PosX, r.Float32)
	protocol.OptionalFunc(r, &x.PosY, r.Float32)
	protocol.OptionalFunc(r, &x.PosZ, r.Float32)
	protocol.OptionalFunc(r, &x.RotX, r.Float32)
	protocol.OptionalFunc(r, &x.RotY, r.Float32)
	if IsProtoGTE(r, ID729) {
		protocol.OptionalFunc(r, &x.RotationSpeed, r.Float32)
		protocol.OptionalFunc(r, &x.SnapToTarget, r.Bool)
	}
	if IsProtoGTE(r, ID748) {
		protocol.OptionalFunc(r, &x.HorizontalRotationLimit, r.Vec2)
		protocol.OptionalFunc(r, &x.VerticalRotationLimit, r.Vec2)
		protocol.OptionalFunc(r, &x.ContinueTargeting, r.Bool)
	}
	if IsProtoGTE(r, ID766) {
		protocol.OptionalFunc(r, &x.TrackingRadius, r.Float32)
	}
	if IsProtoGTE(r, ID776) {
		protocol.OptionalFunc(r, &x.MinYawLimit, r.Float32)
		protocol.OptionalFunc(r, &x.MaxYawLimit, r.Float32)
	}
	protocol.OptionalFunc(r, &x.ViewOffset, r.Vec2)
	if IsProtoGTE(r, ID729) {
		protocol.OptionalFunc(r, &x.EntityOffset, r.Vec3)
	}
	protocol.OptionalFunc(r, &x.Radius, r.Float32)
	protocol.OptionalFunc(r, &x.AudioListener, r.Uint8)
	protocol.OptionalFunc(r, &x.PlayerEffects, r.Bool)
	if IsProtoGTE(r, ID748) && IsProtoLT(r, ID818) {
		protocol.OptionalFunc(r, &x.AlignTargetAndCameraForward, r.Bool)
	}
	if IsProtoGTE(r, ID766) {
		protocol.OptionalMarshaler(r, &x.AimAssist)
	}
	if IsProtoGTE(r, ID800) {
		protocol.OptionalFunc(r, &x.ControlScheme, r.Uint8)
	}
}

// CameraInstructionSet represents a camera instruction that sets the camera to a specified preset and can be extended
// with easing functions and translations to the camera's position and rotation.
type CameraInstructionSet struct {
	// Preset is the index of the preset in the CameraPresets packet sent to the player.
	Preset uint32
	// Ease represents the easing function that is used by the instruction.
	Ease protocol.Optional[protocol.CameraEase]
	// Position represents the position of the camera.
	Position protocol.Optional[mgl32.Vec3]
	// Rotation represents the rotation of the camera.
	Rotation protocol.Optional[mgl32.Vec2]
	// Facing is a vector that the camera will always face towards during the duration of the instruction.
	Facing protocol.Optional[mgl32.Vec3]
	// ViewOffset is an offset based on a pivot point to the player, causing the camera to be shifted in a
	// certain direction.
	ViewOffset protocol.Optional[mgl32.Vec2]
	// EntityOffset is an offset from the entity that the camera should be rendered at.
	EntityOffset protocol.Optional[mgl32.Vec3]
	// Default determines whether the camera is a default camera or not.
	Default protocol.Optional[bool]
	// IgnoreStartingValuesComponent behavior is currently unknown.
	IgnoreStartingValuesComponent bool
}

func (x *CameraInstructionSet) FromLatest(cis protocol.CameraInstructionSet) CameraInstructionSet {
	x.Preset = cis.Preset
	x.Ease = cis.Ease
	x.Position = cis.Position
	x.Rotation = cis.Rotation
	x.Facing = cis.Facing
	x.ViewOffset = cis.ViewOffset
	x.EntityOffset = cis.EntityOffset
	x.Default = cis.Default
	x.IgnoreStartingValuesComponent = cis.IgnoreStartingValuesComponent
	return *x
}

func (x *CameraInstructionSet) ToLatest() protocol.CameraInstructionSet {
	return protocol.CameraInstructionSet{
		Preset:                        x.Preset,
		Ease:                          x.Ease,
		Position:                      x.Position,
		Rotation:                      x.Rotation,
		Facing:                        x.Facing,
		ViewOffset:                    x.ViewOffset,
		EntityOffset:                  x.EntityOffset,
		Default:                       x.Default,
		IgnoreStartingValuesComponent: x.IgnoreStartingValuesComponent,
	}
}

// Marshal encodes/decodes a CameraInstructionSet.
func (x *CameraInstructionSet) Marshal(r protocol.IO) {
	r.Uint32(&x.Preset)
	protocol.OptionalMarshaler(r, &x.Ease)
	protocol.OptionalFunc(r, &x.Position, r.Vec3)
	protocol.OptionalFunc(r, &x.Rotation, r.Vec2)
	protocol.OptionalFunc(r, &x.Facing, r.Vec3)
	protocol.OptionalFunc(r, &x.ViewOffset, r.Vec2)
	protocol.OptionalFunc(r, &x.EntityOffset, r.Vec3)
	protocol.OptionalFunc(r, &x.Default, r.Bool)
	if IsProtoGTE(r, ID818) {
		r.Bool(&x.IgnoreStartingValuesComponent)
	}
}
