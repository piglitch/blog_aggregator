package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

type state struct {
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
	fmt.Println("SetUser: cfgPath =", cfgPath) // Added fmt.Println
	filePath, err := getConfigFilePath(cfgPath)
	if err != nil {
		return err
	}
	cfg.CurrentUser = userName
	newData, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}
	fmt.Println(filePath, "//// 68 ////")
	if err := os.WriteFile(filePath, newData, 0644); err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state, cmd command, cfgPath string) error {
	fmt.Println("handlerLogin: cfgPath =", cfgPath) // Added fmt.Println
	if len(cmd.args) == 0 {
		os.Exit(1)
		return errors.New("no arguments passed in args")
	}
	err := s.cfg.SetUser(cmd.args[0], cfgPath)
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}

func (c *commands) register(name string, f func(*state, command, string) error) {
	cmdMap := c.cmdName
	cmdMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	const configFileName = ".gatorconfig.json"
	err := handlerLogin(s, cmd, configFileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	const configFileName = ".gatorconfig.json"
	readCfg, err := Read(configFileName)
	if err != nil {
		return
	}
	newState := state{
		cfg: &readCfg,
	}
	cmdMap := make(
		map[string]func(s *state, cmd command, cfgPath string) error,
	)
	newCommands := commands{
		cmdName: cmdMap,
	}
	newCommands.register("login", handlerLogin)
	cliArgs := os.Args
	if len(cliArgs) < 2 {
		fmt.Println("insufficient args")
		os.Exit(1)
		return
	}
	newCliCmd := command{
		name: cliArgs[1],
		args: cliArgs[2:],
	}
	fmt.Println(newCliCmd.name, newCliCmd.args)
	err = newCommands.run(&newState, newCliCmd)
	if err != nil {
		fmt.Println(err)
	}
}
