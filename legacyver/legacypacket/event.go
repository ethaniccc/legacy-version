package legacypacket

import (
	"github.com/ethaniccc/legacy-version/internal/typeconf"
	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Event is sent by the server to send an event with additional data. It is typically sent to the client for
// telemetry reasons, much like the SimpleEvent packet.
type Event struct {
	// EntityRuntimeID is the runtime ID of the player. The runtime ID is unique for each world session, and
	// entities are generally identified in packets using this runtime ID.
	EntityRuntimeID int64
	// UsePlayerID ...
	// TODO: Figure out what UsePlayerID is for.
	UsePlayerID bool
	// Event is the event that is transmitted.
	Event protocol.Event
}

// ID ...
func (*Event) ID() uint32 {
	return packet.IDEvent
}

func (pk *Event) Marshal(io protocol.IO) {
	io.Varint64(&pk.EntityRuntimeID)
	io.EventType(&pk.Event)
	if proto.IsProtoGTE(io, proto.ID898) {
		io.Bool(&pk.UsePlayerID)
	} else {
		v := typeconf.BoolToByte(pk.UsePlayerID)
		io.Uint8(&v)
		pk.UsePlayerID = typeconf.ByteToBool(v)
	}
	pk.Event.Marshal(io)
}
