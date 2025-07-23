package bot

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// wake up cykodigo!
func WakeUp() {
	token := os.Getenv("TOKEN")
	sess, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Fatalf("Can't create session: %v", err)
	}

	sess.AddHandler(InterHandler)
	sess.AddHandler(MsgHandler)

	err = sess.Open()

	if err != nil {
		log.Fatalf("Can't open session: %v", err)
	}

	defer sess.Close()

	appID := sess.State.User.ID

	// you know what
	// i hate this loop
	// cus its sooooooooooooooooooooooooooooooooooooooo long
	// like to startup bot on my shitty laptop it takes around a minute
	for _, v := range cmds {
		_, err := sess.ApplicationCommandCreate(appID, "", v)

		if err != nil {
			log.Fatalf("Can't create '%s' command: %v", v.Name, err)
		}

		log.Printf("Created '%s' command", v.Name)
	}

	log.Printf("Online")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	log.Printf("Offline")
}
