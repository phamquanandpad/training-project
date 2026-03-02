package service_test

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway/mock"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
)

type PrepareUserRegisterFields struct {
	ctx                   context.Context
	mockAuthReaderGateway *mock_gateway.MockAuthReaderGateway
}

type UserRegisterArgs struct {
	ctx context.Context
	in  *input.UserRegister
}

type UserRegisterTestcase struct {
	prepare func(f *PrepareUserRegisterFields)
	args    UserRegisterArgs
	wantErr bool
}

func Test_userRegister_Register(t *testing.T) {
	t.Parallel()

	testTables := map[string]UserRegisterTestcase{
		"Register successfully": {
			prepare: func(f *PrepareUserRegisterFields) {
				f.mockAuthReaderGateway.
					EXPECT().
					Register(f.ctx, "testuser", "user@example.com", "password123").
					Return(nil).
					Times(1)
			},
			args: UserRegisterArgs{
				ctx: context.Background(),
				in: &input.UserRegister{
					Username: "testuser",
					Email:    "user@example.com",
					Password: "password123",
				},
			},
			wantErr: false,
		},
		"Fail to register when gateway returns error": {
			prepare: func(f *PrepareUserRegisterFields) {
				f.mockAuthReaderGateway.
					EXPECT().
					Register(f.ctx, "testuser", "existing@example.com", "password123").
					Return(errors.New("gateway error")).
					Times(1)
			},
			args: UserRegisterArgs{
				ctx: context.Background(),
				in: &input.UserRegister{
					Username: "testuser",
					Email:    "existing@example.com",
					Password: "password123",
				},
			},
			wantErr: true,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockAuthReaderGateway := mock_gateway.NewMockAuthReaderGateway(ctrl)

			f := &PrepareUserRegisterFields{
				ctx:                   context.Background(),
				mockAuthReaderGateway: mockAuthReaderGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			svc := service.NewUserRegister(mockAuthReaderGateway)

			err := svc.Register(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("userRegister.Register() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("userRegister.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
