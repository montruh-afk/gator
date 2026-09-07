package network

import (
	"html"
	"sync"
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
	var wg sync.WaitGroup

	wg.Add(2)

	//worker #1
	go func () {
		defer wg.Done()
		r.cleanMain()
	} ()
	
	//worker #2
	go func () {
		defer wg.Done()
		r.cleanTitle()
	} ()
	
	//wait for both workers to execute call to Done()
	wg.Wait()
}

func (r *RSSFeed) cleanMain() {
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)
}

func (r *RSSFeed) cleanTitle() {
	for i := range r.Channel.Item {
		r.Channel.Item[i].Title = html.UnescapeString(r.Channel.Item[i].Title)
		r.Channel.Item[i].Description = html.UnescapeString(r.Channel.Item[i].Description)
	}
}