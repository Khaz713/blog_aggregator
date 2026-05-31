package main

import (
	"errors"
	"log"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login expects a single argument, the username")
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	}
	log.Println("User set")
	return nil
}
