package legacypacket

import (
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// BookEdit is sent by the client when it edits a book.
type BookEdit struct {
	// InventorySlot is the slot in which the edited book may be found.
	InventorySlot int32
	// ActionType is the type of the book edit action.
	ActionType uint32
	// PageNumber is the page number concerned by the action.
	PageNumber int32
	// SecondaryPageNumber is the second page number for swap operations.
	SecondaryPageNumber int32
	// Text is the text written on the page.
	Text string
	// PhotoName is the attached education photo name.
	PhotoName string
	// Title is the title used when signing the book.
	Title string
	// Author is the author used when signing the book.
	Author string
	// XUID is the XBOX Live User ID of the editor.
	XUID string
}

// ID ...
func (*BookEdit) ID() uint32 {
	return packet.IDBookEdit
}

func (pk *BookEdit) Marshal(io protocol.IO) {
	if proto.IsProtoGTE(io, proto.ID924) {
		io.Varint32(&pk.InventorySlot)
		io.Varuint32(&pk.ActionType)
		switch pk.ActionType {
		case packet.BookActionReplacePage, packet.BookActionAddPage:
			io.Varint32(&pk.PageNumber)
			io.String(&pk.Text)
			io.String(&pk.PhotoName)
		case packet.BookActionDeletePage:
			io.Varint32(&pk.PageNumber)
		case packet.BookActionSwapPages:
			io.Varint32(&pk.PageNumber)
			io.Varint32(&pk.SecondaryPageNumber)
		case packet.BookActionSign:
			io.String(&pk.Title)
			io.String(&pk.Author)
			io.String(&pk.XUID)
		default:
			io.UnknownEnumOption(pk.ActionType, "book edit action type")
		}
		return
	}

	actionType := byte(pk.ActionType)
	io.Uint8(&actionType)
	pk.ActionType = uint32(actionType)

	inventorySlot := byte(pk.InventorySlot)
	io.Uint8(&inventorySlot)
	pk.InventorySlot = int32(inventorySlot)

	switch pk.ActionType {
	case packet.BookActionReplacePage, packet.BookActionAddPage:
		pageNumber := byte(pk.PageNumber)
		io.Uint8(&pageNumber)
		pk.PageNumber = int32(pageNumber)
		io.String(&pk.Text)
		io.String(&pk.PhotoName)
	case packet.BookActionDeletePage:
		pageNumber := byte(pk.PageNumber)
		io.Uint8(&pageNumber)
		pk.PageNumber = int32(pageNumber)
	case packet.BookActionSwapPages:
		pageNumber := byte(pk.PageNumber)
		io.Uint8(&pageNumber)
		pk.PageNumber = int32(pageNumber)
		secondaryPageNumber := byte(pk.SecondaryPageNumber)
		io.Uint8(&secondaryPageNumber)
		pk.SecondaryPageNumber = int32(secondaryPageNumber)
	case packet.BookActionSign:
		io.String(&pk.Title)
		io.String(&pk.Author)
		io.String(&pk.XUID)
	default:
		io.UnknownEnumOption(pk.ActionType, "book edit action type")
	}
}
