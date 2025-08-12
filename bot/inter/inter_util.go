package inter

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cyxigo/cykodigo/bot/data"
	"github.com/cyxigo/cykodigo/bot/database"
)

// quick note for myself cus i will forget 100%:
// sender - person who used command
// target - person who was specified in [user] option
//
// also, i dont why but this file is really commented :p

// util function to send interaction responses
func respond(sess *discordgo.Session, inter *discordgo.InteractionCreate, content string, files []*discordgo.File) {
	data := &discordgo.InteractionResponseData{
		Content: content,
		Files:   files,
	}

	sess.InteractionRespond(inter.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

// util function to send interaction responses with embeds
func respondEmbed(sess *discordgo.Session, inter *discordgo.InteractionCreate, content string,
	files []*discordgo.File, embeds []*discordgo.MessageEmbed) {
	data := &discordgo.InteractionResponseData{
		Content: content,
		Files:   files,
		Embeds:  embeds,
	}

	sess.InteractionRespond(inter.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

// util function for getting interaction sender cus yes
func getSender(sess *discordgo.Session, inter *discordgo.InteractionCreate) (*discordgo.User, bool) {
	sender := inter.User

	if sender == nil && inter.Member != nil {
		sender = inter.Member.User
	}

	// how
	if sender == nil {
		respond(sess, inter, "Couldn't get interaction sender", nil)
		return nil, false
	}

	return sender, true
}

// util function for getting interaction [user] in commands like
// /meowat [user]
func getTarget(sess *discordgo.Session, inter *discordgo.InteractionCreate, required bool) (
	*discordgo.User, bool) {
	target := &discordgo.User{}
	options := inter.ApplicationCommandData().Options

	if len(options) == 0 && !required {
		return nil, true // its not required so not having user option is ok
	} // else clause is unreachable

	userID := options[0].Value.(string)

	if user, ok := inter.ApplicationCommandData().Resolved.Users[userID]; ok {
		target = user
	}

	// i genuinely dont know how this can fail but anyways heres check
	// if we somehow didnt find user
	if required && target == nil {
		respond(sess, inter, "Couldn't find target user", nil)
		return nil, false
	}

	return target, true
}

// util function for getting interaction user and sender cus yes
func getSenderAndTarget(sess *discordgo.Session, inter *discordgo.InteractionCreate) (
	*discordgo.User, *discordgo.User, bool) {
	target, ok := getTarget(sess, inter, true)

	if !ok {
		return nil, nil, false
	}

	sender, ok := getSender(sess, inter)
	return sender, target, ok
}

// util function for getting target user from commands with optional [user] option
// like /balance [user] <-- optional
func getOptionalTarget(sess *discordgo.Session, inter *discordgo.InteractionCreate) (
	*discordgo.User, bool) {
	// check for [user] option
	target, ok := getTarget(sess, inter, false)

	// if theres no [user] specified just use sender instead
	if target == nil {
		target, ok = getSender(sess, inter)
	}

	return target, ok
}

// util function for getting an item from command option
// returns item and its price
func getItemFromOption(sess *discordgo.Session, inter *discordgo.InteractionCreate, idx int) (
	string, int64, bool) {
	item := strings.ToLower(inter.ApplicationCommandData().Options[idx].StringValue())
	price, exists := data.ShopItems[item]

	if !exists {
		content := fmt.Sprintf("There's no item **%v**", item)
		respond(sess, inter, content, nil)

		return "", 0, false
	}

	return item, int64(price), true
}

// util function for getting optional item amount
// will return 1 if amount isnt specified
func getItemAmountOption(sess *discordgo.Session, inter *discordgo.InteractionCreate, action string, idx int) (
	int64, bool) {
	options := inter.ApplicationCommandData().Options

	if len(options) <= idx {
		return 1, true
	}

	value := options[idx].IntValue()

	if value < 1 {
		content := fmt.Sprintf("You can't %v less than **1** item", action)
		respond(sess, inter, content, nil)

		return 0, false
	}

	if value > 1000 {
		content := fmt.Sprintf("You can't %v more than **1000** items at a time", action)
		respond(sess, inter, content, nil)

		return 0, false
	}

	return value, true
}

// another util function for commands like
// /meowat [user]
func handleActionOnCmd(sess *discordgo.Session, inter *discordgo.InteractionCreate, what string) {
	sender, target, ok := getSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	content := fmt.Sprintf("%v **%v** %v %v", sender.Mention(), what, target.Mention(), data.EmojiCykodigo)
	respond(sess, inter, content, nil)
}

// util function for handling commands in interactions that send image
// can be also used for sending gifs
func handleImageCmd(sess *discordgo.Session, inter *discordgo.InteractionCreate, content string, imgPath string) {
	file, err := os.Open(imgPath)

	if err != nil {
		log.Printf("Error opening '%v': %v", imgPath, err)
		respond(sess, inter, "Couldn't open image", nil)

		return
	}

	defer file.Close()

	imgName := filepath.Base(imgPath)
	discordFile := &discordgo.File{
		Name:   imgName,
		Reader: file,
	}
	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("**%v** %v", content, data.EmojiCykodigo),
		Color:       data.DefaultEmbedColor,
		Image: &discordgo.MessageEmbedImage{
			URL: "attachment://" + imgName,
		},
	}

	respondEmbed(sess, inter, "", []*discordgo.File{discordFile}, []*discordgo.MessageEmbed{embed})
}

// util function for getting all server members
//
// defined to increase readability in /balanceall handler
func getAllMembers(sess *discordgo.Session, inter *discordgo.InteractionCreate) (
	[]*discordgo.Member, bool) {
	guildID := inter.GuildID
	allMembers := []*discordgo.Member{}
	after := ""

	for {
		members, err := sess.GuildMembers(guildID, after, 1000)

		if err != nil {
			log.Printf("Failed to get %v member list: %v", guildID, err)
			respond(sess, inter, "Failed to get member list", nil)

			return nil, false
		}

		if len(members) == 0 {
			break
		}

		allMembers = append(allMembers, members...)
		after = members[len(members)-1].User.ID
	}

	return allMembers, true
}

const (
	workCooldown     = 10 * 60
	stealCooldown    = 15 * 60
	rouletteCooldown = 5 * 60
)

// util function for checking cooldowns in sql transactions
func txCheckCd(sess *discordgo.Session, inter *discordgo.InteractionCreate, tx *sql.Tx, userID, field string,
	cooldown int64, cmd string) bool {
	lastTime := int64(0)
	err := tx.QueryRow(
		`
		SELECT `+field+` 
		FROM cooldowns 
		WHERE user_id = ?
		`, userID).Scan(&lastTime)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Cooldown check error in /%s: %v", cmd, err)
		respond(sess, inter, "Failed to check cooldown", nil)

		return false
	}

	currentTime := time.Now().Unix()

	if currentTime-lastTime < cooldown {
		remaining := cooldown - (currentTime - lastTime)
		content := fmt.Sprintf("You need to wait **%vm %vs**", remaining/60, remaining%60)

		respond(sess, inter, content, nil)

		return false
	}

	return true
}

// util function for checking if you have enough items in your inventory in sql transactions
// "fatal" parameter means that this function should respond to an interaction with error if not enough items
func txCheckInventory(sess *discordgo.Session, inter *discordgo.InteractionCreate, tx *sql.Tx, userID, item string,
	minAmount int64, fatal bool, cmd string) bool {
	count := int64(0)
	err := tx.QueryRow(
		`
		SELECT amount 
		FROM inventory
		WHERE user_id = ? AND item = ?
		`,
		userID, item).Scan(&count)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Inventory check error in /%v: %v", cmd, err)
		respond(sess, inter, "Failed to check inventory", nil)

		return false
	}

	if count < minAmount {
		if fatal {
			content := fmt.Sprintf("You don't have enough **%s**", item)
			respond(sess, inter, content, nil)
		}

		return false
	}

	return true
}

// util function for checking if you have enough items in your inventory
//
// "fatal" parameter means that this function should respond to an interaction with error if not enough items
func checkInventory(sess *discordgo.Session, inter *discordgo.InteractionCreate, db *sql.DB, userID, item string,
	minAmount int64, fatal bool, cmd string) bool {
	count := int64(0)
	err := db.QueryRow(
		`
		SELECT amount 
		FROM inventory
		WHERE user_id = ? AND item = ?
		`,
		userID, item).Scan(&count)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Inventory check error in /%v: %v", cmd, err)
		respond(sess, inter, "Failed to check inventory", nil)

		return false
	}

	if count < minAmount {
		if fatal {
			content := fmt.Sprintf("You don't have enough **%s**", item)
			respond(sess, inter, content, nil)
		}

		return false
	}

	return true
}

