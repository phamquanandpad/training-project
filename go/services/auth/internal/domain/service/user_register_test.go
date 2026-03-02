package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway/mock"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
)

type PrepareRegisterFields struct {
	ctx                        context.Context
	mockBinder                 *mock_gateway.MockBinder
	mockUserCommandsGateway    *mock_gateway.MockUserCommandsGateway
	mockUserQueriesGateway     *mock_gateway.MockUserQueriesGateway
	mockUserSvcCommandsGateway *mock_gateway.MockUserServiceCommandsGateway
}

type RegisterArgs struct {
	ctx context.Context
	in  *input.UserRegister
}

type RegisterTestcase struct {
	name    string
	prepare func(f *PrepareRegisterFields)
	args    RegisterArgs
	wantErr bool
}

func Test_userRegister_Register(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testTables := map[string]RegisterTestcase{
		"Register successfully": {
			name: "Register successfully",
			prepare: func(f *PrepareRegisterFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByEmail(f.ctx, "newuser@example.com").
					Return(nil, nil).
					Times(1)

				f.mockUserCommandsGateway.
					EXPECT().
					CreateUser(f.ctx, gomock.Any()).
					Return(&auth_models.User{
						ID:        auth_models.UserID(10),
						Email:     "newuser@example.com",
						Password:  "hashedpassword",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil).
					Times(1)

				f.mockUserSvcCommandsGateway.
					EXPECT().
					CreateUser(f.ctx, gomock.Any()).
					Return(&todo.User{
						ID:        todo.UserID(10),
						Username:  "newuser",
						Email:     cast.Ptr("newuser@example.com"),
						CreatedAt: now,
						UpdatedAt: now,
					}, nil).
					Times(1)
			},
			args: RegisterArgs{
				ctx: context.Background(),
				in: &input.UserRegister{
					Username: "newuser",
					Email:    "newuser@example.com",
					Password: "Password123!",
				},
			},
			wantErr: false,
		},
		"Email already exists": {
			name: "Email already exists",
			prepare: func(f *PrepareRegisterFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByEmail(f.ctx, "existing@example.com").
					Return(&auth_models.User{
						ID:    auth_models.UserID(1),
						Email: "existing@example.com",
					}, nil).
					Times(1)
			},
			args: RegisterArgs{
				ctx: context.Background(),
				in: &input.UserRegister{
					Username: "existinguser",
					Email:    "existing@example.com",
					Password: "Password123!",
				},
			},
			wantErr: true,
		},
		"Internal error on GetUserByEmail": {
			name: "Internal error on GetUserByEmail",
			prepare: func(f *PrepareRegisterFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByEmail(f.ctx, "newuser@example.com").
					Return(nil, errors.New("database error")).
					Times(1)
			},
			args: RegisterArgs{
				ctx: context.Background(),
				in: &input.UserRegister{
					Username: "newuser",
					Email:    "newuser@example.com",
					Password: "Password123!",
				},
			},
			wantErr: true,
		},
		"Internal error on CreateUser (auth DB)": {
			name: "Internal error on CreateUser (auth DB)",
			prepare: func(f *PrepareRegisterFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByEmail(f.ctx, "newuser@example.com").
					Return(nil, nil).
					Times(1)

				f.mockUserCommandsGateway.
					EXPECT().
					CreateUser(f.ctx, gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			args: RegisterArgs{
				ctx: context.Background(),
				in: &input.UserRegister{
					Username: "newuser",
					Email:    "newuser@example.com",
					Password: "Password123!",
				},
			},
			wantErr: true,
		},
		"Internal error on CreateUser (todo service)": {
			name: "Internal error on CreateUser (todo service)",
			prepare: func(f *PrepareRegisterFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByEmail(f.ctx, "newuser@example.com").
					Return(nil, nil).
					Times(1)

				f.mockUserCommandsGateway.
					EXPECT().
					CreateUser(f.ctx, gomock.Any()).
					Return(&auth_models.User{
						ID:        auth_models.UserID(10),
						Email:     "newuser@example.com",
						Password:  "hashedpassword",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil).
					Times(1)

				f.mockUserSvcCommandsGateway.
					EXPECT().
					CreateUser(f.ctx, gomock.Any()).
					Return(nil, errors.New("grpc error")).
					Times(1)
			},
			args: RegisterArgs{
				ctx: context.Background(),
				in: &input.UserRegister{
					Username: "newuser",
					Email:    "newuser@example.com",
					Password: "Password123!",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range testTables {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockBinder := mock_gateway.NewMockBinder(ctrl)
			mockUserCommandsGateway := mock_gateway.NewMockUserCommandsGateway(ctrl)
			mockUserQueriesGateway := mock_gateway.NewMockUserQueriesGateway(ctrl)
			mockUserSvcCommandsGateway := mock_gateway.NewMockUserServiceCommandsGateway(ctrl)

			userRegisterService := service.NewUserRegister(
				mockBinder,
				mockUserCommandsGateway,
				mockUserQueriesGateway,
				mockUserSvcCommandsGateway,
			)

			f := &PrepareRegisterFields{
				ctx:                        context.Background(),
				mockBinder:                 mockBinder,
				mockUserCommandsGateway:    mockUserCommandsGateway,
				mockUserQueriesGateway:     mockUserQueriesGateway,
				mockUserSvcCommandsGateway: mockUserSvcCommandsGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			err := userRegisterService.Register(f.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRegisterService.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
