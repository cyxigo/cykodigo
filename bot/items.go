package bot

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	itemCandy  = "candy"
	itemApple  = "apple"
	itemFish   = "fish"
	itemCatnip = "catnip"

	itemKnife = "knife"
	itemGun   = "gun"

	itemDiamond = "diamond"
)

var foodItems = []string{
	itemCandy, itemApple, itemFish, itemCatnip,
}

var weaponItems = []string{
	itemKnife, itemGun,
}

var shopItems = map[string]int{
	itemCandy:   50,
	itemApple:   100,
	itemFish:    75,
	itemCatnip:  150,
	itemKnife:   200,
	itemGun:     500,
	itemDiamond: 1000,
}

func isFood(item string) bool {
	return slices.Contains(foodItems, item)
}

func isWeapon(item string) bool {
	return slices.Contains(weaponItems, item)
}

// util function for getting an item from command option
// returns item and its price
func getItemFromInterOption(sess *discordgo.Session, inter *discordgo.InteractionCreate, idx int) (
	string, int, bool) {
	item := strings.ToLower(inter.ApplicationCommandData().Options[idx].StringValue())
	price, exists := shopItems[item]

	if !exists {
		content := fmt.Sprintf("There's no item **%s**!!!", item)
		respond(sess, inter, content, nil, false)

		return "", 0, false
	}

	return item, price, true
}
