package store

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/param108/grpc-chat/models"
	"os"
)

var (
	testStore *ChatStore
)

func newTestChatStore() (*ChatStore, error) {
	db, err := gorm.Open(os.Getenv("TEST_DB_DRIVER"), testDbDSN())
	if err != nil {
		return nil, err
	}
	return &ChatStore{db}, nil
}

func testDbDSN() string {
	driver := os.Getenv("TEST_DB_DRIVER")

	username := os.Getenv("TEST_DB_USERNAME")
	password := os.Getenv("TEST_DB_PASSWORD")
	host := os.Getenv("TEST_DB_HOST")
	port := os.Getenv("TEST_DB_PORT")
	name := os.Getenv("TEST_DB_NAME")

	dbConfigString := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&timezone=UTC", driver, username,
		password, host, port, name)
	return dbConfigString
}

func init() {
	newTestStore, err := newTestChatStore()
	if err != nil {
		panic(err.Error())
	}
	testStore = newTestStore
}

func teardown() {
	if err := testStore.db.Delete(models.ChatMessage{}).Error; err != nil {
		panic("Failed to clear ChhatMessages db:" + err.Error())
	}
	if err := testStore.db.Delete(models.UserGroup{}).Error; err != nil {
		panic("Failed to clear UserGroup db:" + err.Error())
	}
	if err := testStore.db.Delete(models.ChatGroup{}).Error; err != nil {
		panic("Failed to clear ChatGroup db:" + err.Error())
	}
	if err := testStore.db.Delete(models.UserToken{}).Error; err != nil {
		panic("Failed to clear UserToken db:" + err.Error())
	}
	if err := testStore.db.Delete(models.User{}).Error; err != nil {
		panic("Failed to clear User db:" + err.Error())
	}
}
