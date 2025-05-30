package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func ListChatScopes() []string {
	return []string{
		"https://www.googleapis.com/auth/chat.admin.memberships",
		"https://www.googleapis.com/auth/chat.admin.memberships.readonly",
		"https://www.googleapis.com/auth/chat.admin.spaces",
		"https://www.googleapis.com/auth/chat.admin.spaces.readonly",
		"https://www.googleapis.com/auth/chat.memberships",
		"https://www.googleapis.com/auth/chat.memberships.app",
		"https://www.googleapis.com/auth/chat.memberships.readonly",
		"https://www.googleapis.com/auth/chat.messages",
		"https://www.googleapis.com/auth/chat.messages.create",
		"https://www.googleapis.com/auth/chat.messages.reactions",
		"https://www.googleapis.com/auth/chat.messages.reactions.create",
		"https://www.googleapis.com/auth/chat.messages.reactions.readonly",
		"https://www.googleapis.com/auth/chat.messages.readonly",
		"https://www.googleapis.com/auth/chat.spaces",
		"https://www.googleapis.com/auth/chat.spaces.create",
		"https://www.googleapis.com/auth/chat.spaces.readonly",
		"https://www.googleapis.com/auth/chat.users.readstate",
		"https://www.googleapis.com/auth/chat.users.readstate.readonly",
	}
}
func ListGoogleScopes() []string {
	scopes := []string{
		gmail.GmailLabelsScope,
		gmail.GmailModifyScope,
		gmail.MailGoogleComScope,
		gmail.GmailSettingsBasicScope,
	}
	scopes = append(scopes, ListChatScopes()...)
	return scopes
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	pwd, _ := os.Getwd()
	tokFile := pwd + "/google-credential/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()
	pwd, _ := os.Getwd()
	b, err := os.ReadFile(pwd + "/google-credential/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, ListGoogleScopes()...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
}
