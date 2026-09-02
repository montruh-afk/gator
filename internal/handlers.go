package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/montruh-afk/gator/internal/database"
	"github.com/montruh-afk/gator/internal/network"
	"os"
	"time"
)

type Command struct {
	Name string
	Args []string
}

func HandlerLogin(s *State, cmd Command) error {
	length := len(cmd.Args)

	// Check if the user provided a username
	if length < 1 {
		return fmt.Errorf("Please provide a username to log in\n")
	} else if length > 1 {
		return fmt.Errorf("Username must not contain a whitespace.\n")
	}
	username := cmd.Args[0]
	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println("User not found\n\t", err)
		os.Exit(1)
	}

	if err := s.Configuration.SetUser(user.Name); err != nil {
		return fmt.Errorf("couldn't set current user to %v\n", err)
	}

	fmt.Println("Welcome", username)
	return nil
}

func Register(s *State, cmd Command) error {
	length := len(cmd.Args)

	if length < 1 {
		return fmt.Errorf("Please provide a valid username\n")
	} else if length > 1 {
		return fmt.Errorf("Username must not contain a whitespace.\n")
	}
	name := cmd.Args[0]
	if user, err := s.Db.GetUser(context.Background(), name); err == nil {
		fmt.Println("username:", user.Name, "with ID:", user.ID, "already exists")
		os.Exit(1)
	}
	parameters := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
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
	return Write(*s.Configuration)
}

func Users(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Println("Something went wrong", err)
		os.Exit(1)
	}
	if len(users) < 1 {
		return fmt.Errorf("No users on record\n")
	}
	current := s.Configuration.Current_user_name
	for _, user := range users {
		if user.Name == current {
			fmt.Println("*", user.Name, "(current)")
		} else {
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

func AddFeed(s *State, cmd Command) error {
	length := len(cmd.Args)
	if length != 2 {
		return fmt.Errorf("Invalid format\n\tUsage: addfeed '<feed name>' '<feed url>'\n")
	} else if network.Validateargs("AddFeed", cmd.Args) {
		user, err := s.Db.GetUser(context.Background(), s.Configuration.Current_user_name)
		if err != nil {
			return fmt.Errorf("Something went wrong while fetching user details: %v\n", err)
		}
		feed := database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Args[0],
			Url:       cmd.Args[1],
			UserID:    user.ID,
		}
		newFeed, err := s.Db.CreateFeed(context.Background(), feed)
		if err != nil {
			return fmt.Errorf("Error adding feed: %v\n", err)
		}

		fmt.Printf("Operation completed successfully\n\t- Feed name: %s\n\t- Feed URL: %s\n\t- Created at: %v\n\t- Linked to User ID: %v\n", newFeed.Name, newFeed.Url, newFeed.CreatedAt, newFeed.UserID)
		if err := addFeedFollow(s, cmd); err != nil {
			return fmt.Errorf("Something broke while attempting to follow %s \nPlease try running 'follow' <feed URL>\n", newFeed.Name)
		}
		return nil
	}
	return fmt.Errorf("Something went wrong...")

}

func getUserById(s *State, userid uuid.UUID) (*database.User, error) {
	user, err := s.Db.FetchUser(context.Background(), userid)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func Feeds(s *State, cmd Command) error {
	allFeed, err := s.Db.FetchFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Something went wrong: %s\n", err)
	}
	for _, feed := range allFeed {
		user, err := getUserById(s, feed.UserID)
		if err != nil {
			continue
		}
		fmt.Printf("• %s\n\t- %s Posted by: %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func Follow(s *State, cmd Command) error {
	length := len(cmd.Args)
	if length < 1 {
		return fmt.Errorf("Invalid format\n\tUsage: follow <feed url>\n")
	} else if network.Validateargs("Follow", cmd.Args) {
		userInfo, err := s.Db.GetUser(context.Background(), s.Configuration.Current_user_name)
		if err != nil {
			return fmt.Errorf("Something went wrong in authenticating user\n\t -%v\n", err)
		}
		feedInfo, err_ := s.Db.GetFeedbyURL(context.Background(), cmd.Args[0])
		if err_ != nil {
			return fmt.Errorf("Something broke while fetching results from our records: %v\n", err_)
		}

		newFollow := database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userInfo.ID,
			FeedID:    feedInfo.ID,
		}
		_, err = s.Db.CreateFeedFollow(context.Background(), newFollow)
		if err != nil {
			return fmt.Errorf("Something went wrong on our end: %v\n", err)
		}
		fmt.Printf("%s is now following %s\n", s.Configuration.Current_user_name, feedInfo.Name)
		return nil
	}
	return fmt.Errorf("Something went wrong...\n")
}

func addFeedFollow(s *State, cmd Command) error {
	if network.Validateargs("AddFeed", cmd.Args) {
		userInfo, err := s.Db.GetUser(context.Background(), s.Configuration.Current_user_name)
		if err != nil {
			return fmt.Errorf("Something went wrong in authenticating user\n\t -%v\n", err)
		}
		feedInfo, err_ := s.Db.GetFeedbyURL(context.Background(), cmd.Args[1])
		if err_ != nil {
			return fmt.Errorf("Something broke while fetching results from our records: %v\n", err_)
		}

		newFollow := database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userInfo.ID,
			FeedID:    feedInfo.ID,
		}
		_, err = s.Db.CreateFeedFollow(context.Background(), newFollow)
		if err != nil {
			return fmt.Errorf("Something went wrong on our end: %v\n", err)
		}
		fmt.Printf("%s is now following %s\n", s.Configuration.Current_user_name, feedInfo.Name)
		return nil
	}
	return fmt.Errorf("Something broke while attempting to follow %s\n", cmd.Args[0])
}

func Following(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Invalid format\n\t following takes no arguments\n")
	}
	user := s.Configuration.Current_user_name
	follows, err := s.Db.GetUserFollows(context.Background(), user)
	if err != nil {
		return fmt.Errorf("Something went wrong in fetching results: %v\n", err)
	}
	fmt.Println(user, "is currently following")
	for _, follow := range follows {
		fmt.Printf("\t- %s\n", follow.FeedName)
	}
	return nil
}
