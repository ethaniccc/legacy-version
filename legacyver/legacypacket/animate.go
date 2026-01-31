package legacypacket

import (
	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Animate is sent by the server to send a player animation from one player to all viewers of that player. It
// is used for a couple of actions, such as arm swimming and critical hits.
type Animate struct {
	// ActionType is the ID of the animation action to execute. It is one of the action type constants that
	// may be found above.
	ActionType uint8
	// EntityRuntimeID is the runtime ID of the player that the animation should be played upon. The runtime
	// ID is unique for each world session, and entities are generally identified in packets using this
	// runtime ID.
	EntityRuntimeID uint64
	// Data ...
	Data float32
	// SwingSource is the source for swing actions. It is one of the action type constants that
	// may be found above.
	SwingSource protocol.Optional[string]
}

// ID ...
func (*Animate) ID() uint32 {
	return packet.IDAnimate
}

func (pk *Animate) Marshal(io protocol.IO) {
	if proto.IsProtoGTE(io, proto.ID898) {
		io.Uint8(&pk.ActionType)
	} else {
		v := int32(pk.ActionType)
		io.Varint32(&v)
		pk.ActionType = uint8(v)
	}
	io.Varuint64(&pk.EntityRuntimeID)
	if proto.IsProtoGTE(io, proto.ID859) || (pk.ActionType == 128 || pk.ActionType == 127) {
		io.Float32(&pk.Data)
	}
	if proto.IsProtoGTE(io, proto.ID898) {
		protocol.OptionalFunc(io, &pk.SwingSource, io.String)
	}
}
