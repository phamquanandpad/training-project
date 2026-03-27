package service

import (
	"context"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

type TodoHelper interface {
	CanAccessTodo(
		ctx context.Context,
		userID todo.UserID,
		todoID todo.TodoID,
	) (bool, error)
}
