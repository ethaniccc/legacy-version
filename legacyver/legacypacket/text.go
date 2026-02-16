package legacypacket

import (
	"fmt"
	"strings"
  
	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	TextTypeRaw = iota
	TextTypeChat
	TextTypeTranslation
	TextTypePopup
	TextTypeJukeboxPopup
	TextTypeTip
	TextTypeSystem
	TextTypeWhisper
	TextTypeAnnouncement
	TextTypeObjectWhisper
	TextTypeObject
	TextTypeObjectAnnouncement
)

// Text is sent by the client to the server to send chat messages, and by the server to the client to forward
// or send messages, which may be chat, popups, tips etc.
type Text struct {
	// TextType is the type of the text sent. When a client sends this to the server, it should always be
	// TextTypeChat. If the server sends it, it may be one of the other text types above.
	TextType byte
	// NeedsTranslation specifies if any of the messages need to be translated. It seems that where % is found
	// in translatable text types, these are translated regardless of this bool. Translatable text types
	// include TextTypeTranslation, TextTypeTip, TextTypePopup and TextTypeJukeboxPopup.
	NeedsTranslation bool
	// SourceName is the name of the source of the messages. This source is displayed in text types such as
	// the TextTypeChat and TextTypeWhisper, where typically the username is shown.
	SourceName string
	// Message is the message of the packet. This field is set for each TextType and is the main component of
	// the packet.
	Message string
	// Parameters is a list of parameters that should be filled into the message. These parameters are only
	// written if the type of the packet is TextTypeTranslation, TextTypeTip, TextTypePopup or TextTypeJukeboxPopup.
	Parameters []string
	// XUID is the XBOX Live user ID of the player that sent the message. It is only set for packets of
	// TextTypeChat. When sent to a player, the player will only be shown the chat message if a player with
	// this XUID is present in the player list and not muted, or if the XUID is empty.
	XUID string
	// PlatformChatID is an identifier only set for particular platforms when chatting (presumably only for
	// Nintendo Switch). It is otherwise an empty string, and is used to decide which players are able to
	// chat with each other.
	PlatformChatID string
	// FilteredMessage is a filtered version of Message with all the profanity removed. The client will use
	// this over Message if this field is not empty and they have the "Filter Profanity" setting enabled.
	FilteredMessage protocol.Optional[string]
}

// ID ...
func (*Text) ID() uint32 {
	return packet.IDText
}

func (pk *Text) Marshal(io protocol.IO) {
	if proto.IsProtoLT(io, proto.ID898) {
		io.Uint8(&pk.TextType)
	}
	io.Bool(&pk.NeedsTranslation)
	if proto.IsProtoGTE(io, proto.ID898) {
		var categoryType uint8
		if pk.TextType == TextTypeRaw || pk.TextType == TextTypeTip || pk.TextType == TextTypeSystem || pk.TextType == TextTypeObjectWhisper || pk.TextType == TextTypeObjectAnnouncement || pk.TextType == TextTypeObject {
			categoryType = packet.TextCategoryMessageOnly
		} else if pk.TextType == TextTypeChat || pk.TextType == TextTypeWhisper || pk.TextType == TextTypeAnnouncement {
			categoryType = packet.TextCategoryAuthoredMessage
		} else {
			categoryType = packet.TextCategoryMessageWithParameters
		}
		io.Uint8(&categoryType)
		// Protocols 898-923 include a quirky Mojang text-category string block.
		// 924+ removed this and only keeps the uint8 category and text type.
		if proto.IsProtoLT(io, proto.ID924) {
			for _, v := range textCategoryConstants(categoryType) {
				stringConst(io, v)
			}
		}
		io.Uint8(&pk.TextType)
	}
	switch pk.TextType {
	case TextTypeChat, TextTypeWhisper, TextTypeAnnouncement:
		io.String(&pk.SourceName)
		io.String(&pk.Message)
	case TextTypeRaw, TextTypeTip, TextTypeSystem, TextTypeObject, TextTypeObjectWhisper, TextTypeObjectAnnouncement:
		io.String(&pk.Message)
	case TextTypeTranslation, TextTypePopup, TextTypeJukeboxPopup:
		io.String(&pk.Message)
		protocol.FuncSlice(io, &pk.Parameters, io.String)
	}
	if proto.IsProtoGTE(io, proto.ID898) {
		if len(pk.Message) == 0 {
			io.InvalidValue(pk.Message, "message", "string cannot be empty")
		}
	}
	io.String(&pk.XUID)
	io.String(&pk.PlatformChatID)
	if proto.IsProtoGTE(io, proto.ID685) {
		if proto.IsProtoGTE(io, proto.ID898) {
			protocol.OptionalFunc(io, &pk.FilteredMessage, io.String)
		} else {
			v, _ := pk.FilteredMessage.Value()
			io.String(&v)
			if v != "" {
				pk.FilteredMessage = protocol.Option(v)
			} else {
				pk.FilteredMessage = protocol.Optional[string]{}
			}
		}
	}
}

func textCategoryConstants(categoryType uint8) []string {
	switch categoryType {
	case packet.TextCategoryMessageOnly:
		return []string{"raw", "tip", "systemMessage", "textObjectWhisper", "textObjectAnnouncement", "textObject"}
	case packet.TextCategoryAuthoredMessage:
		return []string{"chat", "whisper", "announcement"}
	default:
		return []string{"translate", "popup", "jukeboxPopup"}
	}
}

func stringConst(io protocol.IO, expected string) {
	if proto.IsReader(io) {
		var got string
		io.String(&got)
		if !strings.EqualFold(got, expected) {
			io.InvalidValue(got, "text category constant", fmt.Sprintf("expected %q", expected))
		}
		return
	}
	v := expected
	io.String(&v)
}
