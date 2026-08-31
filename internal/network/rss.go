package network

import (
	"html"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (r *RSSFeed) cleaner(done chan struct{}) error {
	defer close(done)
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)
	for i := range r.Channel.Item {
		r.Channel.Item[i].Title = html.UnescapeString(r.Channel.Item[i].Title)
		r.Channel.Item[i].Description = html.UnescapeString(r.Channel.Item[i].Description)
	}
	return nil
}
