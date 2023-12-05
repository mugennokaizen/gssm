package utils

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func ReadConfigFromHomeDir() {
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
}
