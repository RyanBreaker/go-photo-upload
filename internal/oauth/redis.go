package oauth

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

const (
	accessToken  = "access_token"
	refreshToken = "refresh_token"
)

var client *redis.Client

func init() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalln("Error parsing redis URL:", err)
	}
	client = redis.NewClient(opt)
}

func GetAccessToken() string {
	var token string
	for {
		var err error
		token, err = getToken(accessToken)
		if errors.Is(err, redis.Nil) {
			log.Println("Refreshing access token from GetAccessToken...")
			RefreshAccessToken()
			continue
		}
		break
	}
	return token
}

func GetRefreshToken() string {
	t, _ := getToken(refreshToken)
	return t
}

func getToken(tokenType string) (string, error) {
	ctx := context.Background()
	val, err := client.Get(ctx, tokenType).Result()
	if err != nil {
		log.Println("Error getting access token:", err)
		return "", err
	}
	return val, nil
}

func setToken(tokenType string, token *oauthTokens) {
	ctx := context.Background()
	var err error
	if tokenType == accessToken {
		//expiration := time.Duration(token.ExpiresIn) * time.Second
		expiration := time.Duration(15) * time.Second
		err = client.Set(ctx, tokenType, token.AccessToken, expiration).Err()
	} else {
		err = client.Set(ctx, tokenType, token.RefreshToken, 0).Err()
	}
	if err != nil {
		log.Println("Error setting token:", err)
	}
}