package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	_ "github.com/lib/pq"
	"main.go/internal/database"
)

type Config struct {
	DbUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

type state struct {
	db *database.Queries
	cfg *Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmdName map[string]func(*state, command, string) error
}

func getConfigFilePath(cfgpath string) (string, error) {
	fmt.Println("getConfigFilePath: cfgpath =", cfgpath) // Added fmt.Println
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("cannot find home directory")
		return "", err
	}
	configPath := filepath.Join(homeDir, cfgpath)
	return configPath, nil
}

func Read(cfgpath string) (Config, error) {
	fmt.Println("Read: cfgpath =", cfgpath) // Added fmt.Println
	configPath, err := getConfigFilePath(cfgpath)
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg *Config) SetUser(userName string, cfgPath string) error {
	filePath, err := getConfigFilePath(cfgPath)
	if err != nil {
		return err
	}
	cfg.CurrentUser = userName
	newData, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, newData, 0644); err != nil {
		return err
	}
	return nil
}

