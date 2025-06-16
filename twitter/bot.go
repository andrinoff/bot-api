package twitter_bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func PostTweet(message string, imageBytes []byte) (string, error) {
	// --- 1. Validate Input ---
	if message == "" {
		return "", fmt.Errorf("tweet message cannot be empty")
	}

	// --- 2. Configure Twitter Client from Environment Variables ---
	apiKey := os.Getenv("TWITTER_API_KEY")
	apiSecret := os.Getenv("TWITTER_API_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	if apiKey == "" || apiSecret == "" || accessToken == "" || accessTokenSecret == "" {
		return "", fmt.Errorf("twitter API credentials are not set in environment variables")
	}

	anaconda.SetConsumerKey(apiKey)
	anaconda.SetConsumerSecret(apiSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	// --- 3. Upload Media using v1.1 (if image exists) ---
	var mediaID string
	if imageBytes != nil {
		fmt.Println("Image data provided, uploading via v1.1 Media API...")
		media, err := api.UploadMedia(string(imageBytes))
		if err != nil {
			return "", fmt.Errorf("twitter v1.1 media upload failed: %w", err)
		}
		mediaID = media.MediaIDString
		fmt.Println("Media uploaded successfully. Media ID:", mediaID)
	}

	// --- 4. Post the Tweet using v2 Tweets API ---
	// Create the JSON payload required by the v2 endpoint.
	tweetPayload := map[string]interface{}{"text": message}
	if mediaID != "" {
		tweetPayload["media"] = map[string][]string{
			"media_ids": {mediaID},
		}
	}
	payloadBytes, err := json.Marshal(tweetPayload)
	if err != nil {
		return "", fmt.Errorf("failed to create v2 tweet payload: %w", err)
	}

	// Manually create and send an authorized request to the v2 tweets endpoint.
	v2Endpoint := "https://api.twitter.com/2/tweets"
	req, err := http.NewRequest("POST", v2Endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create v2 request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// The anaconda library conveniently provides an authenticated http.Client
	// that handles the OAuth1 signing for us.
	fmt.Println("Posting tweet via v2 Tweets API...")
	resp, err := api.HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send v2 tweet request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response from the server.
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("v2 tweet post failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// --- 5. Parse the successful response to get the Tweet ID ---
	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse v2 tweet response: %w", err)
	}

	return result.Data.ID, nil
}