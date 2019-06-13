package store

import (
	"fmt"
	"os"
)

func DbDSN() string {
	driver := os.Getenv("DB_DRIVER")

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dbConfigString := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&timezone=UTC", driver, username,
		password, host, port, name)
	return dbConfigString
}
