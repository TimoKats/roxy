package pkg

import (
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

func (api Api) Stats(idx *Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, header := marshall(idx.Urls, JSON)
		w.Header().Set("Content-Type", header)
		w.Write(data) //nolint:errcheck
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
func (api Api) Get(idx *Index, format Format) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := Query{
			Urls:       getListParam(r.URL, "urls"),
			Categories: getListParam(r.URL, "category"),
			Keywords:   getListParam(r.URL, "keywords"),
			Amount:     getIntParam(r.URL, "amount", 10),
		}
		result := idx.Get(query)
		data, header := marshall(result, format)
		w.Header().Set("Content-Type", header)
		w.Write(data) //nolint:errcheck
	}
}

// refreshes the index by removing all entries, and fetching the feed.
func (api Api) Refresh(idx *Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tempUrls := make([]Url, len(idx.Urls))
		copy(tempUrls, idx.Urls)
		idx.Clear()
		for _, url := range tempUrls {
			log.Printf("refreshing: %s", url.Url)
			err := idx.Add(url.Url, url.Category)
			if err != nil {
				http.Error(w, "fail: "+url.Url, http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
