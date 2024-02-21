package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


func main() {
	var addr = "localhost:8001"
	http.HandleFunc("/ws",WsHandler)
	http.ListenAndServe(addr, nil)
}


func WsHandler (w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error: ", err)
	}
	defer conn.Close()

	// heart beat
	go func() {		
		err = conn.WriteMessage(1, []byte("heart:beat"))
		if err != nil {
			log.Println("write error: ", err)
		}
	}()


	for {
		mt, msg, err := conn.ReadMessage()

		if err != nil {
			log.Println("read error: ", err)
			break
		}

		newMsg := string(msg) + " callback"
		msg = []byte(newMsg)

		err = conn.WriteMessage(mt, msg)
		if err != nil {
			log.Println("write error: ", err)
			break
		}
	}

}