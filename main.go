package main

import (
	"fmt"
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

	router.GET("/ws/games/:id", func(c *gin.Context) {
		gameId := c.Param("id")
		serveWs(hub, c.Writer, c.Request, gameId)
	})

	router.GET("/api/games", func(c *gin.Context) {
		receiveChan := make(chan map[string]*GameState)

		request := &AllGameStatesRequest{
			receive: receiveChan,
		}

		fmt.Println("here here")

		hub.states <- request

		fmt.Println("there")

		val := <-receiveChan

		c.JSON(200, val)
	})

	router.PUT("/api/games/:id", func(c *gin.Context) {
		gameId := c.Param("id")

		bytes, err := ioutil.ReadAll(c.Request.Body)

		if err != nil {
			c.JSON(200, gin.H{"message": "bad"})
		}

		event, err := UnmarshalEvent(bytes)

		if err != nil {
			c.JSON(200, gin.H{"message": "bad"})
		}

		request := &GameEventRequest{
			gameId: gameId,
			event:  event,
		}

		hub.event <- request

		c.JSON(200, gin.H{
			"message": "cool",
		})
	})

	router.GET("/api/games/:id", func(c *gin.Context) {
		gameId := c.Param("id")

		receiveChan := make(chan *GameState)

		request := &GameStateRequest{
			gameId:  gameId,
			receive: receiveChan,
		}

		hub.state <- request

		val := <-receiveChan

		c.JSON(200, val)
	})

	router.Run(":" + port)
}
