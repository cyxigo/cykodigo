package bot

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
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
		fmt.Sprintf("my super duper website %s\n", emojiCykodigo) +
		"P.S for nerds: Terms of Service and Privacy Policy are also there"

	featureEmbed := embedText(features)
	otherFeaturesEmbed := embedText(otherFeatures)
	websiteEmbed := embedText(website)

	content := fmt.Sprintf("**Super-duper manual %s**\n", emojiCykodigo) +
		"-# Copyright (c) 2025 cyxigo"

	interRespondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{
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
	interRespond(sess, inter, content, nil, false)
}

func handleBarkAt(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	handleDoAtCmd(sess, inter, "bark")
}

func handleCat(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	files, err := os.ReadDir("res/cat")

	if err != nil {
		log.Printf("Error reading res/cat: %v", err)
		interRespond(sess, inter, "Couldn't find any cats :<", nil, false)

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
		interRespond(sess, inter, "Couldn't find any cats", nil, false)
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

	interRespond(sess, inter, result, nil, false)
}

func handleAssault(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	db, ok := getDB(inter.GuildID)

	if !ok {
		return
	}

	sender, target, ok := interGetSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	// note: Options[0] is target
	item, _, ok := interGetItemFromOption(sess, inter, 1)

	if !ok {
		return
	}

	if !isWeapon(item) {
		content := fmt.Sprintf("Item **%v** isn't a weapon!!!", item)
		interRespond(sess, inter, content, nil, false)

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
		interRespond(sess, inter, "Failed to check inventory :<", nil, false)

		return
	}

	if count == 0 {
		interRespond(sess, inter, fmt.Sprintf("You don't have **%v** in your inventory!!!", item), nil, false)
		return
	}

	if !ok {
		return
	}

	chance := 0
	content := ""
	successMessage := "killed them!"

	switch item {
	case itemKnife:
		chance = 20
		content = fmt.Sprintf("%v tried to stab %v with a knife and... ", sender.Mention(), target.Mention())
	case itemGun:
		chance = 70
		content = fmt.Sprintf("%v tried to shoot %v with a gun and... ", sender.Mention(), target.Mention())
	case itemBomb:
		chance = 90
		content = fmt.Sprintf("%v threw a bomb at %v and... ", sender.Mention(), target.Mention())
		successMessage = "BOOM!1!11!! 💥💥💥💥💥"
	case itemNuke:
		chance = 100
		content = fmt.Sprintf("%v nuked %v and... ", sender.Mention(), target.Mention())
		successMessage = "OBLITERATED THEM!!! BOOM!1!11!! 💥💥💥💥💥"
	}

	result := "failed! oops"

	if rand.IntN(99) < chance {
		result = successMessage
	}

	content += result
	interRespond(sess, inter, content, nil, true)
}

func handleWork(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, ok := interGetSender(sess, inter)

	if !ok {
		return
	}

	tx, ok := interBeginTx(sess, inter, cmdWork)

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
		interRespond(sess, inter, "Failed to check work cooldown :<", nil, false)

		return
	}

	const cooldown = 10 * 60
	currentTime := time.Now().Unix()

	if currentTime-lastWork < cooldown {
		remaining := cooldown - (currentTime - lastWork)
		content := fmt.Sprintf("You need to wait **%vm %vs** before working again!!!", remaining/60, remaining%60)

		interRespond(sess, inter, content, nil, false)

		return
	}

	isHigh, _ := getUserHighInfo(tx, sender.ID)
	// random number from range 100-200
	money := rand.IntN(100) + 100

	if isHigh {
		// apply 30% reduction if high
		money = int(float64(money) * 0.7)
	}

	if !interTxMoneyOp(sess, inter, tx, sender.ID, money, opAdd, cmdWork) {
		return
	}

	if !interTxUpdateCd(sess, inter, tx, sender.ID, cdWorkField, currentTime, cmdWork) {
		return
	}

	if !interCommitTx(sess, inter, tx, cmdWork) {
		return
	}

	content := ""

	if isHigh {
		content = fmt.Sprintf("You are **high** %s Actually, it's not good... Money from work has decreased!!!\n",
			emojiCatr)
	}

	content += fmt.Sprintf("You worked and got **%v** money!", money)
	interRespond(sess, inter, content, nil, false)
}

func handleBalance(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	db, ok := getDB(inter.GuildID)

	if !ok {
		return
	}

	target, ok := interGetOptionalTarget(sess, inter)

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
		interRespond(sess, inter, "Failed to check balance :<", nil, false)

		return
	}

	content := fmt.Sprintf("%v's balance: **%v** money", target.Mention(), balance)
	interRespond(sess, inter, content, nil, true)
}

