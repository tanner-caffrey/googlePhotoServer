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
type ServerResponse struct {
	Message string `json:"message"`
}

type MediaMetadata struct {
	// CreationTime string `json:"creationTime"`
	Height string `json:"height"`
	// Photo        Photo  `json:"photo"`
	Width string `json:"width"`
}

// type Photo struct {
// 	ApertureFNumber float64 `json:"apertureFNumber"`
// 	CameraMake      string  `json:"cameraMake"`
// 	CameraModel     string  `json:"cameraModel"`
// 	ExposureTime    string  `json:"exposureTime"`
// 	FocalLength     float64 `json:"focalLength"`
// 	IsoEquivalent   int     `json:"isoEquivalent"`
// }

type MediaItem struct {
	BaseUrl       string        `json:"baseUrl"`
	Filename      string        `json:"filename"`
	ID            string        `json:"id"`
	MediaMetadata MediaMetadata `json:"mediaMetadata"`
	// MimeType      string        `json:"mimeType"`
	// ProductUrl    string        `json:"productUrl"`
}

type Response struct {
	MediaItems    []MediaItem `json:"mediaItems"`
	NextPageToken string      `json:"nextPageToken"`
}

// Handler function for the GET endpoint
func photosHandler(w http.ResponseWriter, r *http.Request) {
	response := ServerResponse{Message: "Hello, World!"}
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

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Unable to decode response: %v", err)
	}

	fmt.Println(result)
}
