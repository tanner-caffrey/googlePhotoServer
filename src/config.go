package googlePhotoServer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	SecretsPath string `json:"SecretsPath"`
}

type Secrets struct {
	ClientId                string `json:"client_id"`
	ProjectId               string `json:"project_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientSecret            string `json:"client_secret"`
}

type SecretsWrapper struct {
	Secrets Secrets `json:"web"`
}

func GetConfig() (*Config, error) {
	// Open json file
	byteValue, err := readJsonFile("config/config.json")
	if err != nil {
		fmt.Println("Error reading config", err)
		return nil, err
	}

	// Unmarshal the JSON data into a map
	var result Config
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil, err
	}
	return &result, nil
}

func GetSecrets(config *Config) (*Secrets, error) {
	if config == nil {
		return nil, fmt.Errorf("error reading secrets: Invalid config file (nil)")
	}
	byteValue, err := readJsonFile(config.SecretsPath)
	if err != nil {
		fmt.Println("Error reading config", err)
		return nil, err
	}
	// Unmarshal the JSON data into a map
	var result SecretsWrapper
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil, err
	}

	return &result.Secrets, nil
}

func readJsonFile(path string) ([]byte, error) {
	// Open json file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file", path, "-", err)
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

// Secrets
func GetClient(secretsPath string) (*http.Client, error) {
	ctx := context.Background()

	// Load OAuth 2.0 credentials from a file
	b, err := os.ReadFile(secretsPath)
	if err != nil {
		fmt.Printf("Unable to read client secret file: %v\n", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/photoslibrary.readonly")
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v\n", err)
		return nil, err
	}

	client := getClient(ctx, config)

	return client, nil
}

func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(ctx, tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		fmt.Printf("Unable to read authorization code: %v\n", err)
		return nil
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		fmt.Printf("Unable to retrieve token from web: %v\n", err)
		return nil
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Unable to create file: %v\n", err)
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
