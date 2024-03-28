package oauth

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
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
		slog.Error("Error getting access token", slog.Any("error", err))
		return "", err
	}
	return val, nil
}

func setToken(tokenType string, token *oauthTokens) {
	ctx := context.Background()
	var err error
	if tokenType == accessToken {
		expiration := time.Duration(token.ExpiresIn) * time.Second
		err = client.Set(ctx, tokenType, token.AccessToken, expiration).Err()
	} else {
		err = client.Set(ctx, tokenType, token.RefreshToken, 0).Err()
	}
	if err != nil {
		slog.Error("Error setting token: %s", err.Error())
	}
}
