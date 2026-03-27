package auth

import (
	"strconv"
	"time"
)

type UserID int64

type User struct {
	ID        UserID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type NewUser struct {
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

func ParseUserID(id string) (UserID, error) {
	userIDInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, err
	}
	userID := UserID(userIDInt)
	return userID, nil
}

func (u *User) IsDeleted() bool {
	if u == nil {
		return false
	}

	return u.DeletedAt != nil
}

func (u *User) Delete() {
	if u == nil {
		return
	}

	now := time.Now()
	u.DeletedAt = &now
}