func handleTransfer(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, target, ok := interGetSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	if sender.ID == target.ID {
		interRespond(sess, inter, "You can't transfer money to yourself!!!", nil, false)
		return
	}

	options := inter.ApplicationCommandData().Options
	amount := (int)(options[1].Value.(float64))

	if amount <= 0 {
		interRespond(sess, inter, "Transfer amount must be positive!!!", nil, false)
		return
	}

	tx, ok := interBeginTx(sess, inter, cmdTransfer)

	if !ok {
		return
	}

	defer tx.Rollback()

	balance := getUserBalance(tx, sender.ID)

	if balance < amount {
		interRespond(sess, inter, "You don't have enough money for this transfer!!!", nil, false)
		return
	}

	if !interTxMoneyOp(sess, inter, tx, sender.ID, amount, opSub, cmdTransfer) {
		return
	}

	if !interTxMoneyOp(sess, inter, tx, target.ID, amount, opAdd, cmdTransfer) {
		return
	}

	if !interCommitTx(sess, inter, tx, cmdTransfer) {
		return
	}

	response := fmt.Sprintf("%v transferred %v money to %v", sender.Mention(), amount, target.Mention())
	interRespond(sess, inter, response, nil, true)
}

func handleSteal(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	db, ok := getDB(inter.GuildID)

	if !ok {
		return
	}

	sender, target, ok := interGetSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	if sender.ID == target.ID {
		interRespond(sess, inter, "You can't steal from yourself!!!", nil, false)
		return
	}

	tx, ok := interBeginTx(sess, inter, cmdSteal)

	if !ok {
		return
	}

	defer tx.Rollback()

	targetBalance := getUserBalance(tx, target.ID)

	if targetBalance <= 0 {
		content := fmt.Sprintf("%v is **broke!** Nothing to steal", target.Mention())
		interRespond(sess, inter, content, nil, true)

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
		interRespond(sess, inter, "Failed to check steal cooldown :<", nil, false)

		return
	}

	const cooldown = 20 * 60
	currentTime := time.Now().Unix()

	if currentTime-lastStealFail < cooldown {
		remaining := cooldown - (currentTime - lastStealFail)
		content := fmt.Sprintf("You need to wait **%vm %vs** before stealing again after failure!!!", remaining/60,
			remaining%60)

		interRespond(sess, inter, content, nil, false)

		return
	}

	content := ""
	successChance := 50
	isHigh, _ := getUserHighInfo(tx, sender.ID)

	if isHigh {
		content = fmt.Sprintf("You are **high**, chances of successful steal has increased %s...\n", emojiCatr)
		successChance = 80
	}

	if rand.IntN(100) < successChance {
		targetBalance := getUserBalance(tx, target.ID)

		// max() all stuff so we cant get 0 from stealing lol
		stealPercent := max(1, rand.IntN(51))
		stealAmount := max(1, (stealPercent*targetBalance)/100)

		if !interTxMoneyOp(sess, inter, tx, target.ID, stealAmount, opSub, cmdSteal) {
			return
		}

		if !interTxMoneyOp(sess, inter, tx, sender.ID, stealAmount, opAdd, cmdSteal) {
			return
		}

		content += fmt.Sprintf("You successfully stole **%v** money from %v!", stealAmount, target.Mention())
	} else {
		const penalty = 20

		if !interTxMoneyOp(sess, inter, tx, sender.ID, penalty, opSub, cmdSteal) {
			return
		}

		if !interTxUpdateCd(sess, inter, tx, sender.ID, cdStealField, currentTime, cmdSteal) {
			return
		}

		if isHigh {
			content = fmt.Sprintf("You failed to steal from %v and lost **%v** money! :<\n"+
				"Even being **high** couldn't help you %s",
				target.Mention(), penalty, emojiCatr)
		} else {
			content = fmt.Sprintf("You failed to steal from %v and lost **%v** money! :<\n"+
				"**Pro tip:** being **high** increases chances of a successful steal %s",
				target.Mention(), penalty, emojiCatr)
		}

	}

	if !interCommitTx(sess, inter, tx, cmdSteal) {
		return
	}

	interRespond(sess, inter, content, nil, true)
}

func handleShop(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	builder := strings.Builder{}
	content := fmt.Sprintf("**Shop %s**\n", emojiCykodigo) +
		"-# Use `/buy [item]` to buy something!!!"

	for item, price := range shopItems {
		itemString := fmt.Sprintf("- %v: %v money\n", item, price)
		builder.WriteString(itemString)
	}

	embed := embedText(builder.String())
	interRespondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{embed}, false)
}

func handleBuy(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, ok := interGetSender(sess, inter)

	if !ok {
		return
	}

	item, price, ok := interGetItemFromOption(sess, inter, 0)

	if !ok {
		return
	}

	tx, ok := interBeginTx(sess, inter, cmdBuy)

	if !ok {
		return
	}

	defer tx.Rollback()

	balance := getUserBalance(tx, sender.ID)

	if balance < price {
		content := fmt.Sprintf("Too broke for **%v**, go work!!!", item)
		interRespond(sess, inter, content, nil, false)

		return
	}

	if !interTxMoneyOp(sess, inter, tx, sender.ID, price, opSub, cmdBuy) {
		return
	}

	if _, err := tx.Exec("INSERT INTO inventory(user_id, item) VALUES(?, ?)", sender.ID, item); err != nil {
		log.Printf("Insert error in /inventory: %v", err)
		interRespond(sess, inter, "Failed to add item to inventory :<", nil, false)

		return
	}

	if !interCommitTx(sess, inter, tx, cmdBuy) {
		return
	}

	interRespond(sess, inter, fmt.Sprintf("You bought **%v** for **%v** money", item, price), nil, false)
}

