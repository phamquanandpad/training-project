package datastore

import (
	"context"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

type userWriter struct{}

func NewUserWriter() gateway.UserCommandsGateway {
	return &userWriter{}
}

func (w *userWriter) CreateUser(ctx context.Context, user todo.NewUser) (*todo.User, error) {
	tx, err := ExtractTodoDB(ctx)
	if err != nil {
		return nil, err
	}
	db := tx.WithContext(ctx)

	createdUser := todo.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	err = db.
		Create(&createdUser).
		Error
	if err != nil {
		return nil, err
	}
	return &createdUser, nil
}
