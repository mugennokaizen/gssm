package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type CliConfig struct {
	MasterKeyHash string `toml:"master_key_hash"`
	MasterKeySalt string `toml:"master_key_salt"`
}

func ReadConfigFromHomeDirToViper() {
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
		panic(fmt.Errorf("fatal error config file at path %s. Ensure that you have run gssm init and created config.toml file", gssmPath))
	}
}

func ReadCliConfigFromHomeDirToViper() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	gssmPath := filepath.Join(homeDir, ".gssm")

	if _, err := os.Stat(gssmPath); errors.Is(err, os.ErrNotExist) {
		panic(errors.New("app is not initialized"))
	}

	configPath := filepath.Join(gssmPath, "cli")

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		panic(errors.New("cli config does not exist"))
	}

	viper.SetConfigName("cli")
	viper.SetConfigType("toml")
	viper.AddConfigPath(gssmPath)
}

func ReadCliConfigToStruct() CliConfig {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	gssmPath := filepath.Join(homeDir, ".gssm")

	if _, err := os.Stat(gssmPath); errors.Is(err, os.ErrNotExist) {
		panic(errors.New("app is not initialized"))
	}

	configPath := filepath.Join(gssmPath, "cli")

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		panic(errors.New("cli config does not exist"))
	}

	var cliConfig CliConfig

	_, err = toml.DecodeFile(configPath, &cliConfig)
	if err != nil {
		return CliConfig{}
	}

	return cliConfig
}

func DumpCliConfig(cliConfig *CliConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	gssmPath := filepath.Join(homeDir, ".gssm")

	if _, err := os.Stat(gssmPath); errors.Is(err, os.ErrNotExist) {
		panic(errors.New("app is not initialized"))
	}

	configPath := filepath.Join(gssmPath, "cli")

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		panic(errors.New("cli config does not exist"))
	}

	file, err := os.OpenFile(configPath, os.O_CREATE, 0644)
	if err != nil {
		panic(errors.New("can't open cli config file to dump"))
	}

	enc := toml.NewEncoder(file)

	err = enc.Encode(cliConfig)
	if err != nil {
		return err
	}

	return nil
}
