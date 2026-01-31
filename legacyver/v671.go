package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion671 ...
	ItemVersion671 = 181
	// BlockVersion671 ...
	BlockVersion671 int32 = (1 << 24) | (20 << 16) | (80 << 8)
)

var (
	//go:embed data/required_item_list_671.json
	requiredItemList671 []byte
	//go:embed data/block_states_671.nbt
	blockStateData671 []byte
)

// New671 ...
func New671(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList671, ItemVersion671)
	blockTranslator := lookupOrCreateBlockTranslator(671, BlockVersion671, blockStateData671)
	return &Protocol{
		ver:             "1.20.80",
		id:              proto.ID671,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
