package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/cusackj4/rssAggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	data, err := s.db.GetNextFeed(ctx)
	if err != nil {
		return fmt.Errorf("Error Getting feed")
	}

	err = s.db.SetLastFetched(ctx, data.ID)
	if err != nil {
		log.Printf("SetLastFetched failed: %v", err)
		return fmt.Errorf("Error fetching feed!")
	}
	feedContent, err := fetchFeed(ctx, data.Url)
	if err != nil {
		return fmt.Errorf("Error fetching feed content!")
	}

	for i := range feedContent.Channel.Item {
		//fmt.Printf("%v\n", feedContent.Channel.Item[i].Title)
		fmt.Printf("%v\n", feedContent.Channel.Item[i].Link)
		// fmt.Printf("%v\n", feedContent.Channel.Item[i].PubDate)
		//fmt.Printf("%v\n\n\n", feedContent.Channel.Item[i].Description)

		_, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: feedContent.Channel.Item[i].Title, Valid: feedContent.Channel.Item[i].Title != ""},
			Url:         feedContent.Channel.Item[i].Link,
			Description: sql.NullString{String: feedContent.Channel.Item[i].Description, Valid: feedContent.Channel.Item[i].Description != ""},
			PublishedAt: parsePubTime(feedContent.Channel.Item[i].PubDate),
			FeedID:      data.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			log.Printf("Couldn't create post %v", err)
		}
	}

	return nil
}

func parsePubTime(dateStr string) sql.NullTime {
	if dateStr == "" {
		return sql.NullTime{Valid: false}
	}

	t, err := time.Parse(time.RFC1123Z, dateStr)
	return sql.NullTime{Time: t, Valid: err == nil}
}
