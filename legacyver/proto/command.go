package proto

import (
	"math"
	_ "unsafe"

	"github.com/akmalfairuz/legacy-version/internal/typeconf"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func MarshalCommand(r protocol.IO, c *protocol.Command) {
	r.String(&c.Name)
	r.String(&c.Description)
	r.Uint16(&c.Flags)

	if IsProtoGTE(r, ID898) {
		permLevel := commandPermissionToString(c.PermissionLevel)
		r.String(&permLevel)
		commandPermissionFromString(r, &c.PermissionLevel, permLevel)
	} else {
		r.Uint8(&c.PermissionLevel)
	}
	r.Uint32(&c.AliasesOffset)
	if IsProtoGTE(r, ID898) {
		protocol.FuncSlice(r, &c.ChainedSubcommandOffsets, r.Uint32)
	} else {
		offsets := typeconf.SliceIntToSliceInt[uint32, uint16](c.ChainedSubcommandOffsets)
		protocol.FuncSlice(r, &offsets, r.Uint16)
		c.ChainedSubcommandOffsets = typeconf.SliceIntToSliceInt[uint16, uint32](offsets)
	}
	protocol.FuncIOSlice(r, &c.Overloads, MarshalCommandOverload)
}

func MarshalCommandOverload(r protocol.IO, x *protocol.CommandOverload) {
	r.Bool(&x.Chaining)
	protocol.FuncIOSlice(r, &x.Parameters, MarshalCommandParameter)
}

func MarshalCommandParameter(r protocol.IO, x *protocol.CommandParameter) {
	r.String(&x.Name)
	// Type is written as is for every protocol. The argument type numbering
	// did not change between 1.26.0 and 1.26.44: the gophertunnel constants
	// at protocol 2168 (Int=1, Float=3, RValue=4, ... Command=87) are the
	// wire values older clients use too. gophertunnel briefly renumbered its
	// constants for 1.26.30 and reverted that in #466, and the translation
	// that used to live here was written against the renumbered set.
	//
	// The low 20 bits are only an argument type for basic parameters. For a
	// parameter with CommandArgEnum, CommandArgSoftEnum or CommandArgSuffixed
	// they hold an index into the Enums, DynamicEnums or Suffixes table of
	// the packet, so any translation of the value breaks those parameters.
	r.Uint32(&x.Type)
	r.Bool(&x.Optional)
	r.Uint8(&x.Options)
}

// CommandEnumContext holds context required for encoding command enums.
type CommandEnumContext struct {
	EnumValues []string
}

// Marshal encodes/decodes a CommandEnum.
func (ctx CommandEnumContext) Marshal(r protocol.IO, x *protocol.CommandEnum) {
	r.String(&x.Type)
	if IsProtoGTE(r, ID898) {
		protocol.FuncSlice(r, &x.ValueIndices, r.Uint32)
	} else {
		protocol.FuncIOSlice(r, &x.ValueIndices, ctx.enumOption)
	}
}

// enumOption writes/reads a command enum option as a byte/uint16/uint32,
// depending on the amount of enum values.
func (ctx CommandEnumContext) enumOption(r protocol.IO, opt *uint32) {
	n := len(ctx.EnumValues)
	switch {
	case n <= math.MaxUint8:
		val := byte(*opt)
		r.Uint8(&val)
		*opt = uint32(val)
	case n <= math.MaxUint16:
		val := uint16(*opt)
		r.Uint16(&val)
		*opt = uint32(val)
	default:
		r.Uint32(opt)
	}
}

func MarshalChainedSubcommand(r protocol.IO, x *protocol.ChainedSubcommand) {
	r.String(&x.Name)
	protocol.FuncIOSlice(r, &x.Values, MarshalChainedSubcommandValue)
}

func MarshalChainedSubcommandValue(r protocol.IO, x *protocol.ChainedSubcommandValue) {
	if IsProtoGTE(r, ID898) {
		r.Varuint32(&x.Index)
		r.Varuint32(&x.Value)
	} else {
		v := uint16(x.Index)
		r.Uint16(&v)
		x.Index = uint32(v)
		vv := uint16(x.Value)
		r.Uint16(&vv)
		x.Value = uint32(vv)
	}
}

// CommandOriginData reads/writes a CommandOrigin x using IO r.
func CommandOriginData(r protocol.IO, x *protocol.CommandOrigin) {
	if IsProtoGTE(r, ID898) {
		originStr := commandOriginToString(x.Origin)
		r.String(&originStr)
		commandOriginFromString(r, &x.Origin, originStr)
	} else {
		r.Varuint32(&x.Origin)
	}
	r.UUID(&x.UUID)
	r.String(&x.RequestID)
	if IsProtoGTE(r, ID898) {
		r.Int64(&x.PlayerUniqueID)
	} else if x.Origin == protocol.CommandOriginDevConsole || x.Origin == protocol.CommandOriginTest {
		r.Varint64(&x.PlayerUniqueID)
	}
}

func MarshalCommandOutputMessage(r protocol.IO, x *protocol.CommandOutputMessage) {
	if IsProtoLT(r, ID898) {
		r.Bool(&x.Success)
	}
	r.String(&x.Message)
	if IsProtoGTE(r, ID898) {
		r.Bool(&x.Success)
	}
	protocol.FuncSlice(r, &x.Parameters, r.String)
}

//go:linkname commandPermissionToString github.com/sandertv/gophertunnel/minecraft/protocol.commandPermissionToString
func commandPermissionToString(x byte) string

//go:linkname commandPermissionFromString github.com/sandertv/gophertunnel/minecraft/protocol.commandPermissionFromString
func commandPermissionFromString(r protocol.IO, x *byte, s string)

//go:linkname commandOriginToString github.com/sandertv/gophertunnel/minecraft/protocol.commandOriginToString
func commandOriginToString(x uint32) string

//go:linkname commandOriginFromString github.com/sandertv/gophertunnel/minecraft/protocol.commandOriginFromString
func commandOriginFromString(r protocol.IO, x *uint32, s string)
