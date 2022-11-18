package services

import (
	"balance-service/repositories"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"time"
)

var ErrInsufficientBalance = errors.New("insufficient balance to provide a transaction")
var ErrTransactionAlreadyProcessed = errors.New("transaction already processed")
var ErrTransactionNotFound = errors.New("not found transaction to withdrawal")
var ErrTransactionWrongAmount = errors.New("withdrawal transaction should has same amount, as initial one")
var ErrTransactionAlreadyCancelled = errors.New("transaction already cancelled")

func StoreReplenishmentTransaction(userID int64, amount int64) (*repositories.User, error) {
	if amount <= 0 {
		return nil, errors.New("amount should be positive")
	}

	transaction := repositories.Transaction{
		UserID:           userID,
		Amount:           amount,
		IsReserveAccount: false,
		CreatedAt:        time.Now().UTC(),
	}

	tx := repositories.DB.MustBegin()
	if err := repositories.StoreUserIfNotExists(tx, userID); err != nil {
		return nil, err
	}
	if _, err := repositories.LockUser(tx, userID); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := repositories.StoreTransaction(tx, &transaction); err != nil {
		return nil, err
	}

	user, err := repositories.UpdateUserBalance(tx, userID, amount)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func StoreReservationTransaction(userID int64, amount int64, orderID int64, serviceID int64) (*repositories.User, error) {
	if amount < 0 {
		return nil, errors.New("amount should be either positive or zero")
	}

	withdrawalTransaction := repositories.Transaction{
		UserID:           userID,
		ServiceID:        sql.NullInt64{Int64: serviceID, Valid: true},
		OrderID:          sql.NullInt64{Int64: orderID, Valid: true},
		Amount:           -amount,
		IsReserveAccount: false,
		CreatedAt:        time.Now().UTC(),
	}
	reservationTransaction := repositories.Transaction{
		UserID:           userID,
		ServiceID:        sql.NullInt64{Int64: serviceID, Valid: true},
		OrderID:          sql.NullInt64{Int64: orderID, Valid: true},
		Amount:           amount,
		IsReserveAccount: true,
		CreatedAt:        time.Now().UTC(),
	}

	tx := repositories.DB.MustBegin()

	if user, err := repositories.GetUser(tx, userID); err != nil || user == nil {
		_ = tx.Rollback()

		if err != nil {
			return nil, err
		}

		return nil, ErrInsufficientBalance
	}

	user, err := repositories.LockUser(tx, userID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if user.Balance+withdrawalTransaction.Amount < 0 {
		_ = tx.Rollback()
		return nil, ErrInsufficientBalance
	}

	existingTransaction, err := repositories.GetServiceTransaction(tx, userID, serviceID, orderID, false)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if existingTransaction != nil {
		_ = tx.Rollback()
		return nil, ErrTransactionAlreadyProcessed
	}

	if err := repositories.StoreTransaction(tx, &withdrawalTransaction); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := repositories.StoreTransaction(tx, &reservationTransaction); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	user, err = repositories.UpdateUserBalance(tx, userID, -amount)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return user, nil
}

func StoreWithdrawalTransaction(userID int64, amount int64, orderID int64, serviceID int64) (*repositories.User, error) {
	if amount < 0 {
		return nil, errors.New("amount should be either positive or zero")
	}

	cancelReservationTransaction := repositories.Transaction{
		UserID:           userID,
		ServiceID:        sql.NullInt64{Int64: serviceID, Valid: true},
		OrderID:          sql.NullInt64{Int64: orderID, Valid: true},
		Amount:           -amount,
		IsReserveAccount: true,
		CreatedAt:        time.Now().UTC(),
	}

	tx := repositories.DB.MustBegin()

	if user, err := repositories.GetUser(tx, userID); err != nil || user == nil {
		_ = tx.Rollback()

		if err != nil {
			return nil, err
		}

		return nil, ErrTransactionNotFound
	}

	user, err := repositories.LockUser(tx, userID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	transactionToCancel, err := repositories.GetServiceTransaction(tx, userID, serviceID, orderID, true)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if transactionToCancel == nil {
		_ = tx.Rollback()
		return nil, ErrTransactionNotFound
	}

	if transactionToCancel.Amount != amount {
		_ = tx.Rollback()
		return nil, ErrTransactionWrongAmount
	}

	cancelReservationTransaction.CancelledTransactionId = sql.NullInt64{Int64: transactionToCancel.ID, Valid: true}

	reserved, err := repositories.GetUserReservedAmount(tx, userID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if reserved-amount < 0 {
		_ = tx.Rollback()
		return nil, errors.New("reserved balance cannot be less than 0")
	}

	cancellingTransaction, err := repositories.GetCancellingTransaction(tx, transactionToCancel.ID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if cancellingTransaction != nil {
		_ = tx.Rollback()
		return nil, ErrTransactionAlreadyCancelled
	}

	if err := repositories.StoreTransaction(tx, &cancelReservationTransaction); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return user, nil
}
