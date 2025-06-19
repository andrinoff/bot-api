package handler

import (
	discord_bot "bot-api/discord"
	telegram_bot "bot-api/telegram"
	twitter_bot "bot-api/twitter"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// RequestPayload defines the structure for the incoming JSON.
type RequestPayload struct {
	Password string `json:"password"`
	Content  string `json:"content"`
	Discord string `json:"discord"`
	Telegram string `json:"telegram"`
	Twitter string `json:"twitter"`
}

// ResponsePayload is defined in handler.go and reused here.

// processContent is a placeholder for your logic.
func processContent(content string, discord string, telegram string, twitter string) {
	if twitter == "true" {
		tweet_id, err := twitter_bot.PostTweet(content, nil)
		if err != nil {
			fmt.Println("Error posting tweet:", err)
			return
		}
	
		fmt.Println(tweet_id)
	}
	if telegram == "true" {
		telegram_bot.SendNoImage(content)
	}
	if discord == "true" {
		discord_bot.SendNoImage(content)
	}
}

// Handler for the /api/no-image endpoint.
func NoImageHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS and Content-Type headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Method Not Allowed.", Success: false})
		return
	}

	// Decode the incoming JSON body
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Invalid request body.", Success: false})
		return
	}

	// Compare passwords
	correctPassword := os.Getenv("BACKEND_PASSWORD")
	if payload.Password != correctPassword {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Invalid password.", Success: false})
		return
	}

	// Process the content
	processContent(payload.Content, payload.Discord, payload.Telegram, payload.Twitter)

	// Send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponsePayload{Message: "Content received and processed successfully.", Success: true})
}
