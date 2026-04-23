package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/CusackJ4/gator/internal/config"
	"github.com/CusackJ4/gator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

// logs in a user (user must be registered)
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username Expected!")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Println("User doesn't exist")
		os.Exit(1)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User logged in!: %s!\n", cmd.args[0])

	return nil
}

// Registers a user into the database
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Name expected!")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		fmt.Println("Error at create user (Duplicate name?)")
		os.Exit(1)
	}
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User created! Data: %v\n", user)

	return nil
}

// Resets the DB to a clean slate (retains columns)
func reset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("Reset takes no additional args!")
	}
	err := s.db.ResetDB(context.Background())
	if err != nil {
		fmt.Println("Error resetting!")
		os.Exit(1)
	}
	fmt.Println("table reset!")
	os.Exit(0)

	return nil
}

// Prints the names in the DB
func getNames(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("Users takes no additional args!")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		println("Error retrieving users!")
		os.Exit(1)
	}

	for _, user := range users {
		if s.cfg.CurrentUsername == user {
			fmt.Printf("* %v (current)\n", user)
		} else {
			fmt.Printf("* %v\n", user)
		}
	}

	return nil
}

// add an rss feed
func addFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("addfeed needs two args: name and url")
	}
	/* currentUser := s.cfg.CurrentUsername // get current user
	user, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		return err
	} */

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		fmt.Println("Err creating feed follow!")

	}

	var feedMap map[string]interface{}
	data, _ := json.Marshal(feed)
	json.Unmarshal(data, &feedMap)

	for key, value := range feedMap {
		fmt.Printf("%s: %v\n", key, value)
	}
	return nil

}

// Aggregator command - fetches feeds
func agg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Agg requires a time string!")
	}

	timeSplit, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error parsing time duration. Example usage: 30s, 1h23m, 10ms")
	}
	ticker := time.NewTicker(timeSplit)

	for ; ; <-ticker.C {
		fmt.Printf("Fetching feeds every %s seconds\n", timeSplit)
		err := scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("Scraping error!")
		}
	}

}

// feeds overview
func viewFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("ViewFeeds takes no additional args!")
	}
	items, err := s.db.GetFeeds(context.Background()) // items is a slice of feeds
	if err != nil {
		println("Error retrieving feeds!")
		os.Exit(1)
	}
	for i := range items {
		name, err := s.db.GetUserName(context.Background(), items[i].UserID)
		if err != nil {
			fmt.Println("Username not found!")
			os.Exit(1)
		} else {
			fmt.Printf("Username: %v\n", name)
		}
		fmt.Printf("The feed name is: %v\n", items[i].Name)
		fmt.Printf("The feed url is: %v\n", items[i].Url)
	}
	return nil
}

func follow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow function takes a single url string")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0]) // get feed
	if err != nil {
		return fmt.Errorf("error retrieving URL in follow function: %w\n", err)
	}

	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error at feedFollowRow: %w\n", err)
	}
	fmt.Printf("Feed Name: %v\n, User Name: %v\n", feedFollowRow.FeedName, feedFollowRow.UserName)

	return nil
}

// used to show who the logged-in user is following
func following(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("func following takes no arguments!!")
	}

	data, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name) //returns a slice of GetFeedFollowsForUser rows, which is a struct
	if err != nil {
		return fmt.Errorf("Problem returning feed follows user data!")
	}

	if len(data) == 0 {
		fmt.Printf("%v is not following anyone!\n", user.Name)
		return nil
	}
	fmt.Printf("%v's feeds:\n", data[0].UserName)
	for i := range data {
		fmt.Printf(" - %v\n", data[i].FeedName)
	}

	return nil
}

// used to unfollow
func unfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Function Unfollows takes one argument... a url!")
	}

	err := s.db.Unfollow(context.Background(), database.UnfollowParams{
		Name: user.Name,
		Url:  cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("Error unfollowing!")
	}

	return nil

}

// used to browse posts
func browse(s *state, cmd command, user database.User) error {
	lmt := 2

	if len(cmd.args) > 0 {
		if len(cmd.args) > 1 {
			return fmt.Errorf("The browse function only takes an a single integer as a LIMIT argument.\n")
		}
		if val, err := strconv.Atoi(cmd.args[0]); err == nil {
			lmt = val
		} else {
			return fmt.Errorf("The argument provided, %v, is not a valid integer.\n", cmd.args[0])
		}
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  lmt,
	})
	if err != nil {
		return fmt.Errorf("Error getting posts! %v\n", err)
	}
	for i := range posts {
		fmt.Printf("Feed: %v\n", posts[i].FeedName)
		fmt.Printf("Title: %v\n", posts[i].Title.String)
		fmt.Printf("Publish Date: %v\n", posts[i].PublishedAt.Time)
		fmt.Printf("URL: %v\n", posts[i].Url)
		fmt.Printf("Description: %v\n", posts[i].Description.String)
	}

	return nil
}
