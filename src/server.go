package googlePhotoServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var client *http.Client

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

type Photo struct {
	BaseUrl  string `json:"baseUrl"`
	Filename string `json:"filename"`
}
type PhotoResponse struct {
	Photos []Photo `json:"photos"`
	Size   int     `json:"size"`
}

// Handler function for the GET endpoint
func photosHandler(w http.ResponseWriter, r *http.Request) {
	gpResp := queryPhotos(client)
	photos := convertResponseToPhotoResponse(gpResp)
	// Marshal PhotoResponse to JSON
	jsonResponse, err := json.Marshal(photos)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the JSON response
	w.Write(jsonResponse)
}

// ConvertResponseToPhotoResponse converts a Response struct into a PhotoResponse struct
func convertResponseToPhotoResponse(response Response) PhotoResponse {
	photos := make([]Photo, len(response.MediaItems))

	for i, item := range response.MediaItems {
		photos[i] = Photo{
			BaseUrl:  item.BaseUrl,
			Filename: item.Filename,
		}
	}

	return PhotoResponse{
		Photos: photos,
		Size:   len(photos),
	}
}

func startListeners() {
	http.HandleFunc("/photos", photosHandler)
	log.Fatal(http.ListenAndServe("localhost:443", nil))
}

func StartServer() {
	client = createClient()
	startListeners()
	// fmt.Print(queryPhotos(client)) //remove
	// create a listener
	// make the listener call photoshandler i guess
}

func createClient() *http.Client {
	// config, err := GetConfig()
	// if err != nil {
	// 	fmt.Println("Unable to get config.")
	// 	return nil
	// }
	client := GetAuthClient()

	return client
}

func queryPhotos(client *http.Client) Response {
	// data := map[string]interface{}{
	// 	"filters": map[string]interface{}{
	// 		"contentFilter": map[string]interface{}{
	// 			"includedContentCategories": []string{"PETS"},
	// 		},
	// 	},
	// }

	//getAlbumInfo(client)

	data := map[string]interface{}{
		"albumId":  "ADmVkWt1jie_gCSizAXqKBu0bwgW1iWxkwN5d6Xto0YuxC5yfoVIcKYBOQW2otfbcFnTVcHJBOWu",
		"pageSize": 100,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return Response{}
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

	return result
}

func getAlbumInfo(client *http.Client) {
	type Album struct {
		ID                    string `json:"id"`
		Title                 string `json:"title"`
		ProductURL            string `json:"productUrl"`
		IsWriteable           bool   `json:"isWriteable"`
		MediaItemsCount       string `json:"mediaItemsCount"`
		CoverPhotoBaseURL     string `json:"coverPhotoBaseUrl"`
		CoverPhotoMediaItemID string `json:"coverPhotoMediaItemId"`
	}

	type Albums struct {
		Albums        []Album `json:"albums"`
		NextPageToken string  `json:"nextPageToken"`
	}

	resp, err := client.Get("https://photoslibrary.googleapis.com/v1/albums")
	if err != nil {
		log.Fatalf("Unable to get album info: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read response body: %v", err)
	}

	var result Albums
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Unable to decode response: %v", err)
	}

	for _, album := range result.Albums {
		if album.Title == "Gwynniebook" {
			fmt.Println(album.ID)
		}
	}
}
