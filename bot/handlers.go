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

// returns true if youre dead
func roulette() bool {
	bullet := 3
	return rand.IntN(6) == bullet
}

// flip a coin!
func flipACoin() bool {
	return rand.IntN(2) != 1
}

// util function to get random number in range [min, max] cus go
// for some reason doesnt have it builtin
func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

// util function for getting interaction sender cus yes
func getInterSender(inter *discordgo.InteractionCreate) (*discordgo.User, error) {
	sender := inter.User

	if sender == nil && inter.Member != nil {
		sender = inter.Member.User
	}

	if sender == nil {
		return nil, fmt.Errorf("couldn't get interaction sender :<")
	}

	return sender, nil
}

// util function for getting interaction [member] in commands like
// /meowat [member]
func getInterUser(inter *discordgo.InteractionCreate, required bool) (*discordgo.User, error) {
	var targetUser *discordgo.User
	options := inter.ApplicationCommandData().Options

	if len(options) > 0 && options[0].Type == discordgo.ApplicationCommandOptionUser {
		userID := options[0].Value.(string)

		if user, ok := inter.ApplicationCommandData().Resolved.Users[userID]; ok {
			targetUser = user
		}
	}

	if required && targetUser == nil {
		return nil, fmt.Errorf("couldn't find target user :<")
	}

	return targetUser, nil
}

// util function for getting interaction user and sender cus yes
func getInterSenderAndTargetUser(inter *discordgo.InteractionCreate) (*discordgo.User, *discordgo.User, error) {
	targetUser, err := getInterUser(inter, true)

	if err != nil {
		return nil, nil, err
	}

	sender, err := getInterSender(inter)

	if err != nil {
		return nil, nil, err
	}

	return sender, targetUser, nil
}

// util function for getting user balances in sql transactions
// yes /balance doesnt use it
// cus /balance doesnt need sql transactions since its just one query
func getUserBalance(tx *sql.Tx, userID string) int {
	balance := 0
	err := tx.QueryRow("SELECT balance FROM balances WHERE user_id = ?", userID).Scan(&balance)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Balance query error: %v", err)
	}

	return balance
}

// util function to send interaction responses
func respond(sess *discordgo.Session, inter *discordgo.InteractionCreate, content string, files []*discordgo.File,
	allowMentions bool) {
	data := &discordgo.InteractionResponseData{
		Content: content,
		Files:   files,
	}

	if allowMentions {
		data.AllowedMentions = &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{
				discordgo.AllowedMentionTypeUsers},
		}
	}

	sess.InteractionRespond(inter.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

// another util function for commands like
// /meowat [member]
func handleTargetedCmd(sess *discordgo.Session, inter *discordgo.InteractionCreate,
	contentFunc func(sender *discordgo.User, target *discordgo.User) string) {
	sender, targetUser, err := getInterSenderAndTargetUser(inter)

	if err != nil {
		respond(sess, inter, err.Error(), nil, false)
		return
	}

	content := contentFunc(sender, targetUser)
	respond(sess, inter, content, nil, true)
}

// util function for handling commands that send image like
// /me
func handleImageCmd(sess *discordgo.Session, inter *discordgo.InteractionCreate, content string, imgName string,
	imgPath string) {
	file, err := os.Open(imgPath)

	if err != nil {
		log.Printf("Error opening '%s': %v", imgName, err)
		respond(sess, inter, "Couldn't open image :<", nil, false)

		return
	}

	defer file.Close()

	respond(sess, inter, content, []*discordgo.File{{
		Name:   imgName,
		Reader: file,
	}}, false)
}

func handleMeowat(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	contentFunc := func(sender, target *discordgo.User) string {
		return fmt.Sprintf("%s meows at %s!", sender.Mention(), target.Mention())
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
		return fmt.Sprintf("%s barks at %s!", sender.Mention(), target.Mention())
	}

	handleTargetedCmd(sess, inter, contentFunc)
}

func handleRoulette(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	result := "Victory!!! You're alive!!!"

	if roulette() {
		result = "Sorry, you're dead, better luck next ti- uhh.."
	}

	respond(sess, inter, result, nil, false)
}

func handleAssault(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	contentFunc := func(sender, target *discordgo.User) string {
		result := "killed them!"

		if flipACoin() {
			result = "failed! oops"
		}

		return fmt.Sprintf("%s tried to assault %s and... %s", sender.Mention(), target.Mention(), result)
	}

	handleTargetedCmd(sess, inter, contentFunc)
}

// !!! its a joke command !!!
func handleSexnkill(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	contentFunc := func(sender, target *discordgo.User) string {
		mpreg := "made them pregnant"

		if flipACoin() {
			mpreg = "failed to make them pregnant"
		}

		return fmt.Sprintf("%s had sex with %s, %s and killed them!", sender.Mention(), target.Mention(), mpreg)
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
	filePath := fmt.Sprintf("res/cat/%s", img)
	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("Error opening %s: %v", img, err)
		respond(sess, inter, "Couldn't open cat picture :<", nil, false)

		return
	}

	defer file.Close()

	respond(sess, inter, "Cat!", []*discordgo.File{{
		Name:   img,
		Reader: file,
	}}, false)
}

func handleWork(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, err := getInterSender(inter)

	if err != nil {
		respond(sess, inter, err.Error(), nil, false)
		return
	}

	tx, err := DB.Begin()

	if err != nil {
		log.Printf("Failed to begin transaction in /work: %v", err)
		respond(sess, inter, "Failed to work :<", nil, false)

		return
	}

	defer tx.Rollback()

	var lastWork sql.NullInt64
	err = tx.QueryRow("SELECT last_work FROM balances WHERE user_id = ?", sender.ID).Scan(&lastWork)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Cooldown check error in /work: %v", err)
		respond(sess, inter, "Failed to check work cooldown :<", nil, false)

		return
	}

	const cooldown = 30 * 60 // 30 minutes in seconds
	currentTime := time.Now().Unix()

	if lastWork.Valid && (currentTime-lastWork.Int64) < cooldown {
		remaining := cooldown - (currentTime - lastWork.Int64)
		content := fmt.Sprintf("You need to wait %d minutes before working again!!!", remaining/60)

		respond(sess, inter, content, nil, false)

		return
	}

	money := randRange(100, 200)

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

	if err := tx.Commit(); err != nil {
		log.Printf("Commit error in /work: %v", err)
		respond(sess, inter, "Failed to finalize work :<", nil, false)

		return
	}

	content := fmt.Sprintf("You worked and got %d money!1!11!!", money)
	respond(sess, inter, content, nil, false)
}

