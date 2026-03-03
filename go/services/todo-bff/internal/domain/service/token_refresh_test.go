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

type PrepareTokenRefreshFields struct {
	ctx                   context.Context
	mockAuthReaderGateway *mock_gateway.MockAuthReaderGateway
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

	testTables := map[string]TokenRefreshTestcase{
		"Refresh token successfully": {
			prepare: func(f *PrepareTokenRefreshFields) {
				f.mockAuthReaderGateway.
					EXPECT().
					RefreshToken(f.ctx, "refresh-token-value").
					Return(&auth_model.AccessToken{
						Token:   "new-access-token",
						Expires: 3600,
					}, nil).
					Times(1)
			},
			args: TokenRefreshArgs{
				ctx: context.Background(),
				in: &input.TokenRefresh{
					RefreshToken: "refresh-token-value",
				},
			},
			expected: &output.TokenRefresh{
				AccessToken: &auth_model.AccessToken{
					Token:   "new-access-token",
					Expires: 3600,
				},
			},
			wantErr: false,
		},
		"Fail to refresh token when gateway returns error": {
			prepare: func(f *PrepareTokenRefreshFields) {
				f.mockAuthReaderGateway.
					EXPECT().
					RefreshToken(f.ctx, "invalid-refresh-token").
					Return(nil, errors.New("gateway error")).
					Times(1)
			},
			args: TokenRefreshArgs{
				ctx: context.Background(),
				in: &input.TokenRefresh{
					RefreshToken: "invalid-refresh-token",
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

			f := &PrepareTokenRefreshFields{
				ctx:                   context.Background(),
				mockAuthReaderGateway: mockAuthReaderGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			svc := service.NewTokenRefresh(mockAuthReaderGateway)

			actual, err := svc.RefreshToken(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("tokenRefresh.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("tokenRefresh.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("tokenRefresh.RefreshToken() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
