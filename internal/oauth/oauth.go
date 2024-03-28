package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

var ClientId = os.Getenv("DBX_CLIENT_ID")
var clientSecret = os.Getenv("DBX_CLIENT_SECRET")

var AuthorizeUri string
var redirectUri string

const tokenUrl = "https://api.dropboxapi.com/oauth2/token"

func init() {
	if os.Getenv("ENV") == "production" {
		redirectUri = "https://wedding-photos.breaker.rocks/oauth2/redirect"
	} else {
		redirectUri = "http://localhost:8080/oauth2/redirect"
	}
	AuthorizeUri = fmt.Sprintf(
		"https://www.dropbox.com/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&token_access_type=offline",
		ClientId,
		redirectUri,
	)
}

type oauthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func unmarshalTokens(r io.Reader) (*oauthTokens, error) {
	body, _ := io.ReadAll(r)

	var tokens *oauthTokens
	err := json.Unmarshal(body, &tokens)
	return tokens, err
}

func SetTokens(code string) {
	data := url.Values{}
	data.Add("code", code)
	data.Add("grant_type", "authorization_code")
	data.Add("client_id", ClientId)
	data.Add("client_secret", clientSecret)
	data.Add("redirect_uri", redirectUri)

	res, err := http.PostForm(tokenUrl, data)
	if err != nil {
		slog.Error("Error getting access and refresh tokens", slog.Any("error", err))
		return
	}
	defer res.Body.Close()

	tokens, _ := unmarshalTokens(res.Body)

	setToken(accessToken, tokens)
	setToken(refreshToken, tokens)
}

func RefreshAccessToken() {
	slog.Info("Refreshing access token")

	// TODO: This needs to be rate-limited
	refreshToken := GetRefreshToken()
	if refreshToken == "" {
		slog.Warn("No refresh token")
		return // TODO: error?
	}

	data := url.Values{}
	data.Add("grant_type", "refresh_token")
	data.Add("refresh_token", refreshToken)
	data.Add("client_id", ClientId)
	data.Add("client_secret", clientSecret)

	res, err := http.PostForm(tokenUrl, data)
	if err != nil {
		slog.Error("Error while refreshing access token", slog.Any("error", err))
		return
	}
	defer res.Body.Close()

	tokens, _ := unmarshalTokens(res.Body)

	if tokens.AccessToken == "" {
		slog.Error("Empty access token received")
		return // TODO: error?
	}

	setToken(accessToken, tokens)
	slog.Info("Set new access token")
}
