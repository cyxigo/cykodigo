package data

import (
	"slices"
)

const (
	ItemCandy   = "candy"
	ItemApple   = "apple"
	ItemFish    = "fish"
	ItemCatnip  = "catnip"
	ItemMeth    = "meth"
	ItemCocaine = "cocaine"
	ItemChalk   = "chalk"

	ItemKnife       = "knife"
	ItemGun         = "gun"
	ItemBomb        = "bomb"
	ItemNuke        = "nuke"
	ItemDevilsKnife = "devilsknife"

	ItemDiamond = "diamond"
)

var FoodItems = []string{
	// dont ask me why you can eat knife, bomb and nuke
	// ask yourself why would you eat that lol
	ItemCandy,
	ItemApple,
	ItemFish,
	ItemChalk,
	ItemKnife,
	ItemBomb,
	ItemNuke,
}

var DrugItems = []string{
	ItemCatnip,
	ItemMeth,
	ItemCocaine,
}

var WeaponItems = []string{
	ItemKnife,
	ItemGun,
	ItemBomb,
	ItemNuke,
	ItemDevilsKnife,
}

var ShopItems = map[string]int{
	ItemCandy:   50,
	ItemApple:   100,
	ItemFish:    75,
	ItemCatnip:  150,
	ItemMeth:    350,
	ItemCocaine: 200,
	ItemChalk:   100,

	ItemKnife:       200,
	ItemGun:         500,
	ItemBomb:        700,
	ItemNuke:        2000,
	ItemDevilsKnife: 3000,

	ItemDiamond: 1000,
}

func IsFood(item string) bool {
	return slices.Contains(FoodItems, item) || slices.Contains(DrugItems, item)
}

// note: also returns duration
func IsDrug(item string) (int64, bool) {
	exists := slices.Contains(DrugItems, item)

	if !exists {
		return 0, false
	}

	switch item {
	case ItemCatnip:
		return 1 * 60, true
	case ItemMeth:
		return 5 * 60, true
	case ItemCocaine:
		return 3 * 60, true
	default:
		return 0, false
	}
}

func IsWeapon(item string) bool {
	return slices.Contains(WeaponItems, item)
}
