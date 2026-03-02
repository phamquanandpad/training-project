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

type PrepareTokenVerifyFields struct {
	ctx                    context.Context
	mockBinder             *mock_gateway.MockBinder
	mockJwtVerifyGateway   *mock_gateway.MockJwtVerifyGateway
	mockUserQueriesGateway *mock_gateway.MockUserQueriesGateway
}

type TokenVerifyArgs struct {
	ctx context.Context
	in  *input.TokenVerify
}

type TokenVerifyTestcase struct {
	name     string
	prepare  func(f *PrepareTokenVerifyFields)
	args     TokenVerifyArgs
	expected *output.TokenVerify
	wantErr  bool
}

func Test_tokenVerify_VerifyToken(t *testing.T) {
	t.Parallel()

	now := time.Now()

	existingUser := &auth_models.User{
		ID:        auth_models.UserID(1),
		Email:     "user1@example.com",
		Password:  "hashedpassword",
		CreatedAt: now,
		UpdatedAt: now,
	}

	testTables := map[string]TokenVerifyTestcase{
		"Verify token successfully": {
			name: "Verify token successfully",
			prepare: func(f *PrepareTokenVerifyFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockJwtVerifyGateway.
					EXPECT().
					VerifyAccessToken("valid_access_token").
					Return(auth_models.UserID(1), nil).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByID(f.ctx, auth_models.UserID(1)).
					Return(existingUser, nil).
					Times(1)
			},
			args: TokenVerifyArgs{
				ctx: context.Background(),
				in: &input.TokenVerify{
					AccessToken: "valid_access_token",
				},
			},
			expected: &output.TokenVerify{
				UserID: auth_models.UserID(1),
			},
			wantErr: false,
		},
		"Invalid access token": {
			prepare: func(f *PrepareTokenVerifyFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockJwtVerifyGateway.
					EXPECT().
					VerifyAccessToken("invalid_access_token").
					Return(auth_models.UserID(0), errors.New("token is invalid")).
					Times(1)
			},
			args: TokenVerifyArgs{
				ctx: context.Background(),
				in: &input.TokenVerify{
					AccessToken: "invalid_access_token",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"User not found after token verification": {
			prepare: func(f *PrepareTokenVerifyFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockJwtVerifyGateway.
					EXPECT().
					VerifyAccessToken("valid_access_token").
					Return(auth_models.UserID(999), nil).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByID(f.ctx, auth_models.UserID(999)).
					Return(nil, nil).
					Times(1)
			},
			args: TokenVerifyArgs{
				ctx: context.Background(),
				in: &input.TokenVerify{
					AccessToken: "valid_access_token",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Internal error on GetUserByID": {
			prepare: func(f *PrepareTokenVerifyFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockJwtVerifyGateway.
					EXPECT().
					VerifyAccessToken("valid_access_token").
					Return(auth_models.UserID(1), nil).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUserByID(f.ctx, auth_models.UserID(1)).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			args: TokenVerifyArgs{
				ctx: context.Background(),
				in: &input.TokenVerify{
					AccessToken: "valid_access_token",
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
			mockJwtVerifyGateway := mock_gateway.NewMockJwtVerifyGateway(ctrl)
			mockUserQueriesGateway := mock_gateway.NewMockUserQueriesGateway(ctrl)

			tokenVerifyService := service.NewTokenVerify(
				mockBinder,
				mockJwtVerifyGateway,
				mockUserQueriesGateway,
			)

			f := &PrepareTokenVerifyFields{
				ctx:                    context.Background(),
				mockBinder:             mockBinder,
				mockJwtVerifyGateway:   mockJwtVerifyGateway,
				mockUserQueriesGateway: mockUserQueriesGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			res, err := tokenVerifyService.VerifyToken(f.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokenVerifyService.VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(res, tt.expected); diff != "" {
				t.Errorf("VerifyToken result mismatch (-actual +expected):\n%s", diff)
			}
		})
	}
}
