package controllers

import (
	"balance-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StoreReplenishmentTransactionInput struct {
	UserID int64 `json:"user_id" binding:"required,gt=0"`
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

type StoreReservationTransactionInput struct {
	UserID    int64 `json:"user_id" binding:"required,gt=0"`
	Amount    int64 `json:"amount" binding:"required,gt=0"`
	ServiceID int64 `json:"service_id" binding:"required,gt=0"`
	OrderID   int64 `json:"order_id" binding:"required,gt=0"`
}

type StoreWithdrawalTransactionInput struct {
	UserID    int64 `json:"user_id" binding:"required,gt=0"`
	Amount    int64 `json:"amount" binding:"required,gt=0"`
	ServiceID int64 `json:"service_id" binding:"required,gt=0"`
	OrderID   int64 `json:"order_id" binding:"required,gt=0"`
}

func StoreReplenishmentTransaction(c *gin.Context) {
	var json StoreReplenishmentTransactionInput
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := services.StoreReplenishmentTransaction(json.UserID, json.Amount)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": user.ID,
		"balance": user.Balance,
	})
}

func StoreReservationTransaction(c *gin.Context) {
	var json StoreReservationTransactionInput
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := services.StoreReservationTransaction(json.UserID, json.Amount, json.OrderID, json.ServiceID)

	if err == services.ErrInsufficientBalance || err == services.ErrTransactionAlreadyProcessed {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": user.ID,
		"balance": user.Balance,
	})
}

func StoreWithdrawalTransaction(c *gin.Context) {
	var json StoreWithdrawalTransactionInput
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := services.StoreWithdrawalTransaction(json.UserID, json.Amount, json.OrderID, json.ServiceID)

	if err == services.ErrTransactionNotFound || err == services.ErrTransactionWrongAmount || err == services.ErrTransactionAlreadyCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": user.ID,
		"balance": user.Balance,
	})
}
