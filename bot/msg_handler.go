package bot

import (
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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

func handleMsgNature(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	imgPath := "res/gif/nature.gif"
	file, err := os.Open(imgPath)

	if err != nil {
		log.Printf("Error opening '%v': %v", imgPath, err)
		sess.ChannelMessageSend(msg.ChannelID, "Couldn't open image :<")

		return
	}

	defer file.Close()

	// yeah.. all of this is a bunch of code ripped out from handleImageCmd
	// ¯\_(ツ)_/¯
	imgName := filepath.Base(imgPath)
	discordFile := &discordgo.File{
		Name:   imgName,
		Reader: file,
	}
	embed := &discordgo.MessageEmbed{
		Description: "**RULES OF NATURE!!!**",
		Color:       defaultEmbedColor,
		Image: &discordgo.MessageEmbedImage{
			URL: "attachment://" + imgName,
		},
	}
	msgSend := &discordgo.MessageSend{
		Embed: embed,
		File:  discordFile,
	}

	sess.ChannelMessageSendComplex(msg.ChannelID, msgSend)
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
	case strings.Contains(content, cmdMsgNature):
		handleMsgNature(sess, msg)
	}
}
