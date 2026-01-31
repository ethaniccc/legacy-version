package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion827 ...
	ItemVersion827 = 241
	// BlockVersion827 ...
	BlockVersion827 int32 = (1 << 24) | (21 << 16) | (100 << 8)
)

var (
	//go:embed data/required_item_list_827.json
	requiredItemList827 []byte
	//go:embed data/block_states_827.nbt
	blockStateData827 []byte
)

func New827(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList827, ItemVersion827)
	blockTranslator := lookupOrCreateBlockTranslator(827, BlockVersion827, blockStateData827)
	return &Protocol{
		ver:             "1.21.100",
		id:              proto.ID827,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
