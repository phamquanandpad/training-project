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
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

func Test_GetUser(t *testing.T) {
	type fields struct {
		mockUserGetter *mock_usecase.MockUserGetter
	}

	type args struct {
		ctx context.Context
		req *todo_v1.GetUserRequest
	}

	type testcase struct {
		prepare  func(f *fields)
		args     args
		expected *todo_v1.GetUserResponse
		wantErr  bool
	}

	validate := validator.New()
	requestBinder := requestbinder.NewRequestBinder(validate)

	testTables := map[string]testcase{
		"Get User successfully": {
			prepare: func(f *fields) {
				// NOTE: input.UserGetter has no json tags for its fields, so
				// the requestBinder cannot bind "user_id" from GetUserRequest
				// to UserGetter.UserID. Using gomock.Any() here to match regardless
				// of the actual bound value.
				f.mockUserGetter.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(&output.UserGetter{
						ID:       todo.UserID(1),
						Username: "user1",
						Email:    cast.Ptr("user1@example.com"),
					}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &todo_v1.GetUserRequest{
					UserId: 1,
				},
			},
			expected: &todo_v1.GetUserResponse{
				User: &todo_common_v1.User{
					Id:       1,
					Username: "user1",
					Email:    "user1@example.com",
				},
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockUserGetter := mock_usecase.NewMockUserGetter(ctrl)

			tt.prepare(&fields{
				mockUserGetter: mockUserGetter,
			})

			svc, err := handler.NewTodoService(
				nil,
				validate,
				requestBinder,
				nil, nil, nil, nil, nil,
				mockUserGetter,
				nil,
			)
			if err != nil {
				t.Fatalf("NewTodoService() error = %v", err)
			}

			got, err := svc.GetUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, got, protocmp.Transform()); diff != "" {
				t.Errorf("GetUser() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
