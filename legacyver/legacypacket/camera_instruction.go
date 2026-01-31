package legacypacket

import (
	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// CameraInstruction gives a custom camera specific instructions to operate.
type CameraInstruction struct {
	// Set is a camera instruction that sets the camera to a specified preset.
	Set protocol.Optional[proto.CameraInstructionSet]
	// Clear can be set to true to clear all the current camera instructions.
	Clear protocol.Optional[bool]
	// Fade is a camera instruction that fades the screen to a specified colour.
	Fade protocol.Optional[protocol.CameraInstructionFade]
	// Target is a camera instruction that targets a specific entity.
	Target protocol.Optional[protocol.CameraInstructionTarget]
	// RemoveTarget can be set to true to remove the current aim assist target.
	RemoveTarget protocol.Optional[bool]
	// FieldOfView is a camera instruction that updates the field of view for the camera.
	FieldOfView protocol.Optional[protocol.CameraInstructionFieldOfView]
	// Spline is a camera instruction that creates a spline path for the camera to follow.
	Spline protocol.Optional[protocol.CameraSplineInstruction]
	// AttachToEntity is the entity ID to attach the camera to.
	AttachToEntity protocol.Optional[int64]
	// DetachFromEntity can be set to true to detach the camera from the current entity.
	DetachFromEntity protocol.Optional[bool]
}

// ID ...
func (*CameraInstruction) ID() uint32 {
	return packet.IDCameraInstruction
}

func (pk *CameraInstruction) Marshal(io protocol.IO) {
	protocol.OptionalMarshaler(io, &pk.Set)
	protocol.OptionalFunc(io, &pk.Clear, io.Bool)
	protocol.OptionalMarshaler(io, &pk.Fade)
	if proto.IsProtoGTE(io, proto.ID712) {
		protocol.OptionalMarshaler(io, &pk.Target)
		protocol.OptionalFunc(io, &pk.RemoveTarget, io.Bool)
	}
	if proto.IsProtoGTE(io, proto.ID827) {
		protocol.OptionalMarshaler(io, &pk.FieldOfView)
	}
	if proto.IsProtoGTE(io, proto.ID859) {
		protocol.OptionalMarshaler(io, &pk.Spline)
		protocol.OptionalFunc(io, &pk.AttachToEntity, io.Int64)
		protocol.OptionalFunc(io, &pk.DetachFromEntity, io.Bool)
	}
}
