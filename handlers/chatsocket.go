package handlers

import (
	"ChatApplication/models"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatRoom struct {
	users    map[*websocket.Conn]string
	messages chan models.ChatMessage
	lock     sync.Mutex
}

// New Chat Room
func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		users:    make(map[*websocket.Conn]string),
		messages: make(chan models.ChatMessage),
	}
}

// for read msg and send to all
func (room *ChatRoom) Run() {
	for {
		msg := <-room.messages
		room.lock.Lock()
		for user := range room.users {
			err := user.WriteJSON(msg)
			if err != nil {
				user.Close()
				delete(room.users, user)
			}
		}
		room.lock.Unlock()
	}
}

// updgrade http to websocket
var webSocketUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

//new User to the chat

func (room *ChatRoom) AddNewUser(ctx *gin.Context) {
	username := ctx.Query("username")
	conn, err := webSocketUpgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	room.lock.Lock()
	room.users[conn] = username
	room.lock.Unlock()

	room.messages <- models.ChatMessage{User: "System", Msg: fmt.Sprintf(" %s Joined the Chat", username)}

	for {
		var chat models.ChatMessage
		err := conn.ReadJSON(&chat)
		if err != nil {
			room.lock.Lock()
			delete(room.users, conn)
			room.lock.Unlock()
			room.messages <- models.ChatMessage{User: "System", Msg: fmt.Sprintf("%s ended the chat", username)}
			break
		}
		chat.User = username
		room.messages <- chat
	}

}
