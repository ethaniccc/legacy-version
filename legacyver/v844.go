package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion844 ...
	ItemVersion844 = 251
	// BlockVersion844 ...
	BlockVersion844 int32 = (1 << 24) | (21 << 16) | (111 << 8)
)

var (
	//go:embed data/required_item_list_844.json
	requiredItemList844 []byte
	//go:embed data/block_states_844.nbt
	blockStateData844 []byte
)

func New844(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList844, ItemVersion844)
	blockTranslator := lookupOrCreateBlockTranslator(844, BlockVersion844, blockStateData844)
	return &Protocol{
		ver:             "1.21.111",
		id:              proto.ID844,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
