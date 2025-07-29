package main

import (
	"github.com/cyxigo/cykodigo/bot"
	"github.com/cyxigo/cykodigo/bot/data"
)

// load .env stuff and init economy database
func init() {
	data.InitEnv()
}

func main() {
	bot.WakeUp()
}
