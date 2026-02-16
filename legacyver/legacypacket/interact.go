package legacypacket

import (
	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Interact is sent by the client when it interacts with another entity in some way. It used to be used for
// normal entity and block interaction, but this is no longer the case now.
type Interact struct {
	// Action type is the ID of the action that was executed by the player. It is one of the constants that
	// may be found above.
	ActionType byte
	// TargetEntityRuntimeID is the runtime ID of the entity that the player interacted with. This is empty
	// for the InteractActionOpenInventory action type.
	TargetEntityRuntimeID uint64
	// Position associated with the ActionType above. For the InteractActionMouseOverEntity, this is the
	// position relative to the entity moused over over which the player hovered with its mouse/touch. For the
	// InteractActionLeaveVehicle, this is the position that the player spawns at after leaving the vehicle.
	Position protocol.Optional[mgl32.Vec3]
}

// ID ...
func (*Interact) ID() uint32 {
	return packet.IDInteract
}

func (pk *Interact) Marshal(io protocol.IO) {
	io.Uint8(&pk.ActionType)
	io.Varuint64(&pk.TargetEntityRuntimeID)
	if proto.IsProtoGTE(io, proto.ID898) {
		protocol.OptionalFunc(io, &pk.Position, io.Vec3)
	} else if pk.ActionType == packet.InteractActionMouseOverEntity || pk.ActionType == packet.InteractActionLeaveVehicle {
		pos, _ := pk.Position.Value()
		io.Vec3(&pos)
		if pos != (mgl32.Vec3{}) {
			pk.Position = protocol.Option(pos)
		} else {
			pk.Position = protocol.Optional[mgl32.Vec3]{}
		}
	}
}
