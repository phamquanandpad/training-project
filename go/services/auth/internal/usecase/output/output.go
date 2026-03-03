package output

import (
	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"
)

type UserLogin struct {
	UserID                     auth_models.UserID
	AccessToken                string
	RefreshToken               string
	AccessTokenExpireDuration  int64
	RefreshTokenExpireDuration int64
}

type UserRegister struct{}

type TokenRefresh struct {
	AccessToken               string
	AccessTokenExpireDuration int64
}

type TokenVerify struct {
	UserID auth_models.UserID
}
