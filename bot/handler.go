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

func handleMeowat(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	contentFunc := func(sender, target *discordgo.User) string {
		return fmt.Sprintf("%v meows at %v!", sender.Mention(), target.Mention())
	}

	handleTargetedCmd(sess, inter, contentFunc)
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
	contentFunc := func(sender, target *discordgo.User) string {
		return fmt.Sprintf("%v barks at %v!", sender.Mention(), target.Mention())
	}

	handleTargetedCmd(sess, inter, contentFunc)
}

func handleCat(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	files, err := os.ReadDir("res/cat")

	if err != nil {
		log.Printf("Error reading res/cat: %v", err)
		respond(sess, inter, "Couldn't find any cats :<", nil, false)

		return
	}

	var pngFiles []string

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
	sender, target, ok := getInterSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	// note: Options[0] is target
	item, _, ok := getItemFromInterOption(sess, inter, 1)

	if !ok {
		return
	}

	if !isWeapon(item) {
		content := fmt.Sprintf("Item **%v** isn't a weapon!!!", item)
		respond(sess, inter, content, nil, false)

		return
	}

	count := 0
	err := DB.QueryRow(
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
	}

	result := "failed! oops"

	if rand.IntN(99) < chance {
		result = "killed them!"
	}

	content += result
	respond(sess, inter, content, nil, true)
}

func handleWork(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, ok := getInterSender(sess, inter)

	if !ok {
		return
	}

	tx, ok := interBeginTx(sess, inter, cmdWork)

	if !ok {
		return
	}

	defer tx.Rollback()

	var lastWork sql.NullInt64
	err := tx.QueryRow(
		`
		SELECT last_work 
		FROM balances 
		WHERE user_id = ?
		`, sender.ID).Scan(&lastWork)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to scan row in /work: %v", err)
		respond(sess, inter, "Failed to check work cooldown :<", nil, false)

		return
	}

	const cooldown = 30 * 60 // 30 minutes in seconds
	currentTime := time.Now().Unix()

	if lastWork.Valid && (currentTime-lastWork.Int64) < cooldown {
		remaining := cooldown - (currentTime - lastWork.Int64)
		content := fmt.Sprintf("You need to wait **%v** minutes before working again!!!", remaining/60)

		respond(sess, inter, content, nil, false)

		return
	}

	// random number from range 100-200
	money := rand.IntN(100) + 100

	_, err = tx.Exec(
		`
		INSERT INTO balances(user_id, balance, last_work)
		VALUES(?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET 
			balance = balance + ?,
			last_work = ?
		`,
		sender.ID, money, currentTime, money, currentTime)

	if err != nil {
		log.Printf("Update error in /work: %v", err)
		respond(sess, inter, "Failed to work :<", nil, false)

		return
	}

	if !interCommitTx(sess, inter, tx, cmdWork) {
		return
	}

	content := fmt.Sprintf("You worked and got **%v** money!1!11!!", money)
	respond(sess, inter, content, nil, false)
}

