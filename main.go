package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)


type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

type state struct {
	cfg *Config
}

type command struct {
	name string
	args []string
}

func getConfigFilePath(cfgpath string) (string, error){
	homeDir := "/mnt/f/blog_aggregator/"
	configPath := filepath.Join(homeDir, cfgpath)
	return configPath, nil
}

func Read(cfgpath string) (Config, error) {
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
	if err != nil{
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("No arguments passed in args")
	}
	
}

func main(){
	var cfg Config
	const configFileName = "gatorconfig.json"
	cfg, err := Read(configFileName)
	if err != nil {
		return 
	}
	userName := "Avi Banerjee"
	cfg.SetUser(userName, configFileName)


}