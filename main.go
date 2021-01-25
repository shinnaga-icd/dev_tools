package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	serve := gin.Default()
	serve.Use(gin.Logger())
	serve.LoadHTMLGlob("templates/*.tmpl.html")

	serve.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"message": "GET Reqest",
		})
	})

	serve.POST("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"message": "POST Reqest",
			"bhreq":   c.PostForm("bhreq"),
			"p4":      c.PostForm("p4"),
		})
	})

	serve.Run(":" + port)
}
