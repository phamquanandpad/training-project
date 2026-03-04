package usecase

import (
	"context"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

type TodoGetter interface {
	Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error)
}

type TodoLister interface {
	List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error)
}

type TodoCreator interface {
	Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error)
}

type TodoUpdater interface {
	Update(ctx context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error)
}

type TodoDeleter interface {
	SoftDelete(ctx context.Context, in *input.TodoDeleter) error
}

type UserGetter interface {
	Get(ctx context.Context, in *input.UserGetter) (*output.UserGetter, error)
}

type UserCreator interface {
	Create(ctx context.Context, in *input.UserCreator) (*output.UserCreator, error)
}

type DataValidator interface {
	ValidateUserRequest(
		ctx context.Context,
		in *input.UserRequestValidator,
	) (*output.UserRequestValidator, error)
}
