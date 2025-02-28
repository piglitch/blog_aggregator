package main

import (
	"context"

	"main.go/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command, string) error {
	return func(s1 *state, c command, s2 string) error {
		dbUser, err := s1.db.GetUser(context.Background(), s1.cfg.CurrentUser)
		if err != nil {
			return err
		}
		err = handler(s1, c, dbUser)
		if err != nil {
			return err
		}
		return nil
	}
}