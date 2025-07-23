package bot

import (
	"bytes"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// returns true if you're dead :<
func roulette() bool {
	bullet := 3
	return rand.IntN(6) == bullet
}

// flip a coin!
func flipACoin() bool {
	return rand.IntN(2) != 1
}

// util function for commands like
// meowat [member]
// returns sender and a [member]
func getUserAndSender(inter *discordgo.InteractionCreate) (
	*discordgo.User, *discordgo.User, error,
) {
	var targetUser *discordgo.User
	options := inter.ApplicationCommandData().Options

	if len(options) > 0 &&
		options[0].Type == discordgo.ApplicationCommandOptionUser {
		userID := options[0].Value.(string)

		if user,
			ok := inter.ApplicationCommandData().Resolved.Users[userID]; ok {
			targetUser = user
		}
	}

	if targetUser == nil {
		return nil, nil, fmt.Errorf("couldn't find user :<")
	}

	sender := inter.User

	if sender == nil && inter.Member != nil {
		sender = inter.Member.User
	}

	return targetUser, sender, nil
}

// util function to send interaction responses
func respond(
	sess *discordgo.Session,
	inter *discordgo.InteractionCreate,
	content string,
	files []*discordgo.File,
	allowMentions bool,
) {
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
// meowat [member]
func handleTargetedCmd(
	sess *discordgo.Session,
	inter *discordgo.InteractionCreate,
	contentFunc func(sender *discordgo.User, target *discordgo.User) string,
) {
	targetUser, sender, err := getUserAndSender(inter)

	if err != nil {
		respond(sess, inter, err.Error(), nil, false)
		return
	}

	content := contentFunc(sender, targetUser)
	respond(sess, inter, content, nil, true)
}

func handleMeowat(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	contentFunc := func(sender, target *discordgo.User) string {
		return fmt.Sprintf("%s meows at %s!", sender.Mention(),
			target.Mention())
	}

	handleTargetedCmd(sess, inter, contentFunc)
}

func handleAssault(
	sess *discordgo.Session,
	inter *discordgo.InteractionCreate,
) {
	contentFunc := func(sender, target *discordgo.User) string {
		result := "killed them!"

		if flipACoin() {
			result = "failed! oops"
		}

		return fmt.Sprintf("%s tried to assault %s and... %s",
			sender.Mention(), target.Mention(), result)
	}

	handleTargetedCmd(sess, inter, contentFunc)
}

// !!! its a joke command !!!
func handleSexnkill(
	sess *discordgo.Session,
	inter *discordgo.InteractionCreate,
) {
	contentFunc := func(sender, target *discordgo.User) string {
		mpreg := "made them pregnant"

		if flipACoin() {
			mpreg = "failed to make them pregnant"
		}

		return fmt.Sprintf("%s had sex with %s, %s and killed them!",
			sender.Mention(), target.Mention(), mpreg)
	}

	handleTargetedCmd(sess, inter, contentFunc)
}

// handler for / commands
func InteractionHandler(
	sess *discordgo.Session,
	inter *discordgo.InteractionCreate,
) {
	// we dont need to handle some weird interaction stuff here
	// only commands
	if inter.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch inter.ApplicationCommandData().Name {
	case CmdHelp:
		respond(sess, inter, "I don't know... umm... awkwar...", nil, false)
	case CmdMeow:
		respond(sess, inter, "Meow!", nil, false)
	case CmdMeowat:
		handleMeowat(sess, inter)
	case CmdRoulette:
		result := "Victory!!! You're alive!!!"

		if roulette() {
			result = "Sorry, you're dead, better luck next ti- uhh.."
		}

		respond(sess, inter, result, nil, false)
	case CmdMe:
		imgBytes, err := os.ReadFile("res/me.png")

		if err != nil {
			log.Printf("Cannot read 'me.png': %v", err)
			respond(sess, inter,
				"Oops, I couldn't find my own picture :<", nil, false)

			return
		}

		respond(sess, inter, "", []*discordgo.File{{
			Name:   "me.png",
			Reader: bytes.NewReader(imgBytes),
		}}, false)
	case CmdAssault:
		handleAssault(sess, inter)
	case CmdSexnkill:
		handleSexnkill(sess, inter)
	}
}

// handler for messages content
func MessageHandler(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == sess.State.User.ID {
		return
	}

	if strings.Contains(msg.Content, sess.State.User.Username) {
		msgs := []string{
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
		}

		sess.ChannelMessageSend(msg.ChannelID, msgs[rand.IntN(len(msgs))])
	}
}
