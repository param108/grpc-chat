package models

import (
	"time"
)

type available struct {
	ID        int
	UserID    string
	StartTime time.Time
	EndTime   time.Time
}
