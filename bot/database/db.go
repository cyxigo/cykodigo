package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// maps guild IDs to database connections
var dbCache = make(map[string]*sql.DB)

// databases file directory
// this is where all databases created using getDB() will be stored
const dbDir = "database"

// initializes a database by creating all the needed tables and indexes
func initDB(db *sql.DB, name string) bool {
	_, err := db.Exec(
		`
		CREATE TABLE IF NOT EXISTS balances (
			user_id TEXT PRIMARY KEY,
			balance INTEGER NOT NULL DEFAULT 0
		);
		CREATE TABLE IF NOT EXISTS cooldowns (
			user_id TEXT PRIMARY KEY,
			last_work INTEGER NOT NULL DEFAULT 0,
			last_steal_fail INTEGER NOT NULL DEFAULT 0
		);
		CREATE TABLE IF NOT EXISTS inventory (
			user_id TEXT PRIMARY KEY,
			item TEXT NOT NULL,
			amount INTEGER NOT NULL DEFAULT 1,
			UNIQUE("user_id", "item")
		);
		CREATE TABLE IF NOT EXISTS meth_effects (
			user_id TEXT PRIMARY KEY,
			end_time INTEGER NOT NULL DEFAULT 0
		);
		
		CREATE INDEX IF NOT EXISTS idx_balances_user ON balances(user_id);
		CREATE INDEX IF NOT EXISTS idx_cooldowns_user ON cooldowns(user_id);
		CREATE INDEX IF NOT EXISTS idx_inventory_user ON inventory(user_id);
		CREATE INDEX IF NOT EXISTS idx_meth_effects_user ON meth_effects(user_id);
		`)

	if err != nil {
		log.Printf("Failed to create '%v' tables: %v", name, err)
		return false
	}

	return true
}

// returns the database for the guild with guildID
// creates a database if it doesn't exist
func GetDB(guildID string) (*sql.DB, bool) {
	if db, exists := dbCache[guildID]; exists {
		return db, true
	}

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Printf("Failed to get database for guild '%v': %v", guildID, err)
		return nil, false
	}

	dbPath := filepath.Join(dbDir, guildID+".db")
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Printf("Failed to open database for guild '%v': %v", guildID, err)
		return nil, false
	}

	if !initDB(db, dbPath) {
		return nil, false
	}

	dbCache[guildID] = db
	return db, true
}

// util function to get info about a users meth effect status directly from database (not via transaction)
// used exactly one time but i dont care it shortens the code
func GetUserHighInfo(db *sql.DB, userID string) (bool, int64) {
	endTime := int64(0)
	err := db.QueryRow(
		`
		SELECT end_time 
		FROM meth_effects 
		WHERE user_id = ?
		`, userID).Scan(&endTime)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Query error in TxGetUserHighInfo: %v", err)
	}

	return time.Now().Unix() < endTime, endTime
}
