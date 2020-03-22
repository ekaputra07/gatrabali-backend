package common

import (
	"encoding/json"
	"testing"
)

func TestToFeedTransformation(t *testing.T) {
	blob := []byte(`
	{
		"id": 2,
		"user_id": 1,
		"feed_url": "http://www.example.com/bali/badung/feed",
		"site_url": "http://www.example.com",
		"title": "example.com - Badung",
		"checked_at": "2019-06-16T04:54:24.861894Z",
		"etag_header": "",
		"last_modified_header": "Sun, 16 Jun 2019 04:53:09 GMT",
		"parsing_error_message": "",
		"parsing_error_count": 0,
		"scraper_rules": "",
		"rewrite_rules": "",
		"crawler": true,
		"user_agent": "",
		"username": "",
		"password": "",
		"category": {
		  "id": 2,
		  "title": "Badung",
		  "user_id": 1
		},
		"icon": {
		  "feed_id": 2,
		  "icon_id": 1
		}
	  }
	`)

	var mFeed MFeed
	if err := json.Unmarshal(blob, &mFeed); err != nil {
		t.Error(err)
	}

	feed := mFeed.ToFeed()
	if feed.ID != 2 {
		t.Error("Wrong ID value")
	}
	if feed.UserID != 1 {
		t.Error("Wrong UserID value")
	}
	if feed.SiteURL != "http://www.example.com" {
		t.Error("Wrong SiteURL value")
	}
	if feed.FeedURL != "http://www.example.com/bali/badung/feed" {
		t.Error("Wrong FeedURL value")
	}
	if feed.Title != "example.com - Badung" {
		t.Error("Wrong Title value")
	}
	if feed.CheckedAt != "2019-06-16T04:54:24.861894Z" {
		t.Error("Wrong CheckedAt value")
	}
	if feed.Category != 2 {
		t.Error("Wrong Category value")
	}
}

func TestToEntryTransformation(t *testing.T) {
	blob := []byte(`
	{
		"id": 1000,
		"user_id": 1,
		"feed_id": 15,
		"status": "read",
		"hash": "474a2a1c8ae65b9e591d1b9f16c914b5bb8c411ccb34e8ea5d4b4d15ab47709c",
		"title": "Test Title",
		"url": "http://url/",
		"comments_url": "http://comments/",
		"published_at": "2019-05-04T09:59:25Z",
		"content": "Test Content",
		"author": "Example",
		"starred": false,
		"enclosures": [
		  {
			"id": 0,
			"user_id": 0,
			"entry_id": 0,
			"url": "http://image.jpg",
			"mime_type": "image/jpg",
			"size": 0
		  }
		],
		"feed": {
		  "id": 15,
		  "user_id": 1,
		  "feed_url": "http://feed/",
		  "site_url": "http://site/",
		  "title": " Example.com - Bangli",
		  "checked_at": "2019-06-16T03:54:25.179466Z",
		  "etag_header": "",
		  "last_modified_header": "",
		  "parsing_error_message": "",
		  "parsing_error_count": 0,
		  "scraper_rules": "",
		  "rewrite_rules": "",
		  "crawler": true,
		  "user_agent": "",
		  "username": "",
		  "password": "",
		  "category": {
			"id": 3,
			"title": "Bangli",
			"user_id": 1
		  },
		  "icon": {
			"feed_id": 15,
			"icon_id": 3
		  }
		}
	  }
	`)
	var mEntry MEntry
	if err := json.Unmarshal(blob, &mEntry); err != nil {
		t.Error(err)
	}

	entry, _ := mEntry.ToEntry()
	if entry.ID != 1000 {
		t.Error("Wrong ID value")
	}
	if entry.UserID != 1 {
		t.Error("Wrong UserID value")
	}
	if entry.FeedID != 15 {
		t.Error("Wrong FeedID value")
	}
	if entry.Title != "Test Title" {
		t.Error("Wrong Title value")
	}
	if entry.URL != "http://url/" {
		t.Error("Wrong URL value")
	}
	if entry.Content != "Test Content" {
		t.Error("Wrong Content value")
	}
	if *entry.CommentsURL != "http://comments/" {
		t.Error("Wrong CommentsURL value")
	}
	if *entry.Author != "Example" {
		t.Error("Wrong Author value")
	}
	if len(*entry.Enclosures) != 1 {
		t.Error("Wrong Enclosures value")
	}
	if entry.PublishedAt != 1556963965000 {
		t.Error("Wrong PublishedAt value")
	}
	if len(entry.Categories) != 1 {
		t.Error("Wrong Categories value")
	}
}
