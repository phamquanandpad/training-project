package datavalidator

import (
	"context"
	"encoding/json"
	"strings"

	"google.golang.org/grpc"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
)

type reqData struct {
	UserAttributes todo_v1.UserAttributes `json:"user_attributes"`
}

// TODO Please reconsider whether to disable 'cyclp' and 'gocognit'.
//
//nolint:cyclop,gocognit
func UnaryServerInterceptor(dataValidator usecase.DataValidator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		methodName, _ := grpc.Method(ctx)
		// Check if method name is in whitelist, skip permission checker
		if isSkipDataValidator(methodName) {
			return handler(ctx, req)
		}

		byteData, err := json.Marshal(req)
		if err != nil {
			return nil, app_errors.NewInternalError(
				"dataValidator.json.Marshal",
				err,
			)
		}

		reqData := &reqData{}
		err = json.Unmarshal(byteData, &reqData)
		if err != nil {
			return nil, app_errors.NewInternalError(
				"dataValidator.json.Unmarshal",
				err,
			)
		}

		userID := todo.UserID(reqData.UserAttributes.UserId)

		if userID.Int64() != 0 {
			validatedUser, err := dataValidator.ValidateUserRequest(ctx, &input.UserRequestValidator{UserID: userID})

			if err != nil {
				return nil, err
			}

			if validatedUser != nil {
				ctx = todo.WithUser(ctx, (*todo.User)(validatedUser))
			}
		}

		return handler(ctx, req)
	}
}

func isSkipDataValidator(grpcMethodName string) bool {
	// Skip check permission for requests from user service
	if strings.HasPrefix(grpcMethodName, "/todo.todo.v1.TodoService/PostUser") {
		return true
	}

	// Skip check permission for health check requests
	if strings.HasPrefix(grpcMethodName, "/grpc.health.v1.Health/") {
		return true
	}

	return false
}
