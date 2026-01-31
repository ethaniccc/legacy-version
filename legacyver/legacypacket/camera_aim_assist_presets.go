package legacypacket

import (
	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	CameraAunAssistPresetOperationSet = iota
	CameraAunAssistPresetOperationAddToExisting
)

// CameraAimAssistPresets is sent by the server to the client to provide a list of categories and presets
// that can be used when sending a CameraAimAssist packet or a CameraInstruction including aim assist.
type CameraAimAssistPresets struct {
	// Categories is a list of groups of categories which can be referenced by one of the Presets.
	Categories []protocol.CameraAimAssistCategory
	// Presets is a list of presets which define a base for how aim assist should behave
	Presets []protocol.CameraAimAssistPreset
	// Operation is the operation to perform with the presets. It is one of the constants above.
	Operation byte
}

// ID ...
func (*CameraAimAssistPresets) ID() uint32 {
	return packet.IDCameraAimAssistPresets
}

func (pk *CameraAimAssistPresets) Marshal(io protocol.IO) {
	protocol.Slice(io, &pk.Categories)
	protocol.Slice(io, &pk.Presets)
	if proto.IsProtoGTE(io, proto.ID776) {
		io.Uint8(&pk.Operation)
	}
}
