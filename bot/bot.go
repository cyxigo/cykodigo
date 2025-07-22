package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func Run() {
	token := os.Getenv("TOKEN")
	sess, err := discordgo.New(token)

	if err != nil {
		log.Fatalf("Cannot create session: %v", sess)
	}

	cmds := []*discordgo.ApplicationCommand{
		{
			Name:        "meow",
			Description: "He will meow.",
		},
	}

	sess.AddHandler(Handler)
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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Die.")
}
