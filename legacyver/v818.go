package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion818 ...
	ItemVersion818 = 231
	// BlockVersion818 ...
	BlockVersion818 int32 = (1 << 24) | (21 << 16) | (90 << 8)
)

var (
	//go:embed data/required_item_list_818.json
	requiredItemList818 []byte
	//go:embed data/block_states_818.nbt
	blockStateData818 []byte
)

func New818(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList818, ItemVersion818)
	blockTranslator := lookupOrCreateBlockTranslator(818, BlockVersion818, blockStateData818)
	return &Protocol{
		ver:             "1.21.90",
		id:              proto.ID818,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
