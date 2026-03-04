package handler_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/testing/protocmp"

	app_errors "github.com/phamquanandpad/training-project/go/services/auth/internal/errors"
	mock_usecase "github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/mock"
	auth_v1 "github.com/phamquanandpad/training-project/grpc/go/auth/v1"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/handler"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/handler/requestbinder"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

func Test_RefreshToken(t *testing.T) {
	type fields struct {
		mockTokenRefresh *mock_usecase.MockTokenRefresh
	}

	type args struct {
		ctx context.Context
		req *auth_v1.RefreshTokenRequest
	}

	type testcase struct {
		prepare  func(f *fields)
		args     args
		expected *auth_v1.RefreshTokenResponse
		wantErr  bool
	}

	validate := validator.New()
	requestBinder := requestbinder.NewRequestBinder(validate)
	accessTokenExpireDuration := int64(1 * time.Hour.Seconds())

	testTables := map[string]testcase{
		"Refresh token successfully": {
			prepare: func(f *fields) {
				f.mockTokenRefresh.
					EXPECT().
					RefreshToken(gomock.Any(), &input.TokenRefresh{
						RefreshToken: "refresh-token",
					}).
					Return(&output.TokenRefresh{
						AccessToken:               "new-access-token",
						AccessTokenExpireDuration: accessTokenExpireDuration,
					}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.RefreshTokenRequest{
					RefreshToken: "refresh-token",
				},
			},
			expected: &auth_v1.RefreshTokenResponse{
				AccessToken:               "new-access-token",
				AccessTokenExpireDuration: accessTokenExpireDuration,
			},
			wantErr: false,
		},
		"Missing refresh token returns validation error": {
			prepare: func(f *fields) {
				f.mockTokenRefresh.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).Times(0)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.RefreshTokenRequest{
					RefreshToken: "",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Usecase returns auth error": {
			prepare: func(f *fields) {
				f.mockTokenRefresh.
					EXPECT().
					RefreshToken(gomock.Any(), &input.TokenRefresh{
						RefreshToken: "expired-token",
					}).
					Return(nil, app_errors.NewAuthNError("token is expired or invalid", nil, nil)).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.RefreshTokenRequest{
					RefreshToken: "expired-token",
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

			mockTokenRefresh := mock_usecase.NewMockTokenRefresh(ctrl)

			tt.prepare(&fields{
				mockTokenRefresh: mockTokenRefresh,
			})

			svc, err := handler.NewAuthService(
				nil,
				validate,
				requestBinder,
				nil, nil, nil,
				mockTokenRefresh,
			)
			if err != nil {
				t.Fatalf("NewAuthService() error = %v", err)
			}

			got, err := svc.RefreshToken(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, got,
				protocmp.Transform(),
			); diff != "" {
				t.Errorf("RefreshToken() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
