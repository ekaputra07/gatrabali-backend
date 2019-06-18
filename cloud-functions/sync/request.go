package sync

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/apps4bali/gatrabali-backend/common/model"
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

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("GetCategories() error status code: %v", res.StatusCode)
	}

	var categories model.CategoryList
	json.NewDecoder(res.Body).Decode(&categories)
	return &categories, nil
}

// GetFeedIcon calls /v1/feeds/:FeedID/icon
func GetFeedIcon(feedID int64) (*model.FeedIcon, error) {
	req, err := createGetRequest(fmt.Sprintf("/v1/feeds/%v/icon", feedID))
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("GetFeedIcon(%v) error status code: %v", feedID, res.StatusCode)
	}

	var icon model.FeedIcon
	json.NewDecoder(res.Body).Decode(&icon)
	return &icon, nil
}

// GetFeed calls /v1/feeds/:ID
func GetFeed(id int64) (*model.Feed, error) {
	req, err := createGetRequest(fmt.Sprintf("/v1/feeds/%v", id))
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("GetFeed(%v) error status code: %v", id, res.StatusCode)
	}

	var mFeed model.MFeed
	json.NewDecoder(res.Body).Decode(&mFeed)
	feed := mFeed.ToFeed()

	// get feed icon, if error return feed without icon
	icon, err := GetFeedIcon(id)
	if err != nil {
		return &feed, nil
	}
	feed.SetIcon(icon)
	return &feed, nil
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

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("GetEntry(%v) error status code: %v", id, res.StatusCode)
	}

	var mEntry model.MEntry
	json.NewDecoder(res.Body).Decode(&mEntry)

	entry := mEntry.ToEntry()
	return &entry, nil
}
