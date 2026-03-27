package handler

import (
	"github.com/go-playground/validator/v10"

	auth_v1 "github.com/phamquanandpad/training-project/grpc/go/auth/v1"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/config"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/handler/requestbinder"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase"
)

type authService struct {
	auth_v1.UnimplementedAuthServiceServer

	conf          *config.Config
	validate      *validator.Validate
	requestBinder *requestbinder.RequestBinder
	userLogin     usecase.UserLogin
	userRegister  usecase.UserRegister
	tokenVerify   usecase.TokenVerify
	tokenRefresh  usecase.TokenRefresh
}

func NewAuthService(
	conf *config.Config,
	validate *validator.Validate,
	requestBinder *requestbinder.RequestBinder,
	userLogin usecase.UserLogin,
	userRegister usecase.UserRegister,
	tokenVerify usecase.TokenVerify,
	tokenRefresh usecase.TokenRefresh,
) (auth_v1.AuthServiceServer, error) {
	return &authService{
		conf:          conf,
		validate:      validate,
		requestBinder: requestBinder,
		userLogin:     userLogin,
		userRegister:  userRegister,
		tokenVerify:   tokenVerify,
		tokenRefresh:  tokenRefresh,
	}, nil
}
