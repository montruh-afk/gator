package internal

import (
	"context"
	"fmt"
	"os"
	"time"
	"strconv"
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
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Invalid format\n\tUsage: agg <time period> e.g 1m, 1h, 300ms... this is how frequent the feeds will be updated.\n")
	}
	time_between_reqs := cmd.Args[0]
	reqs, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return fmt.Errorf("Invalid duration %w\n", err)
	}
	fmt.Printf("Collecting feeds every %s...\n", time_between_reqs)
	ticker := time.NewTicker(reqs)
	for ; ; <-ticker.C {
		network.ScrapeFeeds(s.Db)
	}
}

func AddFeed(s *State, cmd Command, user database.User) error {
	length := len(cmd.Args)
	if length != 2 {
		return fmt.Errorf("Invalid format\n\tUsage: addfeed '<feed name>' '<feed url>'\n")
	} else if network.Validateargs("AddFeed", cmd.Args) {
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
		if err := addFeedFollow(s, cmd, user); err != nil {
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

func Follow(s *State, cmd Command, user database.User) error {
	length := len(cmd.Args)
	if length < 1 {
		return fmt.Errorf("Invalid format\n\tUsage: follow <feed url>\n")
	} else if network.Validateargs("Follow", cmd.Args) {
		feedInfo, err_ := s.Db.GetFeedbyURL(context.Background(), cmd.Args[0])
		if err_ != nil {
			return fmt.Errorf("Something broke while fetching results from our records: %v\n", err_)
		}

		newFollow := database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feedInfo.ID,
		}
		_, err := s.Db.CreateFeedFollow(context.Background(), newFollow)
		if err != nil {
			return fmt.Errorf("Something went wrong on our end: %v\n", err)
		}
		fmt.Printf("%s is now following %s\n", user.Name, feedInfo.Name)
		return nil
	}
	return fmt.Errorf("Something went wrong...\n")
}

func addFeedFollow(s *State, cmd Command, user database.User) error {
	if network.Validateargs("AddFeed", cmd.Args) {
		feedInfo, err_ := s.Db.GetFeedbyURL(context.Background(), cmd.Args[1])
		if err_ != nil {
			return fmt.Errorf("Something broke while fetching results from our records: %v\n", err_)
		}

		newFollow := database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feedInfo.ID,
		}
		_, err := s.Db.CreateFeedFollow(context.Background(), newFollow)
		if err != nil {
			return fmt.Errorf("Something went wrong on our end: %v\n", err)
		}
		fmt.Printf("%s is now following %s\n", user.Name, feedInfo.Name)
		return nil
	}
	return fmt.Errorf("Something broke while attempting to follow %s\n", cmd.Args[0])
}

func Following(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Invalid format\n\t following takes no arguments\n")
	}
	follows, err := s.Db.GetUserFollows(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("Something went wrong in fetching results: %v\n", err)
	}
	if len(follows) < 1 {
		fmt.Println("You are not keeping up with any feeds at the moment.")
	} else {
		fmt.Println(user.Name, "is currently following")
		for _, follow := range follows {
			fmt.Printf("\t- %s\n", follow.FeedName)
		}
	}

	return nil
}

func UnFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Invalid format, please provide a valid url\n\t Usage: unfollow <url>\n")
	} else if network.Validateargs("UnFollow", cmd.Args) {
		feed, err := s.Db.GetFeedbyURL(context.Background(), cmd.Args[0])
		if err != nil {
			return err
		}
		unFollowParams := database.UnfollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		}
		if err := s.Db.Unfollow(context.Background(), unFollowParams); err != nil {
			return fmt.Errorf("Something went wrong: %w\n", err)
		}
		fmt.Printf("You are no longer following %s at %s\n", feed.Name, feed.Url)
	}
	return nil
}

func Browse(s *State, cmd Command, user database.User) error {
	var limit int

	//defaults to 2
	limit = 2

	if len(cmd.Args) >= 1 {
		val, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("Please provide a valid number to limit result by: %w\n", err)
		}
		limit = val
	}

	posts, err := s.Db.GetPosts(context.Background(), database.GetPostsParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return err
	}
	fmt.Println("Fetching posts...")
	if len(posts) < 1 {
		fmt.Println("You are not following any feeds")
		return nil
	}
	for _, post := range posts {
		if post.Description.Valid {
			fmt.Printf("- %s @ %s\n\t• %s\n", post.Title, post.Url, post.Description.String)
		} else {
			fmt.Printf("- %s @ %s\n", post.Title, post.Url)
		}
		
	}
	return nil
}