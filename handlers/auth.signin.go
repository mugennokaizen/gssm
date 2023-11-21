package handlers

import (
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"gssm/data"
	"gssm/db"
	"gssm/types"
	"gssm/utils"
)

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *auth) SignIn(c *fiber.Ctx) error {
	ctx := c.UserContext()
	us := do.MustInvoke[*db.UserSource](h.inj)
	t := do.MustInvoke[*data.TokenProcessor](h.inj)
	p := new(signInRequest)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	if !govalidator.IsEmail(p.Email) {
		return types.RnD(c, types.ResultBadEmail)
	}

	if !govalidator.StringMatches(p.Password, govalidator.ASCII) {
		return types.RnD(c, types.ResultBadPasswordAlphabet)
	}

	if len(p.Password) < 8 {
		return types.RnD(c, types.ResultBadPassword)
	}

	if !us.IsUserExist(ctx, p.Email) {
		return types.RnD(c, types.ResultUserNotFound)
	}

	user, _ := us.GetUserByEmail(ctx, p.Email)

	passwordIsCorrect := utils.VerifyPassword(p.Password, user.PasswordHash, user.Salt)

	if !passwordIsCorrect {
		return types.RnD(c, types.ResultWrongPassword)
	}

	accessToken, err := t.GenerateToken(user.Id, t.AccessTokenDuration)
	if err != nil {
		return types.RnD(c, types.ResultTokenGenerationError)
	}

	refreshToken, err := t.GenerateToken(user.Id, t.RefreshTokenDuration)
	if err != nil {
		return types.RnD(c, types.ResultTokenGenerationError)
	}

	c.Cookie(t.GetAccessTokenCookie(accessToken))
	c.Cookie(t.GetRefreshTokenCookie(refreshToken))

	return types.RnD(c, types.ResultOk)
}
