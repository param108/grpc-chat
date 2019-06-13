package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

type ChatStore struct {
	db *gorm.DB
}

func NewChatStore() (*ChatStore, error) {
	db, err := gorm.Open(os.Getenv("DB_DRIVER"), DbDSN())
	if err != nil {
		return nil, err
	}
	return &ChatStore{db}, nil
}
