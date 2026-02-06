package todo

type TodoID int64

type TodoStatus string

const (
	Pending   TodoStatus = "pending"
	InProcess TodoStatus = "in_process"
	Done      TodoStatus = "done"
)

type Todo struct {
	UserID      UserID
	Title       string
	Description string
	Status      TodoStatus
}

func (id *TodoID) Int64() int64 {
	if id == nil {
		return 0
	}
	return int64(*id)
}
