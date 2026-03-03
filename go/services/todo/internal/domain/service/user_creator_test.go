package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway/mock"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

type PrepareCreateUserFields struct {
	ctx                     context.Context
	mockBinder              *mock_gateway.MockBinder
	mockUserCommandsGateway *mock_gateway.MockUserCommandsGateway
}

type CreateUserArgs struct {
	ctx context.Context
	in  *input.UserCreator
}

type CreateUserTestcase struct {
	prepare  func(f *PrepareCreateUserFields)
	args     CreateUserArgs
	expected *output.UserCreator
	wantErr  bool
}

func Test_userCreator_Create(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testTables := map[string]CreateUserTestcase{
		"Create User successfully": {
			prepare: func(f *PrepareCreateUserFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserCommandsGateway.
					EXPECT().
					CreateUser(f.ctx, todo.NewUser{
						ID:        todo.UserID(10),
						Username:  "user10",
						Email:     cast.Ptr("user10@example.com"),
						CreatedAt: now,
						UpdatedAt: now,
					}).
					Return(&todo.User{
						ID:       todo.UserID(10),
						Username: "user10",
						Email:    cast.Ptr("user10@example.com"),
					}, nil).
					Times(1)
			},
			args: CreateUserArgs{
				ctx: context.Background(),
				in: &input.UserCreator{
					User: todo.User{
						ID:        todo.UserID(10),
						Username:  "user10",
						Email:     cast.Ptr("user10@example.com"),
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			expected: &output.UserCreator{
				ID:       todo.UserID(10),
				Username: "user10",
				Email:    cast.Ptr("user10@example.com"),
			},
			wantErr: false,
		},
		"Internal error when creating user": {
			prepare: func(f *PrepareCreateUserFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserCommandsGateway.
					EXPECT().
					CreateUser(f.ctx, gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			args: CreateUserArgs{
				ctx: context.Background(),
				in: &input.UserCreator{
					User: todo.User{
						ID:        todo.UserID(1),
						Username:  "user1",
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockBinder := mock_gateway.NewMockBinder(ctrl)
			mockUserCommandsGateway := mock_gateway.NewMockUserCommandsGateway(ctrl)

			userCreatorService := service.NewUserCreator(
				mockBinder,
				mockUserCommandsGateway,
			)

			f := &PrepareCreateUserFields{
				ctx:                     context.Background(),
				mockBinder:              mockBinder,
				mockUserCommandsGateway: mockUserCommandsGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			res, err := userCreatorService.Create(f.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("userCreatorService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(res, tt.expected); diff != "" {
				t.Errorf("CreateUser result mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}
