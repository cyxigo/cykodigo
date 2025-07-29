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
		fmt.Sprintf("my super duper website %s\n", data.EmojiCykodigo) +
		"P.S for nerds: Terms of Service and Privacy Policy are also there"

	featureEmbed := data.EmbedText(features)
	otherFeaturesEmbed := data.EmbedText(otherFeatures)
	websiteEmbed := data.EmbedText(website)

	content := fmt.Sprintf("**Super-duper manual %s**\n", data.EmojiCykodigo) +
		"-# Copyright (c) 2025 cyxigo"

	respondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{
		featureEmbed,
		otherFeaturesEmbed,
		websiteEmbed,
	}, false)
}

func handleMeowat(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	handleDoAtCmd(sess, inter, "meow")
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

func handleBarkAt(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	handleDoAtCmd(sess, inter, "bark")
}

func handleCat(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	files, err := os.ReadDir("res/cat")

	if err != nil {
		log.Printf("Error reading res/cat: %v", err)
		respond(sess, inter, "Couldn't find any cats :<", nil, false)

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

func handleRoulette(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	result := "**Victory!!!** You're alive!!!"
	bullet := 3

	if rand.IntN(5) == bullet {
		result = "Sorry, you're dead, better luck next ti- uhh.."
	}

	respond(sess, inter, result, nil, false)
}

func handleAssault(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	db, ok := database.GetDB(inter.GuildID)

	if !ok {
		return
	}

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
		content := fmt.Sprintf("Item **%v** isn't a weapon!!!", item)
		respond(sess, inter, content, nil, false)

		return
	}

	count := 0
	err := db.QueryRow(
		`
		SELECT COUNT(*) 
		FROM inventory 
		WHERE user_id = ? AND item = ?
		`,
		sender.ID, item).Scan(&count)

	if err != nil {
		log.Printf("Failed to scan row in /assault: %v", err)
		respond(sess, inter, "Failed to check inventory :<", nil, false)

		return
	}

	if count == 0 {
		respond(sess, inter, fmt.Sprintf("You don't have **%v** in your inventory!!!", item), nil, false)
		return
	}

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
		successMessage = "OBLITERATED THEM!!! BOOM!1!11!! 💥💥💥💥💥"
	}

	result := "failed! oops"

	if rand.IntN(99) < chance {
		result = successMessage
	}

	content += result
	respond(sess, inter, content, nil, true)
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

	lastWork := int64(0)
	err := tx.QueryRow(
		`
		SELECT last_work 
		FROM cooldowns 
		WHERE user_id = ?
		`, sender.ID).Scan(&lastWork)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to scan row in /work: %v", err)
		respond(sess, inter, "Failed to check work cooldown :<", nil, false)

		return
	}

	const cooldown = 10 * 60
	currentTime := time.Now().Unix()

	if currentTime-lastWork < cooldown {
		remaining := cooldown - (currentTime - lastWork)
		content := fmt.Sprintf("You need to wait **%vm %vs** before working again!!!", remaining/60, remaining%60)

		respond(sess, inter, content, nil, false)

		return
	}

	isHigh, _ := database.TxGetUserHighInfo(tx, sender.ID)
	// random number from range 100-200
	money := rand.IntN(100) + 100

	if isHigh {
		// apply 30% reduction if high
		money = int(float64(money) * 0.7)
	}

	if !txMoneyOp(sess, inter, tx, sender.ID, money, "+", data.CmdWork) {
		return
	}

	if !txUpdateCd(sess, inter, tx, sender.ID, "last_work", currentTime, data.CmdWork) {
		return
	}

	if !commitTx(sess, inter, tx, data.CmdWork) {
		return
	}

	content := ""

	if isHigh {
		content = fmt.Sprintf("You are **high** %s Actually, it's not good... Money from work has decreased!!!\n",
			data.EmojiCatr)
	}

	content += fmt.Sprintf("You worked and got **%v** money!", money)
	respond(sess, inter, content, nil, false)
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
		respond(sess, inter, "Failed to check balance :<", nil, false)

		return
	}

	content := fmt.Sprintf("%v's balance: **%v** money", target.Mention(), balance)
	respond(sess, inter, content, nil, true)
}

