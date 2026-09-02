package network

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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
	switch caller {
	case "Follow":
		if len(args) != 1 {
			return false
		}
		holdURL := make(chan struct{})
		go validateURL(args[0], holdURL)
		<- holdURL
		return true
		
	case "AddFeed":
		if len(args) != 2 {
			return false
		}
		holdName := make(chan struct{})
		holdURL := make(chan struct{})

		go validateName(args[0], holdName)
		go validateURL(args[1], holdURL)

		<-holdName
		<-holdURL
		return true
	}
	return false
}

func validateName(name string, hold chan struct{}) error {
	defer close(hold)
	_, err := url.ParseRequestURI(name)
	if err == nil {
		fmt.Println("Invalid format\n\tUsage: addfeed <feed name> <feed url> | follow <feed url> when calling 'follow'")
		os.Exit(1)
	}

	return nil
}

func validateURL(Url string, hold chan struct{}) error {
	defer close(hold)
	_, err := url.ParseRequestURI(Url)
	if err != nil {
		fmt.Println("Invalid format\n\tUsage: addfeed <feed name> <feed url> | follow <feed url> when calling 'follow'")
		os.Exit(1)
	}
	return nil
}
