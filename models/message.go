package models

import (
	"time"
)

type ChatMessage struct {
	ID          string
	ServerTime  time.Time
	ClientTime  time.Time
	Message     string
	UserID      string
	ChatGroupID string
}

type ChatGroup struct {
	ID       string
	ChatName string
}

type ChatMessageArchive struct {
	ID          string
	ServerTime  time.Time
	ClientTime  time.Time
	Message     string
	UserID      string
	ChatGroupID string
}
