package googlePhotoServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Define a struct for the JSON response
type Response struct {
	Message string `json:"message"`
}

// Handler function for the GET endpoint
func photosHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Hello, World!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func StartServer() {
	client := createClient()
	queryPhotos(client) //remove
	// create a listener
	// make the listener call photoshandler i guess
}

func createClient() *http.Client {
	config, err := GetConfig()
	if err != nil {
		fmt.Println("Unable to get config.")
		return nil
	}
	client, err := GetClient(config.SecretsPath)
	if err != nil {
		fmt.Println("Unable to generate client.")
		return nil
	}
	return client
}

func queryPhotos(client *http.Client) {
	data := map[string]interface{}{
		"filters": map[string]interface{}{
			"contentFilter": map[string]interface{}{
				"includedContentCategories": []string{"PETS"},
			},
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := client.Post(
		"https://photoslibrary.googleapis.com/v1/mediaItems:search",
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Unable to retrieve albums: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Unable to decode response: %v", err)
	}

	fmt.Println(result)
}
