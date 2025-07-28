package bot

import (
	"slices"
)

const (
	itemCandy  = "candy"
	itemApple  = "apple"
	itemFish   = "fish"
	itemCatnip = "catnip"
	itemMeth   = "meth"

	itemKnife = "knife"
	itemGun   = "gun"
	itemBomb  = "bomb"
	itemNuke  = "nuke"

	itemDiamond = "diamond"
)

var foodItems = []string{
	// dont ask me why you can eat knife, bomb and nuke
	// ask yourself why would you eat that lol
	itemCandy,
	itemApple,
	itemFish,
	itemCatnip,
	itemMeth,
	itemKnife,
	itemBomb,
	itemNuke,
}

var weaponItems = []string{
	itemKnife,
	itemGun,
	itemBomb,
	itemNuke,
}

var shopItems = map[string]int{
	itemCandy:  50,
	itemApple:  100,
	itemFish:   75,
	itemCatnip: 150,
	itemMeth:   300,

	itemKnife: 200,
	itemGun:   500,
	itemBomb:  700,
	itemNuke:  2000,

	itemDiamond: 1000,
}

func isFood(item string) bool {
	return slices.Contains(foodItems, item)
}

func isWeapon(item string) bool {
	return slices.Contains(weaponItems, item)
}
