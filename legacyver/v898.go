package legacyver

import (
	_ "embed"

	"github.com/ethaniccc/legacy-version/mapping"
)

const (
	// ItemVersion898 ...
	ItemVersion898 = 251
	// BlockVersion898 ...
	BlockVersion898 int32 = (1 << 24) | (21 << 16) | (130 << 8)
)

var (
	//go:embed data/dragonfly_items.json
	dragonflyLatestItemList []byte
	//go:embed data/required_item_list_898.json
	requiredItemList898 []byte
	//go:embed data/block_states_898.nbt
	blockStateData898 []byte

	itemMappingLatestPocketMine = mapping.NewItemMapping(requiredItemList898, ItemVersion898)
	itemMappingLatestDragonfly  = mapping.NewItemMapping(dragonflyLatestItemList, ItemVersion898)
	blockMappingLatest          = mapping.NewBlockMapping(blockStateData898)
)

func itemMappingLatest(dragonflyMapping bool) mapping.Item {
	if dragonflyMapping {
		return itemMappingLatestDragonfly
	}
	return itemMappingLatestPocketMine
}
