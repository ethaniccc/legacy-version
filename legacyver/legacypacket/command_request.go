package legacypacket

import (
	"github.com/ethaniccc/legacy-version/internal/typeconf"
	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// CommandRequest is sent by the client to request the execution of a server-side command. Although some
// servers support sending commands using the Text packet, this packet is guaranteed to have the correct
// result.
type CommandRequest struct {
	// CommandLine is the raw entered command line. The client does no parsing of the command line by itself
	// (unlike it did in the early stages), but lets the server do that.
	CommandLine string
	// CommandOrigin is the data specifying the origin of the command. In other words, the source that the
	// command was from, such as the player itself or a websocket server.
	CommandOrigin protocol.CommandOrigin
	// Internal specifies if the command request internal. Setting it to false seems to work and the usage of
	// this field is not known.
	Internal bool
	// Version is the version of the command that is being executed. This field currently has no purpose or functionality.
	Version string
}

// ID ...
func (*CommandRequest) ID() uint32 {
	return packet.IDCommandRequest
}

func (pk *CommandRequest) Marshal(io protocol.IO) {
	io.String(&pk.CommandLine)
	proto.CommandOriginData(io, &pk.CommandOrigin)
	io.Bool(&pk.Internal)
	if proto.IsProtoGTE(io, proto.ID898) {
		io.String(&pk.Version)
	} else {
		v := typeconf.StringToInt[int32](pk.Version, 0)
		io.Varint32(&v)
		pk.Version = typeconf.IntToString[int32](v, "")
	}
}
