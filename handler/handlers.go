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
}

const kHostUrl = "http://localhost:9808/"

// Given a long URL, create a corresponding short URL and save it in database.
func CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch a short URL till a unique one is fetched.
	shortUrl := shortener.GenerateShortLink(creationRequest.LongUrl)
	for store.RetrieveLongUrl(shortUrl) != "" {
		shortUrl = shortener.GenerateShortLink(creationRequest.LongUrl)
	}

	store.SaveUrlMapping(shortUrl, creationRequest.LongUrl)

	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": kHostUrl + shortUrl,
	})

}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	longUrl := store.RetrieveLongUrl(shortUrl)
	c.Redirect(302, longUrl)
}
