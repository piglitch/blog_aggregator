package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
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

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func handlerLogin(s *state, cmd command, cfgPath string) error {
	if len(cmd.args) == 0 {
		os.Exit(1)
		return errors.New("no arguments passed in args")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Println("user does not exist")
		os.Exit(1)
	}
	err = s.cfg.SetUser(cmd.args[0], cfgPath)
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}

func registerHandler(s *state, cmd command, cfgPath string) error {
	type user struct {
		ID uuid.UUID
		CreatedAt time.Time
		UpdatedAt time.Time
		Name string
	}
	newUser := user{
		Name: cmd.args[0],
	}
	getUser, err := s.db.GetUser(context.Background(), newUser.Name)
	if err != nil {
		println(err, newUser.Name)
	}
	println(newUser.Name, getUser.Name)
	if getUser.Name == newUser.Name {
		fmt.Println("user already exists")
		os.Exit(1)
	}
	dbUser, err := s.db.CreateUser(context.Background(), newUser.Name)
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(dbUser.Name, cfgPath)
	if err != nil {
		return err
	}
	fmt.Printf("User has been set and registered. CurrentUser: %s", dbUser.Name)
	return nil
}

func resetHandler(s *state, cmd command, cfgPath string) error { 
	err := s.db.ResetDb(context.Background())
	if err != err {
		return err
	}
	return nil
}

func getUsersHandler(s *state, cmd command, cfgPath string) error {
	dbUsers, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range dbUsers {
		var userName string
		userName = user.Name
		if user.Name == s.cfg.CurrentUser {
			userName = user.Name + " " + "(current)"
		}
		fmt.Printf("* %s\n", userName)
	}
	return nil
}

func addFeedsHandler(s *state, cmd command, cfgPath string) error {
	println(cmd.args)
	type feed struct {
		Name string
		Url string
		UserID uuid.UUID
	}
	dbUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		return err
	}
	newFeed := feed{
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: dbUser.ID,
	}
	err = s.db.AddFeed(context.Background(), database.AddFeedParams(newFeed))
	if err != nil {
		return err
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	v := RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}	
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}
	v.Channel.Description = html.UnescapeString(v.Channel.Description)
	v.Channel.Title = html.UnescapeString(v.Channel.Title)
	return &v, nil
}

func (c *commands) register(name string, f func(*state, command, string) error) {
	cmdMap := c.cmdName
	cmdMap[name] = f
}

func registerFetch(_ *state, _ command, _ string) error {
	fetchUrl := "https://www.wagslane.dev/index.xml"
	s, err := fetchFeed(context.Background(), fetchUrl)
	if err != nil {
		return err
	}
	for _, str := range s.Channel.Item {
		fmt.Println(str)
	}
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	const configFileName = ".gatorconfig.json"
	err := c.cmdName[cmd.name](s, cmd, configFileName) 
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
	newCommands.register("register", registerHandler)
	newCommands.register("reset", resetHandler)
	newCommands.register("users", getUsersHandler)
	newCommands.register("agg", registerFetch)
	newCommands.register("addfeed", addFeedsHandler)

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
	err = newCommands.run(&newState, newCliCmd)
	if err != nil {
		fmt.Println(err)
	}
}
