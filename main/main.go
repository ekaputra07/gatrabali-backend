package main

import (
	"fmt"

	"gatrabali/sync"
)

// This package is not part of Cloud Function but used to test parts of the project.
// Example usage:
// MINIFLUX_HOST=http://gatrabali.com MINIFLUX_USER=user MINIFLUX_PASS=pass go run main/main.go
func main() {
	// -- Get categories
	cats, err := sync.GetCategories()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cats)

	// -- Get Entry
	entry, err := sync.GetEntry(2000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(entry.Enclosures)
}
