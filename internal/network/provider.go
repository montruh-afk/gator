package network

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"github.com/montruh-afk/gator/internal/database"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func Fetchfeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	feed := &RSSFeed{}
	if err := xml.Unmarshal(data, feed); err != nil {
		return nil, err
	}

	done := make(chan struct{})
	go feed.cleaner(done)

	<-done

	return feed, nil
}

func Validateargs(caller string, args []string) bool {
	var wg sync.WaitGroup
	switch caller {
	case "Follow":
		if len(args) != 1 {
			return false
		}
		wg.Add(1)
		validateURL(args[0], &wg)
		return true

	case "AddFeed":
		if len(args) != 2 {
			return false
		}
		wg.Add(2)
		go validateName(args[0], &wg)
		go validateURL(args[1], &wg)
		wg.Wait()
		return true

	case "UnFollow":
		if len(args) != 1 {
			return false
		}
		wg.Add(1)
		validateURL(args[0], &wg)
		return true
	}
	return false
}

func concurrentPost(s *database.Queries, feedID uuid.UUID, feed *RSSItem) {
	now := time.Now()
	post := database.CreatePostParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Title:     feed.Title,
		Url:       feed.Link,
		FeedID:    feedID,
	}

	getdescription(&post, feed)
	parseTime(&post, feed)

	res, err := s.CreatePost(context.Background(), post)
	if err != nil {
		if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			log.Printf("Couldn't create post: %v", err)
		}
		return
	}
	fmt.Println("Added Title •", res.Title)

}

func getdescription(post *database.CreatePostParams, feed *RSSItem) {
	desc := sql.NullString{
		String: feed.Description,
		Valid:  feed.Description != "",
	}
	post.Description = desc
}

func parseTime(post *database.CreatePostParams, feed *RSSItem) {

	var pubdate time.Time
	var err error

	pubdate, err = time.Parse(time.RFC1123Z, feed.PubDate)
	if err != nil {
		pubdate, err = time.Parse(time.RFC1123, feed.PubDate)
	}
	if err != nil {
		log.Printf("Something went wrong while attempting to set Publish Date for title: %s\n\t %v\n", feed.Title, err)
		return
	}
	post.PublishedAt = sql.NullTime{
		Time:  pubdate,
		Valid: true,
	}
}

func ScrapeFeeds(s *database.Queries) error {
	parentFeed, err := s.GetNextFeed(context.Background())
	if err != nil {
		return fmt.Errorf("Something broke on our end: %w\n", err)
	}
	if err := s.MarkFeed(context.Background(), parentFeed.ID); err != nil {
		return fmt.Errorf("Something went wrong: %w\n", err)
	}

	feeds, err := Fetchfeed(context.Background(), parentFeed.Url)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i := range feeds.Channel.Item {
		wg.Add(1)

		go func() {
			wg.Done()
			concurrentPost(s, parentFeed.ID, &feeds.Channel.Item[i])
		}()

	}
	wg.Wait()
	return nil
}

func validateName(name string, wg *sync.WaitGroup) error {
	defer wg.Done()
	_, err := url.ParseRequestURI(name)
	if err == nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func validateURL(Url string, wg *sync.WaitGroup) error {
	defer wg.Done()
	_, err := url.ParseRequestURI(Url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}
