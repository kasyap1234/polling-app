package main

import (
"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	Send chan []byte 
}

// sending message from the channel to the client ; 
func (c*Client) WritePump(){
	defer c.conn.Close()
	for message :=range c.Send{
		err :=c.conn.WriteMessage(websocket.TextMessage,message)
		if err !=nil {
			log.Println("Websocket write error",err)
		}

	}

}

func (c*Client) ReadPump(){
	
	defer func(){
		c.hub.unregister <- c 
		c.conn.Close()
	}()
	for {
		_,message,err :=c.conn.ReadMessage()
	
		if err !=nil {
			log.Println("websocket read error",err)

		}
		c.hub.HandleVote(message)
	}

}



