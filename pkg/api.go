package pkg

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"
)

type Api struct{}

// healthcheck endpoint, does nothing (useful)
func (api Api) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("roxy is running...")) //nolint:errcheck
	}
}

// endpoint for adding rss feeds through api
func (api Api) Add(idx *Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := getStrParam(r.URL, "category")
		urls := getListParam(r.URL, "urls")
		if len(urls) == 0 {
			http.Error(w, "no url", http.StatusBadRequest)
			return
		}
		for _, url := range urls {
			err := idx.Add(url, category)
			if err != nil {
				http.Error(w, "error: "+url, http.StatusInternalServerError)
				idx.Clear()
				return
			}
		}
		w.Write([]byte("add " + strings.Join(urls, ","))) //nolint:errcheck
	}
}

// query the rss feeds using url parameters, returns xml in the body
func (api Api) Get(idx *Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := Query{
			Urls:     getListParam(r.URL, "urls"),
			Category: getStrParam(r.URL, "category"),
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
			err := idx.Add(url.Url, url.Category)
			if err != nil {
				http.Error(w, "fail: "+url.Url, http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
