package msg

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cyxigo/cykodigo/bot/data"
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
		"Look, whatever you selling, I'm not buying yo.",
		"WHY BE THE [[Little Sponge]] WHO HATES ITS [[$4.99]] LIFE\n" +
			"WHEN YOU CAN BE A\n" +
			"[[BIG SHOT!!!]]\n" +
			"[[BIG SHOT!!!!]]\n" +
			"[[BIG SHOT!!!!!]]\n" +
			"THAT'S RIGHT!! NOW'S YOUR CHANCE TO BE A [[BIG SHOT]]!!",
		"[[Hyperlink Blocked]]",
		"I'm a [[BIG SHOT]]!",
		"I'LL SHOW YOU WHAT IT MEANS TO BE FREE!",
		"I CAN DO ANYTHING!",
		"CHAOS, CHAOS!",
		"DO YOU WANT TO PLAY GAME, GAME?",
		"UEE HEE HEE",
		"I'm old!",
		"Freedom...",
		"True",
		"False",
		"Wow",
		"Um",
		"Uh",
		"Umm",
		"Uhh",
		"Yo",
		"Yes?",
		"Hm",
		"Hmm",
		"Hmmmmm",
		"Hmmmmmmmmmmmmmmmmmmmmmmm",
		"Hmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm",
		"Hm?",
		"Shhhhhhhhhhhhhhhh",
		"Grrrrrrrrrrrrrrrr",
		"Let's dance, use `/dance`",
		"You better go work, use `/work`",
		"",
		"IT'S!! TV!! TIME!!!",
		"I LOVE TV",
		"I'm groovy and never glooby",
	}

	handleMsgReplyRandVariant(sess, msg, variants)
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

	handleMsgReplyRandVariant(sess, msg, variants)
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
	case strings.Contains(content, sess.State.User.Username):
		switch {
		case strings.Contains(content, data.CmdMsgExplodeBalls):
			sess.ChannelMessageSend(msg.ChannelID, "BOOM!1!11!! 💥💥💥💥💥")
		default:
			handleMsgBotUsername(sess, msg)
		}
	case strings.Contains(content, data.CmdMsgMeow):
		handleMsgMeow(sess, msg)
	case strings.Contains(content, data.CmdMsgCrazy):
		sess.ChannelMessageSend(msg.ChannelID,
			"Crazy? I was crazy once. They locked me in a room, a rubber room, a rubber room with rats, "+
				"and rats make me crazy.",
		)
	case strings.Contains(content, data.CmdMsgGlamptastic):
		sess.ChannelMessageSend(msg.ChannelID, "glamptastic!")
	}
}
