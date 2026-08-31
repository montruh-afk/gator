package main

import (
	"fmt"
	"os"

	"github.com/montruh-afk/gator/internal"
)


func start(state *internal.State, cmd *internal.Commands) {
	cmd.Register("login", internal.HandlerLogin)
	cmd.Register("register", internal.Register)
	cmd.Register("reset", internal.Reset)
	cmd.Register("users", internal.Users)
	checker(state, cmd)
}

func checker(state *internal.State, cmd *internal.Commands) error {
	if len(os.Args) < 2 {
		fmt.Printf("Missing arguments needed to run gator\n\t Usage: cli <command> [args...]\n")
		os.Exit(1)
	}
	if len(os.Args) >= 2 {
		args := os.Args[2:]
		command := internal.Command{
		Name: os.Args[1],
		Args: args,
		}
		cmd.Run(state, command)
		return nil
	}
	os.Exit(1)
	return nil
	
}