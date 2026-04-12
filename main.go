package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/Pradhyumna-Joshi/go_mcp_frontend/components"
	"github.com/Pradhyumna-Joshi/go_mcp_frontend/models"
	"github.com/gorilla/websocket"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var upgrader websocket.Upgrader

var (
	serverConf = make([]models.MCPServerConfig, 2)
	modelConf  = make([]models.ModelConfig, 2)
)

func main() {

	titles := []string{
		"What's on your mind?",
		"Where should we begin?",
		"System ready. What’s next?",
		"What's the goal?",
		"How can Nexus help today?",
		"What's on your agenda?",
	}

	llm, err := ollama.New(ollama.WithModel("mistral"))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /getmodelconf", func(w http.ResponseWriter, r *http.Request) {
		components.ConfigMain(modelConf).Render(r.Context(), w)
	})

	mux.HandleFunc("POST /setmodelconf", func(w http.ResponseWriter, r *http.Request) {
		modelConf[0].Model = r.FormValue("model1")
		modelConf[0].APIKey = r.FormValue("api_key1")
		modelConf[0].BaseURL = r.FormValue("base_url1")
		modelConf[0].Name = r.FormValue("name1")
		modelConf[0].Provider = r.FormValue("provider1")

		modelConf[1].Model = r.FormValue("model2")
		modelConf[1].APIKey = r.FormValue("api_key2")
		modelConf[1].BaseURL = r.FormValue("base_url2")
		modelConf[1].Name = r.FormValue("name2")
		modelConf[1].Provider = r.FormValue("provider2")

		modelConf[0].HasSaved = true
		modelConf[1].HasSaved = true
		fmt.Println(modelConf)
		components.ConfigDisabled(modelConf).Render(r.Context(), w)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//components.Home().Render(r.Context(), w)
		components.Login().Render(r.Context(), w)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		fmt.Println(username, password)
		components.Home().Render(r.Context(), w)
	})

	mux.HandleFunc("/newChat", func(w http.ResponseWriter, r *http.Request) {
		components.ChatHome("Pradhyumna", titles[rand.Intn(len(titles))], true).Render(r.Context(), w)

	})

	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {

		if !modelConf[0].HasSaved {
			components.ConfigMain(modelConf).Render(r.Context(), w)
		} else {
			components.ConfigDisabled(modelConf).Render(r.Context(), w)
		}
	})

	mux.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
		if !serverConf[0].HasSaved {
			components.Tools(serverConf).Render(r.Context(), w)
		} else {
			components.ToolsDisabled(serverConf).Render(r.Context(), w)
		}
	})

	mux.HandleFunc("POST /settoolconf", func(w http.ResponseWriter, r *http.Request) {
		serverConf[0].Name = r.FormValue("name1")
		serverConf[0].BaseURL = r.FormValue("base_url1")
		serverConf[0].Transport = r.FormValue("transport1")

		serverConf[1].Name = r.FormValue("name2")
		serverConf[1].BaseURL = r.FormValue("base_url2")
		serverConf[1].Transport = r.FormValue("transport2")

		serverConf[0].HasSaved = true
		serverConf[1].HasSaved = true
		fmt.Println(serverConf)
		components.ToolsDisabled(serverConf).Render(r.Context(), w)
	})

	mux.HandleFunc("GET /gettoolconf", func(w http.ResponseWriter, r *http.Request) {
		components.Tools(serverConf).Render(r.Context(), w)
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
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Println(err)
				}
				break
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
