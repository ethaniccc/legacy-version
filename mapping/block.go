package mapping

import (
	"bytes"
	"sort"

	"github.com/ethaniccc/legacy-version/internal"

	"github.com/df-mc/worldupgrader/blockupgrader"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/segmentio/fasthash/fnv1"
)

type Block interface {
	// StateToRuntimeID converts a block state to a runtime ID.
	StateToRuntimeID(blockupgrader.BlockState) (uint32, bool)
	// RuntimeIDToState converts a runtime ID to a name and its state properties.
	RuntimeIDToState(uint32) (blockupgrader.BlockState, bool)
	// DowngradeBlockActorData downgrades the input sub chunk to a legacy block actor.
	DowngradeBlockActorData(map[string]any)
	// UpgradeBlockActorData upgrades the input sub chunk to the latest block actor.
	UpgradeBlockActorData(map[string]any)
	HashToRuntimeID(hash uint32) (rid uint32, ok bool)
	RuntimeIDToHash(rid uint32) (hash uint32, ok bool)
	// Adjust adjusts the latest mappings to account for custom states.
	Adjust([]protocol.BlockEntry)
	Air() uint32
	InfoUpdate() uint32
}

type DefaultBlockMapping struct {
	// states holds a list of all possible vanilla block states.
	states []blockupgrader.BlockState
	// stateRuntimeIDs holds a map for looking up the runtime ID of a block by the stateHash it produces.
	stateRuntimeIDs map[internal.StateHash]uint32
	// runtimeIDToState holds a map for looking up the blockState of a block by its runtime ID.
	runtimeIDToState     map[uint32]blockupgrader.BlockState
	upgrader, downgrader func(map[string]any) map[string]any

	// airRID is the runtime ID of the air block in the latest version of the game.
	airRID uint32

	// infoUpdateRID is the runtime ID of the info_update block in the latest version of the game.
	infoUpdateBlockRID uint32

	networkhashToRids map[uint32]uint32
	ridsToNetworkhash map[uint32]uint32
}

func NewBlockMapping(raw []byte) *DefaultBlockMapping {
	dec := nbt.NewDecoder(bytes.NewBuffer(raw))

	var states []blockupgrader.BlockState
	stateRuntimeIDs := make(map[internal.StateHash]uint32)
	runtimeIDToState := make(map[uint32]blockupgrader.BlockState)
	var airRID *uint32
	var infoUpdateBlockRID *uint32
	networkhashToRids := make(map[uint32]uint32)
	ridsToNetworkhash := make(map[uint32]uint32)

	for {
		var s blockupgrader.BlockState
		if err := dec.Decode(&s); err != nil {
			break
		}

		rid := uint32(len(states))
		states = append(states, s)
		switch s.Name {
		case "minecraft:air":
			airRID = &rid
		case "minecraft:info_update":
			infoUpdateBlockRID = &rid
		}

		upgradedStates := blockupgrader.Upgrade(s)
		stateRuntimeIDs[internal.HashState(upgradedStates)] = rid
		runtimeIDToState[rid] = s
		networkhashToRids[networkBlockHash(upgradedStates.Name, upgradedStates.Properties)] = rid
		ridsToNetworkhash[rid] = networkBlockHash(s.Name, s.Properties)
	}
	if airRID == nil {
		panic("couldn't find air")
	}
	if infoUpdateBlockRID == nil {
		panic("couldn't find info_update block")
	}

	return &DefaultBlockMapping{
		states:             states,
		stateRuntimeIDs:    stateRuntimeIDs,
		runtimeIDToState:   runtimeIDToState,
		airRID:             *airRID,
		infoUpdateBlockRID: *infoUpdateBlockRID,
		networkhashToRids:  networkhashToRids,
		ridsToNetworkhash:  ridsToNetworkhash,
	}
}

func (m *DefaultBlockMapping) WithBlockActorRemapper(downgrader, upgrader func(map[string]any) map[string]any) *DefaultBlockMapping {
	m.downgrader = downgrader
	m.upgrader = upgrader
	return m
}

func (m *DefaultBlockMapping) StateToRuntimeID(state blockupgrader.BlockState) (uint32, bool) {
	rid, ok := m.stateRuntimeIDs[internal.HashState(blockupgrader.Upgrade(state))]
	return rid, ok
}

func (m *DefaultBlockMapping) RuntimeIDToState(runtimeId uint32) (blockupgrader.BlockState, bool) {
	state, found := m.runtimeIDToState[runtimeId]
	return state, found
}

func (m *DefaultBlockMapping) DowngradeBlockActorData(actorData map[string]any) {
	if m.downgrader != nil {
		m.downgrader(actorData)
	}
}

func (m *DefaultBlockMapping) UpgradeBlockActorData(actorData map[string]any) {
	if m.upgrader != nil {
		m.upgrader(actorData)
	}
}

func (m *DefaultBlockMapping) RuntimeIDToHash(rid uint32) (uint32, bool) {
	hash, found := m.ridsToNetworkhash[rid]
	return hash, found
}

func (m *DefaultBlockMapping) HashToRuntimeID(hash uint32) (uint32, bool) {
	rid, found := m.networkhashToRids[hash]
	return rid, found
}

func (m *DefaultBlockMapping) Adjust(entries []protocol.BlockEntry) {
	if len(entries) == 0 {
		return
	}

	customStates := convert(entries)
	var newStates []blockupgrader.BlockState
	for _, state := range customStates {
		if _, ok := m.StateToRuntimeID(state); !ok {
			newStates = append(newStates, state)
		}
	}
	if len(newStates) == 0 {
		return
	}

	adjustedStates := append(m.states, customStates...)
	sort.SliceStable(adjustedStates, func(i, j int) bool {
		stateOne, stateTwo := adjustedStates[i], adjustedStates[j]
		return stateOne.Name != stateTwo.Name && fnv1.HashString64(stateOne.Name) < fnv1.HashString64(stateTwo.Name)
	})

	m.stateRuntimeIDs = make(map[internal.StateHash]uint32, len(adjustedStates))
	m.runtimeIDToState = make(map[uint32]blockupgrader.BlockState, len(adjustedStates))
	for rid, state := range adjustedStates {
		m.stateRuntimeIDs[internal.HashState(blockupgrader.Upgrade(state))] = uint32(rid)
		m.runtimeIDToState[uint32(rid)] = state
	}
}

func (m *DefaultBlockMapping) Air() uint32 {
	return m.airRID
}

func (m *DefaultBlockMapping) InfoUpdate() uint32 {
	return m.infoUpdateBlockRID
}
