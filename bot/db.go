package bot

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	path, ok := getEnvVariable("DB_PATH")

	if !ok {
		return
	}

	var err error
	DB, err = sql.Open("sqlite3", path)

	if err != nil {
		log.Fatalf("Can't open '%v': %v", path, err)
	}

	_, err = DB.Exec(
		`
		CREATE TABLE IF NOT EXISTS balances (
			user_id TEXT PRIMARY KEY,
			balance INTEGER NOT NULL DEFAULT 0,
			last_work INTEGER NOT NULL DEFAULT 0,
			last_steal_fail INTEGER NOT NULL DEFAULT 0
		);
		CREATE TABLE IF NOT EXISTS inventory (
			user_id TEXT NOT NULL,
			item TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_inventory_user ON inventory(user_id);
		`)

	if err != nil {
		log.Fatalf("Error creating '%v' tables: %v", path, err)
	}
}
