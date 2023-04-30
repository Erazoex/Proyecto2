package server

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}

	fmt.Fprintf(w, "API is working correctly!")
}
