package handler_test

import (
	"context"
	"testing"

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
)

func Test_Register(t *testing.T) {
	type fields struct {
		mockUserRegister *mock_usecase.MockUserRegister
	}

	type args struct {
		ctx context.Context
		req *auth_v1.RegisterRequest
	}

	type testcase struct {
		prepare  func(f *fields)
		args     args
		expected *auth_v1.RegisterResponse
		wantErr  bool
	}

	validate := validator.New()
	requestBinder := requestbinder.NewRequestBinder(validate)

	testTables := map[string]testcase{
		"Register successfully": {
			prepare: func(f *fields) {
				f.mockUserRegister.
					EXPECT().
					Register(gomock.Any(), &input.UserRegister{
						Username: "testuser",
						Email:    "test@example.com",
						Password: "password123",
					}).
					Return(nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.RegisterRequest{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			expected: &auth_v1.RegisterResponse{},
			wantErr:  false,
		},
		"Missing username returns validation error": {
			prepare: func(f *fields) {
				f.mockUserRegister.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.RegisterRequest{
					Username: "",
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Missing email returns validation error": {
			prepare: func(f *fields) {
				f.mockUserRegister.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.RegisterRequest{
					Username: "testuser",
					Email:    "",
					Password: "password123",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Invalid email format returns validation error": {
			prepare: func(f *fields) {
				f.mockUserRegister.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.RegisterRequest{
					Username: "testuser",
					Email:    "not-an-email",
					Password: "password123",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Password too short returns validation error": {
			prepare: func(f *fields) {
				f.mockUserRegister.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.RegisterRequest{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "abc",
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Usecase returns error": {
			prepare: func(f *fields) {
				f.mockUserRegister.
					EXPECT().
					Register(gomock.Any(), &input.UserRegister{
						Username: "testuser",
						Email:    "existing@example.com",
						Password: "password123",
					}).
					Return(app_errors.NewAlreadyExistsError("email already registered", nil, nil)).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &auth_v1.RegisterRequest{
					Username: "testuser",
					Email:    "existing@example.com",
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

			mockUserRegister := mock_usecase.NewMockUserRegister(ctrl)

			tt.prepare(&fields{
				mockUserRegister: mockUserRegister,
			})

			svc, err := handler.NewAuthService(
				nil,
				validate,
				requestBinder,
				nil,
				mockUserRegister,
				nil, nil,
			)
			if err != nil {
				t.Fatalf("NewAuthService() error = %v", err)
			}

			got, err := svc.Register(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, got,
				protocmp.Transform(),
			); diff != "" {
				t.Errorf("Register() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
