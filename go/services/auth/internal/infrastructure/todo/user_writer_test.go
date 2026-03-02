package todo_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"
	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"
	mock_todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1/mock"

	todo_service "github.com/phamquanandpad/training-project/go/services/auth/internal/infrastructure/todo"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/todo"
)

func Test_userWriter_CreateUser(t *testing.T) {
	type fields struct {
		mockTodoServiceClient *mock_todo_v1.MockTodoServiceClient
	}

	type args struct {
		ctx     context.Context
		newUser todo.NewUser
	}

	type expected struct {
		user *todo.User
		err  error
	}

	testTables := map[string]struct {
		prepare func(f *fields)
		args    args
		expect  expected
		wantErr bool
	}{
		"Create User successfully": {
			prepare: func(f *fields) {
				f.mockTodoServiceClient.
					EXPECT().
					PostUser(gomock.Any(),
						&todo_v1.PostUserRequest{
							User: &todo_common_v1.User{
								Id:       10,
								Username: "user10",
								Email:    "user10@example.com",
							},
						},
					).
					Return(&todo_v1.PostUserResponse{}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				newUser: todo.NewUser{
					ID:       todo.UserID(10),
					Username: "user10",
					Email:    cast.Ptr("user10@example.com"),
				},
			},
			expect: expected{
				user: &todo.User{
					ID:       todo.UserID(10),
					Username: "user10",
					Email:    cast.Ptr("user10@example.com"),
				},
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			f := &fields{
				mockTodoServiceClient: mock_todo_v1.NewMockTodoServiceClient(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			userWriter := todo_service.NewUserWriter(f.mockTodoServiceClient)

			got, err := userWriter.CreateUser(tt.args.ctx, tt.args.newUser)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				if diff := cmp.Diff(tt.expect.user, got); diff != "" {
					t.Errorf("CreateUser() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
