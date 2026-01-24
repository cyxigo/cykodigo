package bot

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/cyxigo/cykodigo/bot/data"
	"github.com/cyxigo/cykodigo/bot/inter"
	"github.com/cyxigo/cykodigo/bot/msg"
)

// wake up cykodigo!
func WakeUp() {
	token, ok := data.GetEnvVariable("TOKEN")

	if !ok {
		return
	}

	sess, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Fatalf("Can't create session: %v", err)
	}

	dialer, httpClient, ok := SetupProxy()

	if ok {
		sess.Client = httpClient
		sess.Dialer = dialer
	}

	sess.StateEnabled = true

	sess.AddHandler(inter.InterHandler)
	sess.AddHandler(msg.MsgHandler)

	sess.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages

	err = sess.Open()

	if err != nil {
		log.Fatalf("Can't open session: %v", err)
	}

	defer sess.Close()

	sess.ApplicationCommandBulkOverwrite(sess.State.User.ID, "", data.Cmds)

	log.Printf("Online")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	log.Printf("Offline")
}
