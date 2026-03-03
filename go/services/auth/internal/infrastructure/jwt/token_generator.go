package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/config"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway"
)

type tokenGenerator struct {
	jwt_conf *config.JwtConfig
}

func NewTokenGenerator(
	jwt_conf *config.JwtConfig,
) gateway.JwtGenerateGateway {
	return &tokenGenerator{
		jwt_conf: jwt_conf,
	}
}

func (g *tokenGenerator) GenerateAccessToken(userID auth_models.UserID) (string, int64, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Second * time.Duration(g.jwt_conf.AccessTokenExpireDuration)),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(g.jwt_conf.AccessTokenSecret))
	if err != nil {
		return "", 0, err
	}

	return signedString, int64(g.jwt_conf.AccessTokenExpireDuration), nil
}

func (g *tokenGenerator) GenerateRefreshToken(userID auth_models.UserID) (string, int64, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Second * time.Duration(g.jwt_conf.RefreshTokenExpireDuration)),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(g.jwt_conf.RefreshTokenSecret))
	if err != nil {
		return "", 0, err
	}

	return signedString, int64(g.jwt_conf.RefreshTokenExpireDuration), nil
}
