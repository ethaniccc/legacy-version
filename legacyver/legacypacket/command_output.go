package legacypacket

import (
	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// CommandOutput is sent by the server to the client to send text as output of a command. Most servers do not
// use this packet and instead simply send Text packets, but there is reason to send it.
// If the origin of a CommandRequest packet is not the player itself, but, for example, a websocket server,
// sending a Text packet will not do what is expected: The message should go to the websocket server, not to
// the client's chat. The CommandOutput packet will make sure the messages are relayed to the correct origin
// of the command request.
type CommandOutput struct {
	// CommandOrigin is the data specifying the origin of the command. In other words, the source that the
	// command request was from, such as the player itself or a websocket server. The client forwards the
	// messages in this packet to the right origin, depending on what is sent here.
	CommandOrigin protocol.CommandOrigin
	// OutputType specifies the type of output that is sent. The OutputType sent by vanilla games appears to
	// be 3, which seems to work.
	OutputType string
	// SuccessCount is the amount of times that a command was executed successfully as a result of the command
	// that was requested. For servers, this is usually a rather meaningless fields, but for vanilla, this is
	// applicable for commands created with Functions.
	SuccessCount uint32
	// OutputMessages is a list of all output messages that should be sent to the player. Whether they are
	// shown or not, depends on the type of the messages.
	OutputMessages []proto.CommandOutputMessage
	// DataSet ... TODO: Find out what this is for.
	DataSet protocol.Optional[string]
}

// ID ...
func (*CommandOutput) ID() uint32 {
	return packet.IDCommandOutput
}

func (pk *CommandOutput) Marshal(io protocol.IO) {
	proto.CommandOriginData(io, &pk.CommandOrigin)
	if proto.IsProtoGTE(io, proto.ID898) {
		io.String(&pk.OutputType)
		io.Uint32(&pk.SuccessCount)
	} else {
		v := LegacyCommandOutputType(pk.OutputType)
		io.Uint8(&v)
		pk.OutputType = LatestCommandOutputType(v)
		io.Varuint32(&pk.SuccessCount)
	}
	protocol.Slice(io, &pk.OutputMessages)
	if proto.IsProtoGTE(io, proto.ID898) {
		protocol.OptionalFunc(io, &pk.DataSet, io.String)
	} else if pk.OutputType == "dataset" {
		v, _ := pk.DataSet.Value()
		io.String(&v)
		if v != "" {
			pk.DataSet = protocol.Option(v)
		} else {
			pk.DataSet = protocol.Optional[string]{}
		}
	}
}

var legacyCommandOutputTypeMap = map[byte]string{
	0: "none",       // packet.CommandOutputTypeNone
	1: "lastoutput", // packet.CommandOutputTypeLastOutput
	2: "silent",     // packet.CommandOutputTypeSilent
	3: "alloutput",  // packet.CommandOutputTypeAllOutput
	4: "dataset",    // packet.CommandOutputTypeDataSet
}

func LegacyCommandOutputType(outputType string) byte {
	for k, v := range legacyCommandOutputTypeMap {
		if v == outputType {
			return k
		}
	}
	return 3 // Default to all output.
}

func LatestCommandOutputType(outputType byte) string {
	if v, ok := legacyCommandOutputTypeMap[outputType]; ok {
		return v
	}
	return "alloutput" // packet.CommandOutputTypeAllOutput
}
