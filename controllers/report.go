package controllers

import (
	"balance-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StoreReportInput struct {
	Year  int `json:"year" binding:"required,min=0,max=9999"`
	Month int `json:"month" binding:"required,min=1,max=12"`
}

func StoreReport(c *gin.Context) {
	var json StoreReportInput
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	filePath, err := services.StoreReport(json.Month, json.Year)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	url := scheme + "://" + c.Request.Host + "/" + filePath

	c.JSON(http.StatusCreated, gin.H{
		"url": url,
	})
}
