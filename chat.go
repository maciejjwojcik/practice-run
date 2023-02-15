package main

import "sync"

type Chat struct { //switch to map
	mu    sync.Mutex
	rooms []string
}

type Room struct {
	messages map[int]string
}

func (c *Chat) CreateRoom(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rooms = append(c.rooms, name)
}

//join room

//leave room

//send message to channel
