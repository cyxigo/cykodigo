package main

import (
	"log"

	"github.com/cyxigo/cykodigo/bot"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	bot.Run()
}
