package googlePhotoServer

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetAuthClient() *http.Client {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	refreshToken := os.Getenv("GOOGLE_REFRESH_TOKEN")

	var oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"), // Set this environment variable
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
	}

	// TokenSource to refresh tokens
	tokenSource := oauthConfig.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: refreshToken,
	})

	// Get a new token
	token, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}

	client := oauthConfig.Client(context.Background(), token)

	return client
}
