package msg

import (
	"math/rand/v2"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cyxigo/cykodigo/bot/data"
)

func handleMsgReplyRandVariant(sess *discordgo.Session, msg *discordgo.MessageCreate, variants []string) {
	sess.MessageReactionAdd(msg.ChannelID, msg.ID, data.EmojiReactCykodigo)

	chosen := variants[rand.IntN(len(variants))]
	lines := strings.Split(chosen, "\n")

	for i := range lines {
		lines[i] = lines[i] + " " + data.EmojiCykodigo
	}

	content := strings.Join(lines, "\n")

	sess.ChannelMessageSendReply(msg.ChannelID, content, msg.Reference())
}
