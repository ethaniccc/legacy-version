package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion859 ...
	ItemVersion859 = 251
	// BlockVersion859 ...
	BlockVersion859 int32 = (1 << 24) | (21 << 16) | (120 << 8)
)

var (
	//go:embed data/required_item_list_859.json
	requiredItemList859 []byte
	//go:embed data/block_states_859.nbt
	blockStateData859 []byte
)

func New859(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList859, ItemVersion859)
	blockTranslator := lookupOrCreateBlockTranslator(859, BlockVersion859, blockStateData859)
	return &Protocol{
		ver:             "1.21.120",
		id:              proto.ID859,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
