package internal

import (
	"fmt"
)

type Commands struct {
	Handler map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if exists, ok := c.Handler[cmd.Name]; ok {
		return exists(s, cmd)
	}
	return fmt.Errorf("%s failed: Nonexistent handler present for the provided command.\n", cmd.Name)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if _, ok := c.Handler[name]; ok {
		fmt.Printf("Operation nullified: %s is a process still in service\n", name)
		return
	} else {
		c.Handler[name] = f
	}
}