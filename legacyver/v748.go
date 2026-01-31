package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion748 ...
	ItemVersion748 = 221
	// BlockVersion748 ...
	BlockVersion748 int32 = (1 << 24) | (21 << 16) | (40 << 8)
)

var (
	//go:embed data/required_item_list_748.json
	requiredItemList748 []byte
	//go:embed data/block_states_748.nbt
	blockStateData748 []byte
)

func New748(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList748, ItemVersion748)
	blockTranslator := lookupOrCreateBlockTranslator(748, BlockVersion748, blockStateData748)
	return &Protocol{
		ver:             "1.21.40",
		id:              proto.ID748,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
