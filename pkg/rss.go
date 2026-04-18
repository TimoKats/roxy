package pkg

import (
	"regexp"
	"slices"
	"strings"
)

// extracts keywords from title, used for querying based on keywords
func (item *Item) Keywords() []string {
	re := regexp.MustCompile(`[a-zA-Z]+`)
	words := re.FindAllString(strings.ToLower(item.Title), -1)
	keywords := []string{}
	for _, w := range words {
		if len(w) >= 4 {
			keywords = append(keywords, w)
		}
	}
	return keywords
}

// for each query aspect (url, tags, keywords), check if it fails
// ps, King Terry said case/switch are devine, hence the choice
func (item *Item) QueryMatch(query Query) bool {
	switch {
	case len(query.Urls) > 0 && !slices.Contains(query.Urls, item.parentFeed.Url):
		return false
	case len(query.Tags) > 0 && !overlap(query.Tags, item.parentFeed.Tags):
		return false
	case len(query.Keywords) > 0 && !overlap(query.Keywords, item.Keywords()):
		return false
	default:
		return true
	}
}

// converts all time strings to datetime objects
func (feed *Feed) ParseTime() {
	feed.Channel.timestamp = parsePubDate(feed.Channel.PubDate)
	for i, item := range feed.Channel.Items {
		timestamp := parsePubDate(item.PubDate)
		feed.Channel.Items[i].timestamp = timestamp
	}
}
