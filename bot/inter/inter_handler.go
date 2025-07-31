package inter

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cyxigo/cykodigo/bot/data"
	"github.com/cyxigo/cykodigo/bot/database"
)

func handleHelp(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	features := "**Features**\n" +
		"This bot isn't so serious. But he has money system!\n" +
		"You can work using `/work` command, buy stuff using `/buy`, check your balance using `/balance` " +
		"and check stuff that you bought using `/inventory`\n" +
		"You can also steal money from people using `/steal` and transfer money using `/transfer`"
	otherFeatures := "**Other features**\n" +
		"cykodigo has many other commands and stuff, I didn't describe all his features hehe\n" +
		"I think it would be more fun if you found out all this by yourself\n" +
		"Meow-meow!"
	website := "**My website!**\n" +
		"I don't think anyone will look in here, but [here](https://cyxigo.github.io/cykodigo-io/) " +
		fmt.Sprintf("my super duper website %v\n", data.EmojiCykodigo) +
		"P.S for nerds: Terms of Service and Privacy Policy are also there"
	waitWhat := "**Wait what? I have no money on another server, this bot is stupid!**\n" +
		"Well, it's not a bot, it's you. cykodigo is smart, cykodigo thinks it will be more fun if each server had " +
		"their own local money"

	featureEmbed := data.EmbedText(features)
	otherFeaturesEmbed := data.EmbedText(otherFeatures)
	websiteEmbed := data.EmbedText(website)
	waitWhatEmbed := data.EmbedText(waitWhat)

	content := fmt.Sprintf("**Super-duper manual %v**\n", data.EmojiCykodigo) +
		"-# Copyright (c) 2025 cyxigo"

	respondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{
		featureEmbed,
		otherFeaturesEmbed,
		websiteEmbed,
		waitWhatEmbed,
	}, false)
}

func handleBark(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	content := "I'm a cat, I can't bark you "
	compliments := []string{
		"idiot",
		"dumbass",
		"dog",
		"stupid",
	}

	content += compliments[rand.IntN(len(compliments))]
	respond(sess, inter, content, nil, false)
}

func handleCat(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	files, err := os.ReadDir("res/cat")

	if err != nil {
		log.Printf("Error reading res/cat: %v", err)
		respond(sess, inter, "Couldn't find any cats", nil, false)

		return
	}

	pngFiles := []string{}

	for _, file := range files {
		if !file.IsDir() &&
			strings.HasSuffix(strings.ToLower(file.Name()), ".png") {
			pngFiles = append(pngFiles, file.Name())
		}
	}

	if len(pngFiles) == 0 {
		respond(sess, inter, "Couldn't find any cats", nil, false)
		return
	}

	img := pngFiles[rand.IntN(len(pngFiles))]
	path := fmt.Sprintf("res/cat/%v", img)

	handleImageCmd(sess, inter, "Cat!", path)
}

func handleAssault(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, target, ok := getSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	// note: Options[0] is target
	item, _, ok := getItemFromOption(sess, inter, 1)

	if !ok {
		return
	}

	if !data.IsWeapon(item) {
		content := fmt.Sprintf("Item **%v** isn't a weapon", item)
		respond(sess, inter, content, nil, false)

		return
	}

	tx, ok := beginTx(sess, inter, data.CmdAssault)

	if !ok || !checkInventory(sess, inter, tx, sender.ID, item, 1, data.CmdAssault) {
		return
	}

	defer tx.Rollback()

	if !ok {
		return
	}

	chance := 0
	content := ""
	successMessage := "killed them!"

	switch item {
	case data.ItemKnife:
		chance = 20
		content = fmt.Sprintf("%v tried to stab %v with a knife and... ", sender.Mention(), target.Mention())
	case data.ItemGun:
		chance = 70
		content = fmt.Sprintf("%v tried to shoot %v with a gun and... ", sender.Mention(), target.Mention())
	case data.ItemBomb:
		chance = 90
		content = fmt.Sprintf("%v threw a bomb at %v and... ", sender.Mention(), target.Mention())
		successMessage = "BOOM!1!11!! 💥💥💥💥💥"
	case data.ItemNuke:
		chance = 100
		content = fmt.Sprintf("%v nuked %v and... ", sender.Mention(), target.Mention())
		successMessage = "**OBLITERATED** THEM!!! BOOM!1!11!! 💥💥💥💥💥"
	}

	result := "failed! oops"

	if rand.IntN(99) < chance {
		result = successMessage
	}

	content += result
	respond(sess, inter, content, nil, true)
}

