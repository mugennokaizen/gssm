package data

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"gssm/types"
	"time"
)

type TokenProcessor struct {
	RefreshCookieName    string
	RefreshTokenDuration time.Duration
	AccessCookieName     string
	AccessTokenDuration  time.Duration
	IsSecure             bool
	SecretKey            string
}

func NewTokenProcessor(_ *do.Injector) (*TokenProcessor, error) {
	return &TokenProcessor{
		RefreshCookieName:    viper.GetString("jwt.refresh_cookie_name"),
		RefreshTokenDuration: viper.GetDuration("jwt.refresh_token_duration"),
		AccessCookieName:     viper.GetString("jwt.access_cookie_name"),
		AccessTokenDuration:  viper.GetDuration("jwt.access_token_duration"),
		IsSecure:             viper.GetBool("jwt.secure_token"),
		SecretKey:            viper.GetString("jwt.secret_key"),
	}, nil
}

func (t *TokenProcessor) VerifyToken(token string) (bool, *types.Claims, error) {
	claims := &types.Claims{}

	jwtKey := []byte(t.SecretKey)

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return false, nil, fiber.ErrUnauthorized
		}
		return false, nil, fiber.ErrBadRequest
	}
	if !tkn.Valid {
		return false, nil, fiber.ErrUnauthorized
	}

	return true, claims, nil
}

func (t *TokenProcessor) GenerateToken(id types.ULID, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)

	claims := types.Claims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &claims)
	tokenString, err := token.SignedString([]byte(t.SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *TokenProcessor) GetRefreshTokenCookie(refreshToken string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = t.RefreshCookieName
	cookie.Value = refreshToken
	cookie.HTTPOnly = true
	cookie.Secure = t.IsSecure
	cookie.Expires = time.Now().Add(t.RefreshTokenDuration)

	return cookie
}

func (t *TokenProcessor) GetAccessTokenCookie(accessToken string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = t.AccessCookieName
	cookie.Value = accessToken
	cookie.HTTPOnly = true
	cookie.Secure = t.IsSecure
	cookie.Expires = time.Now().Add(t.AccessTokenDuration)

	return cookie
}
