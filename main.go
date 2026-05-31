package main

import (
	"blog_aggregator/internal/config"
	"log"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	currentState := &state{
		config: &cfg,
	}

	availableCommands := commands{}
	availableCommands.register("login", handlerLogin)

	arguments := os.Args[1:]
	log.Println(arguments)
	newCommand := command{
		name: arguments[0],
		args: arguments[1:],
	}

	err = availableCommands.run(currentState, newCommand)
	if err != nil {
		log.Fatal(err)
	}

}
