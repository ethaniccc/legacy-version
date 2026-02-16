package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion649 ...
	ItemVersion649 = 161
	// BlockVersion649 ...
	BlockVersion649 int32 = (1 << 24) | (20 << 16) | (60 << 8)
)

var (
	//go:embed data/required_item_list_649.json
	requiredItemList649 []byte
	//go:embed data/block_states_649.nbt
	blockStateData649 []byte
)

// New649 ...
func New649(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList649, ItemVersion649)
	blockTranslator := lookupOrCreateBlockTranslator(649, BlockVersion649, blockStateData649)
	return &Protocol{
		ver:             "1.20.60",
		id:              proto.ID649,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
