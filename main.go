package main

import (
	"fmt"
	"os"

	"github.com/GE1S7/gator/internal/config"
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

	cmds.register("login", handlerLogin)

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
