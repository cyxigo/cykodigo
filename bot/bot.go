package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// wake up cykodigo!
func WakeUp() {
	token := os.Getenv("TOKEN")
	sess, err := discordgo.New(token)

	if err != nil {
		log.Fatalf("Cannot create session: %v", sess)
	}

	// slash commands, looooooooooong list of them
	cmds := []*discordgo.ApplicationCommand{
		{
			Name:        CmdHelp,
			Description: "This... will not help actually",
		},
		{
			Name:        CmdMeow,
			Description: "He will meow",
		},
		{
			Name:        CmdMeowat,
			Description: "Meow at someone!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "member",
					Description: "Person to meow at",
					Required:    true,
				},
			},
		},
		{
			Name:        CmdRoulette,
			Description: "Why don't we play a little game?",
		},
		{
			Name:        CmdMe,
			Description: "Send picture of me! Nah, not you of course",
		},
		{
			Name:        CmdAssault,
			Description: "Try to assault someone... shh...",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "member",
					Description: "Person to try to assault",
					Required:    true,
				},
			},
		},
		{
			// !!! its a joke command !!!
			Name:        CmdSexnkill,
			Description: "I don't know why you would do that.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "member",
					Description: "Person to... you know",
					Required:    true,
				},
			},
		},
	}

	sess.AddHandler(InteractionHandler)
	sess.AddHandler(MessageHandler)

	err = sess.Open()

	for _, v := range cmds {
		_, err := sess.ApplicationCommandCreate(sess.State.User.ID, "", v)

		if err != nil {
			log.Fatalf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	if err != nil {
		log.Fatalf("Cannot open session: %v", err)
	}

	defer sess.Close()

	fmt.Println("Good morning cykodigo :)")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	fmt.Println("Die.")
}
