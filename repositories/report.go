package repositories

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type Report struct {
	ID                int64         `db:"id"`
	Month             int           `db:"month"`
	Year              int           `db:"year"`
	CreatedAt         time.Time     `db:"created_at"`
	FilePath          string        `db:"file_path"`
	LastTransactionID sql.NullInt64 `db:"last_transaction_id"`
}

func FindReport(month int, year int, lastTransactionID sql.NullInt64) (string, error) {
	var filePath string

	var err error
	if lastTransactionID.Valid {
		err = DB.QueryRow("SELECT file_path FROM reports WHERE year=$1 AND month=$2 AND last_transaction_id=$3 LIMIT 1", year, month, lastTransactionID).Scan(&filePath)
	} else {
		err = DB.QueryRow("SELECT file_path FROM reports WHERE year=$1 AND month=$2 AND last_transaction_id IS NULL LIMIT 1", year, month).Scan(&filePath)
	}

	if err != nil {
		return "", err
	}

	return filePath, nil
}

func StoreReport(report *Report) error {
	insertQuery := "INSERT INTO reports (month, year, created_at, file_path, last_transaction_id) VALUES (:month, :year, :created_at, :file_path, :last_transaction_id)"
	_, err := DB.NamedExec(insertQuery, report)
	return err
}
