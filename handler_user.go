package main

import (
	"blog_aggregator/internal/database"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("login expects a single argument, the username")
	}
	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}
	err = s.config.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	
	log.Printf("User %s logged in", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("register expects a single argument, the username")
	}
	arg := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	user, err := s.db.CreateUser(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	if err = s.config.SetUser(arg.Name); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	log.Printf("User created: %v", user.Name)
	return nil

}
