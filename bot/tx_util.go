package bot

import (
	"database/sql"
	"log"
	"time"
)

// util function for getting user balances in sql transactions
// yes /balance doesnt use it
// cus /balance doesnt need sql transactions since its just one query
func getUserBalance(tx *sql.Tx, userID string) int {
	balance := 0
	err := tx.QueryRow(
		`
		SELECT balance
		FROM balances 
		WHERE user_id = ?
		`, userID).Scan(&balance)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Query error in getUserBalance: %v", err)
	}

	return balance
}

// util function to get info about a users meth effect status in sql transactions
// returns whether the user is currently high and the effect end time
func getUserHighInfo(tx *sql.Tx, userID string) (bool, int64) {
	endTime := int64(0)
	err := tx.QueryRow(
		`
		SELECT end_time 
		FROM meth_effects 
		WHERE user_id = ?
		`, userID).Scan(&endTime)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Query error in isUserHigh: %v", err)
	}

	return time.Now().Unix() < endTime, endTime
}
