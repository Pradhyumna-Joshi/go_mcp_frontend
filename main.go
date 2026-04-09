package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Pradhyumna-Joshi/go_mcp_frontend/components"
	"github.com/gorilla/websocket"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var upgrader websocket.Upgrader

func main() {

	llm, err := ollama.New(ollama.WithModel("mistral"))
	if err != nil {
		log.Fatal(err)
	}

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

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
			}

			log.Println(string(msg))
			llm.Call(r.Context(), string(msg), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				conn.WriteMessage(websocket.TextMessage, chunk)
				return nil
			}))

			conn.WriteMessage(websocket.TextMessage, []byte("[DONE]"))
		}

	})

	log.Println("Server running on port 8000")
	http.ListenAndServe(":8000", mux)
}
