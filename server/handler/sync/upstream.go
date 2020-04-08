package sync

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"server/common/types"
	"server/config"
)

// create client with default 5secs timeout
var client = &http.Client{Timeout: time.Duration(5) * time.Second}

// create a request object which includes Basic Authentication
func createGetRequest(path string) (*http.Request, error) {
	url := config.MinifluxHost + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(config.MinifluxUser, config.MinifluxPass)
	return req, nil
}

// getCategories calls /v1/categories
func getCategories() (*types.CategoryList, error) {
	req, err := createGetRequest("/v1/categories")
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if http.StatusOK != res.StatusCode {
		return nil, fmt.Errorf("getCategories() error status code: %v", res.StatusCode)
	}

	var categories types.CategoryList
	json.NewDecoder(res.Body).Decode(&categories)
	return &categories, nil
}

// getFeedIcon calls /v1/feeds/:FeedID/icon
func getFeedIcon(feedID int64) (*types.FeedIcon, error) {
	req, err := createGetRequest(fmt.Sprintf("/v1/feeds/%v/icon", feedID))
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if http.StatusOK != res.StatusCode {
		return nil, fmt.Errorf("getFeedIcon(%v) error status code: %v", feedID, res.StatusCode)
	}

	var icon types.FeedIcon
	json.NewDecoder(res.Body).Decode(&icon)
	return &icon, nil
}

// getFeed calls /v1/feeds/:ID
func getFeed(id int64) (*types.Feed, error) {
	req, err := createGetRequest(fmt.Sprintf("/v1/feeds/%v", id))
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if http.StatusOK != res.StatusCode {
		return nil, fmt.Errorf("getFeed(%v) error status code: %v", id, res.StatusCode)
	}

	var mFeed types.MFeed
	json.NewDecoder(res.Body).Decode(&mFeed)
	feed := mFeed.ToFeed()

	// get feed icon, if error return feed without icon
	icon, err := getFeedIcon(id)
	if err != nil {
		return &feed, nil
	}
	feed.SetIcon(icon)
	return &feed, nil
}

// getEntry calls /v1/entries/:ID
func getEntry(id int64) (*types.Entry, error) {
	req, err := createGetRequest(fmt.Sprintf("/v1/entries/%v", id))
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if http.StatusOK != res.StatusCode {
		return nil, fmt.Errorf("getEntry(%v) error status code: %v", id, res.StatusCode)
	}

	var mEntry types.MEntry
	json.NewDecoder(res.Body).Decode(&mEntry)

	entry, err := mEntry.ToEntry()
	if err != nil {
		return nil, err
	}
	return entry, nil
}
