package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r*http.Request)bool{
		return true 
	},
}

func ServeWS(hub *Hub,w http.ResponseWriter,r *http.Request){
	// this function is being used to initate the client and register it to the hub ; 

	conn,err :=upgrader.Upgrade(w,r,nil)
	if err !=nil {
		log.Println("Error upgrading to websocket",conn)
		return 
	}
	client :=&Client{
		hub : hub, 
		conn : conn , 
		Send : make(chan[]byte), 
	}
	hub.register <-client 
	// these go routines are used to send and receive data 
	go client.WritePump()
	go client.ReadPump()

}


func main(){
	// initiate the hub 
	hub :=NewHub()
	// separate goroutine to ensure that the websocket does not block the rest of the process 
	go hub.Run()

	 r :=chi.NewRouter()
	 r.Get("/",serveHome)
	 r.Get("/ws",func(w http.ResponseWriter,r*http.Request)
{
	ServeWS(hub,w,r); 
})
log.Println("Server started on port :8080")
log.Fatal(http.ListenAndServe(":8080",r))


}