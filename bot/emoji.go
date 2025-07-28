package bot

import "fmt"

var (
	// i cant set all of this as const
	// so just please dont change them
	emojiCykodigo = emoji("cykodigo", "1399481322866741388")
	emojiCatr     = emoji("catr", "1399481457977856172")
)

func emoji(name string, id string) string {
	return fmt.Sprintf("<:%s:%s>", name, id)
}
