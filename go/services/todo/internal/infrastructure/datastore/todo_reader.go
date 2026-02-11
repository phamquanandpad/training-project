package datastore

import (
	"context"

	"errors"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"gorm.io/gorm"
)

type todoReader struct{}

func NewTodoReader() gateway.TodoQueriesGateway {
	return &todoReader{}
}

func (r *todoReader) GetTodo(context context.Context, id todo.TodoID) (*todo.Todo, error) {
	tx, err := ExtractTodoDB(context)
	if err != nil {
		return nil, err
	}
	db := tx.WithContext(context)

	todo := new(todo.Todo)
	err = db.
		Where("id = ? AND deleted_at IS NULL", id).
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

func (r *todoReader) ListTodos(context context.Context) ([]*todo.Todo, error) {
	tx, er := ExtractTodoDB(context)
	if er != nil {
		return nil, er
	}
	db := tx.WithContext(context)

	var todos []*todo.Todo
	err := db.
		Where("deleted_at IS NULL").
		Find(&todos).
		Error
	if err != nil {
		return nil, err
	}

	return todos, nil
}
