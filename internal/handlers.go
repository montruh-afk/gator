package internal

import (
	"context"
	"fmt"
	"os"
	"time"
	"github.com/google/uuid"
	"github.com/montruh-afk/gator/internal/database"
	"github.com/montruh-afk/gator/internal/network"
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

func Reset(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		fmt.Println("Reset must be called with no arguments")
		os.Exit(1)
	}
	if err := s.Db.DeleteUsers(context.Background()); err != nil {
		fmt.Println("Something went wrong\n\t", err)
		os.Exit(1)
	}
	s.Configuration.Current_user_name = ""
	fmt.Println("Operation successful")
	return nil
}

func Users(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Println("Something went wrong", err)
		os.Exit(1)
	}

	current := s.Configuration.Current_user_name
	for _, user := range users {
		if user.Name == current {
			fmt.Println("*", user.Name, "(current)")
		} else{
			fmt.Println("*", user.Name)
		}
		
	}

	return nil
}

func Agg(s *State, cmd Command) error {
	feed, err := network.Fetchfeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}