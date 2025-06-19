package handler

import (
	discord_bot "bot-api/discord"
	telegram_bot "bot-api/telegram"
	twitter_bot "bot-api/twitter"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// ResponsePayload defines the structure for our JSON responses.
type ResponsePayload struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func processImage(ctx context.Context, fileBytes []byte, filename string, message string, discord string, telegram string, twitter string) error {
	fmt.Printf("Processing image: %s, Size: %d bytes\n", filename, len(fileBytes))

	if twitter == "true" {
		tweet_id, err := twitter_bot.PostTweet(ctx, message, fileBytes)
		if err != nil {
			return fmt.Errorf("error posting tweet: %w", err)
		}
		fmt.Println(tweet_id)
	}

	if telegram == "true" {
		if err := telegram_bot.Send(ctx, message, bytes.NewReader(fileBytes)); err != nil {
			return fmt.Errorf("error sending Telegram message: %w", err)
		}
	}

	if discord == "true" {
		if err := discord_bot.Send(ctx, message, bytes.NewReader(fileBytes)); err != nil {
			return fmt.Errorf("error sending Discord message: %w", err)
		}
	}

	return nil
}

func errorMiddleware(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := next(w, r)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}

		if err != nil {
			log.Printf("%s %s %s ERROR %s", r.Method, r.URL.Path, r.RemoteAddr, err)
		} else {
			log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		}

	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	errorMiddleware(handler)(w, r)
}

func handler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	// --- 1. Set CORS Headers ---
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// --- 2. Handle Preflight OPTIONS Request ---
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	// --- 3. Handle HTTP Method ---
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ResponsePayload{Message: "Method Not Allowed. Please use POST.", Success: false})
		return nil
	}

	// --- 4. Parse Multipart Form ---
	// We set a max memory limit for the form parts. 10MB in this case.
	// Larger files will be temporarily stored on disk.
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return err
	}

	// --- 5. Get Password from Form Value ---
	password := r.FormValue("password")
	correctPassword := os.Getenv("BACKEND_PASSWORD")

	// --- 6. Compare Passwords ---
	if password != correctPassword {
		return err
	}

	// --- 8 . Get what social media to post to

	// --- 7. Get the Image File from the Form ---
	file, handler, err := r.FormFile("image")
	if err != nil {
		return err
	}
	defer file.Close()

	// --- 8. Read the file into a byte slice ---
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	fmt.Print(r.FormValue("discord"), r.FormValue("telegram"), r.FormValue("twitter"))

	// --- 9. Process the Image ---
	if err := processImage(ctx, fileBytes, handler.Filename, r.FormValue("message"), r.FormValue("discord"), r.FormValue("telegram"), r.FormValue("twitter")); err != nil {
		return err
	}

	// --- 10. Send Success Response ---
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponsePayload{Message: "Image received and processed successfully.", Success: true})
	return nil
}
