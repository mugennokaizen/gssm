package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"gssm/data"
	"gssm/db"
	"gssm/handlers"
	"log"
	"os"
	"os/signal"
)

func main() {
	app := NewFiberServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":9000"); err != nil {
		log.Panic(err)
	}
}

func NewFiberServer() *fiber.App {
	viper.SetConfigName("/cmd/front/config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config fileddd: %w", err))
	}

	injector := do.New()
	do.Provide(injector, db.NewDatabase)
	do.Provide(injector, db.NewUserSource)
	do.Provide(injector, data.NewTokenProcessor)

	app := fiber.New()

	authHandler := handlers.NewAuthHandler(injector)
	authGroup := app.Group(authHandler.GetGroup())
	authGroup.Post("/sign-in", authHandler.SignIn)
	authGroup.Post("/sign-up", authHandler.SignUp)
	authGroup.Post("/refresh", authHandler.Refresh)

	//app.Use(middlewares.NewJwt(middlewares.JwtConfig{
	//	Filter:            nil,
	//	RefreshCookieName: viper.GetString("jwt.refresh_cookie_name"),
	//	AccessCookieName:  viper.GetString("jwt.access_cookie_name"),
	//	SecretKey:         viper.GetString("jwt.secret_key"),
	//}))

	return app
}
