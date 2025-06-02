package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/GE1S7/gator/internal/config"
	"github.com/GE1S7/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
	conf, err := config.Read()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := state{
		conf: &conf,
	}

	cmds := commands{
		commandFunctions: make(map[string]func(state *state, cmd command) error),
	}

	db, err := sql.Open("postgres", conf.DbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s.db = dbQueries

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
