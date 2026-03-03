package handler

import (
	"context"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/mapper"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
)

func (h *todoService) GetUser(
	ctx context.Context,
	req *todo_v1.GetUserRequest,
) (*todo_v1.GetUserResponse, error) {
	var in input.UserGetter

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	out, err := h.userGetter.Get(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &todo_v1.GetUserResponse{
		User: mapper.ToUserGRPCResponse((*todo.User)(out)),
	}, nil
}
