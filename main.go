package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "5002"
	}

	hub := newHub()

	go hub.run()

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")

	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "game.tmpl.html", nil)
	})

	router.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})

	router.PUT("/api/game", func(c *gin.Context) {
		bytes, err := ioutil.ReadAll(c.Request.Body)

		if err != nil {
			c.JSON(200, gin.H{"message": "bad"})
		}

		event, err := UnmarshalEvent(bytes)

		if err != nil {
			c.JSON(200, gin.H{"message": "bad"})
		}

		hub.event <- event

		c.JSON(200, gin.H{
			"message": "cool",
		})
	})

	router.GET("/api/game", func(c *gin.Context) {
		receiveChan := make(chan GameState)

		hub.state <- receiveChan

		val := <-receiveChan

		c.JSON(200, val)
	})

	router.Run(":" + port)
}
