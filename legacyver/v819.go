package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion819 ...
	ItemVersion819 = 231
	// BlockVersion819 ...
	BlockVersion819 int32 = (1 << 24) | (21 << 16) | (93 << 8)
)

var (
	//go:embed data/required_item_list_819.json
	requiredItemList819 []byte
	//go:embed data/block_states_819.nbt
	blockStateData819 []byte
)

func New819(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList819, ItemVersion819)
	blockTranslator := lookupOrCreateBlockTranslator(819, BlockVersion819, blockStateData819)
	return &Protocol{
		ver:             "1.21.90",
		id:              proto.ID819,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
