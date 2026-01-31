package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion686 ...
	ItemVersion686 = 191
	// BlockVersion686 ...
	BlockVersion686 int32 = (1 << 24) | (21 << 16) | (2 << 8)
)

var (
	//go:embed data/required_item_list_686.json
	requiredItemList686 []byte
	//go:embed data/block_states_686.nbt
	blockStateData686 []byte
)

func New686(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList686, ItemVersion686)
	blockTranslator := lookupOrCreateBlockTranslator(686, BlockVersion686, blockStateData686)
	return &Protocol{
		ver:             "1.21.2",
		id:              proto.ID686,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
