package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion729 ...
	ItemVersion729 = 211
	// BlockVersion729 ...
	BlockVersion729 int32 = (1 << 24) | (21 << 16) | (30 << 8)
)

var (
	//go:embed data/required_item_list_729.json
	requiredItemList729 []byte
	//go:embed data/block_states_729.nbt
	blockStateData729 []byte
)

func New729(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList729, ItemVersion729)
	blockTranslator := lookupOrCreateBlockTranslator(729, BlockVersion729, blockStateData729)
	return &Protocol{
		ver:             "1.21.30",
		id:              proto.ID729,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
