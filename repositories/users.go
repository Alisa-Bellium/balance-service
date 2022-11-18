package repositories

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID      int64 `db:"id"`
	Balance int64 `db:"balance"`
}

func LockUser(tx *sqlx.Tx, userID int64) (*User, error) {
	var user User

	// Preventing race condition by locking user record
	lockUserQuery := "SELECT * FROM users WHERE id=$1 FOR UPDATE"
	if err := tx.Get(&user, lockUserQuery, userID); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return &user, nil
}

func GetUser(tx *sqlx.Tx, ID int64) (*User, error) {
	var user User

	var err error
	if tx == nil {
		err = DB.Get(&user, "SELECT * FROM users WHERE id=$1 LIMIT 1", ID)
	} else {
		err = tx.Get(&user, "SELECT * FROM users WHERE id=$1 LIMIT 1", ID)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, nil
}

func StoreUser(tx *sqlx.Tx, ID int64) error {
	user := User{ID: ID, Balance: 0}
	_, err := tx.NamedExec("INSERT INTO users (id, balance) VALUES (:id, :balance)", &user)

	return err
}

func StoreUserIfNotExists(tx *sqlx.Tx, userID int64) error {
	user, err := GetUser(tx, userID)

	if err != nil {
		return err
	}

	if user == nil {
		if err = StoreUser(tx, userID); err != nil {
			return err
		}
	}

	return nil
}

func UpdateUserBalance(tx *sqlx.Tx, userID int64, amount int64) (*User, error) {
	var user User
	updateUserQuery := "UPDATE users SET balance = balance + $1 WHERE id=$2 RETURNING *"

	if err := tx.QueryRowx(updateUserQuery, amount, userID).StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
