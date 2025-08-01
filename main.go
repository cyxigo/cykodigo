package main

import (
	"github.com/cyxigo/cykodigo/bot"
	"github.com/cyxigo/cykodigo/bot/data"
)

func init() {
	data.InitEnv()
}

func main() {
	bot.WakeUp()
}
