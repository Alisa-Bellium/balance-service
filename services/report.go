package services

import (
	"balance-service/repositories"
	"database/sql"
	"encoding/csv"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"time"
)

func StoreReport(month int, year int) (string, error) {
	lastID, err := repositories.GetLastTransactionIDForReport(nil, month, year)
	if err != nil {
		return "", err
	}

	filePath, err := repositories.FindReport(month, year, lastID)

	if err == nil {
		return filePath, nil
	}

	if err != sql.ErrNoRows {
		return "", err
	}

	transactionReports, err := repositories.GetReport(nil, month, year)
	if err != nil {
		return "", err
	}

	fileName, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	filePath = "data/" + fileName.String() + ".csv"
	csvRecords := getCsvRecords(transactionReports)
	f, err := os.Create(filePath)
	defer f.Close()

	if err != nil {
		return "", err
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(csvRecords)

	if err != nil {
		return "", err
	}

	maxID := getMaxIdFromTransactionReports(transactionReports)

	report := repositories.Report{
		Month:             month,
		Year:              year,
		CreatedAt:         time.Now().UTC(),
		FilePath:          filePath,
		LastTransactionID: maxID,
	}
	err = repositories.StoreReport(&report)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func getMaxIdFromTransactionReports(transactionReports []repositories.TransactionReport) sql.NullInt64 {
	if len(transactionReports) == 0 {
		return sql.NullInt64{}
	}

	maxID := transactionReports[0].LastTransactionID

	for _, e := range transactionReports {
		if maxID < e.LastTransactionID {
			maxID = e.LastTransactionID
		}
	}

	return sql.NullInt64{Int64: maxID, Valid: true}
}

func getCsvRecords(transactionReports []repositories.TransactionReport) [][]string {
	var records [][]string

	for _, e := range transactionReports {
		record := []string{strconv.FormatInt(e.ServiceID, 10), strconv.FormatInt(e.Total, 10)}
		records = append(records, record)
	}

	return records

}
