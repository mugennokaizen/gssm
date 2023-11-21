package db

import (
	"github.com/samber/do"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDatabase(inj *do.Injector) (*gorm.DB, error) {
	connectionString := viper.GetString("connection_string")
	return gorm.Open(postgres.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: false,
	})
}
