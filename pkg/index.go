package pkg

import (
	"bufio"
	"encoding/xml"
	"log"
	"net/http"
	"os"
)

// gets rss feed from a url and adds it to the index, and parses the pubdates
func (idx *Index) Add(url string, category string) error {
	var feed Feed
	resp, err := http.Get(url) //nolint:errcheck
	if err != nil {
		return err
	}
	if err = xml.NewDecoder(resp.Body).Decode(&feed); err == nil {
		feed.Category = category
		feed.Url = url
		feed.ParseTime()
		for _, item := range feed.Channel.Items {
			item.parentFeed = &feed
			idx.Rank = insertSorted(idx.Rank, &item)
		}
		idx.Urls = append(idx.Urls, Url{url, category, len(feed.Channel.Items)})
		log.Printf("added to feed: '%s' %v", url, category)
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
func (idx *Index) Load(filename string) error {
	if filename == "" {
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("can't open: %s", filename)
		return err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if url, category := parseLine(line); len(url) > 0 {
			if err := idx.Add(url, category); err != nil {
				log.Printf("error adding url: '%s'", url)
			}
		}
	}
	return nil
}

// servers all api endpoints for an index instance
func (idx *Index) Serve(port string) {
	api := Api{}
	log.Printf("serving on: http://localhost%s", port)
	http.HandleFunc("/", api.Ping())
	http.HandleFunc("/stats", api.Stats(idx))
	http.HandleFunc("/add", api.Add(idx))
	http.HandleFunc("/xml", api.Get(idx, XML))
	http.HandleFunc("/json", api.Get(idx, JSON))
	http.HandleFunc("/refresh", api.Refresh(idx))
	log.Fatal(http.ListenAndServe(port, nil))
}

// removes all entries from rank, return copy of urls for re-create.
func (idx *Index) Clear() {
	idx.Rank = nil
	idx.Urls = nil
}

// initiate rss feed index class (enforce singleton?)
func NewIndex() *Index {
	log.Println("starting roxy...")
	return &Index{}
}
