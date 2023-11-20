package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"gssm/data"
	"gssm/types"
	"time"
)

func (a *auth) Refresh(c *fiber.Ctx) error {

	t := do.MustInvoke[*data.TokenProcessor](a.inj)

	refreshTokenCookie := c.Cookies(t.RefreshCookieName)

	tokenIsValid, claims, err := t.VerifyToken(refreshTokenCookie)
	if err != nil || !tokenIsValid {
		return err
	}

	if claims.ExpiresAt.Time.Sub(time.Now()) <= 0 {
		return fiber.ErrUnauthorized
	}

	accessTokenCookie := c.Cookies(t.AccessCookieName)

	if len(accessTokenCookie) > 0 {
		tokenIsValid, claims, err = t.VerifyToken(accessTokenCookie)
		if err != nil || !tokenIsValid {
			return err
		}
		if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
			return fiber.ErrBadRequest
		}
	}

	accessToken, err := t.GenerateToken(claims.Id, t.AccessTokenDuration)
	if err != nil {
		return types.RnD(c, types.ResultTokenGenerationError)
	}

	refreshToken, err := t.GenerateToken(claims.Id, t.RefreshTokenDuration)
	if err != nil {
		return types.RnD(c, types.ResultTokenGenerationError)
	}

	c.Cookie(t.GetAccessTokenCookie(accessToken))
	c.Cookie(t.GetRefreshTokenCookie(refreshToken))

	return types.RnD(c, types.ResultOk)
}
