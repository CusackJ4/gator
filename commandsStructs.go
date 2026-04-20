package main

import "fmt"

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	/* f is the function being retrieved and run by
	this method. */
	f, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("Command not found!")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state,
	command) error) {

	c.cmds[name] = f

}
