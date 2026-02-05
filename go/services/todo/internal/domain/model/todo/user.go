package model

import (
	"strconv"
	"time"
)

type UserID int64

type User struct {
	ID        UserID
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type NewUser struct {
	Username string
	Email    string
	Password string
}

func (id *UserID) Int64() int64 {
	if id == nil {
		return 0
	}
	return int64(*id)
}

func (id *UserID) String() string {
	if id == nil {
		return ""
	}
	return strconv.FormatInt(int64(*id), 10)
}

func NewUserID(id int64) *UserID {
	userID := UserID(id)
	return &userID
}
