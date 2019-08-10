package main

import (
	"fmt"

	"function/sync"
)

// This package is not part of Cloud Function but used to test parts of the project.
// Example usage:
// MINIFLUX_HOST=http://gatrabali.com MINIFLUX_USER=user MINIFLUX_PASS=pass go run main/main.go
func main() {
	// -- Get categories
	cats, err := sync.GetCategories()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(cats)
	}

	// -- Get Entry
	entry, err := sync.GetEntry(2000)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(entry.Enclosures)
	}

	// -- Get Feed
	feed, err := sync.GetFeed(15)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(feed)
	}

	// -- Get Feed Icon
	icon, err := sync.GetFeedIcon(15)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Icon: {%v %v}\n", icon.ID, icon.MimeType)
	}

}
