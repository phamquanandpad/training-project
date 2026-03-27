package jwt_test

import (
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/config"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/infrastructure/jwt"
)

var testJwtConfig = &config.JwtConfig{
	AccessTokenSecret:          "test_access_secret_key_32bytes!",
	RefreshTokenSecret:         "test_refresh_secret_key_32bytes!",
	AccessTokenExpireDuration:  int64(15 * time.Minute.Seconds()),
	RefreshTokenExpireDuration: int64(24 * time.Hour.Seconds()),
}

func Test_tokenGenerator_GenerateAccessToken(t *testing.T) {
	t.Parallel()

	type args struct {
		userID auth_models.UserID
	}

	type testcase struct {
		args    args
		wantErr bool
	}

	testTables := map[string]testcase{
		"Generate access token for user 1": {
			args:    args{userID: auth_models.UserID(1)},
			wantErr: false,
		},
		"Generate access token for user 99": {
			args:    args{userID: auth_models.UserID(99)},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			generator := jwt.NewTokenGenerator(testJwtConfig)

			tokenStr, expiresSecond, err := generator.GenerateAccessToken(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tokenStr == "" {
				t.Error("GenerateAccessToken() returned empty token string")
			}

			if expiresSecond != testJwtConfig.AccessTokenExpireDuration {
				t.Errorf("GenerateAccessToken() expiresSecond = %d, want %d", expiresSecond, testJwtConfig.AccessTokenExpireDuration)
			}
		})
	}
}

func Test_tokenGenerator_GenerateRefreshToken(t *testing.T) {
	t.Parallel()

	type args struct {
		userID auth_models.UserID
	}

	type testcase struct {
		args    args
		wantErr bool
	}

	testTables := map[string]testcase{
		"Generate refresh token for user 1": {
			args:    args{userID: auth_models.UserID(1)},
			wantErr: false,
		},
		"Generate refresh token for user 99": {
			args:    args{userID: auth_models.UserID(99)},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			generator := jwt.NewTokenGenerator(testJwtConfig)

			tokenStr, expiresSecond, err := generator.GenerateRefreshToken(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tokenStr == "" {
				t.Error("GenerateRefreshToken() returned empty token string")
			}

			if expiresSecond != int64(testJwtConfig.RefreshTokenExpireDuration) {
				t.Errorf("GenerateRefreshToken() expiresSecond = %d, want %d", expiresSecond, testJwtConfig.RefreshTokenExpireDuration)
			}
		})
	}
}

func Test_tokenVerifier_VerifyAccessToken(t *testing.T) {
	t.Parallel()

	generator := jwt.NewTokenGenerator(testJwtConfig)
	verifier := jwt.NewTokenVerifier(testJwtConfig)

	t.Run("Verify valid access token", func(t *testing.T) {
		t.Parallel()

		expectedUserID := auth_models.UserID(42)
		tokenStr, _, err := generator.GenerateAccessToken(expectedUserID)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		gotUserID, err := verifier.VerifyAccessToken(tokenStr)
		if err != nil {
			t.Errorf("VerifyAccessToken() unexpected error = %v", err)
			return
		}

		if gotUserID != expectedUserID {
			t.Errorf("VerifyAccessToken() userID = %v, want %v", gotUserID, expectedUserID)
		}
	})

	t.Run("Verify invalid access token", func(t *testing.T) {
		t.Parallel()

		_, err := verifier.VerifyAccessToken("this.is.invalid")
		if err == nil {
			t.Error("VerifyAccessToken() expected error for invalid token, got nil")
		}
	})

	t.Run("Verify access token signed with wrong secret", func(t *testing.T) {
		t.Parallel()

		wrongConfig := &config.JwtConfig{
			AccessTokenSecret:         "wrong_secret_key_that_is_different",
			AccessTokenExpireDuration: int64(15 * time.Minute.Seconds()),
		}
		wrongGenerator := jwt.NewTokenGenerator(wrongConfig)

		tokenStr, _, err := wrongGenerator.GenerateAccessToken(auth_models.UserID(1))
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		_, err = verifier.VerifyAccessToken(tokenStr)
		if err == nil {
			t.Error("VerifyAccessToken() expected error for wrong secret, got nil")
		}
	})

	t.Run("Verify expired access token", func(t *testing.T) {
		t.Parallel()

		// Build an already-expired token manually.
		claims := map[string]interface{}{
			"user_id": "1",
			"iat":     time.Now().Add(-2 * time.Hour).Unix(),
			"exp":     time.Now().Add(-1 * time.Hour).Unix(),
		}
		token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, gojwt.MapClaims(claims))
		tokenStr, err := token.SignedString([]byte(testJwtConfig.AccessTokenSecret))
		if err != nil {
			t.Fatalf("failed to sign expired token: %v", err)
		}

		_, err = verifier.VerifyAccessToken(tokenStr)
		if err == nil {
			t.Error("VerifyAccessToken() expected error for expired token, got nil")
		}
	})
}

func Test_tokenVerifier_VerifyRefreshToken(t *testing.T) {
	t.Parallel()

	generator := jwt.NewTokenGenerator(testJwtConfig)
	verifier := jwt.NewTokenVerifier(testJwtConfig)

	t.Run("Verify valid refresh token", func(t *testing.T) {
		t.Parallel()

		expectedUserID := auth_models.UserID(7)
		tokenStr, _, err := generator.GenerateRefreshToken(expectedUserID)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		gotUserID, err := verifier.VerifyRefreshToken(tokenStr)
		if err != nil {
			t.Errorf("VerifyRefreshToken() unexpected error = %v", err)
			return
		}

		if gotUserID != expectedUserID {
			t.Errorf("VerifyRefreshToken() userID = %v, want %v", gotUserID, expectedUserID)
		}
	})

	t.Run("Verify invalid refresh token", func(t *testing.T) {
		t.Parallel()

		_, err := verifier.VerifyRefreshToken("this.is.invalid")
		if err == nil {
			t.Error("VerifyRefreshToken() expected error for invalid token, got nil")
		}
	})

	t.Run("Verify refresh token signed with wrong secret", func(t *testing.T) {
		t.Parallel()

		wrongConfig := &config.JwtConfig{
			RefreshTokenSecret:         "wrong_secret_key_that_is_different",
			RefreshTokenExpireDuration: int64(24 * time.Hour.Seconds()),
		}
		wrongGenerator := jwt.NewTokenGenerator(wrongConfig)

		tokenStr, _, err := wrongGenerator.GenerateRefreshToken(auth_models.UserID(1))
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		_, err = verifier.VerifyRefreshToken(tokenStr)
		if err == nil {
			t.Error("VerifyRefreshToken() expected error for wrong secret, got nil")
		}
	})

	t.Run("Access token cannot be used as refresh token", func(t *testing.T) {
		t.Parallel()

		// Access and refresh tokens use different secrets, so cross-use should fail.
		accessToken, _, err := generator.GenerateAccessToken(auth_models.UserID(1))
		if err != nil {
			t.Fatalf("failed to generate access token: %v", err)
		}

		_, err = verifier.VerifyRefreshToken(accessToken)
		if err == nil {
			t.Error("VerifyRefreshToken() expected error when using access token, got nil")
		}
	})
}
