package legacypacket

import (
	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// SetHud is sent by the server to set the visibility of individual HUD elements on the client. It is
// important to note that the client does not reset the state of the HUD elements after it leaves a server,
// meaning they can leak into sessions on different servers. To be safe, you should reset the visibility of
// all HUD elements when a player connects.
type SetHud struct {
	// Elements is a list of HUD elements that are being modified. The values can be any of the HudElement
	// constants above.
	Elements []int32
	// Visibility represents the new visibility of the specified Elements. It can be any of the HudVisibility
	// constants above.
	Visibility int32
}

// ID ...
func (*SetHud) ID() uint32 {
	return packet.IDSetHud
}

func (pk *SetHud) Marshal(io protocol.IO) {
	if proto.IsProtoGTE(io, proto.ID786) {
		protocol.FuncSlice(io, &pk.Elements, io.Varint32)
		io.Varint32(&pk.Visibility)
	} else {
		elements := make([]uint8, len(pk.Elements))
		for i, v := range pk.Elements {
			elements[i] = uint8(v)
		}
		protocol.FuncSlice(io, &elements, io.Uint8)
		for i, v := range elements {
			pk.Elements[i] = int32(v)
		}

		visibility := uint8(pk.Visibility)
		io.Uint8(&visibility)
		pk.Visibility = int32(visibility)
	}
}
