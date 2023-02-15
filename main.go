package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var chat = &Chat{}

type Request struct {
	Action   string `json:"action"`
	RoomName string `json:"room_name"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if messageType != websocket.TextMessage {
			log.Println("invalid message type")
			continue
		}

		var request Request
		if err := json.Unmarshal(raw, &request); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("incorrect request"))
		}

		switch action := request.Action; action {
		case "CreateRoom":
			if err := chat.CreateRoom(request.RoomName); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				return
			}
			resMessage := fmt.Sprintf("created room with name %s", string(raw))
			conn.WriteMessage(websocket.TextMessage, []byte(resMessage))
		case "SendMessage":
			var message Message
			message.username = request.Username
			message.message = request.Message

			if err := chat.SendMessage(request.RoomName, message); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				return
			}
			resMessage := fmt.Sprintf("sent message with text: %s to room: %s", string(request.Message), request.RoomName)
			if err := conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
				log.Println(err)
				return
			}
		case "JoinRoom":
			var user User
			request.Username = user.name

			if err := chat.JoinRoom(request.RoomName, user); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			}
			resMessage := fmt.Sprintf("added user: %s to room: %s", string(request.Username), request.RoomName)
			if err := conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
				log.Println(err)
				return
			}
		case "LeaveRoom":
			var user User
			request.Username = user.name

			chat.LeaveRoom(request.RoomName, user)
			resMessage := fmt.Sprintf("user: %s left room: %s", string(request.Username), request.RoomName)
			if err := conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
				log.Println(err)
				return
			}
		}
	}

}
