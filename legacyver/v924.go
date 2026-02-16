package legacyver

import (
	_ "embed"

	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion924 ...
	ItemVersion924 = 251
	// BlockVersion924 ...
	BlockVersion924 int32 = (1 << 24) | (26 << 16) | (0 << 8)
)

var (
	//go:embed data/dragonfly_items.json
	dragonflyLatestItemList []byte
	//go:embed data/required_item_list_924.json
	requiredItemList924 []byte
	//go:embed data/block_states_924.nbt
	blockStateData924 []byte

	itemMappingLatestPocketMine = mapping.NewItemMapping(requiredItemList924, ItemVersion924)
	itemMappingLatestDragonfly  = mapping.NewItemMapping(dragonflyLatestItemList, ItemVersion924)
	blockMappingLatest          = mapping.NewBlockMapping(blockStateData924)
)

func itemMappingLatest(dragonflyMapping bool) mapping.Item {
	if dragonflyMapping {
		return itemMappingLatestDragonfly
	}
	return itemMappingLatestPocketMine
}
