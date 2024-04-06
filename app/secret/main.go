package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"gssm/config"
	"gssm/db"
	"gssm/handlers"
	"gssm/handlers/middlewares"
	"gssm/immu"
	"log"
	"os"
	"os/signal"
)

func main() {
	app := NewFiberServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

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

	config.ReadConfigFromHomeDirToViper()
	injector := do.New()
	do.Provide(injector, db.NewDatabase)

	do.Provide(injector, immu.NewDatabase)
	do.Provide(injector, immu.NewManager)

	app := fiber.New()

	app.Use(middlewares.NewAccessKey(middlewares.InjectorConfig{
		Filter:   nil,
		Injector: injector,
	}))

	authHandler := handlers.NewAuthHandler(injector)
	authGroup := app.Group(authHandler.GetGroup())
	authGroup.Post("/sign-in", authHandler.SignIn)
	authGroup.Post("/sign-up", authHandler.SignUp)
	authGroup.Post("/refresh", authHandler.Refresh)

	return app
}
