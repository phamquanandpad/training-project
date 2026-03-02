package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway/mock"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

type PrepareGetUserFields struct {
	ctx                    context.Context
	mockBinder             *mock_gateway.MockBinder
	mockUserQueriesGateway *mock_gateway.MockUserQueriesGateway
}

type GetUserArgs struct {
	ctx context.Context
	in  *input.UserGetter
}

type GetUserTestcase struct {
	name     string
	prepare  func(f *PrepareGetUserFields)
	args     GetUserArgs
	expected *output.UserGetter
	wantErr  bool
}

func Test_userGetter_Get(t *testing.T) {
	t.Parallel()

	testTables := map[string]GetUserTestcase{
		"Get User successfully": {
			name: "Get User successfully",
			prepare: func(f *PrepareGetUserFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUser(f.ctx, todo.UserID(1)).
					Return(&todo.User{
						ID:       todo.UserID(1),
						Username: "user1",
						Email:    cast.Ptr("user1@example.com"),
					}, nil).
					Times(1)
			},
			args: GetUserArgs{
				ctx: context.Background(),
				in: &input.UserGetter{
					UserID: todo.UserID(1),
				},
			},
			expected: &output.UserGetter{
				ID:       todo.UserID(1),
				Username: "user1",
				Email:    cast.Ptr("user1@example.com"),
			},
			wantErr: false,
		},
		"User not found": {
			name: "User not found",
			prepare: func(f *PrepareGetUserFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUser(f.ctx, todo.UserID(999)).
					Return(nil, nil).
					Times(1)
			},
			args: GetUserArgs{
				ctx: context.Background(),
				in: &input.UserGetter{
					UserID: todo.UserID(999),
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Internal error when getting user": {
			name: "Internal error when getting user",
			prepare: func(f *PrepareGetUserFields) {
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockUserQueriesGateway.
					EXPECT().
					GetUser(f.ctx, todo.UserID(1)).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			args: GetUserArgs{
				ctx: context.Background(),
				in: &input.UserGetter{
					UserID: todo.UserID(1),
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
			mockUserQueriesGateway := mock_gateway.NewMockUserQueriesGateway(ctrl)

			userGetterService := service.NewUserGetter(
				mockBinder,
				mockUserQueriesGateway,
			)

			f := &PrepareGetUserFields{
				ctx:                    context.Background(),
				mockBinder:             mockBinder,
				mockUserQueriesGateway: mockUserQueriesGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			res, err := userGetterService.Get(f.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("userGetterService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(res, tt.expected); diff != "" {
				t.Errorf("GetUser result mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}
