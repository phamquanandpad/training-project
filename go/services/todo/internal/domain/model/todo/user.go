package todo

import (
	"strconv"
	"time"
)

type UserID int64

type User struct {
	ID        UserID
	Username  string
	Email     *string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type NewUser struct {
	Username string
	Email    *string
	Password string
}

type UserOverview struct {
	ID       UserID
	Username string
	Email    *string
}

type UserWithTodos struct {
	ID       UserID
	Username string
	Email    *string
	Todos    []*Todo `gorm:"foreignKey:UserID"`
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

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

func (u *User) Delete() {
	now := time.Now()
	u.DeletedAt = &now
}
