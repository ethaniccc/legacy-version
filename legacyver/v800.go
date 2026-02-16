package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion800 ...
	ItemVersion800 = 231
	// BlockVersion800 ...
	BlockVersion800 int32 = (1 << 24) | (21 << 16) | (80 << 8)
)

var (
	//go:embed data/required_item_list_800.json
	requiredItemList800 []byte
	//go:embed data/block_states_800.nbt
	blockStateData800 []byte
)

func New800(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList800, ItemVersion800)
	blockTranslator := lookupOrCreateBlockTranslator(800, BlockVersion800, blockStateData800)
	return &Protocol{
		ver:             "1.21.80",
		id:              proto.ID800,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
