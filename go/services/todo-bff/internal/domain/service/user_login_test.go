package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway/mock"
	auth_model "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/output"
)

type PrepareUserLoginFields struct {
	ctx                   context.Context
	mockAuthReaderGateway *mock_gateway.MockAuthReaderGateway
}

type UserLoginArgs struct {
	ctx context.Context
	in  *input.UserLogin
}

type UserLoginTestcase struct {
	prepare  func(f *PrepareUserLoginFields)
	args     UserLoginArgs
	expected *output.UserLogin
	wantErr  bool
}

func Test_userLogin_Login(t *testing.T) {
	t.Parallel()

	testTables := map[string]UserLoginTestcase{
		"Login successfully": {
			prepare: func(f *PrepareUserLoginFields) {
				f.mockAuthReaderGateway.
					EXPECT().
					Login(f.ctx, "user@example.com", "password123").
					Return(&auth_model.Tokens{
						AccessToken: auth_model.AccessToken{
							Token:   "access-token",
							Expires: 3600,
						},
						RefreshToken: auth_model.RefreshToken{
							Token:   "refresh-token",
							Expires: 86400,
						},
					}, nil).
					Times(1)
			},
			args: UserLoginArgs{
				ctx: context.Background(),
				in: &input.UserLogin{
					Email:    "user@example.com",
					Password: "password123",
				},
			},
			expected: &output.UserLogin{
				Tokens: &auth_model.Tokens{
					AccessToken: auth_model.AccessToken{
						Token:   "access-token",
						Expires: 3600,
					},
					RefreshToken: auth_model.RefreshToken{
						Token:   "refresh-token",
						Expires: 86400,
					},
				},
			},
			wantErr: false,
		},
		"Fail to login when gateway returns error": {
			prepare: func(f *PrepareUserLoginFields) {
				f.mockAuthReaderGateway.
					EXPECT().
					Login(f.ctx, "user@example.com", "wrong-password").
					Return(nil, errors.New("gateway error")).
					Times(1)
			},
			args: UserLoginArgs{
				ctx: context.Background(),
				in: &input.UserLogin{
					Email:    "user@example.com",
					Password: "wrong-password",
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

			mockAuthReaderGateway := mock_gateway.NewMockAuthReaderGateway(ctrl)

			f := &PrepareUserLoginFields{
				ctx:                   context.Background(),
				mockAuthReaderGateway: mockAuthReaderGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			svc := service.NewUserLogin(mockAuthReaderGateway)

			actual, err := svc.Login(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("userLogin.Login() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("userLogin.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("userLogin.Login() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
