package pkg

import (
	"encoding/xml"
	"time"
)

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Guid        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	// generated
	timestamp  time.Time
	parentFeed *Feed
}

type Channel struct {
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	Items       []Item   `xml:"item"`
	PubDate     string   `xml:"pubDate"`
	Category    []string `xml:"category"`
	Generator   string   `xml:"generator"`
	// generated
	timestamp time.Time
}

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
	// generated
	Category string
	Url      string
}

type Index struct {
	Rank []*Item
	Urls []struct {
		Url      string
		Category string
	}
}

type Query struct {
	Urls     []string
	Keywords []string
	Category string
	Amount   int
}

type Result struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Items   []*Item  `xml:"channel>item"`
}
