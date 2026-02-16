package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion712 ...
	ItemVersion712 = 201
	// BlockVersion712 ...
	BlockVersion712 int32 = (1 << 24) | (21 << 16) | (20 << 8)
)

var (
	//go:embed data/required_item_list_712.json
	requiredItemList712 []byte
	//go:embed data/block_states_712.nbt
	blockStateData712 []byte
)

func New712(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList712, ItemVersion712)
	blockTranslator := lookupOrCreateBlockTranslator(712, BlockVersion712, blockStateData712)
	return &Protocol{
		ver:             "1.21.20",
		id:              proto.ID712,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
