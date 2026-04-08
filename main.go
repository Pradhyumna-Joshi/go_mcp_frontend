package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Pradhyumna-Joshi/go_mcp_frontend/components"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		components.Base().Render(r.Context(), w)
	})
	mux.HandleFunc("/newChat", func(w http.ResponseWriter, r *http.Request) {
		components.ChatHome("Pradhyumna").Render(r.Context(), w)
	})

	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		components.ConfigMain().Render(r.Context(), w)
	})

	mux.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		components.Settings().Render(r.Context(), w)
	})

	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		components.About().Render(r.Context(), w)
	})

	mux.HandleFunc("/chatMessage", func(w http.ResponseWriter, r *http.Request) {
		input := r.FormValue("input")

		time.Sleep(2 * time.Second)
		components.AIMsgBubble(input).Render(r.Context(), w)
	})
	log.Println("Server running on port 8000")
	http.ListenAndServe(":8000", mux)
}
