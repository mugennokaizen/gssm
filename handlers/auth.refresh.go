package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"gssm/data"
	"gssm/types"
	"time"
)

func (a *auth) Refresh(c *fiber.Ctx) error {

	d := do.MustInvoke[data.TokenProcessor](a.inj)

	refreshTokenCookie := c.Cookies(viper.GetString("jwt.refresh_cookie_name"))

	tokenIsValid, claims, err := d.VerifyToken(refreshTokenCookie)
	if err != nil || !tokenIsValid {
		return err
	}

	if claims.ExpiresAt.Time.Sub(time.Now()) <= 0 {
		return fiber.ErrUnauthorized
	}

	accessTokenCookie := c.Cookies(viper.GetString("jwt.access_cookie_name"))

	if len(accessTokenCookie) > 0 {
		tokenIsValid, claims, err = d.VerifyToken(accessTokenCookie)
		if err != nil || !tokenIsValid {
			return err
		}
		if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
			return fiber.ErrBadRequest
		}
	}

	accessToken, err := d.GenerateToken(claims.Id, viper.GetDuration("jwt.access_token_duration"))
	if err != nil {
		return c.JSON(types.ResponseNoData{
			Code: types.ResultOk,
		})
	}

	refreshToken, err := d.GenerateToken(claims.Id, viper.GetDuration("jwt.refresh_token_duration"))
	if err != nil {
		return c.JSON(types.ResponseNoData{
			Code: types.ResultTokenGenerationError,
		})
	}

	if err != nil {
		return c.JSON(types.ResponseNoData{
			Code: types.ResultTokenGenerationError,
		})
	}

	c.Cookie(d.GetAccessTokenCookie(accessToken))
	c.Cookie(d.GetRefreshTokenCookie(refreshToken))

	return c.JSON(types.ResponseNoData{
		Code: types.ResultOk,
	})
}
