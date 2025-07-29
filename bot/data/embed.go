package data

import "github.com/bwmarrin/discordgo"

const (
	DefaultEmbedColor = 0xFF7B00
)

// util function for creating embeds
func EmbedText(content string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Description: content,
		Color:       DefaultEmbedColor,
	}

	return embed
}
