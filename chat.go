package main

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Chat struct {
	mu    sync.Mutex
	rooms map[string]Room
}

type Room struct {
	messages []Message
	users    map[string]User
}

type User struct {
	name string
	conn *websocket.Conn
}

type Message struct {
	username string
	message  string
}

var (
	ErrUserNotInChat     = errors.New("user is not in the chat")
	ErrSendMessage       = errors.New("error when sending message")
	ErrRoomAlreadyExists = errors.New("room with that name already exists")
	ErrNoRoomWithName    = errors.New("room with that name doesn't exist")
)

func (c *Chat) CreateRoom(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.rooms[name]; exists {
		return ErrRoomAlreadyExists
	}

	var emptyRoom Room
	c.rooms[name] = emptyRoom
	return nil
}

func (c *Chat) SendMessage(roomName string, message Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.rooms[roomName].users[message.username]; !exists {
		return ErrUserNotInChat
	}

	if room, ok := c.rooms[roomName]; ok {
		room.messages = append(room.messages, message)
	} else {
		return ErrSendMessage
	}
	return nil
}

func (c *Chat) JoinRoom(roomName string, user User) error {
	if room, exists := c.rooms[roomName]; exists {
		room.users[user.name] = user
	} else {
		return ErrNoRoomWithName
	}
	return nil
}

func (c *Chat) LeaveRoom(roomName string, user User) {
	delete(c.rooms[roomName].users, user.name)
}
