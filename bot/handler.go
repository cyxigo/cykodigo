package bot

import "github.com/bwmarrin/discordgo"

func Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
	content := ""

	switch inter.ApplicationCommandData().Name {
	case "meow":
		content = "Meow!"
	}

	sess.InteractionRespond(inter.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}
