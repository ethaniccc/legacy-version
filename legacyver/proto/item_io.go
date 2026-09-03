package proto

import (
	"bytes"
	"fmt"

	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	legacyItemDescriptorInvalid = iota
	legacyItemDescriptorDefault
	legacyItemDescriptorMoLang
	legacyItemDescriptorItemTag
	legacyItemDescriptorDeferred
	legacyItemDescriptorComplexAlias
)

func (r *Reader) ItemDescriptorCount(x *protocol.ItemDescriptorCount) {
	if r.protocolID >= ID2168 {
		r.Reader.ItemDescriptorCount(x)
		return
	}
	readLegacyItemDescriptorCount(r, x)
}

func (w *Writer) ItemDescriptorCount(x *protocol.ItemDescriptorCount) {
	if w.protocolID >= ID2168 {
		w.Writer.ItemDescriptorCount(x)
		return
	}
	writeLegacyItemDescriptorCount(w, x)
}

func readLegacyItemDescriptorCount(r *Reader, x *protocol.ItemDescriptorCount) {
	var kind uint8
	r.Uint8(&kind)
	switch kind {
	case legacyItemDescriptorInvalid:
		x.Descriptor = &protocol.InvalidItemDescriptor{}
	case legacyItemDescriptorDefault:
		var networkID, metadata int16
		r.Int16(&networkID)
		if networkID != 0 {
			r.Int16(&metadata)
		}
		// Numeric recipe descriptors no longer have a numeric field in the
		// current API. Retain their metadata; server-bound recipe data is not
		// used by the game.
		x.Descriptor = &protocol.DefaultItemDescriptor{MetadataValue: int32(metadata)}
	case legacyItemDescriptorMoLang:
		var expression string
		var version uint8
		r.String(&expression)
		r.Uint8(&version)
		x.Descriptor = &protocol.MoLangItemDescriptor{Expression: expression, Version: int16(version)}
	case legacyItemDescriptorItemTag:
		var tag string
		r.String(&tag)
		x.Descriptor = &protocol.ItemTagItemDescriptor{Tag: tag}
	case legacyItemDescriptorDeferred:
		var name string
		var metadata int16
		r.String(&name)
		r.Int16(&metadata)
		x.Descriptor = &protocol.DefaultItemDescriptor{Name: name, MetadataValue: int32(metadata)}
	case legacyItemDescriptorComplexAlias:
		var name string
		r.String(&name)
		x.Descriptor = &protocol.DefaultItemDescriptor{Name: name}
	default:
		r.UnknownEnumOption(kind, "item descriptor type")
		return
	}
	r.Varint32(&x.Count)
}

func writeLegacyItemDescriptorCount(w *Writer, x *protocol.ItemDescriptorCount) {
	var kind uint8
	switch x.Descriptor.(type) {
	case nil, *protocol.InvalidItemDescriptor:
		kind = legacyItemDescriptorInvalid
	case *protocol.DefaultItemDescriptor:
		// The deferred legacy variant is name based, so it can represent the
		// current descriptor without an item-runtime-ID mapping.
		kind = legacyItemDescriptorDeferred
	case *protocol.MoLangItemDescriptor:
		kind = legacyItemDescriptorMoLang
	case *protocol.ItemTagItemDescriptor:
		kind = legacyItemDescriptorItemTag
	default:
		w.UnknownEnumOption(fmt.Sprintf("%T", x.Descriptor), "item descriptor type")
		return
	}
	w.Uint8(&kind)
	switch descriptor := x.Descriptor.(type) {
	case nil, *protocol.InvalidItemDescriptor:
	case *protocol.DefaultItemDescriptor:
		w.String(&descriptor.Name)
		metadata := int16(descriptor.MetadataValue)
		w.Int16(&metadata)
	case *protocol.MoLangItemDescriptor:
		w.String(&descriptor.Expression)
		version := uint8(descriptor.Version)
		w.Uint8(&version)
	case *protocol.ItemTagItemDescriptor:
		w.String(&descriptor.Tag)
	}
	w.Varint32(&x.Count)
}

// ItemInstance preserves the item-instance format used before 1.26.30. The
// 1.26.30 format is handled by ItemInstanceNew and 1.26.40 uses gophertunnel's
// current ItemInstance implementation.
func (r *Reader) ItemInstance(x *protocol.ItemInstance) {
	if r.protocolID >= ID2168 {
		r.Reader.ItemInstance(x)
		return
	}
	readLegacyItemInstance(r, x, false)
}

func (w *Writer) ItemInstance(x *protocol.ItemInstance) {
	if w.protocolID >= ID2168 {
		w.Writer.ItemInstance(x)
		return
	}
	writeLegacyItemInstance(w, x, false)
}

// ItemInstanceNew reads/writes the transitional item format used from
// 1.26.30 up to (but excluding) 1.26.40.
func ItemInstanceNew(io protocol.IO, x *protocol.ItemInstance) {
	if IsProtoGTE(io, ID2168) {
		io.ItemInstance(x)
		return
	}
	if r, ok := io.(*Reader); ok {
		readLegacyItemInstance(r, x, true)
		return
	}
	writeLegacyItemInstance(io.(*Writer), x, true)
}

