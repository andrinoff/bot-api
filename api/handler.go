package handler

import (
	"encoding/json"
	"net/http"
	"os"

	discord_bot "bot-api/discord"
	telegram_bot "bot-api/telegram"
)

// RequestPayload defines the structure of the JSON body we expect from the frontend.
type RequestPayload struct {
	Password string `json:"password"`
	Content  string `json:"content"`
}

// ResponsePayload defines the structure for our JSON responses.
type ResponsePayload struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// processContent is a placeholder function that gets called if the password is correct.
// In a real application, this is where you would put your logic to handle the text,


// Handler is the main entry point for the Vercel serverless function.
// It handles the HTTP request, checks the password, and processes the content.
func Handler(w http.ResponseWriter, r *http.Request) {
	// --- 1. Set Headers ---
	// Set the content type to application/json for all responses.
	w.Header().Set("Content-Type", "application/json")

	// --- 2. Handle HTTP Method ---
	// We only want to allow POST requests to this endpoint.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Method Not Allowed. Please use POST.", Success: false})
		return
	}

	// --- 3. Decode Incoming JSON ---
	var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		// If there's an error decoding the JSON, send a bad request response.
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Invalid request body.", Success: false})
		return
	}

	// --- 4. Get Password from Environment Variables ---

	correctPassword := os.Getenv("BACKEND_PASSWORD")

	// --- 5. Compare Passwords ---
	if payload.Password != correctPassword {
		// If the password does not match, return an unauthorized error.
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Invalid password.", Success: false})
		return
	}

	// --- 6. Process the Content ---
	// If the password is correct, call the function to process the content.
	telegram_bot.Send(payload.Content)
	discord_bot.Send(payload.Content)

	// --- 7. Send Success Response ---
	// Let the frontend know that everything was successful.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponsePayload{Message: "Content received and processed successfully.", Success: true})
}

