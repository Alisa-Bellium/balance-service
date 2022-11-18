package main

import (
	"balance-service/controllers"
	"balance-service/repositories"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Set up server
	if err := repositories.CreateConnection(); err != nil {
		log.Fatal(err)
		return
	}

	r := gin.Default()
	r.Static("/data", "./data")

	v1 := r.Group("/v1")
	v1.POST("/transactions/replenish", controllers.StoreReplenishmentTransaction)
	v1.POST("/transactions/reserve", controllers.StoreReservationTransaction)
	v1.POST("/transactions/withdraw", controllers.StoreWithdrawalTransaction)

	v1.GET("/users", controllers.GetUserBalance)

	v1.POST("/report", controllers.StoreReport)

	err := r.Run(":8080")

	if err != nil {
		log.Fatal(err)
		return
	}
}
