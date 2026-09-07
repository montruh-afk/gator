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
	cmd.Register("agg", internal.Agg)
	cmd.Register("addfeed", middlewareLoggedIn(internal.AddFeed))
	cmd.Register("feeds", internal.Feeds)
	cmd.Register("follow", middlewareLoggedIn(internal.Follow))
	cmd.Register("following", middlewareLoggedIn(internal.Following))
	cmd.Register("unfollow", middlewareLoggedIn(internal.UnFollow))
	cmd.Register("browse", middlewareLoggedIn(internal.Browse))
	if err := checker(state, cmd); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func checker(state *internal.State, cmd *internal.Commands) error {
	if len(os.Args) < 2 {
		fmt.Println("Missing arguments needed to run gator\n\t Usage: cli <command> [args...]")
		os.Exit(1)
	} else if _, ok := cmd.Handlers[os.Args[1]]; !ok {
		fmt.Println("Unknown Command")
		os.Exit(1)
	}
	if len(os.Args) >= 2 {
		args := os.Args[2:]
		command := internal.Command{
		Name: os.Args[1],
		Args: args,
		}
		if err := cmd.Run(state, command); err != nil {
			return err
		}
		return nil
	}
	os.Exit(1)
	return nil
	
}

