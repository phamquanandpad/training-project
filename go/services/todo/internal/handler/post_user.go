package handler

import (
	"context"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
)

func (h *todoService) PostUser(
	ctx context.Context,
	req *todo_v1.PostUserRequest,
) (*todo_v1.PostUserResponse, error) {
	var in input.UserCreator

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	_, err := h.userCreator.Create(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &todo_v1.PostUserResponse{}, nil
}

// func toPostUserInput(in *todo_v1.PostUserRequest) *input.UserCreator {
// 	userInput := in.GetUser()
// 	return &input.UserCreator{
// 		ID:        todo.UserID(userInput.GetId()),
// 		Username:  userInput.GetUsername(),
// 		Email:     cast.Ptr(userInput.GetEmail()),
// 		CreatedAt: userInput.GetCreatedAt().AsTime(),
// 		UpdatedAt: userInput.GetUpdatedAt().AsTime(),
// 	}
// }
