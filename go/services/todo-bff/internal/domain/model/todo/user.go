package todo

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