// util function for adding item to the user inventory in sql transactions
//
// remove items via negative amount
//
// deletes the item row if the number of items is zero
func txAddItem(sess *discordgo.Session, inter *discordgo.InteractionCreate, tx *sql.Tx, userID, item string,
	amount int64, cmd string) bool {
	_, err := tx.Exec(
		`
		INSERT INTO inventory (user_id, item, amount) 
		VALUES (?, ?, ?) 
		ON CONFLICT(item) 
		DO UPDATE SET amount = amount + ?
		`,
		userID, item, amount, amount)

	if err != nil {
		log.Printf("Insert error in /%v: %v", cmd, err)
		respond(sess, inter, "Failed to update inventory", nil)

		return false
	}

	_, err = tx.Exec(
		`
		DELETE FROM inventory 
		WHERE user_id = ? AND item = ? AND amount <= 0
		`,
		userID, item)

	if err != nil {
		log.Printf("Delete error in /%v: %v", cmd, err)
		respond(sess, inter, "Failed to update inventory", nil)

		return false
	}

	return true
}

// util function for beginning sql transactions in interactions
//
// note: dont forget to "defer tx.Rollback()"
func beginTx(sess *discordgo.Session, inter *discordgo.InteractionCreate, cmd string) (*sql.Tx, bool) {
	db, ok := database.GetDB(inter.GuildID)

	if !ok {
		return nil, false
	}

	tx, err := db.Begin()

	if err != nil {
		log.Printf("Failed to begin transaction in /%v: %v", cmd, err)

		content := fmt.Sprintf("Failed to start /%v", cmd)
		respond(sess, inter, content, nil)

		return nil, false
	}

	return tx, true
}

