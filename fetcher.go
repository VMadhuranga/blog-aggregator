package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/VMadhuranga/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func startFetchingFeeds(limit int, interval time.Duration, dbq *database.Queries) {
	log.Printf("Fetching feeds on %v goroutines every %v", limit, interval)
	ticker := time.NewTicker(interval)

	for ; ; <-ticker.C {
		feeds, err := dbq.GetNextFeedsToFetch(context.Background(), int32(limit))

		if err != nil {
			log.Printf("Error getting next feed to fetch: %s", err)
			continue
		}

		wg := sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)

			go func() {
				defer wg.Done()
				_, err := dbq.MarkFeedAsFetched(context.Background(), feed.ID)

				if err != nil {
					log.Printf("Error marking feed as fetched: %s", err)
					return
				}

				rssFeed, err := parseXmlToRssFeed(feed.Url)

				if err != nil {
					log.Printf("Error parsing rss feed: %s", err)
					return
				}

				for _, item := range rssFeed.Channel.Item {
					description := sql.NullString{}

					if item.Description != "" {
						description.String = item.Description
						description.Valid = true
					}

					pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)

					if err != nil {
						log.Printf("Error parsing pub date: %s", err)
					}

					_, err = dbq.CreatePost(context.Background(), database.CreatePostParams{
						ID:          uuid.New(),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Title:       item.Title,
						Url:         item.Link,
						Description: description,
						PublishedAt: pubDate,
						FeedID:      feed.ID,
					})

					if err != nil {
						log.Printf("Error creating post: %s", err)
					}
				}

				log.Printf("%v posts collected from %s", len(rssFeed.Channel.Item), feed.Name)
			}()
		}

		wg.Wait()
	}
}