func handleBalance(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	db, ok := database.GetDB(inter.GuildID)

	if !ok {
		return
	}

	target, ok := getOptionalTarget(sess, inter)

	if !ok {
		return
	}

	balance := 0
	err := db.QueryRow(
		`
		SELECT balance 
		FROM balances 
		WHERE user_id = ?
		`, target.ID).Scan(&balance)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to scan row in /balance: %v", err)
		respond(sess, inter, "Failed to check balance", nil, false)

		return
	}

	content := fmt.Sprintf("%v's balance: **%v** money %v", target.Mention(), balance, data.EmojiCykodigo)
	respond(sess, inter, content, nil, true)
}

func handleShop(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	builder := strings.Builder{}
	content := fmt.Sprintf("**Shop %v**\n", data.EmojiCykodigo) +
		"-# Use `/buy [item]` to buy something"

	for item, price := range data.ShopItems {
		itemString := fmt.Sprintf("- %v: %v money\n", item, price)
		builder.WriteString(itemString)
	}

	embed := data.EmbedText(builder.String())
	respondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{embed}, false)
}

func handleInventory(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	db, ok := database.GetDB(inter.GuildID)

	if !ok {
		return
	}

	target, ok := getOptionalTarget(sess, inter)

	if !ok {
		return
	}

	rows, err := db.Query(
		`
		SELECT item, amount
		FROM inventory 
		WHERE user_id = ?
		`, target.ID)

	if err != nil {
		log.Printf("Query error in /inventory: %v", err)
		respond(sess, inter, "Failed to check inventory", nil, false)

		return
	}

	defer rows.Close()

	// int is amount of items user has, so its like
	// item: count
	items := make(map[string]int)

	for rows.Next() {
		item := ""
		amount := 0

		if err := rows.Scan(&item, &amount); err != nil {
			log.Printf("Failed to scan row in /inventory: %v", err)
			continue
		}

		items[item] = amount
	}

	if len(items) == 0 {
		content := fmt.Sprintf("%v inventory: oops! such an empty", target.Mention())
		respond(sess, inter, content, nil, true)

		return
	}

	builder := strings.Builder{}
	content := fmt.Sprintf("%v inventory:\n", target.Mention())

	for item, count := range items {
		itemString := fmt.Sprintf("- %v ×%v\n", item, count)
		builder.WriteString(itemString)
	}

	embed := data.EmbedText(builder.String())
	respondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{embed}, false)
}

func handleLeaderboard(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	db, ok := database.GetDB(inter.GuildID)

	if !ok {
		return
	}

	rows, err := db.Query(
		`
        SELECT user_id, amount
        FROM inventory
        WHERE item = ?
        ORDER BY amount DESC
        LIMIT 10
        `, data.ItemDiamond)

	if err != nil {
		log.Printf("Query error in /leaderboard: %v", err)
		respond(sess, inter, "Failed to check leaderboard", nil, false)

		return
	}

	defer rows.Close()

	leaderboard := []string{}
	position := 1

	for rows.Next() {
		userID := ""
		count := 0

		if err := rows.Scan(&userID, &count); err != nil {
			log.Printf("Failed to scan row in /leaderboard: %v", err)
			continue
		}

		user, err := sess.User(userID)
		displayName := "Unknown user"

		if err == nil {
			displayName = user.Username
		} else {
			log.Printf("Failed to get display name for '%v' in /leaderboard: %v", userID, err)
		}

		entry := fmt.Sprintf(
			"%v. **%v** - `%v` diamonds",
			position,
			displayName,
			count,
		)

		leaderboard = append(leaderboard, entry)
		position++
	}

	if len(leaderboard) == 0 {
		content := "No one has diamonds yet\n" +
			"Such an empty leaderboard\n" +
			"Buy some with `/buy diamond`"
		respond(sess, inter, content, nil, false)

		return
	}

	content := fmt.Sprintf("**Diamond Leaderboard** %v\n", data.EmojiCykodigo) +
		"-# Buy some with `/buy diamond`"
	embed := data.EmbedText(strings.Join(leaderboard, "\n"))
	respondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{embed}, false)
}

