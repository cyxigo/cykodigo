package data

import (
	"slices"
)

const (
	ItemCandy  = "candy"
	ItemApple  = "apple"
	ItemFish   = "fish"
	ItemCatnip = "catnip"
	ItemMeth   = "meth"
	ItemChalk  = "chalk"

	ItemKnife = "knife"
	ItemGun   = "gun"
	ItemBomb  = "bomb"
	ItemNuke  = "nuke"

	ItemDiamond = "diamond"
)

var FoodItems = []string{
	// dont ask me why you can eat knife, bomb and nuke
	// ask yourself why would you eat that lol
	ItemCandy,
	ItemApple,
	ItemFish,
	ItemCatnip,
	ItemMeth,
	ItemChalk,
	ItemKnife,
	ItemBomb,
	ItemNuke,
}

var WeaponItems = []string{
	ItemKnife,
	ItemGun,
	ItemBomb,
	ItemNuke,
}

var ShopItems = map[string]int{
	ItemCandy:  50,
	ItemApple:  100,
	ItemFish:   75,
	ItemCatnip: 150,
	ItemMeth:   300,
	ItemChalk:  100,

	ItemKnife: 200,
	ItemGun:   500,
	ItemBomb:  700,
	ItemNuke:  2000,

	ItemDiamond: 1000,
}

func IsFood(item string) bool {
	return slices.Contains(FoodItems, item)
}

func IsWeapon(item string) bool {
	return slices.Contains(WeaponItems, item)
}
