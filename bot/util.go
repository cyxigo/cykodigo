package bot

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

// quick note for myself cus i will forget 100%:
// sender - person who used command
// target - person who was specified in [member] option
//
// also, i dont why but this file is really commented :p

const (
	defaultEmbedColor = 0xFF7B00
)

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
				discordgo.AllowedMentionTypeUsers,
			},
		}
	}

	sess.InteractionRespond(inter.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

// util function to send interaction responses with embeds
func respondEmbed(sess *discordgo.Session, inter *discordgo.InteractionCreate, content string,
	files []*discordgo.File, embeds []*discordgo.MessageEmbed, allowMentions bool) {
	data := &discordgo.InteractionResponseData{
		Content: content,
		Files:   files,
		Embeds:  embeds,
	}

	if allowMentions {
		data.AllowedMentions = &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{
				discordgo.AllowedMentionTypeUsers,
			},
		}
	}

	sess.InteractionRespond(inter.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

// util function for creating embeds
func embedContent(content string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Description: content,
		Color:       defaultEmbedColor,
	}

	return embed
}

// util function for getting interaction sender cus yes
func getInterSender(sess *discordgo.Session, inter *discordgo.InteractionCreate) (*discordgo.User, bool) {
	sender := inter.User

	if sender == nil && inter.Member != nil {
		sender = inter.Member.User
	}

	// how
	if sender == nil {
		respond(sess, inter, "Couldn't get interaction sender :<", nil, false)
		return nil, false
	}

	return sender, true
}

// util function for getting interaction [member] in commands like
// /meowat [member]
func getInterTarget(sess *discordgo.Session, inter *discordgo.InteractionCreate, required bool) (
	*discordgo.User, bool) {
	var target *discordgo.User
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
		respond(sess, inter, "Couldn't find target user :<", nil, false)
		return nil, false
	}

	return target, true
}

// util function for getting interaction user and sender cus yes
func getInterSenderAndTarget(sess *discordgo.Session, inter *discordgo.InteractionCreate) (
	*discordgo.User, *discordgo.User, bool) {
	target, ok := getInterTarget(sess, inter, true)

	if !ok {
		return nil, nil, false
	}

	sender, ok := getInterSender(sess, inter)
	return sender, target, ok
}

// util function for getting target user from commands with optional [member] option
// like /balance [member] <-- optional
func getInterOptionalTarget(sess *discordgo.Session, inter *discordgo.InteractionCreate) (
	*discordgo.User, bool) {
	// check for [member] option
	target, ok := getInterTarget(sess, inter, false)

	// if theres no [member] specified just use sender instead
	if target == nil {
		target, ok = getInterSender(sess, inter)
	}

	return target, ok
}

// another util function for commands like
// /meowat [member]
func handleTargetedCmd(sess *discordgo.Session, inter *discordgo.InteractionCreate,
	contentFunc func(sender *discordgo.User, target *discordgo.User) string) {
	sender, target, ok := getInterSenderAndTarget(sess, inter)

	if !ok {
		return
	}

	content := contentFunc(sender, target)
	respond(sess, inter, content, nil, true)
}

// util function for handling commands that send image like
// /me
// can be also used for sending gifs
func handleImageCmd(sess *discordgo.Session, inter *discordgo.InteractionCreate, content string, imgPath string) {
	file, err := os.Open(imgPath)

	if err != nil {
		log.Printf("Error opening '%v': %v", imgPath, err)
		respond(sess, inter, "Couldn't open image :<", nil, false)

		return
	}

	defer file.Close()

	imgName := filepath.Base(imgPath)
	discordFile := &discordgo.File{
		Name:   imgName,
		Reader: file,
	}
	description := fmt.Sprintf("**%v**", content)
	embed := &discordgo.MessageEmbed{
		Description: description,
		Color:       defaultEmbedColor,
		Image: &discordgo.MessageEmbedImage{
			URL: "attachment://" + imgName,
		},
	}

	respondEmbed(sess, inter, "", []*discordgo.File{discordFile}, []*discordgo.MessageEmbed{embed}, false)
}

// util function for getting user balances in sql transactions
// yes /balance doesnt use it
// cus /balance doesnt need sql transactions since its just one query
func getUserBalance(tx *sql.Tx, userID string) int {
	balance := 0
	err := tx.QueryRow(
		`
		SELECT balance
		FROM balances 
		WHERE user_id = ?
		`, userID).Scan(&balance)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Query error in getUserBalance: %v", err)
	}

	return balance
}

// util function for money deduction in interactions sql transactions
//
// turns out its very common operation ¯\_(ツ)_/¯
//
// note: you should check for balance being less than amount
func deductMoney(sess *discordgo.Session, inter *discordgo.InteractionCreate, tx *sql.Tx, userID string,
	amount int, cmd string) bool {
	_, err := tx.Exec("UPDATE balances SET balance = balance - ? WHERE user_id = ?", amount, userID, amount)

	if err != nil {
		log.Printf("Deduction error in /%v: %v", cmd, err)
		respond(sess, inter, "Failed to deduct money :<", nil, false)

		return false
	}

	return true
}

// util function for beginning sql transactions in interactions
//
// note: dont forget to "defer tx.Rollback()"
func interBeginTx(sess *discordgo.Session, inter *discordgo.InteractionCreate, cmd string) (*sql.Tx, bool) {
	tx, err := DB.Begin()

	if err != nil {
		log.Printf("Failed to begin transaction in /%v: %v", cmd, err)

		content := fmt.Sprintf("Failed to start /%v :<", cmd)
		respond(sess, inter, content, nil, false)

		return nil, false
	}

	return tx, true
}

// util function for commiting sql transactions in interactions
func interCommitTx(sess *discordgo.Session, inter *discordgo.InteractionCreate, tx *sql.Tx, cmd string) bool {
	if err := tx.Commit(); err != nil {
		log.Printf("Commit error in /%v: %v", cmd, err)

		content := fmt.Sprintf("Failed to finish /%v :<", cmd)
		respond(sess, inter, content, nil, false)

		return false
	}

	return true
}
