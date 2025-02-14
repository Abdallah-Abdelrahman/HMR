package main

import (
	"HMR/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

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
		log.Fatal("Please provide the directory of the app to watch")
	}

	appDir = os.Args[1]
	// Start watching files for changes
	go utils.WatchFiles(utils.ExtractHTMLFiles(appDir), func(a, b, c string) {})

	// Start the websocket server
	fmt.Println("Websocket server started on localhost:8080/ws")
	http.HandleFunc("/ws", wsHandler)

	// Error handling for the server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err) // Log the error and exit
	}
}

// Handle websocket connections and messages
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err) // log instead of panic
		return
	}

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
