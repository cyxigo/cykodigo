package bot

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}
}

// i COULD just return empty string as a sign of failure but you know
// i like returning "ok" bool
func getEnvVariable(name string) (string, bool) {
	variable := os.Getenv(name)

	if variable == "" {
		log.Printf("Couldn't find variable '%v' in .env file", name)
		return "", false
	}

	return variable, true
}
