package handler_test

import (
	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/testing/protocmp"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"
	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	mock_usecase "github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/mock"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/requestbinder"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

func Test_PostUser(t *testing.T) {
	type fields struct {
		mockUserCreator *mock_usecase.MockUserCreator
	}

	type args struct {
		ctx context.Context
		req *todo_v1.PostUserRequest
	}

	type testcase struct {
		prepare  func(f *fields)
		args     args
		expected *todo_v1.PostUserResponse
		wantErr  bool
	}

	validate := validator.New()
	requestBinder := requestbinder.NewRequestBinder(validate)

	testTables := map[string]testcase{
		"Post User successfully": {
			prepare: func(f *fields) {
				f.mockUserCreator.
					EXPECT().
					Create(gomock.Any(), &input.UserCreator{
						User: todo.User{
							ID:       todo.UserID(1),
							Username: "user1",
							Email:    cast.Ptr("user1@example.com"),
						},
					}).
					Return(&output.UserCreator{
						ID:       todo.UserID(1),
						Username: "user1",
						Email:    cast.Ptr("user1@example.com"),
					}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &todo_v1.PostUserRequest{
					User: &todo_common_v1.User{
						Id:       1,
						Username: "user1",
						Email:    "user1@example.com",
					},
				},
			},
			expected: &todo_v1.PostUserResponse{},
			wantErr:  false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockUserCreator := mock_usecase.NewMockUserCreator(ctrl)

			tt.prepare(&fields{
				mockUserCreator: mockUserCreator,
			})

			svc, err := handler.NewTodoService(
				nil,
				validate,
				requestBinder,
				nil, nil, nil, nil, nil, nil,
				mockUserCreator,
			)
			if err != nil {
				t.Fatalf("NewTodoService() error = %v", err)
			}

			got, err := svc.PostUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, got, protocmp.Transform()); diff != "" {
				t.Errorf("PostUser() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
