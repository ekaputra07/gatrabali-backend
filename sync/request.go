package sync

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"gatrabali/model"
)

// create client with default 5secs timeout
var client = &http.Client{Timeout: time.Duration(5) * time.Second}

// create a request object which includes Basic Authentication
func createGetRequest(path string) (*http.Request, error) {
	url := os.Getenv("MINIFLUX_HOST") + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(os.Getenv("MINIFLUX_USER"), os.Getenv("MINIFLUX_PASS"))
	return req, nil
}

// GetCategories calls /v1/categories
func GetCategories() (*model.CategoryList, error) {
	req, err := createGetRequest("/v1/categories")
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var categories model.CategoryList
	json.NewDecoder(res.Body).Decode(&categories)
	return &categories, nil
}

// GetIcon calls /v1/feeds/:FeedID/icon
func GetIcon(id int64) (*model.Feed, error) {
	return nil, nil
}

// GetFeed calls /v1/feeds/:ID
func GetFeed(id int64) (*model.Feed, error) {
	return nil, nil
}

// GetEntry calls /v1/entries/:ID
func GetEntry(id int64) (*model.Entry, error) {
	req, err := createGetRequest(fmt.Sprintf("/v1/entries/%v", id))
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var mEntry model.MEntry
	json.NewDecoder(res.Body).Decode(&mEntry)

	entry := mEntry.ToEntry()
	return &entry, nil
}
