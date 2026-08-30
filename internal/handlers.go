package internal

import (
	"context"
	"fmt"
	"time"
	"os"
	"github.com/google/uuid"
	"github.com/montruh-afk/gator/internal/database"
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
		return fmt.Errorf("Please provide a username to log in\n")
	} else if length > 1 {
		return fmt.Errorf("Username must not contain a whitespace.\n")
	}
	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println("User not found\n\t", err)
		os.Exit(1)
	}
	fmt.Println("username:", user.Name, "with ID:", user.ID, "already exists")
	
	if err := s.Configuration.SetUser(username); err != nil {
		return fmt.Errorf("couldn't set current user: %w\n", err)
	}
	
	fmt.Println("Welcome", username)
	return nil
}

func Register(s *State, cmd Command) error {
	length := len(cmd.Args)
	name := cmd.Args[0]
	if length  < 1 {
		return fmt.Errorf("Please provide a valid username\n")
	} else if length > 1 {
		return fmt.Errorf("Username must not contain a whitespace.\n")
	}

	if user, err := s.Db.GetUser(context.Background(), name); err == nil {
		fmt.Println("username:", user.Name, "with ID:", user.ID, "already exists")
		os.Exit(1)
	}
	parameters := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	}
	user, err := s.Db.CreateUser(context.Background(), parameters)
	if err != nil {
		return fmt.Errorf("Something went wrong: %v", err)
	} else {
		if err := s.Configuration.SetUser(user.Name); err != nil {
			return fmt.Errorf("couldn't set current user: %w\n", err)
		}
		fmt.Printf("User details: \n\t- Name: %s\n\t- ID: %v\n\t- Time created: %v\n", user.Name, user.ID, user.CreatedAt)
	}
	

	return nil
}