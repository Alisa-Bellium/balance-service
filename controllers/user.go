package controllers

import (
	"balance-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetUserBalanceInput struct {
	ID int64 `form:"id" binding:"required,gt=0"`
}

func GetUserBalance(c *gin.Context) {
	var input GetUserBalanceInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, reserved, err := services.GetUserBalance(input.ID)

	if err == services.ErrUserNotExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"balance":  user.Balance,
		"reserved": reserved,
	})
}
