package main

import (
	"context"
	"fmt"
	"strconv"

	"main.go/internal/database"
)

func browsePosts(s *state, cmd command, user database.User) error {
	dbPosts, err := s.db.GetPostsFromUser(context.Background(), user.ID) 
	if err != nil {
		return err
	}
	fmt.Println(cmd.args, len(dbPosts))
	var limit int
	if len(cmd.args) < 1 {
		cmd.args = []string{"2"}
	}
	limit, err = strconv.Atoi(cmd.args[0])
	if err != nil {
		return err
	}
	if limit > len(dbPosts) {
		limit = len(dbPosts)
	}
	for _, post := range dbPosts[:limit]{
		fmt.Println(post.Title)
		fmt.Println(post.Description)
	}
	fmt.Println(dbPosts, user.ID)
	return nil
}