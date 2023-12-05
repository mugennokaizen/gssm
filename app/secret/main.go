package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"gssm/db"
	"gssm/handlers"
	"gssm/handlers/middlewares"
	"gssm/immu"
	"log"
	"os"
	"os/signal"
	"path/filepath"
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

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	gssmPath := filepath.Join(homeDir, ".gssm")

	if _, err := os.Stat(gssmPath); errors.Is(err, os.ErrNotExist) {
		panic(errors.New("app is not initialized"))
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(gssmPath)

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file at path %s. Ensure that you have run gssm init", gssmPath))
	}

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config fileddd: %w", err))
	}

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
