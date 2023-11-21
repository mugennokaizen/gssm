package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"gssm/db"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config fileddd: %w", err))
	}

	injector := do.New()
	do.Provide(injector, db.NewDatabase)
	do.Provide(injector, db.NewUserSource)

	us := do.MustInvoke[*db.UserSource](injector)

	maxGoroutines := 50
	guard := make(chan struct{}, maxGoroutines)

	for i := 0; i < 100000; i++ {
		guard <- struct{}{}

		go func(n int) {
			user, err := us.CreateUser(context.Background(), gofakeit.Email(), gofakeit.Password(true, true, true, false, false, 20))
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println(user)
			<-guard
		}(i)

	}

}
