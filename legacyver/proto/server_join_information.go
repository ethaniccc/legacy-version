package proto

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func MarshalGatheringJoinInfo(r protocol.IO, x *protocol.GatheringJoinInfo) {
	if IsProtoLT(r, ID944) {
		uuidStr := x.ExperienceID.String()
		r.String(&uuidStr)
		if uuidStr != "" {
			x.ExperienceID = uuid.MustParse(uuidStr)
		}
	} else {
		r.UUID(&x.ExperienceID)
	}
	r.String(&x.ExperienceName)
	if IsProtoGTE(r, ID2168) {
		protocol.OptionalFunc(r, &x.ExperienceWorldID, r.UUID)
		protocol.OptionalFunc(r, &x.ExperienceWorldName, r.String)
	} else if IsProtoLT(r, ID944) {
		experienceWorldID, _ := x.ExperienceWorldID.Value()
		uuidStr := experienceWorldID.String()
		r.String(&uuidStr)
		if uuidStr != "" {
			x.ExperienceWorldID = protocol.Option(uuid.MustParse(uuidStr))
		}
		experienceWorldName, _ := x.ExperienceWorldName.Value()
		r.String(&experienceWorldName)
		x.ExperienceWorldName = protocol.Option(experienceWorldName)
	} else {
		experienceWorldID, _ := x.ExperienceWorldID.Value()
		r.UUID(&experienceWorldID)
		x.ExperienceWorldID = protocol.Option(experienceWorldID)
		experienceWorldName, _ := x.ExperienceWorldName.Value()
		r.String(&experienceWorldName)
		x.ExperienceWorldName = protocol.Option(experienceWorldName)
	}
	r.String(&x.CreatorID)
	if IsProtoGTE(r, ID2168) {
		protocol.OptionalFunc(r, &x.TargetID, r.UUID)
		protocol.OptionalFunc(r, &x.ScenarioID, r.String)
		protocol.OptionalFunc(r, &x.ServerID, r.String)
	} else if IsProtoGTE(r, ID944) {
		targetID, _ := x.TargetID.Value()
		r.UUID(&targetID)
		x.TargetID = protocol.Option(targetID)
		// ScenarioID became a string in 1.26.20 (proto 975). In this
		// environment the protocol-944 client (game 1.26.13) already ships
		// the 1.26.20 GatheringJoinInfo, so the string form applies from 944.
		if IsProtoLT(r, ID944) {
			scenarioID, _ := x.ScenarioID.Value()
			sid, _ := uuid.Parse(scenarioID)
			r.UUID(&sid)
			x.ScenarioID = protocol.Option(sid.String())
		} else {
			scenarioID, _ := x.ScenarioID.Value()
			r.String(&scenarioID)
			x.ScenarioID = protocol.Option(scenarioID)
		}
		serverID, _ := x.ServerID.Value()
		r.String(&serverID)
		x.ServerID = protocol.Option(serverID)
	} else {
		storeId := ""
		r.String(&storeId)
	}
}

func MarshalServerJoinInformation(r protocol.IO, x *protocol.ServerJoinInformation) {
	protocol.OptionalFunc(r, &x.GatheringJoinInfo, func(x *protocol.GatheringJoinInfo) { MarshalGatheringJoinInfo(r, x) })
	if IsProtoGTE(r, ID944) {
		protocol.OptionalFunc(r, &x.StoreEntryPointInfo, func(x *protocol.StoreEntryPointInfo) {
			r.String(&x.StoreID)
			r.String(&x.StoreName)
		})
		protocol.OptionalFunc(r, &x.PresenceInfo, func(x *protocol.PresenceInfo) { MarshalPresenceInfo(r, x) })
	}
}

func MarshalPresenceInfo(r protocol.IO, x *protocol.PresenceInfo) {
	if IsProtoGTE(r, ID2168) {
		protocol.OptionalFunc(r, &x.RichPresenceID, r.String)
		return
	}
	if IsProtoGTE(r, ID1001) {
		var experienceName, worldName protocol.Optional[string]
		protocol.OptionalFunc(r, &experienceName, r.String)
		protocol.OptionalFunc(r, &worldName, r.String)
		richPresenceID, _ := x.RichPresenceID.Value()
		r.String(&richPresenceID)
		x.RichPresenceID = protocol.Option(richPresenceID)
		return
	}
	var experienceName, worldName string
	r.String(&experienceName)
	r.String(&worldName)
}
