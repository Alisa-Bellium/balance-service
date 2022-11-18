package repositories

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type Transaction struct {
	ID                     int64         `db:"id"`
	UserID                 int64         `db:"user_id"`
	Amount                 int64         `db:"amount"`
	ServiceID              sql.NullInt64 `db:"service_id"`
	OrderID                sql.NullInt64 `db:"order_id"`
	IsReserveAccount       bool          `db:"is_reserve_account"`
	CreatedAt              time.Time     `db:"created_at"`
	CancelledTransactionId sql.NullInt64 `db:"canceled_transaction_id"`
}

type TransactionReport struct {
	ServiceID         int64 `db:"service_id"`
	Total             int64 `db:"total"`
	LastTransactionID int64 `db:"last_transaction_id"`
}

func StoreTransaction(tx *sqlx.Tx, transaction *Transaction) error {
	insertTransactionQuery := "INSERT INTO transactions (user_id, created_at, amount, service_id, order_id, is_reserve_account, canceled_transaction_id) VALUES (:user_id, :created_at, :amount, :service_id, :order_id, :is_reserve_account, :canceled_transaction_id)"
	_, err := tx.NamedExec(insertTransactionQuery, transaction)
	return err
}

func GetServiceTransaction(tx *sqlx.Tx, userID int64, serviceID int64, orderID int64, isReserveAccount bool) (*Transaction, error) {
	var transaction Transaction
	transactionQuery := "SELECT * FROM transactions WHERE user_id=$1 and service_id=$2 and order_id=$3 and is_reserve_account=$4 LIMIT 1"

	err := tx.Get(&transaction, transactionQuery, userID, serviceID, orderID, isReserveAccount)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &transaction, nil
}

func GetCancellingTransaction(tx *sqlx.Tx, canceledTransactionID int64) (*Transaction, error) {
	var transaction Transaction
	transactionQuery := "SELECT * FROM transactions WHERE canceled_transaction_id=$1 LIMIT 1"
	err := tx.Get(&transaction, transactionQuery, canceledTransactionID)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &transaction, nil
}

func GetUserReservedAmount(tx *sqlx.Tx, userID int64) (int64, error) {
	var sum int64
	selectReservedQuery := "SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id=$1 and is_reserve_account=true"

	var err error
	if tx == nil {
		err = DB.QueryRow(selectReservedQuery, userID).Scan(&sum)
	} else {
		err = tx.QueryRow(selectReservedQuery, userID).Scan(&sum)
	}

	if err != nil {
		return 0, err
	}

	return sum, nil
}

func GetLastTransactionIDForReport(tx *sqlx.Tx, month int, year int) (sql.NullInt64, error) {
	var lastID sql.NullInt64
	selectLastIDQuery := `select max(id)
			from transactions t
			where t.is_reserve_account = false
			  and t.service_id is not null
			  and t.order_id is not null
			  and amount < 0
			  and date_part('year', created_at)=$1 
			  and date_part('month', created_at)=$2
			  and exists(select 1
						 from transactions t2
						 where t.order_id = t2.order_id
						   and t.service_id = t2.service_id
						   and t2.is_reserve_account = true
						   and t2.canceled_transaction_id is not null)`

	var err error
	if tx == nil {
		err = DB.QueryRow(selectLastIDQuery, year, month).Scan(&lastID)
	} else {
		err = tx.QueryRow(selectLastIDQuery, year, month).Scan(&lastID)
	}

	if err != nil {
		return sql.NullInt64{}, err
	}

	return lastID, nil
}

func GetReport(tx *sqlx.Tx, month int, year int) ([]TransactionReport, error) {
	var transactionReports []TransactionReport
	reportQuery := `select t.service_id, sum(abs(t.amount)) as total, max(t.id) as last_transaction_id
			from transactions t
			where t.is_reserve_account = false
			  and t.service_id is not null
			  and t.order_id is not null
			  and amount < 0
			  and date_part('year', created_at)=$1 
			  and date_part('month', created_at)=$2
			  and exists(select 1
						 from transactions t2
						 where t.order_id = t2.order_id
						   and t.service_id = t2.service_id
						   and t2.is_reserve_account = true
						   and t2.canceled_transaction_id is not null)
			group by t.service_id`

	var err error
	if tx == nil {
		err = DB.Select(&transactionReports, reportQuery, year, month)
	} else {
		err = tx.Select(&transactionReports, reportQuery, year, month)
	}

	if err != nil {
		return nil, err
	}

	return transactionReports, nil
}
