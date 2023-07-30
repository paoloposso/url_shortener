package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	db_url "github.com/paoloposso/url_shrt/db/url"
	"github.com/paoloposso/url_shrt/url"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

func init() {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func main() {
	log.Println(os.Getenv("MONGO_URI"))

	repo, err := db_url.NewRepository(os.Getenv("MONGO_URI"), "url_shrt")

	if err != nil {
		log.Fatal(err)
	}

	service := url.NewService(repo)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.POST("/shorten", func(c *gin.Context) {
		var req ShortenRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		shortUrl, err := service.ShortenURL(req.URL)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"shortUrl": shortUrl,
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Service Started")
	}
}