func handleBalance(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	target, ok := getInterOptionalTarget(sess, inter)

	if !ok {
		return
	}

	balance := 0
	err := DB.QueryRow(
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

	content := fmt.Sprintf("%v's balance: **%v** money!1!11!!", target.Mention(), balance)
	respond(sess, inter, content, nil, true)
}

func handleTransfer(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, target, ok := getInterSenderAndTarget(sess, inter)

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

	tx, ok := interBeginTx(sess, inter, cmdTransfer)

	if !ok {
		return
	}

	defer tx.Rollback()

	balance := getUserBalance(tx, sender.ID)

	if balance < amount {
		respond(sess, inter, "You don't have enough money for this transfer!!!", nil, false)
		return
	}

	if !deductMoney(sess, inter, tx, sender.ID, amount, cmdTransfer) {
		return
	}

	_, err := tx.Exec(
		`
        INSERT INTO balances(user_id, balance) 
        VALUES(?, ?) 
        ON CONFLICT(user_id) DO UPDATE SET 
            balance = balance + ?
        `,
		target.ID, amount, amount)

	if err != nil {
		log.Printf("Insert error in /transfer: %v", err)
		respond(sess, inter, "Failed to add money to recipient :<", nil, false)

		return
	}

	if !interCommitTx(sess, inter, tx, cmdTransfer) {
		return
	}

	response := fmt.Sprintf("%v transferred %v money to %v!", sender.Mention(), amount, target.Mention())
	respond(sess, inter, response, nil, true)
}

func handleSteal(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, target, ok := getInterSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	if sender.ID == target.ID {
		respond(sess, inter, "You can't steal from yourself!!!", nil, false)
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
		respond(sess, inter, content, nil, true)

		return
	}

	var lastStealFail sql.NullInt64
	err := DB.QueryRow(
		`
		SELECT last_steal_fail 
		FROM balances 
		WHERE user_id = ?
		`, sender.ID).Scan(&lastStealFail)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to scan row in /steal: %v", err)
		respond(sess, inter, "Failed to check steal cooldown :<", nil, false)

		return
	}

	const cooldown = 60 * 60 // 1 hour in seconds
	currentTime := time.Now().Unix()

	if lastStealFail.Valid && (currentTime-lastStealFail.Int64) < cooldown {
		remaining := cooldown - (currentTime - lastStealFail.Int64)
		content := fmt.Sprintf("You need to wait **%v** minutes before stealing again after failure!!!",
			remaining/60)

		respond(sess, inter, content, nil, false)

		return
	}

	content := ""

	// 20% success chance
	if rand.IntN(100) < 20 {
		targetBalance := getUserBalance(tx, target.ID)
		stealPercent := rand.IntN(51) // random percentage (0-50%)
		stealAmount := (stealPercent * targetBalance) / 100

		if !deductMoney(sess, inter, tx, target.ID, stealAmount, cmdSteal) {
			return
		}

		_, err = tx.Exec(
			`
			INSERT INTO balances(user_id, balance) VALUES(?, ?) 
			ON CONFLICT(user_id) DO UPDATE SET 
				balance = balance + ?
			`,
			sender.ID, stealAmount, stealAmount,
		)

		if err != nil {
			log.Printf("Insert error in /steal: %v", err)
			respond(sess, inter, "Failed to add money to your account :<", nil, false)

			return
		}

		content = fmt.Sprintf("You successfully stole **%v** money from %v!", stealAmount, target.Mention())
	} else {
		const penalty = 20
		_, err := tx.Exec(
			`
			INSERT INTO balances (user_id, balance)
			VALUES (?, ?)
			ON CONFLICT(user_id) DO UPDATE SET 
				balance = balance - ?,
				last_steal_fail = ?
			`,
			sender.ID, -penalty, penalty, currentTime)

		if err != nil {
			log.Printf("Penalty error in /steal: %v", err)
			respond(sess, inter, "Failed to steal :<", nil, false)

			return
		}

		content = fmt.Sprintf("You failed to steal from %v and lost **%v** money! :<", target.Mention(), penalty)
	}

	if !interCommitTx(sess, inter, tx, cmdSteal) {
		return
	}

	respond(sess, inter, content, nil, true)
}

func handleShop(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	var builder strings.Builder
	content := "**Shop!1!11!!**\n" +
		"-# Use `/buy [item]` to buy something!!!"

	for item, price := range shopItems {
		itemString := fmt.Sprintf("- %v: %v money\n", item, price)
		builder.WriteString(itemString)
	}

	embed := embedContent(builder.String())
	respondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{embed}, false)
}

func handleBuy(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, ok := getInterSender(sess, inter)

	if !ok {
		return
	}

	item, price, ok := getItemFromInterOption(sess, inter, 0)

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
		respond(sess, inter, content, nil, false)

		return
	}

	if !deductMoney(sess, inter, tx, sender.ID, price, cmdBuy) {
		return
	}

	_, err := tx.Exec("INSERT INTO inventory(user_id, item) VALUES(?, ?)", sender.ID, item)

	if err != nil {
		log.Printf("Insert error in /inventory: %v", err)
		respond(sess, inter, "Failed to add item to inventory :<", nil, false)

		return
	}

	if !interCommitTx(sess, inter, tx, cmdBuy) {
		return
	}

	respond(sess, inter, fmt.Sprintf("You bought **%v** for **%v** money!", item, price), nil, false)
}

