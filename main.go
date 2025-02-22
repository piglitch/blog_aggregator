package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"main.go/internal/config"
	"main.go/internal/database"
)

type Config struct {
	DbUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

type state struct {
	db *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmdName map[string]func(*state, command, string) error
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
	err := godotenv.Load()
	if err != nil {
		return
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return 
	}
	dbQueries := database.New(db)
	const configFileName = ".gatorconfig.json"
	readCfg, err := config.Read(configFileName)
	if err != nil {
		return
	}
	newState := state{
		cfg: &readCfg,
		db: dbQueries,
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
