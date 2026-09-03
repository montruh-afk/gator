package main

import (
	"github.com/montruh-afk/gator/internal"
	"github.com/montruh-afk/gator/internal/database"
	"context"
)

func middlewareLoggedIn(handler func(s *internal.State, cmd internal.Command, user database.User) error) func(*internal.State, internal.Command) error {
	return func(s *internal.State, c internal.Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Configuration.Current_user_name)
		if err != nil {
			return err
		}
		return handler(s, c, user)
	}
}
