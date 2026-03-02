package handler_test

import (
	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/testing/protocmp"

	auth_v1 "github.com/phamquanandpad/training-project/grpc/go/auth/v1"

	mock_usecase "github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/mock"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/handler"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/handler/requestbinder"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

func Test_VerifyToken(t *testing.T) {
	type fields struct {
		mockTokenVerify *mock_usecase.MockTokenVerify
	}

	type args struct {
		ctx context.Context
		req *auth_v1.VerifyTokenRequest
	}

	type testcase struct {
		prepare  func(f *fields)
		args     args
		expected *auth_v1.VerifyTokenResponse
		wantErr  bool
	}

	validate := validator.New()
	requestBinder := requestbinder.NewRequestBinder(validate)

	testTables := map[string]testcase{
		"Verify token successfully": {
			prepare: func(f *fields) {
				f.mockTokenVerify.
					EXPECT().
					VerifyToken(gomock.Any(), &input.TokenVerify{
						AccessToken: "access-token",
					}).
					Return(&output.TokenVerify{
						UserID: auth_models.UserID(1),
					}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.VerifyTokenRequest{
					AccessToken: "access-token",
				},
			},
			expected: &auth_v1.VerifyTokenResponse{
				UserId: 1,
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockTokenVerify := mock_usecase.NewMockTokenVerify(ctrl)

			tt.prepare(&fields{
				mockTokenVerify: mockTokenVerify,
			})

			svc, err := handler.NewAuthService(
				nil,
				validate,
				requestBinder,
				nil, nil,
				mockTokenVerify,
				nil,
			)
			if err != nil {
				t.Fatalf("NewAuthService() error = %v", err)
			}

			got, err := svc.VerifyToken(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, got,
				protocmp.Transform(),
			); diff != "" {
				t.Errorf("VerifyToken() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
