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
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type NewUser struct {
	ID        UserID
	Username  string
	Email     *string
	CreatedAt time.Time
	UpdatedAt time.Time
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
