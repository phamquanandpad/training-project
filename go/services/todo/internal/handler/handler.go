package handler

import (
	"github.com/go-playground/validator/v10"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/config"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/requestbinder"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase"
)

type todoService struct {
	todo_v1.UnimplementedTodoServiceServer

	conf          *config.Config
	validate      *validator.Validate
	requestBinder *requestbinder.RequestBinder
	todoGetter    usecase.TodoGetter
	todoLister    usecase.TodoLister
	todoCreator   usecase.TodoCreator
	todoUpdater   usecase.TodoUpdater
	todoDeleter   usecase.TodoDeleter
	userGetter    usecase.UserGetter
	userCreator   usecase.UserCreator
}

func NewTodoService(
	conf *config.Config,
	validate *validator.Validate,
	requestBinder *requestbinder.RequestBinder,

	todoGetter usecase.TodoGetter,
	todoLister usecase.TodoLister,
	todoCreator usecase.TodoCreator,
	todoUpdater usecase.TodoUpdater,
	todoDeleter usecase.TodoDeleter,
	userGetter usecase.UserGetter,
	userCreator usecase.UserCreator,
) (todo_v1.TodoServiceServer, error) {
	return &todoService{
		conf:          conf,
		validate:      validate,
		requestBinder: requestBinder,
		todoGetter:    todoGetter,
		todoLister:    todoLister,
		todoCreator:   todoCreator,
		todoUpdater:   todoUpdater,
		todoDeleter:   todoDeleter,
		userGetter:    userGetter,
		userCreator:   userCreator,
	}, nil
}
