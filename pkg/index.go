package pkg

import (
	"bufio"
	"encoding/xml"
	"log"
	"net/http"
	"os"
)

// gets rss feed from a url and adds it to the index, and parses the pubdates
func (idx *Index) Add(url string, tags []string) error {
	var feed Feed
	resp, err := http.Get(url) //nolint:errcheck
	if err != nil {
		return err
	}
	if err = xml.NewDecoder(resp.Body).Decode(&feed); err == nil {
		feed.ParseTime()
		feed.Tags = tags
		feed.Url = url
		for _, item := range feed.Channel.Items {
			item.parentFeed = &feed
			idx.Rank = insertSorted(idx.Rank, &item)
		}
		idx.Urls = append(idx.Urls, url)
		log.Printf("added to feed: '%s' %v", url, tags)
	}
	return err
}

// queries index and returns result object (list of items)
func (idx *Index) Get(query Query) Result {
	result := []*Item{}
	hits := 0
	for _, item := range idx.Rank {
		if item.QueryMatch(query) {
			result = append(result, item)
			hits += 1
		}
		if hits >= query.Amount {
			break
		}
	}
	return Result{
		XMLName: xml.Name{Space: "", Local: "rss"},
		Version: "2.0",
		Items:   result,
	}
}

// loads (newsboat) file rss feeds into the index
func (idx *Index) Load(filename string) {
	if filename == "" {
		return
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("can't open: %s", filename)
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if url, tags := parseLine(line); len(url) > 0 {
			if err := idx.Add(url, tags); err != nil {
				log.Printf("error adding url: '%s'", url)
			}
		}
	}
}

// servers all api endpoints for an index instance
func (idx *Index) Serve(port string) {
	api := Api{}
	log.Printf("serving on: http://localhost%s", port)
	http.HandleFunc("/", api.Ping())
	http.HandleFunc("/add", api.Add(idx))
	http.HandleFunc("/get", api.Get(idx))
	http.HandleFunc("/refresh", api.Refresh(idx))
	log.Fatal(http.ListenAndServe(port, nil))
}

// removes all entries from rank.
func (idx *Index) Clear() {
	idx.Rank = nil
}

// initiate rss feed index class (enforce singleton?)
func NewIndex() *Index {
	log.Println("starting roxy...")
	return &Index{}
}