// util function for commiting sql transactions in interactions
func commitTx(sess *discordgo.Session, inter *discordgo.InteractionCreate, tx *sql.Tx, cmd string) bool {
	if err := tx.Commit(); err != nil {
		log.Printf("Commit error in /%v: %v", cmd, err)

		content := fmt.Sprintf("Failed to finish /%v", cmd)
		respond(sess, inter, content, nil)

		return false
	}

	return true
}

// util function for updating cooldown in interactions sql transactions
//
// e.g work cooldown
func txUpdateCd(sess *discordgo.Session, inter *discordgo.InteractionCreate, tx *sql.Tx, userID string,
	field string, cmd string) bool {
	value := time.Now().Unix()
	_, err := tx.Exec(
		`
		INSERT INTO cooldowns(user_id, `+field+`) 
		VALUES(?, ?) 
		ON CONFLICT(user_id) 
		DO UPDATE SET `+field+` = ?
		`,
		userID, value, value)

	if err != nil {
		log.Printf("Cooldown update error in /%v: %v", cmd, err)
		respond(sess, inter, "Failed to update cooldown", nil)

		return false
	}

	return true
}

// util function for money addition in interactions sql transactions
//
// remove money via negative amount
func txAddMoney(sess *discordgo.Session, inter *discordgo.InteractionCreate, tx *sql.Tx, userID string, amount int64,
	cmd string) bool {
	_, err := tx.Exec(
		`
		INSERT INTO balances(user_id, balance)
		VALUES(?, ?)
		ON CONFLICT(user_id) 
		DO UPDATE SET balance = balance + ?
		`, userID, amount, amount)

	if err != nil {
		log.Printf("Addition error in /%v: %v", cmd, err)
		respond(sess, inter, "Failed to add money", nil)

		return false
	}

	return true
}

// util function for updating high effect in sql transactions
//
// defined to increase readability in /eat handler
//
// returns updated end time on success
func txUpdateHighEffect(sess *discordgo.Session, inter *discordgo.InteractionCreate, tx *sql.Tx, userID string,
	duration int64) (int64, bool) {
	currentTime := time.Now().Unix()
	newEndTime := currentTime + duration

	_, err := tx.Exec(
		`
			INSERT INTO effects(user_id, high_end_time)
			VALUES(?, ?)
			ON CONFLICT(user_id) 
			DO UPDATE SET high_end_time = MAX(?, high_end_time) + ?
			`,
		userID, newEndTime, currentTime, duration)

	if err != nil {
		log.Printf("Failed to update meth effect in /eat: %v", err)
		respond(sess, inter, "Failed to get **high**", nil)

		return 0, false
	}

	updatedEndTime := int64(0)
	err = tx.QueryRow(
		`
			SELECT high_end_time 
			FROM effects 
			WHERE user_id = ?
			`,
		userID).Scan(&updatedEndTime)

	if err != nil {
		log.Printf("Failed to get updated end time in /eat: %v", err)
		respond(sess, inter, "Failed to get **high**", nil)

		return 0, false
	}

	return updatedEndTime, true
}
