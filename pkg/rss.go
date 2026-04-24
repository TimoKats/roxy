package pkg

import (
	"regexp"
	"strings"
)

// extracts keywords from title, used for querying based on keywords
func (item *Item) Keywords() []string {
	re := regexp.MustCompile(`[a-zA-Z]+`)
	words := re.FindAllString(strings.ToLower(item.Description), -1)
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
	case !contains(query.Urls, item.parentFeed.Url):
		return false
	case !contains(query.Categories, item.parentFeed.Category):
		return false
	case !overlap(query.Keywords, item.Keywords()):
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
