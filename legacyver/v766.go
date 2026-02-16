package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion766 ...
	ItemVersion766 = 231
	// BlockVersion766 ...
	BlockVersion766 int32 = (1 << 24) | (21 << 16) | (50 << 8)
)

var (
	//go:embed data/required_item_list_766.json
	requiredItemList766 []byte
	//go:embed data/block_states_766.nbt
	blockStateData766 []byte
)

func New766(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList766, ItemVersion766)
	blockTranslator := lookupOrCreateBlockTranslator(766, BlockVersion766, blockStateData766)
	return &Protocol{
		ver:             "1.21.50",
		id:              proto.ID766,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
