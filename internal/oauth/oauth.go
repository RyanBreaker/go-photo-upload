package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

var ClientId = os.Getenv("DBX_CLIENT_ID")
var clientSecret = os.Getenv("DBX_CLIENT_SECRET")

var RedirectUri string
var AuthorizeUrl = fmt.Sprintf(
	"https://www.dropbox.com/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&token_access_type=offline",
	ClientId,
	RedirectUri,
)

const tokenUrl = "https://api.dropboxapi.com/oauth2/token"

func init() {
	if os.Getenv("ENV") == "production" {
		RedirectUri = "https://wedding-photos.breaker.rocks/oauth2/redirect" // TODO: confirm correct?
	} else {
		RedirectUri = "http://localhost:8080/oauth2/redirect"
	}
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
	data.Add("redirect_uri", RedirectUri)

	res, err := http.PostForm(tokenUrl, data)
	if err != nil {
		log.Println("Error getting access and refresh tokens:", err)
		return
	}
	defer res.Body.Close()

	tokens, _ := unmarshalTokens(res.Body)

	setToken(accessToken, tokens)
	setToken(refreshToken, tokens)
}

func RefreshAccessToken() {
	log.Println("Refreshing access token")

	refreshToken := GetRefreshToken()
	if refreshToken == "" {
		log.Println("No refresh token")
		return // TODO: error?
	}

	data := url.Values{}
	data.Add("grant_type", "refresh_token")
	data.Add("refresh_token", refreshToken)
	data.Add("client_id", ClientId)
	data.Add("client_secret", clientSecret)

	res, err := http.PostForm(tokenUrl, data)
	if err != nil {
		log.Println("Error while refreshing access token:", err)
		return
	}
	defer res.Body.Close()

	tokens, _ := unmarshalTokens(res.Body)

	if tokens.AccessToken == "" {
		log.Println("Empty access token received")
		return // TODO: error?
	}

	setToken(accessToken, tokens)
	log.Println("Set new access token")
}
