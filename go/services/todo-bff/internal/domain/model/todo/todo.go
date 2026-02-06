package todo

import (
	"strconv"
	"time"
)

type TodoID int64

type TodoStatus string

const (
	Pending   TodoStatus = "pending"
	InProcess TodoStatus = "in_process"
	Done      TodoStatus = "done"
)

func (ts TodoStatus) IsValid() bool {
	switch ts {
	case Pending, InProcess, Done:
		return true
	default:
		return false
	}
}

type Todo struct {
	ID          TodoID
	UserID      UserID
	Task        string
	Description string
	Status      TodoStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (id *TodoID) Int64() int64 {
	if id == nil {
		return 0
	}

	return int64(*id)
}

func (id *TodoID) String() string {
	if id == nil {
		return ""
	}

	return strconv.FormatInt(int64(*id), 10)
}

func NewTodoID(id int64) *TodoID {
	todoID := TodoID(id)
	return &todoID
}

func (t *Todo) IsDeleted() bool {
	if t == nil {
		return false
	}

	return t.DeletedAt != nil
}
