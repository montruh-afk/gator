package network

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
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

