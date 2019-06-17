package gatrabali

import (
	"fmt"
	"net/http"
	// "encoding/json"
	// "html"
)

// Hello return JSON encoded hello world
func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}
