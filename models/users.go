package models

import ()

type User struct {
	ID            string
	Name          string
	Role          string
	Password      string
	FirebaseToken string
}

type UserToken struct {
	ID            int
	UserToken     string
	UserID        string
	UserRole      string
	FirebaseToken string
}

type UserGroup struct {
	ID          string
	ChatGroupID string
	UserID      string
}
