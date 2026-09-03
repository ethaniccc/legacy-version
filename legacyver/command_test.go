package legacyver

import (
	"bytes"
	"reflect"
	"testing"

	legacyproto "github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// TestAvailableCommandsParameterTypes checks that command parameter types
// survive the legacy encoding at every supported 1.26 protocol. The packet
// mirrors what play.galaxite.net sends: enum parameters whose low bits are a
// table index (index 0 included), a basic type outside the old whitelist
// (MessageRoot, 68), a float, a soft enum and a suffixed integer. The old
// translation panicked on the enum index 0 and on type 68.
func TestAvailableCommandsParameterTypes(t *testing.T) {
	pk := &packet.AvailableCommands{
		EnumValues: []string{"en_US", "de_DE", "true", "false"},
		Suffixes:   []string{"L", "s"},
		Enums: []protocol.CommandEnum{
			{Type: "language", ValueIndices: []uint32{0, 1}},
			{Type: "bool", ValueIndices: []uint32{2, 3}},
		},
		DynamicEnums: []protocol.DynamicEnum{{Type: "names", Values: []string{"a"}}},
		Commands: []protocol.Command{
			{
				Name: "language", Description: "d", AliasesOffset: 0xffffffff,
				Overloads: []protocol.CommandOverload{{Parameters: []protocol.CommandParameter{
					{Name: "language", Type: protocol.CommandArgValid | protocol.CommandArgEnum | 0},
					{Name: "flag", Type: protocol.CommandArgValid | protocol.CommandArgEnum | 1, Optional: true},
				}}},
			},
			{
				Name: "party", Description: "d", AliasesOffset: 0xffffffff,
				Overloads: []protocol.CommandOverload{{Parameters: []protocol.CommandParameter{
					{Name: "player", Type: protocol.CommandArgValid | protocol.CommandArgTypeMessageRoot},
					{Name: "amount", Type: protocol.CommandArgValid | protocol.CommandArgTypeFloat},
					{Name: "name", Type: protocol.CommandArgValid | protocol.CommandArgSoftEnum | 0},
					{Name: "time", Type: protocol.CommandArgValid | protocol.CommandArgSuffixed | 1},
				}}},
			},
		},
	}
	marshalFn, ok := packets[pk.ID()]
	if !ok {
		t.Fatal("AvailableCommands is not registered for compatibility marshaling")
	}
	for _, id := range []int32{legacyproto.ID924, legacyproto.ID944, legacyproto.ID975, legacyproto.ID1001, legacyproto.ID2168} {
		var buf bytes.Buffer
		marshalFn(legacyproto.NewWriter(protocol.NewWriter(&buf, 0), id, 0), pk)

		got := &packet.AvailableCommands{}
		marshalFn(legacyproto.NewReader(protocol.NewReader(&buf, 0, false), id, 0, false), got)
		if buf.Len() != 0 {
			t.Fatalf("protocol %d: %d bytes left after decode", id, buf.Len())
		}
		for i := range pk.Commands {
			want, have := pk.Commands[i].Overloads[0].Parameters, got.Commands[i].Overloads[0].Parameters
			if !reflect.DeepEqual(want, have) {
				t.Fatalf("protocol %d: /%s parameters differ\nwant: %+v\nhave: %+v", id, pk.Commands[i].Name, want, have)
			}
		}
	}
}
