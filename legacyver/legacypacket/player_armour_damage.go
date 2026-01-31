package legacypacket

import (
	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// PlayerArmourDamage is sent by the server to damage the armour of a player. It is a very efficient packet,
// but generally it's much easier to just send a slot update for the damaged armour.
type PlayerArmourDamage struct {
	// List ...
	List []protocol.PlayerArmourDamageEntry
}

// ID ...
func (pk *PlayerArmourDamage) ID() uint32 {
	return packet.IDPlayerArmourDamage
}

func (pk *PlayerArmourDamage) Marshal(io protocol.IO) {
	if proto.IsProtoLT(io, proto.ID844) {
		var bitset uint8
		flags := []uint8{packet.PlayerArmourDamageFlagHelmet, packet.PlayerArmourDamageFlagChestplate, packet.PlayerArmourDamageFlagLeggings, packet.PlayerArmourDamageFlagBoots}
		if proto.IsProtoGTE(io, proto.ID712) {
			flags = append(flags, packet.PlayerArmourDamageFlagBody)
		}
		if proto.IsReader(io) {
			io.Uint8(&bitset)
			pk.List = make([]protocol.PlayerArmourDamageEntry, 0, len(flags))
			for _, flag := range flags {
				if bitset&flag != 0 {
					v := protocol.PlayerArmourDamageEntry{
						ArmourSlot: int32(flag),
					}
					v2 := int32(0)
					io.Varint32(&v2)
					v.Damage = int16(v2)
					pk.List = append(pk.List, v)
				}
			}
			return
		}

		for _, entry := range pk.List {
			if entry.ArmourSlot == packet.PlayerArmourDamageFlagBody && !proto.IsProtoGTE(io, proto.ID712) {
				continue
			}
			bitset |= 1 << entry.ArmourSlot
		}
		io.Uint8(&bitset)
		for _, flag := range flags {
			if bitset&flag != 0 {
				damage := int32(0)
				for _, entry := range pk.List {
					if entry.ArmourSlot == int32(flag) {
						damage = int32(entry.Damage)
						break
					}
				}
				io.Varint32(&damage)
			}
		}
		return
	}
	protocol.Slice(io, &pk.List)
}
