package bot

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"os"

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

// quick note for myself cus i will forget 100%:
// sender - person who used command
// target - person who was specified in [member] option

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
func getInterTarget(inter *discordgo.InteractionCreate, required bool) (*discordgo.User, error) {
	var target *discordgo.User
	options := inter.ApplicationCommandData().Options

	if len(options) == 0 && !required {
		return nil, nil
	} // else clause is unreachable

	userID := options[0].Value.(string)

	if user, ok := inter.ApplicationCommandData().Resolved.Users[userID]; ok {
		target = user
	}

	// i genuinely dont know how this can fail but anyways heres check
	// if we somehow didnt find user
	if required && target == nil {
		return nil, fmt.Errorf("couldn't find target user :<")
	}

	return target, nil
}

// util function for getting interaction user and sender cus yes
func getInterSenderAndTarget(inter *discordgo.InteractionCreate) (*discordgo.User, *discordgo.User, error) {
	target, err := getInterTarget(inter, true)

	// this check is for golang to shut up because "erm actually value of err never used"
	// i hate it
	if err != nil {
		return nil, nil, err
	}

	sender, err := getInterSender(inter)
	return sender, target, err
}

// util function for getting target user from commands with optional [member] option
// like /balance [member] <-- optional
func getInterOptionalTarget(inter *discordgo.InteractionCreate) (*discordgo.User, error) {
	// check for [member] option
	target, err := getInterTarget(inter, false)

	// if theres no [member] specified just use sender instead
	if target == nil {
		target, err = getInterSender(inter)
	}

	return target, err
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
	sender, target, err := getInterSenderAndTarget(inter)

	if err != nil {
		respond(sess, inter, err.Error(), nil, false)
		return
	}

	content := contentFunc(sender, target)
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

// util function for money deduction in transactions
//
// turns out its very common operation ¯\_(ツ)_/¯
//
// note: you should check for balance being less than amount
func deductMoney(tx *sql.Tx, userID string, amount int) (sql.Result, error) {
	res, err := tx.Exec("UPDATE balances SET balance = balance - ? WHERE user_id = ?", amount, userID, amount)
	return res, err
}
