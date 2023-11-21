package handlers

import (
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"gssm/data"
	"gssm/db"
	"gssm/types"
)

type signUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *auth) SignUp(c *fiber.Ctx) error {
	p := new(signUpRequest)
	ctx := c.UserContext()
	userSource := do.MustInvoke[*db.UserSource](h.inj)
	t := do.MustInvoke[*data.TokenProcessor](h.inj)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	if !govalidator.IsEmail(p.Email) {
		return types.RnD(c, types.ResultBadPassword)
	}

	if !govalidator.StringMatches(p.Password, govalidator.ASCII) {
		return types.RnD(c, types.ResultBadPasswordAlphabet)
	}

	if len(p.Password) < 8 {
		return types.RnD(c, types.ResultBadPassword)
	}

	if userSource.IsUserExist(ctx, p.Email) {
		return types.RnD(c, types.ResultUserAlreadyExist)
	}

	id, err := userSource.CreateUser(ctx, p.Email, p.Password)
	if err != nil {
		return types.RnD(c, types.ResultCreationUserError)
	}

	accessToken, err := t.GenerateToken(id, t.AccessTokenDuration)
	if err != nil {
		return types.RnD(c, types.ResultTokenGenerationError)
	}

	refreshToken, err := t.GenerateToken(id, t.RefreshTokenDuration)
	if err != nil {
		return types.RnD(c, types.ResultTokenGenerationError)
	}

	c.Cookie(t.GetAccessTokenCookie(accessToken))
	c.Cookie(t.GetRefreshTokenCookie(refreshToken))

	return types.RnD(c, types.ResultOk)
}
