package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var (
	shortURLPrefix = "http://localhost:8080/"
	urlCache       = cache.New(5*time.Minute, 10*time.Minute)
)

func generateShortURL() string {
	return shortURLPrefix + randomString(6)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func shortenURL(c *gin.Context) {
	longURL := c.PostForm("url")
	if longURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No URL provided"})
		return
	}

	shortURL := generateShortURL()
	urlCache.Set(shortURL, longURL, cache.DefaultExpiration)

	c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}

func redirect(c *gin.Context) {
	shortURL := c.Query("url")
	longURL, found := urlCache.Get(shortURL)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	fmt.Println(longURL)
	c.Redirect(http.StatusTemporaryRedirect, longURL.(string))
}

func main() {
	router := gin.Default()

	router.POST("/shorten", shortenURL)
	router.GET("/:short", redirect)

	router.Run(":8080")
}
