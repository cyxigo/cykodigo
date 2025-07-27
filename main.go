package main

import (
	"github.com/cyxigo/cykodigo/bot"
)

// load .env stuff and init economy database
func init() {
	bot.InitEnv()
}

func main() {
	bot.WakeUp()
}
