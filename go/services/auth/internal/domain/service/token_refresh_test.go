package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway/mock"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

type PrepareTokenRefreshFields struct {
	ctx                    context.Context
	mockBinder             *mock_gateway.MockBinder
	mockJwtGenerateGateway *mock_gateway.MockJwtGenerateGateway
	mockJwtVerifyGateway   *mock_gateway.MockJwtVerifyGateway
	mockUserQueriesGateway *mock_gateway.MockUserQueriesGateway
}

type TokenRefreshArgs struct {
	ctx context.Context
	in  *input.TokenRefresh
}

type TokenRefreshTestcase struct {
	prepare  func(f *PrepareTokenRefreshFields)
	args     TokenRefreshArgs
	expected *output.TokenRefresh
	wantErr  bool
}

func Test_tokenRefresh_RefreshToken(t *testing.T) {
	t.Parallel()

	now := time.Now()
	accessTokenExpireDuration := int64(1 * time.Hour.Seconds())

	existingUser := &auth_models.User{
		ID:        auth_models.UserID(1),
		Email:     "user1@example.com",
		Password:  "hashedpassword",
		CreatedAt: now,
		UpdatedAt: now,
	}

	testTables := map[string]TokenRefreshTestcase{
		"Refresh token successfully": {
			prepare: func(f *PrepareTokenRefreshFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockJwtVerifyGateway.
					EXPECT().
					VerifyRefreshToken("valid_refresh_token").
					Return(auth_models.UserID(1), nil).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByID(f.ctx, auth_models.UserID(1)).
					Return(existingUser, nil).
					Times(1)

				f.mockJwtGenerateGateway.
					EXPECT().
					GenerateAccessToken(auth_models.UserID(1)).
					Return("new_access_token", accessTokenExpireDuration, nil).
					Times(1)
			},
			args: TokenRefreshArgs{
				ctx: context.Background(),
				in: &input.TokenRefresh{
					RefreshToken: "valid_refresh_token",
				},
			},
			expected: &output.TokenRefresh{
				AccessToken:               "new_access_token",
				AccessTokenExpireDuration: accessTokenExpireDuration,
			},
			wantErr: false,
		},
		"Invalid refresh token": {
			prepare: func(f *PrepareTokenRefreshFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockJwtVerifyGateway.
					EXPECT().
					VerifyRefreshToken("invalid_refresh_token").
					Return(auth_models.UserID(0), errors.New("token is invalid")).
					Times(1)
			},
			args: TokenRefreshArgs{
				ctx: context.Background(),
				in: &input.TokenRefresh{
					RefreshToken: "invalid_refresh_token",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"User not found after token verification": {
			prepare: func(f *PrepareTokenRefreshFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockJwtVerifyGateway.
					EXPECT().
					VerifyRefreshToken("valid_refresh_token").
					Return(auth_models.UserID(999), nil).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByID(f.ctx, auth_models.UserID(999)).
					Return(nil, nil).
					Times(1)
			},
			args: TokenRefreshArgs{
				ctx: context.Background(),
				in: &input.TokenRefresh{
					RefreshToken: "valid_refresh_token",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Internal error on GetUserByID": {
			prepare: func(f *PrepareTokenRefreshFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockJwtVerifyGateway.
					EXPECT().
					VerifyRefreshToken("valid_refresh_token").
					Return(auth_models.UserID(1), nil).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByID(f.ctx, auth_models.UserID(1)).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			args: TokenRefreshArgs{
				ctx: context.Background(),
				in: &input.TokenRefresh{
					RefreshToken: "valid_refresh_token",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Internal error on GenerateAccessToken": {
			prepare: func(f *PrepareTokenRefreshFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockJwtVerifyGateway.
					EXPECT().
					VerifyRefreshToken("valid_refresh_token").
					Return(auth_models.UserID(1), nil).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByID(f.ctx, auth_models.UserID(1)).
					Return(existingUser, nil).
					Times(1)

				f.mockJwtGenerateGateway.
					EXPECT().
					GenerateAccessToken(auth_models.UserID(1)).
					Return("", int64(0), errors.New("jwt error")).
					Times(1)
			},
			args: TokenRefreshArgs{
				ctx: context.Background(),
				in: &input.TokenRefresh{
					RefreshToken: "valid_refresh_token",
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
			mockJwtGenerateGateway := mock_gateway.NewMockJwtGenerateGateway(ctrl)
			mockJwtVerifyGateway := mock_gateway.NewMockJwtVerifyGateway(ctrl)
			mockUserQueriesGateway := mock_gateway.NewMockUserQueriesGateway(ctrl)

			tokenRefreshService := service.NewTokenRefresh(
				mockBinder,
				mockJwtGenerateGateway,
				mockJwtVerifyGateway,
				mockUserQueriesGateway,
			)

			f := &PrepareTokenRefreshFields{
				ctx:                    context.Background(),
				mockBinder:             mockBinder,
				mockJwtGenerateGateway: mockJwtGenerateGateway,
				mockJwtVerifyGateway:   mockJwtVerifyGateway,
				mockUserQueriesGateway: mockUserQueriesGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			res, err := tokenRefreshService.RefreshToken(f.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokenRefreshService.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(res, tt.expected); diff != "" {
				t.Errorf("RefreshToken result mismatch (-actual +expected):\n%s", diff)
			}
		})
	}
}
