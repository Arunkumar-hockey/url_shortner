package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"time"
)

var (
	shortURLPrefix = "http://localhost:8080/"
	longURLCollection *mongo.Collection
)

func initMongoDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	longURLCollection = client.Database("urlshortener").Collection("urls")

	return client
}

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

	_, err := longURLCollection.InsertOne(context.TODO(), bson.M{
		"short_url": shortURL,
		"long_url": longURL,
	})

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}

func redirect(c *gin.Context) {
	fmt.Println("true..")
	shortURL := c.Query("url")

	var result bson.M
	//x := "http://localhost:8080/3BNiMR"
	err := longURLCollection.FindOne(context.TODO(), bson.M{"short_url": shortURL}).Decode(&result)
	fmt.Println("true...")
	if err != nil {
		fmt.Println("Error..", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}
	fmt.Println("==========true=========")
	longURL := result["long_url"].(string)
	fmt.Println(longURL)
	c.Redirect(http.StatusTemporaryRedirect, longURL)
}

func main() {
	initMongoDB()
	router := gin.Default()

	router.POST("/shorten", shortenURL)
	router.GET("/short", redirect)

	router.Run(":8080")
}
