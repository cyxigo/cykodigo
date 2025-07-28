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

// util function for getting .env variables
// logs if the variable doesn't exist
func getEnvVariable(name string) (string, bool) {
	variable, ok := os.LookupEnv(name)

	if !ok {
		log.Printf("Couldn't find variable '%v' in .env file", name)
	}

	return variable, ok
}
