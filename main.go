package main

import (
	"blog_aggregator/internal/config"
	"blog_aggregator/internal/database"
	"database/sql"
	"log"
	"os"
)

import _ "github.com/lib/pq"

type state struct {
	config *config.Config
	db     *database.Queries
}

//const connectionString := "postgres://postgres:3323@localhost:5432/gator"

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	currentState := &state{
		config: &cfg,
	}

	db, err := sql.Open("postgres", currentState.config.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	currentState.db = dbQueries

	cmds := commands{}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	arguments := os.Args[1:]
	log.Printf("Command: %v", arguments)
	newCommand := command{
		name: arguments[0],
		args: arguments[1:],
	}

	err = cmds.run(currentState, newCommand)
	if err != nil {
		log.Fatal(err)
	}

}