func handleRoulette(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, ok := getSender(sess, inter)

	if !ok {
		return
	}

	options := inter.ApplicationCommandData().Options
	bet := options[0].IntValue()

	if bet < 100 {
		respond(sess, inter, "Minimum bet is 100 money", nil, false)
		return
	}

	tx, ok := beginTx(sess, inter, data.CmdRoulette)

	if !ok {
		return
	}

	defer tx.Rollback()

	balance := database.TxGetUserBalance(tx, sender.ID)

	if bet > 50000 {
		content := fmt.Sprintf("Hey don't bet this much %v", data.EmojiCykodigo)
		respond(sess, inter, content, nil, false)

		return
	}

	if balance < bet {
		respond(sess, inter, "You don't have enough money for this bet, go work", nil, false)
		return
	}

	if !txMoneyOp(sess, inter, tx, sender.ID, bet, "-", data.CmdRoulette) {
		return
	}

	content := ""
	successChance := 30
	isHigh, _ := database.TxGetUserHighInfo(tx, sender.ID)

	if isHigh {
		content = fmt.Sprintf("You are **high** %v, chances of successful bet has increased\n", data.EmojiCatr)
		successChance = 50
	}

	if rand.IntN(100) < successChance {
		winnings := int64(float64(bet) * 1.5)

		if !txMoneyOp(sess, inter, tx, sender.ID, winnings, "+", data.CmdRoulette) {
			return
		}

		content += fmt.Sprintf("**JACKPOT** You won **%v** money %v", winnings, data.EmojiCykodigo)
	} else {
		if isHigh {
			content = fmt.Sprintf("You lost your bet of **%v** money. Better luck next time\n"+
				"Even being **high** couldn't help you %v",
				bet, data.EmojiCatr)
		} else {
			content = fmt.Sprintf("You lost your bet of **%v** money. Better luck next time\n"+
				"**Pro tip:** being **high** increases chances of a successful bet %v",
				bet, data.EmojiCatr)
		}
	}

	if !commitTx(sess, inter, tx, data.CmdRoulette) {
		return
	}

	respond(sess, inter, content, nil, false)
}

func handleWork(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, ok := getSender(sess, inter)

	if !ok {
		return
	}

	tx, ok := beginTx(sess, inter, data.CmdWork)

	if !ok {
		return
	}

	defer tx.Rollback()

	if !checkCooldown(sess, inter, tx, sender.ID, "last_work", 10*60, data.CmdWork) {
		return
	}

	isHigh, _ := database.TxGetUserHighInfo(tx, sender.ID)
	// random number from range 100-300
	money := rand.Int64N(200) + 100

	if isHigh {
		// apply 30% reduction if high
		money = int64(float64(money) * 0.7)
	}

	if !txMoneyOp(sess, inter, tx, sender.ID, money, "+", data.CmdWork) {
		return
	}

	if !txUpdateCd(sess, inter, tx, sender.ID, "last_work", time.Now().Unix(), data.CmdWork) {
		return
	}

	if !commitTx(sess, inter, tx, data.CmdWork) {
		return
	}

	content := ""

	if isHigh {
		content = fmt.Sprintf("You are **high** %v, actually, it's not good... Money from work has decreased\n",
			data.EmojiCatr)
	}

	content += fmt.Sprintf("You worked and got **%v** money %v", money, data.EmojiCykodigo)
	respond(sess, inter, content, nil, false)
}

