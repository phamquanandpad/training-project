package datastore

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

type todoReader struct{}

func NewTodoReader() gateway.TodoQueriesGateway {
	return &todoReader{}
}

func (r *todoReader) GetTodo(
	ctx context.Context,
	todoID todo.TodoID,
	userID todo.UserID,
) (*todo.Todo, error) {
	tx, err := ExtractTodoDB(ctx)
	if err != nil {
		return nil, err
	}
	db := tx.WithContext(ctx)

	todo := new(todo.Todo)
	err = db.
		Where("id = ? AND deleted_at IS NULL", todoID).
		Where("user_id = ?", userID).
		First(&todo).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return todo, nil
}

func (r *todoReader) ListTodos(
	ctx context.Context,
	userID todo.UserID,
) ([]*todo.Todo, int, error) {
	tx, er := ExtractTodoDB(ctx)
	if er != nil {
		return nil, 0, er
	}
	db := tx.WithContext(ctx)

	var todos []*todo.Todo
	err := db.
		Where("deleted_at IS NULL").Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&todos).
		Error
	if err != nil {
		return nil, 0, err
	}

	total := len(todos)

	return todos, total, nil
}
