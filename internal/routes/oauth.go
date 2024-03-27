package routes

import (
	"github.com/RyanBreaker/go-photo-upload/internal/oauth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func OauthRoutes(router *gin.Engine) {
	oauthGroup := router.Group("/oauth2")
	{
		oauthGroup.GET("/authorize", func(c *gin.Context) {
			c.Redirect(http.StatusFound, oauth.AuthorizeUrl)
		})

		oauthGroup.GET("/redirect", func(c *gin.Context) {
			code := c.Query("code")
			if code == "" {
				log.Println("code is empty")
			}
			oauth.SetTokens(code)
			c.Redirect(http.StatusFound, "/")
		})
	}
}
