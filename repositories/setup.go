package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

var DB *sqlx.DB

func CreateConnection() error {
	if DB != nil {
		return nil
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	var err error
	DB, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
			host,
			port,
			database,
			user,
			password,
		),
	)

	return err
}
