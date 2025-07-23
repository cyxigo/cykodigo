package main

import (
	"log"

	"github.com/cyxigo/cykodigo/bot"
	"github.com/joho/godotenv"
)

// load .env stuff and init economy database
func init() {
	if err := godotenv.Load(); err != nil {
		// fatal cus bot cant start without a token
		// and token is stored in .env
		log.Fatalf("No .env file found")
	}

	bot.InitDB()
}

func main() {
	bot.WakeUp()
}
