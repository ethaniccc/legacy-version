package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/legacyver/proto"
	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion685 ...
	ItemVersion685 = 191
	// BlockVersion685 ...
	BlockVersion685 int32 = (1 << 24) | (21 << 16) | (0 << 8)
)

// New685 uses same data as 686
func New685(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList686, ItemVersion685)
	blockTranslator := lookupOrCreateBlockTranslator(685, BlockVersion685, blockStateData686)
	return &Protocol{
		ver:             "1.21.0",
		id:              proto.ID685,
		blockTranslator: blockTranslator,
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockTranslator.BlockMapping(), blockMappingLatest),
	}
}