func handleBalance(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	// check for [member] option
	targetUser, err := getInterUser(inter, false)

	if err != nil {
		respond(sess, inter, err.Error(), nil, false)
		return
	}

	// if theres no [member] specified just use sender instead
	if targetUser == nil {
		targetUser, err = getInterSender(inter)

		if err != nil {
			respond(sess, inter, err.Error(), nil, false)
			return
		}
	}

	balance := 0
	err = DB.QueryRow("SELECT balance FROM balances WHERE user_id = ?", targetUser.ID).Scan(&balance)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Database error in /balance: %v", err)
		respond(sess, inter, "Failed to check balance :<", nil, false)

		return
	}

	content := fmt.Sprintf("%s's balance: %d money!1!11!!", targetUser.Mention(), balance)
	respond(sess, inter, content, nil, true)
}

func handleTransfer(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, targetUser, err := getInterSenderAndTargetUser(inter)

	if err != nil {
		respond(sess, inter, err.Error(), nil, false)
		return
	}

	if sender.ID == targetUser.ID {
		respond(sess, inter, "You can't transfer money to yourself!!!", nil, false)
		return
	}

	options := inter.ApplicationCommandData().Options
	amount := 0

	for _, option := range options {
		if option.Name == "amount" {
			amount = (int)(option.Value.(float64))
			break
		}
	}

	if amount <= 0 {
		respond(sess, inter, "Transfer amount must be positive!!!", nil, false)
		return
	}

	tx, err := DB.Begin()

	if err != nil {
		log.Printf("Failed to begin transaction in /transfer: %v", err)
		respond(sess, inter, "Failed to start transfer :<", nil, false)

		return
	}

	defer tx.Rollback()

	res, err := tx.Exec(
		`
        UPDATE balances 
        SET balance = balance - ? 
        WHERE user_id = ? AND balance >= ?
        `,
		amount, sender.ID, amount)

	if err != nil {
		log.Printf("Deduction error in /transfer: %v", err)
		respond(sess, inter, "Failed to deduct from your account :<", nil, false)

		return
	}

	// i think it will return error if database driver doesnt support RowsAffected
	// so i dont do check
	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		respond(sess, inter, "You don't have enough money for this transfer!!!", nil, false)
		return
	}

	_, err = tx.Exec(
		`
        INSERT INTO balances(user_id, balance) 
        VALUES(?, ?) 
        ON CONFLICT(user_id) DO UPDATE SET 
            balance = balance + ?
        `,
		targetUser.ID, amount, amount)

	if err != nil {
		log.Printf("Add balance error in /transfer: %v", err)
		respond(sess, inter, "Failed to add money to recipient :<", nil, false)

		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Commit error in /transfer: %v", err)
		respond(sess, inter, "Failed to finalize transfer :<", nil, false)

		return
	}

	response := fmt.Sprintf("%s transferred %d money to %s!", sender.Mention(), amount, targetUser.Mention())
	respond(sess, inter, response, nil, true)
}