func handleTransfer(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, target, ok := getSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	if sender.ID == target.ID {
		respond(sess, inter, "You can't transfer money to yourself!!!", nil, false)
		return
	}

	options := inter.ApplicationCommandData().Options
	amount := (int)(options[1].Value.(float64))

	if amount <= 0 {
		respond(sess, inter, "Transfer amount must be positive!!!", nil, false)
		return
	}

	tx, ok := beginTx(sess, inter, data.CmdTransfer)

	if !ok {
		return
	}

	defer tx.Rollback()

	balance := database.TxGetUserBalance(tx, sender.ID)

	if balance < amount {
		respond(sess, inter, "You don't have enough money for this transfer!!!", nil, false)
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

	response := fmt.Sprintf("%v transferred %v money to %v", sender.Mention(), amount, target.Mention())
	respond(sess, inter, response, nil, true)
}

func handleSteal(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	db, ok := database.GetDB(inter.GuildID)

	if !ok {
		return
	}

	sender, target, ok := getSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	if sender.ID == target.ID {
		respond(sess, inter, "You can't steal from yourself!!!", nil, false)
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

	lastStealFail := int64(0)
	err := db.QueryRow(
		`
		SELECT last_steal_fail 
		FROM cooldowns 
		WHERE user_id = ?
		`, sender.ID).Scan(&lastStealFail)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to scan row in /steal: %v", err)
		respond(sess, inter, "Failed to check steal cooldown :<", nil, false)

		return
	}

	const cooldown = 20 * 60
	currentTime := time.Now().Unix()

	if currentTime-lastStealFail < cooldown {
		remaining := cooldown - (currentTime - lastStealFail)
		content := fmt.Sprintf("You need to wait **%vm %vs** before stealing again after failure!!!", remaining/60,
			remaining%60)

		respond(sess, inter, content, nil, false)

		return
	}

	content := ""
	successChance := 50
	isHigh, _ := database.TxGetUserHighInfo(tx, sender.ID)

	if isHigh {
		content = fmt.Sprintf("You are **high**, chances of successful steal has increased %s...\n", data.EmojiCatr)
		successChance = 80
	}

	if rand.IntN(100) < successChance {
		targetBalance := database.TxGetUserBalance(tx, target.ID)

		// max() all stuff so we cant get 0 from stealing lol
		// some of this may not be necessary but whatever
		stealPercent := max(1, rand.IntN(51))
		stealAmount := max(1, (stealPercent*targetBalance)/100)

		if !txMoneyOp(sess, inter, tx, target.ID, stealAmount, "-", data.CmdSteal) {
			return
		}

		if !txMoneyOp(sess, inter, tx, sender.ID, stealAmount, "+", data.CmdSteal) {
			return
		}

		content += fmt.Sprintf("You successfully stole **%v** money from %v!", stealAmount, target.Mention())
	} else {
		const penalty = 20

		if !txMoneyOp(sess, inter, tx, sender.ID, penalty, "-", data.CmdSteal) {
			return
		}

		if !txUpdateCd(sess, inter, tx, sender.ID, "last_steal_fail", currentTime, data.CmdSteal) {
			return
		}

		if isHigh {
			content = fmt.Sprintf("You failed to steal from %v and lost **%v** money! :<\n"+
				"Even being **high** couldn't help you %s",
				target.Mention(), penalty, data.EmojiCatr)
		} else {
			content = fmt.Sprintf("You failed to steal from %v and lost **%v** money! :<\n"+
				"**Pro tip:** being **high** increases chances of a successful steal %s",
				target.Mention(), penalty, data.EmojiCatr)
		}

	}

	if !commitTx(sess, inter, tx, data.CmdSteal) {
		return
	}

	respond(sess, inter, content, nil, true)
}

func handleShop(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	builder := strings.Builder{}
	content := fmt.Sprintf("**Shop %s**\n", data.EmojiCykodigo) +
		"-# Use `/buy [item]` to buy something!!!"

	for item, price := range data.ShopItems {
		itemString := fmt.Sprintf("- %v: %v money\n", item, price)
		builder.WriteString(itemString)
	}

	embed := data.EmbedText(builder.String())
	respondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{embed}, false)
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

	tx, ok := beginTx(sess, inter, data.CmdBuy)

	if !ok {
		return
	}

	defer tx.Rollback()

	balance := database.TxGetUserBalance(tx, sender.ID)

	if balance < price {
		content := fmt.Sprintf("Too broke for **%v**, go work!!!", item)
		respond(sess, inter, content, nil, false)

		return
	}

	if !txMoneyOp(sess, inter, tx, sender.ID, price, "-", data.CmdBuy) {
		return
	}

	if _, err := tx.Exec("INSERT INTO inventory(user_id, item) VALUES(?, ?)", sender.ID, item); err != nil {
		log.Printf("Insert error in /inventory: %v", err)
		respond(sess, inter, "Failed to add item to inventory :<", nil, false)

		return
	}

	if !commitTx(sess, inter, tx, data.CmdBuy) {
		return
	}

	respond(sess, inter, fmt.Sprintf("You bought **%v** for **%v** money", item, price), nil, false)
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
		SELECT item 
		FROM inventory 
		WHERE user_id = ?
		`, target.ID)

	if err != nil {
		log.Printf("Query error in /inventory: %v", err)
		respond(sess, inter, "Failed to check inventory :<", nil, false)

		return
	}

	defer rows.Close()

	// int is amount of items user has, so its like
	// item: count
	items := make(map[string]int)

	for rows.Next() {
		item := ""

		if err := rows.Scan(&item); err != nil {
			log.Printf("Failed to scan row in /inventory: %v", err)
			continue
		}

		items[item]++
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
		SELECT user_id, COUNT(*) AS diamond_count
		FROM inventory
		WHERE item = ?
		GROUP BY user_id
		ORDER BY diamond_count DESC
		LIMIT 10
		`, data.ItemDiamond)

	if err != nil {
		log.Printf("Query error in /leaderboard: %v", err)
		respond(sess, inter, "Failed to check leaderboard :<", nil, false)

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
			"Buy some with `/buy diamond`!!!"
		respond(sess, inter, content, nil, false)

		return
	}

	content := "**Diamond Leaderboard**\n" +
		"-# Buy some with `/buy diamond`!!!"
	embed := data.EmbedText(strings.Join(leaderboard, "\n"))
	respondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{embed}, false)
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
		content := fmt.Sprintf("You can't eat **%v**!!!", item)
		respond(sess, inter, content, nil, false)

		return
	}

	tx, ok := beginTx(sess, inter, data.CmdEat)

	if !ok {
		return
	}

	defer tx.Rollback()

	count := 0
	err := tx.QueryRow(
		`
		SELECT COUNT(*) 
		FROM inventory
		WHERE user_id = ? AND item = ?
		`,
		sender.ID, item).Scan(&count)

	if err != nil {
		log.Printf("Failed to scan row in /eat: %v", err)
		respond(sess, inter, "Failed to check inventory :<", nil, false)

		return
	}

	if count == 0 {
		respond(sess, inter, fmt.Sprintf("You don't have **%v** in your inventory!!!", item), nil, false)
		return
	}

	_, err = tx.Exec(
		`
		DELETE FROM inventory 
		WHERE rowid = (
			SELECT rowid 
			FROM inventory 
			WHERE user_id = ? AND item = ? 
			LIMIT 1
		)
		`,
		sender.ID, item)

	if err != nil {
		log.Printf("Delete error in /eat: %v", err)
		respond(sess, inter, "Failed to get inventory :<", nil, false)

		return
	}

	duration := int64(0)

	if item == data.ItemMeth {
		const effectDuration = 5 * 60
		currentTime := time.Now().Unix()
		newEndTime := currentTime + effectDuration

		_, err = tx.Exec(
			`
			INSERT INTO meth_effects(user_id, end_time)
			VALUES(?, ?)
			ON CONFLICT(user_id) DO UPDATE SET
				end_time = MAX(?, end_time) + ?
			`,
			sender.ID, newEndTime, currentTime, effectDuration)

		if err != nil {
			log.Printf("Failed to update meth effect in /eat: %v", err)
			respond(sess, inter, "Failed to get **high** :<", nil, false)

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
			respond(sess, inter, "Failed to get **high** :<", nil, false)

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
			"You're now high for %vm %vs %s", item, duration/60, duration%60, data.EmojiCatr)
		handleImageCmd(sess, inter, content, "res/gif/spin.gif")
	} else {
		message := "Yummy! " + data.EmojiCykodigo

		switch item {
		case data.ItemKnife:
			message = "Oh wait why did you do that??? You're dead from several internal bleeds."
		case data.ItemBomb:
		case data.ItemNuke:
			message = "BOOM!1!11!!💥💥💥💥💥"
		}

		content := fmt.Sprintf("You ate **%v**! %s", item, message)
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
		content = fmt.Sprintf("%s is **high**! Wowowowowowowo...\nTime remaining: **%vm %vs**%s",
			target.Mention(), remaining/60, remaining%60, data.EmojiCatr)
	} else {
		content = fmt.Sprintf("%s isn't **high**! Lame", target.Mention())
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
		handleMeowat(sess, inter)
	case data.CmdBark:
		handleBark(sess, inter)
	case data.CmdBarkat:
		handleBarkAt(sess, inter)
	case data.CmdDoflip:
		handleImageCmd(sess, inter, "Woah!", "res/flip.png")
	case data.CmdExplode:
		handleImageCmd(sess, inter, "WHAAAAAAAA-", "res/boom.png")
	case data.CmdSpin:
		// uhh yes im using handleImageCmd for sending a gif
		// and what? what you gonna do?
		handleImageCmd(sess, inter, "Wooooooo", "res/gif/spin.gif")
	case data.CmdCat:
		handleCat(sess, inter)
	case data.CmdCart:
		handleImageCmd(sess, inter, "Cart!", "res/cart.png")
	case data.CmdRoulette:
		handleRoulette(sess, inter)
	case data.CmdAssault:
		handleAssault(sess, inter)
	case data.CmdWork:
		handleWork(sess, inter)
	case data.CmdBalance:
		handleBalance(sess, inter)
	case data.CmdTransfer:
		handleTransfer(sess, inter)
	case data.CmdSteal:
		handleSteal(sess, inter)
	case data.CmdShop:
		handleShop(sess, inter)
	case data.CmdBuy:
		handleBuy(sess, inter)
	case data.CmdInventory:
		handleInventory(sess, inter)
	case data.CmdLeaderboard:
		handleLeaderboard(sess, inter)
	case data.CmdEat:
		handleEat(sess, inter)
	case data.CmdHigh:
		handleHigh(sess, inter)
	}
}