func handleTransfer(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, target, ok := getSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	if sender.ID == target.ID {
		respond(sess, inter, "You can't transfer money to yourself", nil, false)
		return
	}

	options := inter.ApplicationCommandData().Options
	amount := options[1].IntValue()

	if amount <= 0 {
		respond(sess, inter, "Transfer amount must be positive", nil, false)
		return
	}

	tx, ok := beginTx(sess, inter, data.CmdTransfer)

	if !ok {
		return
	}

	defer tx.Rollback()

	balance := database.TxGetUserBalance(tx, sender.ID)

	if balance < amount {
		respond(sess, inter, "You don't have enough money for this transfer", nil, false)
		return
	}

	if !txMoneyOp(sess, inter, tx, sender.ID, amount, "-", data.CmdTransfer) {
		return
	}

	if !txMoneyOp(sess, inter, tx, target.ID, amount, "+", data.CmdTransfer) {
		return
	}

	if !commitTx(sess, inter, tx, data.CmdTransfer) {
		return
	}

	response := fmt.Sprintf("%v transferred %v money to %v %v", sender.Mention(), amount, target.Mention(),
		data.EmojiCykodigo)
	respond(sess, inter, response, nil, true)
}

func handleSteal(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, target, ok := getSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	if sender.ID == target.ID {
		respond(sess, inter, "You can't steal from yourself", nil, false)
		return
	}

	tx, ok := beginTx(sess, inter, data.CmdSteal)

	if !ok {
		return
	}

	defer tx.Rollback()

	targetBalance := database.TxGetUserBalance(tx, target.ID)

	if targetBalance <= 0 {
		content := fmt.Sprintf("%v is **broke!** Nothing to steal", target.Mention())
		respond(sess, inter, content, nil, true)

		return
	}

	if !checkCooldown(sess, inter, tx, sender.ID, "last_steal_fail", 20*60, data.CmdSteal) {
		return
	}

	content := ""
	successChance := 50
	isHigh, _ := database.TxGetUserHighInfo(tx, sender.ID)

	if isHigh {
		content = fmt.Sprintf("You are **high** %v, chances of successful steal has increased\n", data.EmojiCatr)
		successChance = 80
	}

	if rand.IntN(100) < successChance {
		targetBalance := database.TxGetUserBalance(tx, target.ID)

		// max() all stuff so we cant get 0 from stealing lol
		// some of this may not be necessary but whatever
		stealPercent := max(1, rand.Int64N(51))
		stealAmount := max(1, (stealPercent*targetBalance)/100)

		if !txMoneyOp(sess, inter, tx, target.ID, stealAmount, "-", data.CmdSteal) {
			return
		}

		if !txMoneyOp(sess, inter, tx, sender.ID, stealAmount, "+", data.CmdSteal) {
			return
		}

		content += fmt.Sprintf("You successfully stole **%v** money from %v %v", stealAmount, target.Mention(),
			data.EmojiCykodigo)
	} else {
		const penalty = 20

		if !txMoneyOp(sess, inter, tx, sender.ID, penalty, "-", data.CmdSteal) {
			return
		}

		if !txUpdateCd(sess, inter, tx, sender.ID, "last_steal_fail", time.Now().Unix(), data.CmdSteal) {
			return
		}

		if isHigh {
			content = fmt.Sprintf("You failed to steal from %v and lost **%v** money\n"+
				"Even being **high** couldn't help you %v",
				target.Mention(), penalty, data.EmojiCatr)
		} else {
			content = fmt.Sprintf("You failed to steal from %v and lost **%v** money\n"+
				"**Pro tip:** being **high** increases chances of a successful steal %v",
				target.Mention(), penalty, data.EmojiCatr)
		}
	}

	if !commitTx(sess, inter, tx, data.CmdSteal) {
		return
	}

	respond(sess, inter, content, nil, true)
}

