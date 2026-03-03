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

type PrepareTokenVerifyFields struct {
	ctx                   context.Context
	mockAuthReaderGateway *mock_gateway.MockAuthReaderGateway
}

type TokenVerifyArgs struct {
	ctx context.Context
	in  *input.TokenVerify
}

type TokenVerifyTestcase struct {
	prepare  func(f *PrepareTokenVerifyFields)
	args     TokenVerifyArgs
	expected *output.TokenVerify
	wantErr  bool
}

func Test_tokenVerify_VerifyToken(t *testing.T) {
	t.Parallel()

	testTables := map[string]TokenVerifyTestcase{
		"Verify token successfully": {
			prepare: func(f *PrepareTokenVerifyFields) {
				f.mockAuthReaderGateway.
					EXPECT().
					VerifyToken(f.ctx, "valid-access-token").
					Return(auth_model.NewUserID(1), nil).
					Times(1)
			},
			args: TokenVerifyArgs{
				ctx: context.Background(),
				in: &input.TokenVerify{
					AccessToken: "valid-access-token",
				},
			},
			expected: &output.TokenVerify{
				UserID: auth_model.NewUserID(1),
			},
			wantErr: false,
		},
		"Fail to verify token when gateway returns error": {
			prepare: func(f *PrepareTokenVerifyFields) {
				f.mockAuthReaderGateway.
					EXPECT().
					VerifyToken(f.ctx, "invalid-access-token").
					Return(nil, errors.New("gateway error")).
					Times(1)
			},
			args: TokenVerifyArgs{
				ctx: context.Background(),
				in: &input.TokenVerify{
					AccessToken: "invalid-access-token",
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

			f := &PrepareTokenVerifyFields{
				ctx:                   context.Background(),
				mockAuthReaderGateway: mockAuthReaderGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			svc := service.NewTokenVerify(mockAuthReaderGateway)

			actual, err := svc.VerifyToken(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("tokenVerify.VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("tokenVerify.VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("tokenVerify.VerifyToken() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
