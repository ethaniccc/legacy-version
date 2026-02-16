package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion898 ...
	ItemVersion898 = 251
	// BlockVersion898 ...
	BlockVersion898 int32 = (1 << 24) | (21 << 16) | (130 << 8)
)

var (
	//go:embed data/required_item_list_898.json
	requiredItemList898 []byte
	//go:embed data/block_states_898.nbt
	blockStateData898 []byte
)

func New898(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList898, ItemVersion898)
	blockTranslator := lookupOrCreateBlockTranslator(898, BlockVersion898, blockStateData898)
	return &Protocol{
		ver:             "1.21.130",
		id:              proto.ID898,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
