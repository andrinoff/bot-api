package twitter_bot

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strconv"

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

	// Initialize the anaconda Twitter API client
	anaconda.SetConsumerKey(apiKey)
	anaconda.SetConsumerSecret(apiSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	// --- 3. Upload Image (if provided) ---
	v := url.Values{}
	if imageBytes != nil {
		fmt.Println("Image data provided, uploading to Twitter...")
		// Upload image to Twitter's media endpoint
		encoded := base64.StdEncoding.EncodeToString(imageBytes)
		media, err := api.UploadMedia(encoded)
		if err != nil {
			return "", fmt.Errorf("twitter media upload failed: %w", err)
		}
		
		fmt.Println("Media uploaded successfully. Media ID:", media.MediaIDString)

		// Add the media ID to the tweet parameters
		v.Add("media_ids", strconv.FormatInt(media.MediaID, 10))
	}

	// --- 4. Post the Tweet ---
	fmt.Println("Posting tweet...")
	tweet, err := api.PostTweet(message, v)
	if err != nil {
		return "", fmt.Errorf("failed to post tweet: %w", err)
	}

	// --- 5. Return the Tweet ID on Success ---
	return tweet.IdStr, nil
}
