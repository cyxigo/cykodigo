package main

import (
	"github.com/cyxigo/cykodigo/bot"
)

// load .env stuff and init economy database
func init() {
	bot.InitEnv()
	bot.InitDB()
}

func main() {
	bot.WakeUp()
}
