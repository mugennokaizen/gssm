package types

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"time"
)

type ResultCode int

const (
	ResultOk                   ResultCode = 0
	ResultBadEmail                        = 1001
	ResultUserAlreadyExist                = 1002
	ResultUserNotFound                    = 1003
	ResultBadPassword                     = 1004
	ResultTokenGenerationError            = 1005
	ResultBadPasswordAlphabet             = 1006
	ResultCreationUserError               = 1007
	ResultWrongPassword                   = 1008
)

type ResponseNoData struct {
	Code ResultCode `json:"code"`
}

type Response[T any] struct {
	Code ResultCode `json:"code"`
	Data T          `json:"data"`
}

func RnD(c *fiber.Ctx, code ResultCode) error {
	return c.JSON(ResponseNoData{Code: code})
}

func R[T any](c *fiber.Ctx, code ResultCode, data T) error {
	return c.JSON(Response[T]{
		Code: code,
		Data: data,
	})
}

type ULID string

func (i ULID) ToUnix() int64 {
	return int64(ulid.MustParse(string(i)).Time())
}

func (i ULID) ToTime() time.Time {
	return time.UnixMilli(i.ToUnix())
}

type FiberFunc func(c *fiber.Ctx) error

type Claims struct {
	Id ULID
	jwt.RegisteredClaims
}
