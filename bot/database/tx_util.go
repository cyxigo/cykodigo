package database

import (
	"database/sql"
	"log"
	"time"
)

// util function for getting user balances in sql transactions
func TxGetUserBalance(tx *sql.Tx, userID string) int64 {
	balance := int64(0)
	err := tx.QueryRow(
		`
		SELECT balance
		FROM balances 
		WHERE user_id = ?
		`, userID).Scan(&balance)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Query error in TxGetUserBalance: %v", err)
	}

	return balance
}

// util function for getting user balances (in bank) in sql transactions
func TxGetUserBankBalance(tx *sql.Tx, userID string) int64 {
	balance := int64(0)
	err := tx.QueryRow(
		`
		SELECT bank
		FROM balances 
		WHERE user_id = ?
		`, userID).Scan(&balance)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Query error in TxGetUserBankBalance: %v", err)
	}

	return balance
}

// util function to get info about a user meth effect status in sql transactions
// returns whether the user is currently high and the effect end time
func TxGetUserHighInfo(tx *sql.Tx, userID string) (bool, int64) {
	endTime := int64(0)
	err := tx.QueryRow(
		`
		SELECT high_end_time 
		FROM effects 
		WHERE user_id = ?
		`, userID).Scan(&endTime)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Query error in TxGetUserHighInfo: %v", err)
	}

	return time.Now().Unix() < endTime, endTime
}
