package proto

import (
	"math"

	"github.com/ethaniccc/legacy-version/internal/typeconf"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func legacyCommandPermToNew(x byte) string {
	switch x {
	case protocol.CommandPermissionLevelAny:
		return "any"
	case protocol.CommandPermissionLevelGameDirectors:
		return "gamedirectors"
	case protocol.CommandPermissionLevelAdmin:
		return "admin"
	case protocol.CommandPermissionLevelHost:
		return "host"
	case protocol.CommandPermissionLevelOwner:
		return "owner"
	case protocol.CommandPermissionLevelInternal:
		return "internal"
	default:
		return "unknown"
	}
}

func newCommandPermToLegacy(x string) byte {
	switch x {
	case "any":
		return protocol.CommandPermissionLevelAny
	case "gamedirectors":
		return protocol.CommandPermissionLevelGameDirectors
	case "admin":
		return protocol.CommandPermissionLevelAdmin
	case "host":
		return protocol.CommandPermissionLevelHost
	case "owner":
		return protocol.CommandPermissionLevelOwner
	case "internal":
		return protocol.CommandPermissionLevelInternal
	default:
		return protocol.CommandPermissionLevelAny
	}
}

// Command holds the data that a command requires to be shown to a player client-side. The command is shown in
// the /help command and auto-completed using this data.
type Command struct {
	// Name is the name of the command. The command may be executed using this name, and will be shown in the
	// /help list with it. It currently seems that the client crashes if the Name contains uppercase letters.
	Name string
	// Description is the description of the command. It is shown in the /help list and when starting to write
	// a command.
	Description string
	// Flags is a combination of flags not currently known. Leaving the Flags field empty appears to work.
	Flags uint16
	// PermissionLevel is the command permission level that the player required to execute this command. The
	// field no longer seems to serve a purpose, as the client does not handle the execution of commands
	// anymore: The permissions should be checked server-side.
	PermissionLevel string
	// AliasesOffset is the offset to a CommandEnum that holds the values that
	// should be used as aliases for this command.
	AliasesOffset uint32
	// ChainedSubcommandOffsets is a slice of offsets that all point to a different ChainedSubcommand from the
	// ChainedSubcommands slice in the AvailableCommands packet.
	ChainedSubcommandOffsets []uint32
	// Overloads is a list of command overloads that specify the ways in which a command may be executed. The
	// overloads may be completely different.
	Overloads []protocol.CommandOverload
}

func (c *Command) Marshal(r protocol.IO) {
	r.String(&c.Name)
	r.String(&c.Description)
	r.Uint16(&c.Flags)

	if IsProtoGTE(r, ID898) {
		r.String(&c.PermissionLevel)
	} else {
		permLevel := byte(0)
		r.Uint8(&permLevel)
		c.PermissionLevel = "any"
	}
	r.Uint32(&c.AliasesOffset)
	if IsProtoGTE(r, ID898) {
		protocol.FuncSlice(r, &c.ChainedSubcommandOffsets, r.Uint32)
	} else {
		offsets := typeconf.SliceIntToSliceInt[uint32, uint16](c.ChainedSubcommandOffsets)
		protocol.FuncSlice(r, &offsets, r.Uint16)
		c.ChainedSubcommandOffsets = typeconf.SliceIntToSliceInt[uint16, uint32](offsets)
	}
	protocol.Slice(r, &c.Overloads)
}

func (c *Command) ToLatest() protocol.Command {
	return protocol.Command{
		Name:                     c.Name,
		Description:              c.Description,
		Flags:                    c.Flags,
		PermissionLevel:          newCommandPermToLegacy(c.PermissionLevel),
		AliasesOffset:            c.AliasesOffset,
		ChainedSubcommandOffsets: c.ChainedSubcommandOffsets,
		Overloads:                c.Overloads,
	}
}

func (c *Command) FromLatest(latest protocol.Command) Command {
	c.Name = latest.Name
	c.Description = latest.Description
	c.Flags = latest.Flags
	c.PermissionLevel = legacyCommandPermToNew(latest.PermissionLevel)
	c.AliasesOffset = latest.AliasesOffset
	c.ChainedSubcommandOffsets = latest.ChainedSubcommandOffsets
	c.Overloads = latest.Overloads
	return *c
}

// CommandEnum represents an enum in a command usage. The enum typically has a type and a set of options that
// are valid. A value that is not one of the options results in a failure during execution.
type CommandEnum struct {
	// Type is the type of the command enum. The type will show up in the command usage as the type of the
	// argument if it has a certain amount of arguments, or when Options is set to true in the
	// command holding the enum.
	Type string
	// ValueIndices holds a list of indices that point to the EnumValues slice in the
	// AvailableCommandsPacket. These represent the options of the enum.
	ValueIndices []uint32
}

// CommandEnumContext holds context required for encoding command enums.
type CommandEnumContext struct {
	EnumValues []string
}

// Marshal encodes/decodes a CommandEnum.
func (ctx CommandEnumContext) Marshal(r protocol.IO, x *CommandEnum) {
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

// ToLatest ...
func (ce *CommandEnum) ToLatest() protocol.CommandEnum {
	return protocol.CommandEnum{
		Type:         ce.Type,
		ValueIndices: ce.ValueIndices,
	}
}

// FromLatest ...
func (ce *CommandEnum) FromLatest(latest protocol.CommandEnum) CommandEnum {
	ce.Type = latest.Type
	ce.ValueIndices = latest.ValueIndices
	return *ce
}

// ChainedSubcommand represents a subcommand that can have chained commands, such as /execute which allows you to run
// another command as another entity or at a different position etc.
type ChainedSubcommand struct {
	// Name is the name of the chained subcommand and shows up in the list as a regular subcommand enum.
	Name string
	// Values contains the index and parameter type of the chained subcommand.
	Values []ChainedSubcommandValue
}

func (x *ChainedSubcommand) Marshal(r protocol.IO) {
	r.String(&x.Name)
	protocol.Slice(r, &x.Values)
}

// ToLatest ...
func (x *ChainedSubcommand) ToLatest() protocol.ChainedSubcommand {
	values := make([]protocol.ChainedSubcommandValue, len(x.Values))
	for i, v := range x.Values {
		values[i] = v.ToLatest()
	}
	return protocol.ChainedSubcommand{
		Name:   x.Name,
		Values: values,
	}
}

// FromLatest ...
func (x *ChainedSubcommand) FromLatest(latest protocol.ChainedSubcommand) ChainedSubcommand {
	x.Name = latest.Name
	x.Values = make([]ChainedSubcommandValue, len(latest.Values))
	for i, v := range latest.Values {
		x.Values[i] = (&ChainedSubcommandValue{}).FromLatest(v)
	}
	return *x
}

// ChainedSubcommandValue represents the value for a chained subcommand argument.
type ChainedSubcommandValue struct {
	// Index is the index of the argument in the ChainedSubcommandValues slice from the AvailableCommands packet. This is
	// then used to set the type specified by the Value field below.
	Index uint32
	// Value is a combination of the flags above and specified the type of argument. Unlike regular parameter types,
	// this should NOT contain any of the special flags (valid, enum, suffixed or soft enum) but only the basic types.
	Value uint32
}

func (x *ChainedSubcommandValue) Marshal(r protocol.IO) {
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

// ToLatest ...
func (x *ChainedSubcommandValue) ToLatest() protocol.ChainedSubcommandValue {
	return protocol.ChainedSubcommandValue{
		Index: x.Index,
		Value: x.Value,
	}
}

// FromLatest ...
func (x *ChainedSubcommandValue) FromLatest(latest protocol.ChainedSubcommandValue) ChainedSubcommandValue {
	x.Index = latest.Index
	x.Value = latest.Value
	return *x
}

var (
	commandOrigins = []string{
		"player",                   // protocol.CommandOriginPlayer
		"commandblock",             // protocol.CommandOriginBlock
		"minecartcommandblock",     // protocol.CommandOriginMinecartBlock
		"devconsole",               // protocol.CommandOriginDevConsole
		"test",                     // protocol.CommandOriginTest
		"automationplayer",         // protocol.CommandOriginAutomationPlayer
		"clientautomation",         // protocol.CommandOriginClientAutomation
		"dedicatedserver",          // protocol.CommandOriginDedicatedServer
		"entity",                   // protocol.CommandOriginEntity
		"virtual",                  // protocol.CommandOriginVirtual
		"gameargument",             // protocol.CommandOriginGameArgument
		"entityserver",             // protocol.CommandOriginEntityServer
		"precompiled",              // protocol.CommandOriginPrecompiled
		"gamedirectorentityserver", // protocol.CommandOriginGameDirectorEntityServer
		"scripting",                // protocol.CommandOriginScript
		"executecontext",           // protocol.CommandOriginExecutor
	}
	originToLegacyMap = map[string]uint32{
		"player":                   protocol.CommandOriginPlayer,
		"commandblock":             protocol.CommandOriginBlock,
		"minecartcommandblock":     protocol.CommandOriginMinecartBlock,
		"devconsole":               protocol.CommandOriginDevConsole,
		"test":                     protocol.CommandOriginTest,
		"automationplayer":         protocol.CommandOriginAutomationPlayer,
		"clientautomation":         protocol.CommandOriginClientAutomation,
		"dedicatedserver":          protocol.CommandOriginDedicatedServer,
		"entity":                   protocol.CommandOriginEntity,
		"virtual":                  protocol.CommandOriginVirtual,
		"gameargument":             protocol.CommandOriginGameArgument,
		"entityserver":             protocol.CommandOriginEntityServer,
		"precompiled":              protocol.CommandOriginPrecompiled,
		"gamedirectorentityserver": protocol.CommandOriginGameDirectorEntityServer,
		"scripting":                protocol.CommandOriginScript,
		"executecontext":           protocol.CommandOriginExecutor,
	}
)

func commandOriginFromLegacy(legacy uint32) string {
	if int(legacy) < len(commandOrigins) {
		return commandOrigins[legacy]
	}
	return "player" // default to protocol.CommandOriginPlayer
}

func commandOriginToLegacy(origin string) uint32 {
	legacy, ok := originToLegacyMap[origin]
	if ok {
		return legacy
	}
	return protocol.CommandOriginPlayer
}

// CommandOriginData reads/writes a CommandOrigin x using IO r.
func CommandOriginData(r protocol.IO, x *protocol.CommandOrigin) {
	originV898 := commandOriginFromLegacy(x.Origin)
	if IsProtoGTE(r, ID898) {
		r.String(&originV898)
		x.Origin = commandOriginToLegacy(originV898)
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

// CommandOutputMessage represents a message sent by a command that holds the output of one of the commands
// executed.
type CommandOutputMessage struct {
	// Success indicates if the output message was one of a successful command execution. If set to true, the
	// output message is by default coloured white, whereas if set to false, the message is by default
	// coloured red.
	Success bool
	// Message is the message that is sent to the client in the chat window. It may either be simply a
	// message or a translated built-in string like 'commands.tp.success.coordinates', combined with specific
	// parameters below.
	Message string
	// Parameters is a list of parameters that serve to supply the message sent with additional information,
	// such as the position that a player was teleported to or the effect that was applied to an entity.
	// These parameters only apply for the Minecraft built-in command output.
	Parameters []string
}

// Marshal encodes/decodes a CommandOutputMessage.
func (x *CommandOutputMessage) Marshal(r protocol.IO) {
	if IsProtoLT(r, ID898) {
		r.Bool(&x.Success)
	}
	r.String(&x.Message)
	if IsProtoGTE(r, ID898) {
		r.Bool(&x.Success)
	}
	protocol.FuncSlice(r, &x.Parameters, r.String)
}

// ToLatest ...
func (x *CommandOutputMessage) ToLatest() protocol.CommandOutputMessage {
	return protocol.CommandOutputMessage{
		Success:    x.Success,
		Message:    x.Message,
		Parameters: x.Parameters,
	}
}

// FromLatest ...
func (x *CommandOutputMessage) FromLatest(latest protocol.CommandOutputMessage) CommandOutputMessage {
	x.Success = latest.Success
	x.Message = latest.Message
	x.Parameters = latest.Parameters
	return *x
}
