package datastore

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

type todoWriter struct{}

func NewTodoWriter() gateway.TodoCommandsGateway {
	return &todoWriter{}
}

func (w *todoWriter) CreateTodo(
	ctx context.Context,
	newTodo todo.NewTodo,
) (*todo.Todo, error) {
	tx, err := ExtractTodoDB(ctx)
	if err != nil {
		return nil, err
	}

	db := tx.WithContext(ctx)
	createdTodo := todo.Todo{
		UserID:      newTodo.UserID,
		Task:        newTodo.Task,
		Description: newTodo.Description,
		Status:      newTodo.Status,
	}

	if err := db.
		Create(&createdTodo).
		Error; err != nil {
		return nil, err
	}
	return &createdTodo, nil
}

func (w *todoWriter) UpdateTodo(
	ctx context.Context,
	todoID todo.TodoID,
	userID todo.UserID,
	updateTodo todo.UpdateTodo,
) (*todo.Todo, error) {
	tx, err := ExtractTodoDB(ctx)
	if err != nil {
		return nil, err
	}

	db := tx.WithContext(ctx)

	var t todo.Todo
	if err := db.
		Where("id = ? AND deleted_at IS NULL", todoID).
		Where("user_id = ?", userID).
		First(&t).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if updateTodo.Task != nil {
		t.Task = *updateTodo.Task
	}
	if updateTodo.Status != nil {
		t.Status = *updateTodo.Status
	}
	if updateTodo.Description != nil {
		t.Description = updateTodo.Description
	}

	if err := db.Save(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (w *todoWriter) SoftDeleteTodo(
	ctx context.Context,
	todoID todo.TodoID,
	userID todo.UserID,
) error {
	tx, err := ExtractTodoDB(ctx)
	if err != nil {
		return err
	}

	db := tx.WithContext(ctx)

	var t todo.Todo
	if err := db.
		Where("id = ? AND deleted_at IS NULL", todoID).
		Where("user_id = ?", userID).
		First(&t).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	t.DeletedAt = cast.Ptr(time.Now())
	if err := db.Save(&t).Error; err != nil {
		return err
	}
	return nil
}
