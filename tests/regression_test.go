package tests

import (
	"encoding/xml"
	"net/http"
	"testing"
	"time"

	"github.com/TimoKats/roxy/pkg"
)

func TestAdd(t *testing.T) {
	idx := pkg.NewIndex()
	url := "https://timokats.xyz/feed/website.xml"
	tags := []string{"example", "test"}
	err := idx.Add(url, tags)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(idx.Rank) == 0 {
		t.Error("Rank should not be empty after adding an item")
	}
}

func TestGet(t *testing.T) {
	idx := pkg.NewIndex()
	url := "https://timokats.xyz/feed/website.xml"
	tags := []string{"example", "test"}
	err := idx.Add(url, tags)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	query := pkg.Query{Amount: 10}
	result := idx.Get(query)
	if len(result.Items) == 0 {
		t.Error("Result should contain items after querying")
	}
}

func TestServe(t *testing.T) {
	idx := pkg.NewIndex()
	url := "https://timokats.xyz/feed/website.xml"
	tags := []string{"example", "test"}
	err := idx.Add(url, tags)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	go idx.Serve(":8080")
	time.Sleep(2 * time.Second)

	resp, err := http.Get("http://localhost:8080/get?amount=10")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	result := pkg.Result{}
	err = xml.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result.Items) == 0 {
		t.Error("Result should contain items after querying via API")
	}
}

func TestClear(t *testing.T) {
	idx := pkg.NewIndex()
	url := "https://timokats.xyz/feed/website.xml"
	tags := []string{"example", "test"}
	err := idx.Add(url, tags)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	idx.Clear()
	if len(idx.Rank) > 0 {
		t.Error("Rank should be empty after clearing")
	}
}
