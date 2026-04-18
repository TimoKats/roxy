package pkg

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"
)

type Api struct{}

// healthcheck endpoint
func (api Api) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("roxy is running...")) //nolint:errcheck
	}
}

// endpoint for adding rss feeds through api
func (api Api) Add(idx *Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls := getListParam(r.URL, "urls")
		tags := getListParam(r.URL, "tags")
		if len(urls) == 0 {
			http.Error(w, "no url", http.StatusBadRequest)
			return
		}
		for _, url := range urls {
			err := idx.Add(url, tags)
			if err != nil {
				http.Error(w, "error: "+url, http.StatusInternalServerError)
				idx.Clear()
				return
			}
		}
		w.Write([]byte("added " + strings.Join(urls, ", "))) //nolint:errcheck
	}
}

// /get endpoint to query the rss feeds, returns xml in the body
func (api Api) Get(idx *Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := Query{
			Urls:     getListParam(r.URL, "urls"),
			Tags:     getListParam(r.URL, "tags"),
			Keywords: getListParam(r.URL, "keywords"),
			Amount:   getIntParam(r.URL, "amount", 10),
		}
		result := idx.Get(query)
		xmlData, _ := xml.MarshalIndent(result, "", "\t")
		w.Header().Set("Content-Type", "application/xml")
		w.Write(xmlData) //nolint:errcheck
	}
}

// refreshes the index by removing all entries, and fetching the feed.
func (api Api) Refresh(idx *Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idx.Clear()
		for _, url := range idx.Urls {
			log.Printf("refreshing: %s", url)
			err := idx.Add(url, []string{}) // TODO! RESET TAGS!
			if err != nil {
				http.Error(w, "can't refresh: "+url, http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
