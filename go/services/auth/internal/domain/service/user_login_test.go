package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway/mock"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

type PrepareLoginFields struct {
	ctx                    context.Context
	mockBinder             *mock_gateway.MockBinder
	mockUserQueriesGateway *mock_gateway.MockUserQueriesGateway
	mockJwtGenerateGateway *mock_gateway.MockJwtGenerateGateway
}

type LoginArgs struct {
	ctx context.Context
	in  *input.UserLogin
}

type LoginTestcase struct {
	name     string
	prepare  func(f *PrepareLoginFields)
	args     LoginArgs
	expected *output.UserLogin
	wantErr  bool
}

func Test_userLogin_Login(t *testing.T) {
	t.Parallel()

	now := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

	existingUser := &auth_models.User{
		ID:        auth_models.UserID(1),
		Email:     "user1@example.com",
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	testTables := map[string]LoginTestcase{
		"Login successfully": {
			name: "Login successfully",
			prepare: func(f *PrepareLoginFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByEmail(f.ctx, "user1@example.com").
					Return(existingUser, nil).
					Times(1)

				f.mockJwtGenerateGateway.
					EXPECT().
					GenerateAccessToken(auth_models.UserID(1)).
					Return("access_token_value", int64(900), nil).
					Times(1)

				f.mockJwtGenerateGateway.
					EXPECT().
					GenerateRefreshToken(auth_models.UserID(1)).
					Return("refresh_token_value", int64(86400), nil).
					Times(1)
			},
			args: LoginArgs{
				ctx: context.Background(),
				in: &input.UserLogin{
					Email:    "user1@example.com",
					Password: "Password123!",
				},
			},
			expected: &output.UserLogin{
				UserID:                    auth_models.UserID(1),
				AccessToken:               "access_token_value",
				AccessTokenExpiresSecond:  900,
				RefreshToken:              "refresh_token_value",
				RefreshTokenExpiresSecond: 86400,
			},
			wantErr: false,
		},
		"User not found": {
			name: "User not found",
			prepare: func(f *PrepareLoginFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByEmail(f.ctx, "notfound@example.com").
					Return(nil, nil).
					Times(1)
			},
			args: LoginArgs{
				ctx: context.Background(),
				in: &input.UserLogin{
					Email:    "notfound@example.com",
					Password: "Password123!",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Wrong password": {
			name: "Wrong password",
			prepare: func(f *PrepareLoginFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByEmail(f.ctx, "user1@example.com").
					Return(existingUser, nil).
					Times(1)
			},
			args: LoginArgs{
				ctx: context.Background(),
				in: &input.UserLogin{
					Email:    "user1@example.com",
					Password: "WrongPassword!",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Internal error on GetUserByEmail": {
			name: "Internal error on GetUserByEmail",
			prepare: func(f *PrepareLoginFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByEmail(f.ctx, "user1@example.com").
					Return(nil, errors.New("database error")).
					Times(1)
			},
			args: LoginArgs{
				ctx: context.Background(),
				in: &input.UserLogin{
					Email:    "user1@example.com",
					Password: "Password123!",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Internal error on GenerateAccessToken": {
			name: "Internal error on GenerateAccessToken",
			prepare: func(f *PrepareLoginFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByEmail(f.ctx, "user1@example.com").
					Return(existingUser, nil).
					Times(1)

				f.mockJwtGenerateGateway.
					EXPECT().
					GenerateAccessToken(auth_models.UserID(1)).
					Return("", int64(0), errors.New("jwt error")).
					Times(1)
			},
			args: LoginArgs{
				ctx: context.Background(),
				in: &input.UserLogin{
					Email:    "user1@example.com",
					Password: "Password123!",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Internal error on GenerateRefreshToken": {
			name: "Internal error on GenerateRefreshToken",
			prepare: func(f *PrepareLoginFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByEmail(f.ctx, "user1@example.com").
					Return(existingUser, nil).
					Times(1)

				f.mockJwtGenerateGateway.
					EXPECT().
					GenerateAccessToken(auth_models.UserID(1)).
					Return("access_token_value", int64(900), nil).
					Times(1)

				f.mockJwtGenerateGateway.
					EXPECT().
					GenerateRefreshToken(auth_models.UserID(1)).
					Return("", int64(0), errors.New("jwt error")).
					Times(1)
			},
			args: LoginArgs{
				ctx: context.Background(),
				in: &input.UserLogin{
					Email:    "user1@example.com",
					Password: "Password123!",
				},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range testTables {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockBinder := mock_gateway.NewMockBinder(ctrl)
			mockUserQueriesGateway := mock_gateway.NewMockUserQueriesGateway(ctrl)
			mockJwtGenerateGateway := mock_gateway.NewMockJwtGenerateGateway(ctrl)

			userLoginService := service.NewUserLogin(
				mockBinder,
				mockUserQueriesGateway,
				mockJwtGenerateGateway,
			)

			f := &PrepareLoginFields{
				ctx:                    context.Background(),
				mockBinder:             mockBinder,
				mockUserQueriesGateway: mockUserQueriesGateway,
				mockJwtGenerateGateway: mockJwtGenerateGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			res, err := userLoginService.Login(f.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("userLoginService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(res, tt.expected); diff != "" {
				t.Errorf("Login result mismatch (-actual +expected):\n%s", diff)
			}
		})
	}
}
