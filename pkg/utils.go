package pkg

import (
	"net/url"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

// check if two lists have any overlap. Useful for querying.
func overlap[Type comparable](a []Type, b []Type) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	for _, aItem := range a {
		if slices.Contains(b, aItem) {
			return true
		}
	}
	return false
}

// tries all mentally sane rss datetime formats and returns time object
func parsePubDate(s string) time.Time {
	var rssDateFormats = []string{
		"02 Jan 2006 15:04:05 MST",
		"02 Jan 2006 15:04 MST",
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}
	for _, format := range rssDateFormats {
		if t, err := time.Parse(format, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// parses a newsboat URL line and gets the url and tags
func parseLine(line string) (string, []string) {
	parts := strings.Split(line, " ")
	url, err := url.Parse(parts[0])
	tags := ""
	if err == nil {
		if len(parts) > 1 {
			tags = parts[1]
		}
		return url.String(), []string{tags}
	}
	return "", []string{} // no valid URL found
}

// inserts an item in sorted order. Also returns index it was inserted at.
func insertSorted(ranking []*Item, item *Item) []*Item {
	i := sort.Search(len(ranking), func(i int) bool {
		return ranking[i].timestamp.Before(item.timestamp)
	})
	ranking = append(ranking, nil)
	copy(ranking[i+1:], ranking[i:])
	ranking[i] = item
	return ranking
}

// gets list params from url, and takes out bad values
func getListParam(url *url.URL, param string) []string {
	params := strings.Split(url.Query().Get(param), ",")
	filteredParams := make([]string, 0)
	for _, p := range params {
		if p != "" {
			filteredParams = append(filteredParams, strings.ToLower(p))
		}
	}
	return filteredParams
}

// gets integer param from url, and takes out bad values
func getIntParam(url *url.URL, param string, fallback int) int {
	strValue := url.Query().Get(param)
	if strValue == "" {
		return fallback
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return fallback
	}
	return intValue
}