func handleSteal(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	sender, target, err := getInterSenderAndTargetUser(inter)

	if err != nil {
		respond(sess, inter, err.Error(), nil, false)
		return
	}

	if sender.ID == target.ID {
		respond(sess, inter, "You can't steal from yourself!!!", nil, false)
		return
	}

	tx, err := DB.Begin()

	if err != nil {
		log.Printf("Failed to start transaction in /steal: %v", err)
		respond(sess, inter, "Failed to start stealing :<", nil, false)

		return
	}

	defer tx.Rollback()

	targetBalance := getUserBalance(tx, target.ID)

	if targetBalance <= 0 {
		content := fmt.Sprintf("%s is broke! Nothing to steal.", target.Mention())
		respond(sess, inter, content, nil, true)

		return
	}

	var lastStealFail sql.NullInt64
	err = DB.QueryRow("SELECT last_steal_fail FROM balances WHERE user_id = ?", sender.ID).Scan(&lastStealFail)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Cooldown check error in /steal: %v", err)
		respond(sess, inter, "Failed to check steal cooldown :<", nil, false)

		return
	}

	const cooldown = 60 * 60 // 1 hour in seconds
	currentTime := time.Now().Unix()

	if lastStealFail.Valid && (currentTime-lastStealFail.Int64) < cooldown {
		remaining := cooldown - (currentTime - lastStealFail.Int64)
		content := fmt.Sprintf("You need to wait %d minutes before stealing again after failure!", remaining/60)

		respond(sess, inter, content, nil, false)

		return
	}

	content := ""

	// 20% success chance
	if rand.IntN(100) < 20 {
		targetBalance := getUserBalance(tx, target.ID)
		stealPercent := rand.IntN(51) // random percentage (0-50%)
		stealAmount := (stealPercent * targetBalance) / 100

		_, err := tx.Exec(
			`
			UPDATE balances SET balance = balance - ? WHERE user_id = ? AND balance >= ?
			`,
			stealAmount, target.ID, stealAmount,
		)

		if err != nil {
			log.Printf("Deduction error in /steal: %v", err)
			respond(sess, inter, "Failed to steal :<", nil, false)

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
			log.Printf("Add balance error: %v", err)
			respond(sess, inter, "Steal failed :<", nil, false)

			return
		}

		content = fmt.Sprintf("You successfully stole %d money from %s!", stealAmount, target.Mention())
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

		content = fmt.Sprintf("You failed to steal from %s and lost %d money! :<", target.Mention(), penalty)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Commit error in /steal: %v", err)
		respond(sess, inter, "Failed to finalize steal :<", nil, false)

		return
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
	case CmdHelp:
		respond(sess, inter, "I don't know........ ummmm...... awkwar.......", nil, false)
	case CmdMeow:
		respond(sess, inter, "Meow!", nil, false)
	case CmdMeowat:
		handleMeowat(sess, inter)
	case CmdBark:
		handleBark(sess, inter)
	case CmdBarkat:
		handleBarkAt(sess, inter)
	case CmdRoulette:
		handleRoulette(sess, inter)
	case CmdMe:
		handleImageCmd(sess, inter, "", "me.png", "res/me.png")
	case CmdAssault:
		handleAssault(sess, inter)
	case CmdSexnkill:
		handleSexnkill(sess, inter)
	case CmdCat:
		handleCat(sess, inter)
	case CmdCart:
		handleImageCmd(sess, inter, "Cart!", "cart.png", "res/cart.png")
	case CmdDoflip:
		handleImageCmd(sess, inter, "Woah!", "flip.png", "res/flip.png")
	case CmdWork:
		handleWork(sess, inter)
	case CmdBalance:
		handleBalance(sess, inter)
	case CmdTransfer:
		handleTransfer(sess, inter)
	case CmdSteal:
		handleSteal(sess, inter)
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
	case strings.Contains(content, sess.State.User.Username):
		handleMsgBotUsername(sess, msg)
	case strings.Contains(content, CmdMsgMeow):
		handleMsgMeow(sess, msg)
	case strings.Contains(content, CmdMsgCrazy):
		sess.ChannelMessageSend(msg.ChannelID,
			"Crazy? I was crazy once, They locked me in a room, a rubber room, a rubber room with rats, "+
				"and rats make me crazy.",
		)
	}
}
