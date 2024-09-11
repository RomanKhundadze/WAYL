package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
	"local.wayl/config"
)

var (
	oauth2Config *oauth2.Config
	token        *oauth2.Token
	path         *string
	authPath     *string
)

type PlaybackState struct {
	Item struct {
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
	} `json:"item"`
	IsPlaying bool `json:"is_playing"`
}

func Root(pathP *string, authPathP *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path = pathP
		authPath = authPathP
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func Login(conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauth2Config = &oauth2.Config{
			ClientID:     conf.ClientID,
			ClientSecret: conf.ClientSecret,
			RedirectURL:  "http://localhost" + conf.Port + "/callback",
			Scopes:       []string{"user-read-currently-playing"},
			Endpoint:     spotify.Endpoint,
		}

		refPath := *authPath + "/refreshtoken"

		if _, err := os.Stat(refPath); errors.Is(err, os.ErrNotExist) {
			url := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			fileRefreshToken, err := os.ReadFile(refPath)
			token = &oauth2.Token{
				RefreshToken: string(fileRefreshToken),
			}
			token, err = refreshToken()
			if err != nil {
				log.Println("Error refreshing token:", err)
				return
			}

			http.Redirect(w, r, "/playback", http.StatusSeeOther)
		}
	}
}

func Callback(w http.ResponseWriter, r *http.Request) {
	var err error

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}

	token, err = oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	refPath := *authPath + "/refreshtoken"
	file, err := os.Create(refPath)
	if err != nil {
		log.Println("Error creating refreshtoken file:", err)
		return
	}
	file.WriteString(token.RefreshToken)

	http.Redirect(w, r, "/playback", http.StatusSeeOther)
}

func Playback(w http.ResponseWriter, r *http.Request) {
	if path == nil {
		log.Println("No Path defined going back to root")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	indexPath := *path + "/Website/playback.html"

	tmpl := template.Must(template.ParseFiles(indexPath))
	tmpl.Execute(w, nil)
}

func getCurrentPlaybackState() PlaybackState {
	var state PlaybackState

	if token == nil || token.Expiry.Before(time.Now()) {
		log.Println("Token is missing or expired, refreshing...")
		var err error
		token, err = refreshToken()
		if err != nil {
			log.Println("Error refreshing token:", err)
			return state
		}
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return state
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching playback state:", err)
		return state
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return state
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d\nResponse: %s", resp.StatusCode, body)
		return state
	}

	err = json.Unmarshal(body, &state)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
	}

	return state
}

func refreshToken() (*oauth2.Token, error) {
	if token == nil || token.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
	}

	tokenSource := oauth2Config.TokenSource(context.Background(), token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

func HandleGetPlaybackData(w http.ResponseWriter, r *http.Request) {
	playback := getCurrentPlaybackState()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(playback); err != nil {
		http.Error(w, "Failed to encode playback data", http.StatusInternalServerError)
	}
}
