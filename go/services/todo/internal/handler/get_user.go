package handler

import (
	"context"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"
	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

func (h *todoService) GetUser(
	ctx context.Context,
	req *todo_v1.GetUserRequest,
) (*todo_v1.GetUserResponse, error) {
	var in input.UserGetter

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	user, err := h.userGetter.Get(ctx, &in)
	if err != nil {
		return nil, err
	}

	return toGetUserResponse(user), nil
}

func toGetUserResponse(out *output.UserGetter) *todo_v1.GetUserResponse {
	return &todo_v1.GetUserResponse{
		User: &todo_common_v1.User{
			Id:       int64(out.ID),
			Username: out.Username,
			Email:    cast.Value(out.Email),
		},
	}
}
