package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Constants for the Twitter API base URL
const baseURL = "https://api.twitter.com/2"

// Struct for tweet data
type Tweet struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// Response structure for recent tweets API
type TweetsResponse struct {
	Data []Tweet `json:"data"`
}

// Function to create HTTP headers using Bearer Token
func createHeaders(bearerToken string) http.Header {
	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	headers.Set("Content-Type", "application/json")
	return headers
}

// Function to retrieve the last 10 retweets
func getRecentRetweets(bearerToken string) ([]Tweet, error) {
	// URL for retrieving recent retweets (modify if necessary)
	url := fmt.Sprintf("%s/users/me/tweets?tweet.fields=id&expansions=referenced_tweets.id&max_results=10", baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers with the provided Bearer Token
	req.Header = createHeaders(bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var tweetsResponse TweetsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tweetsResponse); err != nil {
		return nil, err
	}

	return tweetsResponse.Data, nil
}

// Function to unretweet a tweet by its ID
func unretweet(tweetID string, bearerToken string) error {
	url := fmt.Sprintf("%s/tweets/%s/unretweet", baseURL, tweetID)

	// Create DELETE request to unretweet
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	// Set headers with the Bearer Token
	req.Header = createHeaders(bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to unretweet: %v", resp.Status)
	}

	fmt.Printf("Successfully unretweeted tweet with ID: %s\n", tweetID)
	return nil
}

func main() {
	// Retrieve API credentials from the local environment
	apiKey := os.Getenv("TWITTER_API_KEY")
	apiSecret := os.Getenv("TWITTER_API_SECRET")
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")

	if apiKey == "" || apiSecret == "" || bearerToken == "" {
		log.Fatal("Error: Missing required environment variables. Please set TWITTER_API_KEY, TWITTER_API_SECRET, and TWITTER_BEARER_TOKEN.")
	}

	// Step 1: Get the last 10 retweets
	retweets, err := getRecentRetweets(bearerToken)
	if err != nil {
		log.Fatalf("Error retrieving retweets: %v", err)
	}

	// Step 2: Undo each retweet
	for _, tweet := range retweets {
		err := unretweet(tweet.ID, bearerToken)
		if err != nil {
			log.Printf("Error unretweeting tweet with ID %s: %v", tweet.ID, err)
		}
	}
}
