package internal


import (
	"fmt"
)

type Command struct {
	Name string
	Args []string
}

func HandlerLogin(s *State, cmd Command) error {
	length := len(cmd.Args)
	username := cmd.Args[0]

	// Check if the user provided a username
	if length < 1 {
		return fmt.Errorf("Please provide a username to Sign in")
	} else if length > 1 {
		return fmt.Errorf("Username must not contain a whitespace.")
	}

	err := s.Configuration.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	
	fmt.Println("Username has been set: you are ", username)
	return nil
}