func handleBuy(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, ok := getSender(sess, inter)

	if !ok {
		return
	}

	item, price, ok := getItemFromOption(sess, inter, 0)

	if !ok {
		return
	}

	amount, ok := getItemAmountOption(sess, inter, "buy", 1)

	if !ok {
		return
	}

	price *= amount

	tx, ok := beginTx(sess, inter, data.CmdBuy)

	if !ok {
		return
	}

	defer tx.Rollback()

	balance := database.TxGetUserBalance(tx, sender.ID)

	if balance < price {
		content := fmt.Sprintf("Too broke for **%v** (x%v), go work", item, amount)
		respond(sess, inter, content, nil, false)

		return
	}

	if !txMoneyOp(sess, inter, tx, sender.ID, price, "-", data.CmdBuy) {
		return
	}

	if !updateInventory(sess, inter, tx, sender.ID, item, amount, data.CmdBuy) {
		return
	}

	if !commitTx(sess, inter, tx, data.CmdBuy) {
		return
	}

	content := fmt.Sprintf("You bought **%v** (x%v) for **%v** money %v", item, amount, price, data.EmojiCykodigo)
	respond(sess, inter, content, nil, false)
}

func handleGive(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, target, ok := getSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	item, _, ok := getItemFromOption(sess, inter, 1)

	if !ok {
		return
	}

	amount, ok := getItemAmountOption(sess, inter, "give", 2)

	if !ok {
		return
	}

	tx, ok := beginTx(sess, inter, data.CmdGive)

	if !ok {
		return
	}

	defer tx.Rollback()

	if !checkInventory(sess, inter, tx, sender.ID, item, amount, data.CmdGive) {
		return
	}

	if !updateInventory(sess, inter, tx, sender.ID, item, -amount, data.CmdGive) {
		return
	}

	if !updateInventory(sess, inter, tx, target.ID, item, amount, data.CmdGive) {
		return
	}

	if !commitTx(sess, inter, tx, data.CmdGive) {
		return
	}

	content := fmt.Sprintf("%v gave **%v** (x%v) to %v %v", sender.Mention(), item, amount, target.Mention(),
		data.EmojiCykodigo)

	respond(sess, inter, content, nil, true)
}

func handleEat(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, ok := getSender(sess, inter)

	if !ok {
		return
	}

	item, _, ok := getItemFromOption(sess, inter, 0)

	if !ok {
		return
	}

	if !data.IsFood(item) {
		content := fmt.Sprintf("You can't eat **%v**", item)
		respond(sess, inter, content, nil, false)

		return
	}

	tx, ok := beginTx(sess, inter, data.CmdEat)

	if !ok {
		return
	}

	defer tx.Rollback()

	if !checkInventory(sess, inter, tx, sender.ID, item, 1, data.CmdEat) {
		return
	}

	if !updateInventory(sess, inter, tx, sender.ID, item, -1, data.CmdEat) {
		return
	}

	duration := int64(0)

	if item == data.ItemMeth {
		const effectDuration = 5 * 60
		currentTime := time.Now().Unix()
		newEndTime := currentTime + effectDuration

		_, err := tx.Exec(
			`
			INSERT INTO meth_effects(user_id, end_time)
			VALUES(?, ?)
			ON CONFLICT(user_id) 
			DO UPDATE SET end_time = GREATEST(?, end_time) + ?
			`,
			sender.ID, newEndTime, currentTime, effectDuration)

		if err != nil {
			log.Printf("Failed to update meth effect in /eat: %v", err)
			respond(sess, inter, "Failed to get **high**", nil, false)

			return
		}

		updatedEndTime := int64(0)
		err = tx.QueryRow(
			`
			SELECT end_time 
			FROM meth_effects 
			WHERE user_id = ?
			`,
			sender.ID).Scan(&updatedEndTime)

		if err != nil {
			log.Printf("Failed to get updated end time in /eat: %v", err)
			respond(sess, inter, "Failed to get **high**", nil, false)

			return
		}

		duration = updatedEndTime - currentTime
	} else {
		_, endTime := database.TxGetUserHighInfo(tx, sender.ID)
		duration = endTime - time.Now().Unix()
	}

	if !commitTx(sess, inter, tx, data.CmdEat) {
		return
	}

	if item == data.ItemMeth {
		content := fmt.Sprintf("You ate %v! Wowowowowowowowowowo the world is spinning wowowowowowo...\n\n"+
			"You're now high for %vm %vs %v", item, duration/60, duration%60, data.EmojiCatr)
		handleImageCmd(sess, inter, content, "res/gif/spin.gif")
	} else {
		message := "Yummy! " + data.EmojiCykodigo

		switch item {
		case data.ItemKnife:
			message = "Oh wait why did you do that??? You're dead from several internal bleeds."
		case data.ItemBomb:
		case data.ItemNuke:
			message = "BOOM!1!11!!💥💥💥💥💥"
		case data.ItemChalk:
			message = "Crunchy dammit!"
		}

		content := fmt.Sprintf("You ate **%v**! %v", item, message)
		respond(sess, inter, content, nil, false)
	}
}

