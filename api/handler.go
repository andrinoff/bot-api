package handler

import (
	discord_bot "bot-api/discord"
	telegram_bot "bot-api/telegram"
	twitter_bot "bot-api/twitter"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ResponsePayload defines the structure for our JSON responses.
type ResponsePayload struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}


func processImage(fileBytes []byte, filename string, message string) {
	fmt.Printf("Processing image: %s, Size: %d bytes\n", filename, len(fileBytes))
	os.WriteFile("/tmp/"+filename, fileBytes, 0644)
	twitter_bot.PostTweet(message, fileBytes)
	telegram_bot.Send(message, "/tmp/"+filename)
	discord_bot.Send(message, "/tmp/"+filename)
	
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// --- 1. Set CORS Headers ---
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// --- 2. Handle Preflight OPTIONS Request ---
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// --- 3. Handle HTTP Method ---
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Method Not Allowed. Please use POST.", Success: false})
		return
	}

	// --- 4. Parse Multipart Form ---
	// We set a max memory limit for the form parts. 10MB in this case.
	// Larger files will be temporarily stored on disk.
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Could not parse multipart form.", Success: false})
		return
	}

	// --- 5. Get Password from Form Value ---
	password := r.FormValue("password")
	correctPassword := os.Getenv("BACKEND_PASSWORD")

	// --- 6. Compare Passwords ---
	if password != correctPassword {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Invalid password.", Success: false})
		return
	}

	// --- 7. Get the Image File from the Form ---
	file, handler, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Invalid image file in form.", Success: false})
		return
	}
	defer file.Close()

	// --- 8. Read the file into a byte slice ---
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Could not read uploaded file.", Success: false})
		return
	}

	// --- 9. Process the Image ---
	processImage(fileBytes, handler.Filename, r.FormValue("message"))

	// --- 10. Send Success Response ---
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponsePayload{Message: "Image received and processed successfully.", Success: true})
}
