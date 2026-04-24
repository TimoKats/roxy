package pkg

import (
	"encoding/xml"
	"time"
)

type Format string

const (
	JSON Format = "json"
	XML  Format = "xml"
)

type Item struct {
	Title       string `xml:"title" json:"title"`
	Description string `xml:"description" json:"description"`
	Link        string `xml:"link" json:"link"`
	Guid        string `xml:"guid" json:"guid"`
	PubDate     string `xml:"pubDate" json:"pubDate"`
	// generated
	timestamp  time.Time
	parentFeed *Feed
}

type Channel struct {
	Title       string   `xml:"title" json:"title"`
	Description string   `xml:"description" json:"description"`
	Link        string   `xml:"link" json:"link"`
	Items       []Item   `xml:"item"`
	PubDate     string   `xml:"pubDate" json:"pubDate"`
	Category    []string `xml:"category" json:"category"`
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
		Size     int
	}
}

type Query struct {
	Urls       []string
	Keywords   []string
	Categories []string
	Amount     int
}

type Result struct {
	XMLName xml.Name `xml:"rss" json:"-"`
	Version string   `xml:"version,attr" json:"-"`
	Items   []*Item  `xml:"channel>item" json:"items"`
}
