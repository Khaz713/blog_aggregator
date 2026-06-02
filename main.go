package main

import (
	"blog_aggregator/internal/config"
	"blog_aggregator/internal/database"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

//"postgres://postgres:3323@localhost:5432/gator"

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	log.Printf("config: %+v", cfg)

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	db.Close()
	dbQueries := database.New(db)

	currentState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)

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
