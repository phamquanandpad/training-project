package todo

import "strconv"

type UserID int64

type User struct {
	Username string
	Email    *string
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
