package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion786 ...
	ItemVersion786 = 231
	// BlockVersion786 ...
	BlockVersion786 int32 = (1 << 24) | (21 << 16) | (70 << 8)
)

var (
	//go:embed data/required_item_list_786.json
	requiredItemList786 []byte
	//go:embed data/block_states_786.nbt
	blockStateData786 []byte
)

func New786(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList786, ItemVersion786)
	blockTranslator := lookupOrCreateBlockTranslator(786, BlockVersion786, blockStateData786)
	return &Protocol{
		ver:             "1.21.70",
		id:              proto.ID786,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
