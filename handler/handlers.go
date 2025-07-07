package handler

import (
	"github.com/Nilesh015/go-url-shortener/shortener"
	"github.com/Nilesh015/go-url-shortener/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UrlCreationRequest defines the structure for the URL creation request
type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

// CreateShortUrl handles the creation of a short URL
// It expects a JSON body with "long_url" and "user_id" fields
// It generates a short URL using the shortener package and saves the mapping in the store
// It responds with the full short URL
// Example request body: {"long_url": "https://example.com", "user_id": "12345"}
// Example response: {"message": "short url created successfully", "short_url": "http://localhost:9808/jTa4L57P"}
func CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl := shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)
	store.SaveUrlMapping(shortUrl, creationRequest.LongUrl, creationRequest.UserId)

	host := "http://localhost:9808/"
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + shortUrl,
	})

}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	initialUrl := store.RetrieveInitialUrl(shortUrl)
	c.Redirect(302, initialUrl)
}
