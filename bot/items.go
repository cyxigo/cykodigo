package bot

import (
	"slices"
)

const (
	itemCandy  = "candy"
	itemApple  = "apple"
	itemFish   = "fish"
	itemCatnip = "catnip"

	itemKnife = "knife"
	itemGun   = "gun"
	itemBomb  = "bomb"

	itemDiamond = "diamond"
)

var foodItems = []string{
	itemCandy, itemApple, itemFish, itemCatnip,
}

var weaponItems = []string{
	itemKnife, itemGun, itemBomb,
}

var shopItems = map[string]int{
	itemCandy:  50,
	itemApple:  100,
	itemFish:   75,
	itemCatnip: 150,

	itemKnife: 200,
	itemGun:   500,
	itemBomb:  700,

	itemDiamond: 1000,
}

func isFood(item string) bool {
	return slices.Contains(foodItems, item)
}

func isWeapon(item string) bool {
	return slices.Contains(weaponItems, item)
}
