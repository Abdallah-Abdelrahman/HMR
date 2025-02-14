package main

import (
	"HMR/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type Data struct {
	File     string `json:"file"`
	Selector string `json:"selector"`
	Fragment string `json:"fragment"`
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var clientsLock sync.Mutex                   // Ensure concurrent access to clients map is safe

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	var appDir string

	if len(os.Args) < 2 {
		log.Fatal("Please provide the directory path of the app to watch!")
	}

	appDir = os.Args[1]
	// Start watching files for changes
	go utils.WatchFiles(utils.ExtractHTMLFiles(appDir), notifyClient)

	// Start the websocket server
	fmt.Println("Websocket server started on localhost:8080/ws")
	http.HandleFunc("/ws", wsHandler)

	// Error handling for the server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err) // Log the error and exit
	}
}

// notifyClient sends the updated selector and fragment to all connected clients
func notifyClient(file, selector, fragment string) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	data := Data{
		File:     file,
		Selector: selector,
		Fragment: fragment,
	}

	msg, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}

	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Error sending message to client:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// Handle websocket connections and messages
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err) // log instead of panic
		return
	}

	// Add the new client to the map
	clientsLock.Lock()
	clients[conn] = true
	clientsLock.Unlock()

	defer func() {
		// Remove the client when they disconnect
		clientsLock.Lock()
		delete(clients, conn)
		clientsLock.Unlock()
		conn.Close()
	}()

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage error:", err)
			break
		}

		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		if err := conn.WriteMessage(msgType, []byte(fmt.Sprintf("Server successfully recieved %s", msg))); err != nil {
			log.Println("WriteMessage error:", err)
			break
		}
	}
}