func handleInventory(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	db, ok := getDB(inter.GuildID)

	if !ok {
		return
	}

	target, ok := interGetOptionalTarget(sess, inter)

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
		interRespond(sess, inter, "Failed to check inventory :<", nil, false)

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
		interRespond(sess, inter, content, nil, true)

		return
	}

	builder := strings.Builder{}
	content := fmt.Sprintf("%v inventory:\n", target.Mention())

	for item, count := range items {
		itemString := fmt.Sprintf("- %v ×%v\n", item, count)
		builder.WriteString(itemString)
	}

	embed := embedText(builder.String())
	interRespondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{embed}, false)
}

func handleLeaderboard(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	db, ok := getDB(inter.GuildID)

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
		`, itemDiamond)

	if err != nil {
		log.Printf("Query error in /leaderboard: %v", err)
		interRespond(sess, inter, "Failed to check leaderboard :<", nil, false)

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
		interRespond(sess, inter, content, nil, false)

		return
	}

	content := "**Diamond Leaderboard**\n" +
		"-# Buy some with `/buy diamond`!!!"
	embed := embedText(strings.Join(leaderboard, "\n"))
	interRespondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{embed}, false)
}

func handleEat(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, ok := interGetSender(sess, inter)

	if !ok {
		return
	}

	item, _, ok := interGetItemFromOption(sess, inter, 0)

	if !ok {
		return
	}

	if !isFood(item) {
		content := fmt.Sprintf("You can't eat **%v**!!!", item)
		interRespond(sess, inter, content, nil, false)

		return
	}

	tx, ok := interBeginTx(sess, inter, cmdEat)

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
		interRespond(sess, inter, "Failed to check inventory :<", nil, false)

		return
	}

	if count == 0 {
		interRespond(sess, inter, fmt.Sprintf("You don't have **%v** in your inventory!!!", item), nil, false)
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
		interRespond(sess, inter, "Failed to get inventory :<", nil, false)

		return
	}

	duration := int64(0)

	if item == itemMeth {
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
			interRespond(sess, inter, "Failed to get **high** :<", nil, false)

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
			interRespond(sess, inter, "Failed to get **high** :<", nil, false)

			return
		}

		duration = updatedEndTime - currentTime
	} else {
		_, endTime := getUserHighInfo(tx, sender.ID)
		duration = endTime - time.Now().Unix()
	}

	if !interCommitTx(sess, inter, tx, cmdEat) {
		return
	}

	if item == itemMeth {
		content := fmt.Sprintf("You ate %v! Wowowowowowowowowowo the world is spinning wowowowowowo...\n\n"+
			"You're now high for %vm %vs %s", item, duration/60, duration%60, emojiCatr)
		handleImageCmd(sess, inter, content, "res/gif/spin.gif")
	} else {
		message := "Yummy! " + emojiCykodigo

		switch item {
		case itemKnife:
			message = "Oh wait why did you do that??? You're dead from several internal bleeds."
		case itemBomb:
		case itemNuke:
			message = "BOOM!1!11!!💥💥💥💥💥"
		}

		content := fmt.Sprintf("You ate **%v**! %s", item, message)
		interRespond(sess, inter, content, nil, false)
	}
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
	case cmdHelp:
		handleHelp(sess, inter)
	case cmdMe:
		handleImageCmd(sess, inter, "Me!", "res/me.png")
	case cmdMeow:
		interRespond(sess, inter, "Meow!", nil, false)
	case cmdMeowat:
		handleMeowat(sess, inter)
	case cmdBark:
		handleBark(sess, inter)
	case cmdBarkat:
		handleBarkAt(sess, inter)
	case cmdDoflip:
		handleImageCmd(sess, inter, "Woah!", "res/flip.png")
	case cmdExplode:
		handleImageCmd(sess, inter, "WHAAAAAAAA-", "res/boom.png")
	case cmdSpin:
		// uhh yes im using handleImageCmd for sending a gif
		// and what? what you gonna do?
		handleImageCmd(sess, inter, "Wooooooo", "res/gif/spin.gif")
	case cmdCat:
		handleCat(sess, inter)
	case cmdCart:
		handleImageCmd(sess, inter, "Cart!", "res/cart.png")
	case cmdRoulette:
		handleRoulette(sess, inter)
	case cmdAssault:
		handleAssault(sess, inter)
	case cmdWork:
		handleWork(sess, inter)
	case cmdBalance:
		handleBalance(sess, inter)
	case cmdTransfer:
		handleTransfer(sess, inter)
	case cmdSteal:
		handleSteal(sess, inter)
	case cmdShop:
		handleShop(sess, inter)
	case cmdBuy:
		handleBuy(sess, inter)
	case cmdInventory:
		handleInventory(sess, inter)
	case cmdLeaderboard:
		handleLeaderboard(sess, inter)
	case cmdEat:
		handleEat(sess, inter)
	}
}
