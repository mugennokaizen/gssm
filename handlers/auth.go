package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
)

type AuthHandler interface {
	GroupHandler
	SignIn(ctx *fiber.Ctx) error
	//SignUp(ctx *fiber.Ctx) error
	Refresh(ctx *fiber.Ctx) error
}

type auth struct {
	inj *do.Injector
}

func NewAuthHandler(inj *do.Injector) AuthHandler {
	return &auth{
		inj: inj,
	}
}

func (_ *auth) GetGroup() string {
	return "auth"
}
