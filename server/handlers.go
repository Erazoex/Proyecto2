package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/erazoex/proyecto2/analizador"
	"github.com/erazoex/proyecto2/consola"
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}

	fmt.Fprintf(w, "API is working correctly!")
}

func getConsole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newContent := &Response{}
	err := json.NewDecoder(r.Body).Decode(newContent)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}
	// aqui mandar a el analizador todo lo leido en newContent
	var analizador analizador.Analyzer
	analizador.Exe(newContent.Content)
	json.NewEncoder(w).Encode(consola.GetConsole())
}
