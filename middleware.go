package main

import (
	"context"

	"github.com/cusackj4/gator/internal/database"
)

func middlewareFunc(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {

	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}

}
