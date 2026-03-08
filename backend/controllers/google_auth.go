package controllers

import (
	"net/http"

	"learningai/config"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(c *gin.Context) {

	url := config.GoogleConfig.AuthCodeURL("random-state")

	c.Redirect(http.StatusTemporaryRedirect, url)

}

func GoogleCallback(c *gin.Context) {

	code := c.Query("code")

	token, err := config.GoogleConfig.Exchange(c, code)

	if err != nil {
		c.JSON(500, gin.H{"error": "Token exchange failed"})
		return
	}

	client := config.GoogleConfig.Client(c, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user info"})
		return
	}

	defer resp.Body.Close()

	c.JSON(200, gin.H{
		"message": "Login Google berhasil",
	})
}
