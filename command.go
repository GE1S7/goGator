package main

import (
	"fmt"

	"github.com/GE1S7/gator/internal/config"
)

type state struct {
	conf *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Not enough arguments: username required.")
	}
	err := s.conf.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("Username has been set")
	return nil
}

type commands struct {
	commandFunctions map[string]func(state *state, cmd command) error
}

func (c *commands) run(s *state, cmd command) error {
	// check if command with state s exist and runs it
	_, ok := c.commandFunctions[cmd.name]
	if !ok {
		return fmt.Errorf("Command %v does not exist.", cmd.name)
	}

	err := c.commandFunctions[cmd.name](s, cmd)

	return err

}

func (c *commands) register(name string, f func(state *state, cmd command) error) {
	// register a new handler function for a command name
	c.commandFunctions[name] = f
}
