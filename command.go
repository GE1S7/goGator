package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/GE1S7/gator/internal/config"
	"github.com/GE1S7/gator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db   *database.Queries
	conf *config.Config
}

type command struct {
	name string
	args []string
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.conf.UserName)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		return handler(s, cmd, user)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Not enough arguments: username required.")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("No user under such name extists")
	}

	err = s.conf.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("Username has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Not enough arguments: username required.")
	}
	usrData := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	_, err := s.db.GetUser(context.Background(), usrData.Name)
	if err == nil {
		fmt.Println(fmt.Errorf("Usr with name %s already exists", usrData.Name))
		os.Exit(1)
	}

	_, err = s.db.CreateUser(context.Background(), usrData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = s.conf.SetUser(usrData.Name)
	if err != nil {
		return err
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return nil

}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	logged, err := s.db.GetUser(context.Background(), s.conf.UserName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for _, e := range users {
		if e == logged {
			fmt.Println("*", e.Name, "(current)")
			continue
		}
		fmt.Println("*", e.Name)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Feeds:")
	for _, e := range feeds {
		fmt.Println("username:", e.Name_2.String, "feedname:", e.Name, "url:", e.Url)
	}

	return nil
}

func handlerAgg(state *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	decodeHtmlEntities(feed)

	fmt.Println(feed)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		fmt.Println("Error: addfeed takes exactly two arguments")
		os.Exit(1)
	}

	feedData := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	_, err := s.db.CreateFeed(context.Background(), feedData)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	newfeedfollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedData.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), newfeedfollow)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return nil
}

func handlerFollow(s *state, cmd command, current_user database.User) error {
	if len(cmd.args) != 1 {
		fmt.Println("Error: follow takes one url argument")
		os.Exit(1)
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])

	newfeedfollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    current_user.ID,
		FeedID:    feed.ID,
	}
	feedFollowFmtNames, err := s.db.CreateFeedFollow(context.Background(), newfeedfollow)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(feedFollowFmtNames.Name)
	fmt.Println(feedFollowFmtNames.Name_2)

	return nil
}

func handlerFollowing(s *state, cmd command, current_user database.User) error {
	feed_follow, err := s.db.GetFeedFollowsForUser(context.Background(), current_user.ID)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for _, e := range feed_follow {
		feed, err := s.db.GetFeedByID(context.Background(), e.FeedID)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println(feed.Name)
	}

	return nil
}

type commands struct {
	commandFunctions map[string]func(state *state, cmd command) error
}

func handlerUnfollow(s *state, cmd command, current_user database.User) error {
	if len(cmd.args) != 1 {
		fmt.Println("Error: unfollow takes exatcly one argument")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	feedfollow := database.DeleteFeedFollowParams{
		UserID: current_user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), feedfollow)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return nil
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
