package data

import "fmt"

var (
	// i cant set all of this as const
	// so just please dont change them
	EmojiCykodigo = emoji("cykodigo", "1399481322866741388")
	EmojiCatr     = emoji("catr", "1400227174975410327")

	EmojiReactCykodigo = reactionEmoji("cykodigo", "1399481322866741388")
	EmojiReactCatr     = reactionEmoji("catr", "1400227174975410327")
)

func reactionEmoji(name string, id string) string {
	return fmt.Sprintf("%v:%v", name, id)
}

func emoji(name string, id string) string {
	return fmt.Sprintf("<:%v:%v>", name, id)
}
