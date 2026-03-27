package handler_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/testing/protocmp"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"
	app_errors "github.com/phamquanandpad/training-project/go/services/auth/internal/errors"
	mock_usecase "github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/mock"
	auth_v1 "github.com/phamquanandpad/training-project/grpc/go/auth/v1"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/handler"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/handler/requestbinder"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

func Test_Login(t *testing.T) {
	type fields struct {
		mockUserLogin *mock_usecase.MockUserLogin
	}

	type args struct {
		ctx context.Context
		req *auth_v1.LoginRequest
	}

	type testcase struct {
		prepare  func(f *fields)
		args     args
		expected *auth_v1.LoginResponse
		wantErr  bool
	}

	validate := validator.New()
	requestBinder := requestbinder.NewRequestBinder(validate)
	accessTokenExpireDuration := int64(1 * time.Hour.Seconds())
	refreshTokenExpireDuration := int64(24 * time.Hour.Seconds())

	testTables := map[string]testcase{
		"Login successfully": {
			prepare: func(f *fields) {
				f.mockUserLogin.
					EXPECT().
					Login(gomock.Any(), &input.UserLogin{
						Email:    "test@example.com",
						Password: "password123",
					}).
					Return(&output.UserLogin{
						UserID:                     auth_models.UserID(1),
						AccessToken:                "access-token",
						RefreshToken:               "refresh-token",
						AccessTokenExpireDuration:  accessTokenExpireDuration,
						RefreshTokenExpireDuration: refreshTokenExpireDuration,
					}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.LoginRequest{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			expected: &auth_v1.LoginResponse{
				AccessToken:                "access-token",
				RefreshToken:               "refresh-token",
				AccessTokenExpireDuration:  accessTokenExpireDuration,
				RefreshTokenExpireDuration: refreshTokenExpireDuration,
			},
			wantErr: false,
		},
		"Missing email returns validation error": {
			prepare: func(f *fields) {
				f.mockUserLogin.EXPECT().Login(gomock.Any(), gomock.Any()).Times(0)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.LoginRequest{
					Email:    "",
					Password: "password123",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Invalid email format returns validation error": {
			prepare: func(f *fields) {
				f.mockUserLogin.EXPECT().Login(gomock.Any(), gomock.Any()).Times(0)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.LoginRequest{
					Email:    "not-an-email",
					Password: "password123",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Missing password returns validation error": {
			prepare: func(f *fields) {
				f.mockUserLogin.EXPECT().Login(gomock.Any(), gomock.Any()).Times(0)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.LoginRequest{
					Email:    "test@example.com",
					Password: "",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Usecase returns error": {
			prepare: func(f *fields) {
				f.mockUserLogin.
					EXPECT().
					Login(gomock.Any(), &input.UserLogin{
						Email:    "test@example.com",
						Password: "password123",
					}).
					Return(nil, app_errors.NewAuthNError("invalid credentials", nil, nil)).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.LoginRequest{
					Email:    "test@example.com",
					Password: "password123",
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

			mockUserLogin := mock_usecase.NewMockUserLogin(ctrl)

			tt.prepare(&fields{
				mockUserLogin: mockUserLogin,
			})

			svc, err := handler.NewAuthService(
				nil,
				validate,
				requestBinder,
				mockUserLogin,
				nil, nil, nil,
			)
			if err != nil {
				t.Fatalf("NewAuthService() error = %v", err)
			}

			got, err := svc.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, got,
				protocmp.Transform(),
			); diff != "" {
				t.Errorf("Login() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