func handleHigh(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	target, ok := getOptionalTarget(sess, inter)

	if !ok {
		return
	}

	db, ok := database.GetDB(inter.GuildID)

	if !ok {
		return
	}

	isHigh, endTime := database.GetUserHighInfo(db, target.ID)
	currentTime := time.Now().Unix()
	content := ""

	if isHigh {
		remaining := endTime - currentTime
		content = fmt.Sprintf("%v is **high**! Wowowowowowowo...\nTime remaining: **%vm %vs**%v",
			target.Mention(), remaining/60, remaining%60, data.EmojiCatr)
	} else {
		content = fmt.Sprintf("%v isn't **high**! Lame", target.Mention())
	}

	respond(sess, inter, content, nil, true)
}

// handler for interactions
// used only for / commands
func InterHandler(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	// we dont need to handle some weird interaction stuff here
	// only commands
	if inter.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch inter.ApplicationCommandData().Name {
	case data.CmdHelp:
		handleHelp(sess, inter)
	case data.CmdMe:
		handleImageCmd(sess, inter, "Me!", "res/me.png")
	case data.CmdMeow:
		respond(sess, inter, "Meow!", nil, false)
	case data.CmdMeowat:
		handleActionOnCmd(sess, inter, "meows at")
	case data.CmdBark:
		handleBark(sess, inter)
	case data.CmdBarkat:
		handleActionOnCmd(sess, inter, "barks at")
	case data.CmdFlip:
		handleImageCmd(sess, inter, "Woah!", "res/flip.png")
	case data.CmdExplode:
		handleImageCmd(sess, inter, "WHAAAAAAAA-", "res/boom.png")
	case data.CmdSpin:
		handleImageCmd(sess, inter, "Wooooooo", "res/gif/spin.gif")
	case data.CmdCat:
		handleCat(sess, inter)
	case data.CmdCart:
		handleImageCmd(sess, inter, "Cart!", "res/cart.png")
	case data.CmdHate:
		handleActionOnCmd(sess, inter, "HATES")
	case data.CmdAssault:
		handleAssault(sess, inter)
	case data.CmdBalance:
		handleBalance(sess, inter)
	case data.CmdShop:
		handleShop(sess, inter)
	case data.CmdInventory:
		handleInventory(sess, inter)
	case data.CmdLeaderboard:
		handleLeaderboard(sess, inter)
	case data.CmdRoulette:
		handleRoulette(sess, inter)
	case data.CmdWork:
		handleWork(sess, inter)
	case data.CmdTransfer:
		handleTransfer(sess, inter)
	case data.CmdSteal:
		handleSteal(sess, inter)
	case data.CmdBuy:
		handleBuy(sess, inter)
	case data.CmdGive:
		handleGive(sess, inter)
	case data.CmdEat:
		handleEat(sess, inter)
	case data.CmdHigh:
		handleHigh(sess, inter)
	}
}
