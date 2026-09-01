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

func (r *RSSFeed) cleaner(done chan struct{}) {
	defer close(done)
	hold := make(chan struct{})
	hold_ := make(chan struct{})
	go r.cleanMain(hold)
	go r.cleanTitle(hold_)
}

func (r *RSSFeed) cleanMain(hold chan struct{}) {
	defer close(hold)
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)
}

func (r *RSSFeed) cleanTitle(hold chan struct{}) {
	defer close(hold)
	for i := range r.Channel.Item {
		r.Channel.Item[i].Title = html.UnescapeString(r.Channel.Item[i].Title)
		r.Channel.Item[i].Description = html.UnescapeString(r.Channel.Item[i].Description)
	}
}