func handleInventory(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	target, ok := getInterOptionalTarget(sess, inter)

	if !ok {
		return
	}

	rows, err := DB.Query(
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

	var builder strings.Builder
	content := fmt.Sprintf("%v inventory:\n", target.Mention())

	for item, count := range items {
		itemString := fmt.Sprintf("- %v ×%v\n", item, count)
		builder.WriteString(itemString)
	}

	embed := embedContent(builder.String())
	respondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{embed}, false)
}

func handleLeaderboard(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	rows, err := DB.Query(
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
		respond(sess, inter, "Failed to check leaderboard :<", nil, false)

		return
	}

	defer rows.Close()

	leaderboard := []string{}
	position := 1

	for rows.Next() {
		var userID string
		var count int

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

	content := "**Diamond Leaderboard!1!11!!**\n" +
		"-# Buy some with `/buy diamond`!!!"
	embed := embedContent(strings.Join(leaderboard, "\n"))
	respondEmbed(sess, inter, content, nil, []*discordgo.MessageEmbed{embed}, false)
}

func handleEat(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, ok := getInterSender(sess, inter)

	if !ok {
		return
	}

	item, _, ok := getItemFromInterOption(sess, inter, 0)

	if !ok {
		return
	}

	if !isFood(item) {
		content := fmt.Sprintf("You can't eat **%v**!!!", item)
		respond(sess, inter, content, nil, false)

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

	if !interCommitTx(sess, inter, tx, cmdEat) {
		return
	}

	content := fmt.Sprintf("You ate **%v**! Yummy!", item)
	respond(sess, inter, content, nil, false)
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
		respond(sess, inter, "I don't know........ ummmm...... awkwar.......", nil, false)
	case cmdMe:
		handleImageCmd(sess, inter, "Me!", "res/me.png")
	case cmdMeow:
		respond(sess, inter, "Meow!", nil, false)
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
		handleImageCmd(sess, inter, "Wooooooo", "res/spin.gif")
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

// this function doesnt handle some nerdy stuff
// it only handles cases where the message contains the bot name
func handleMsgBotUsername(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	variants := []string{
		"What?",
		"Did someone call me?",
		"Meow?",
		"Hi",
		"I'm a cat!",
		"I'm an orange cat!",
		"Boo",
		"Huh",
		"Huh?",
		"?",
		"???",
		"?????????????????????????????????????",
		"Hey",
		"Hello",
		"Why?",
		"No. Just no.",
		"Well?",
		"So?",
		"I don't know why you said my name",
		":3",
		"Meow!",
		"Wowowowowowowo",
		"OwO",
		"O.o",
		"O.O",
		"Never gonna give you up, never gonna let you down\nNever gonna run around and desert you\n" +
			"Never gonna make you cry, never gonna say goodbye\nNever gonna tell a lie and hurt you",
	}

	sess.ChannelMessageSend(msg.ChannelID, variants[rand.IntN(len(variants))])
}

func handleMsgMeow(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	variants := []string{
		"Meow!",
		"Hi!!! :3",
		":3",
		"I heard a meow!",
		"Meow :3",
		"Meow",
		"jrimbayum",
	}

	sess.ChannelMessageSend(msg.ChannelID, variants[rand.IntN(len(variants))])
}

// handler for messages content
func MsgHandler(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == sess.State.User.ID {
		return
	}

	// convert message content to lowercase so we can understand stuff
	content := strings.ToLower(msg.Content)

	switch {
	// CmdMsgExplodeBalls has 'cykodigo' in it, so check if we didnt get a conflict
	case strings.Contains(content, sess.State.User.Username) && !strings.Contains(content, cmdMsgExplodeBalls):
		handleMsgBotUsername(sess, msg)
	case strings.Contains(content, cmdMsgMeow):
		handleMsgMeow(sess, msg)
	case strings.Contains(content, cmdMsgCrazy):
		sess.ChannelMessageSend(msg.ChannelID,
			"Crazy? I was crazy once, They locked me in a room, a rubber room, a rubber room with rats, "+
				"and rats make me crazy.",
		)
	case strings.Contains(content, cmdMsgExplodeBalls):
		sess.ChannelMessageSend(msg.ChannelID, "BOOM!1!11!! 💥💥💥💥💥")
	case strings.Contains(content, cmdMsgGlamptastic):
		sess.ChannelMessageSend(msg.ChannelID, "glamptastic!")
	}
}