func readLegacyItemInstance(r *Reader, i *protocol.ItemInstance, transitional bool) {
	x := &i.Stack
	if transitional {
		var id int16
		r.Int16(&id)
		x.NetworkID = int32(id)
	} else {
		r.Varint32(&x.NetworkID)
		if x.NetworkID == 0 {
			x.MetadataValue, x.Count, x.BlockRuntimeID, i.StackNetworkID = 0, 0, 0, 0
			x.NBTData, x.CanBePlacedOn, x.CanBreak, x.BlockingTick = nil, nil, nil, 0
			return
		}
	}
	r.Uint16(&x.Count)
	r.Varuint32(&x.MetadataValue)
	var hasNetworkID bool
	r.Bool(&hasNetworkID)
	if hasNetworkID {
		if transitional {
			var reserved uint32
			r.Varuint32(&reserved)
		}
		r.Varint32(&i.StackNetworkID)
	} else {
		i.StackNetworkID = 0
	}
	if transitional {
		protocol.IntegerFunc(&x.BlockRuntimeID, r.Varuint32)
	} else {
		r.Varint32(&x.BlockRuntimeID)
	}
	readLegacyItemUserData(r, x)
}

func writeLegacyItemInstance(w *Writer, i *protocol.ItemInstance, transitional bool) {
	x := &i.Stack
	if transitional {
		id := int16(x.NetworkID)
		w.Int16(&id)
	} else {
		w.Varint32(&x.NetworkID)
		if x.NetworkID == 0 {
			return
		}
	}
	w.Uint16(&x.Count)
	w.Varuint32(&x.MetadataValue)
	hasNetworkID := i.StackNetworkID != 0
	w.Bool(&hasNetworkID)
	if hasNetworkID {
		if transitional {
			var reserved uint32
			w.Varuint32(&reserved)
		}
		w.Varint32(&i.StackNetworkID)
	}
	if transitional {
		protocol.IntegerFunc(&x.BlockRuntimeID, w.Varuint32)
	} else {
		w.Varint32(&x.BlockRuntimeID)
	}
	writeLegacyItemUserData(w, x, transitional && x.NetworkID == 0)
}

func (r *Reader) Item(x *protocol.ItemStack) {
	if r.protocolID >= ID2168 {
		r.Reader.Item(x)
		return
	}
	r.Varint32(&x.NetworkID)
	if x.NetworkID == 0 {
		x.MetadataValue, x.Count, x.BlockRuntimeID = 0, 0, 0
		x.NBTData, x.CanBePlacedOn, x.CanBreak, x.BlockingTick = nil, nil, nil, 0
		return
	}
	r.Uint16(&x.Count)
	r.Varuint32(&x.MetadataValue)
	r.Varint32(&x.BlockRuntimeID)
	readLegacyItemUserData(r, x)
}

func (w *Writer) Item(x *protocol.ItemStack) {
	if w.protocolID >= ID2168 {
		w.Writer.Item(x)
		return
	}
	w.Varint32(&x.NetworkID)
	if x.NetworkID == 0 {
		return
	}
	w.Uint16(&x.Count)
	w.Varuint32(&x.MetadataValue)
	w.Varint32(&x.BlockRuntimeID)
	writeLegacyItemUserData(w, x, false)
}

func readLegacyItemUserData(r *Reader, x *protocol.ItemStack) {
	var extra []byte
	r.ByteSlice(&extra)
	if len(extra) == 0 {
		x.NBTData, x.CanBePlacedOn, x.CanBreak, x.BlockingTick = nil, nil, nil, 0
		return
	}
	// protocol.NewReader in the oomph gophertunnel fork asserts a
	// *bytes.Buffer, a *bytes.Reader panics.
	buf := protocol.NewReader(bytes.NewBuffer(extra), r.shieldID, r.limits)
	var length int16
	buf.Int16(&length)
	if length == -1 {
		var version uint8
		buf.Uint8(&version)
		if version == 1 {
			buf.NBT(&x.NBTData, nbt.LittleEndian)
		} else {
			buf.UnknownEnumOption(version, "item user data version")
			return
		}
	} else if length > 0 {
		buf.NBT(&x.NBTData, nbt.LittleEndian)
	} else {
		x.NBTData = nil
	}
	protocol.FuncSliceUint32Length(buf, &x.CanBePlacedOn, buf.StringUTF)
	protocol.FuncSliceUint32Length(buf, &x.CanBreak, buf.StringUTF)
	if x.NetworkID == r.shieldID {
		buf.Int64(&x.BlockingTick)
	} else {
		x.BlockingTick = 0
	}
}

func writeLegacyItemUserData(w *Writer, x *protocol.ItemStack, empty bool) {
	if empty {
		var extra []byte
		w.ByteSlice(&extra)
		return
	}
	var b bytes.Buffer
	buf := protocol.NewWriter(&b, w.shieldID)
	if len(x.NBTData) != 0 {
		length := int16(-1)
		version := uint8(1)
		buf.Int16(&length)
		buf.Uint8(&version)
		buf.NBT(&x.NBTData, nbt.LittleEndian)
	} else {
		var length int16
		buf.Int16(&length)
	}
	protocol.FuncSliceUint32Length(buf, &x.CanBePlacedOn, buf.StringUTF)
	protocol.FuncSliceUint32Length(buf, &x.CanBreak, buf.StringUTF)
	if x.NetworkID == w.shieldID {
		buf.Int64(&x.BlockingTick)
	}
	extra := b.Bytes()
	w.ByteSlice(&extra)
}
