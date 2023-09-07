package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"goblin/internal/paste"
)

func NewRouter() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", paste.CreatePasteHandler).Methods("POST")
	r.HandleFunc("/{id}", paste.GetPasteHandler).Methods("GET")

	return r
}
