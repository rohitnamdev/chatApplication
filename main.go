package main

import (
	"ChatApplication/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	//Gin Router
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	//render the html page - Chat API
	router.GET("/chat", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	//build chat room for chat
	room := handlers.NewChatRoom()
	go room.Run()
	//for websocket
	router.GET("/ws", func(ctx *gin.Context) {
		room.AddNewUser(ctx)
	})

	//Run Server
	router.Run(":8080")
}
