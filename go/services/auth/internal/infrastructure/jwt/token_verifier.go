package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/config"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway"
)

type tokenVerifier struct {
	jwt_conf *config.JwtConfig
}

func NewTokenVerifier(jwt_conf *config.JwtConfig) gateway.JwtVerifyGateway {
	return &tokenVerifier{
		jwt_conf: jwt_conf,
	}
}

func (g *tokenVerifier) VerifyAccessToken(tokenStr string) (auth_models.UserID, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(g.jwt_conf.AccessTokenSecret), nil
		},
	)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.UserID

		userIDInt, err := auth_models.ParseUserID(userID)
		if err != nil {
			return 0, err
		}

		return userIDInt, nil
	} else {
		return 0, jwt.ErrInvalidKey
	}
}

func (g *tokenVerifier) VerifyRefreshToken(tokenStr string) (auth_models.UserID, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(g.jwt_conf.RefreshTokenSecret), nil
		},
	)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.UserID

		userIDInt, err := auth_models.ParseUserID(userID)
		if err != nil {
			return 0, err
		}

		return userIDInt, nil
	} else {
		return 0, jwt.ErrInvalidKey
	}
}
