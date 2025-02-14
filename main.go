package main

import (
	"HMR/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var oldContent string

func main() {
	// Start watching files for changes
	go watchFiles([]string{"ws.html"}, func(a, b, c string) {})

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

// Detect changes between old and new content in a slice of watched files
func watchFiles(files []string, callback func(string, string, string)) {
	lastModTimes := make(map[string]time.Time)
	_ = callback

	for {
		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				fmt.Println("Error stating file:", err)
				continue
			}

			lastModTime := info.ModTime()
			if lastMod, ok := lastModTimes[file]; !ok || lastModTime.After(lastMod) {
				lastModTimes[file] = lastModTime
				content, err := os.ReadFile(file)
				if err != nil {
					fmt.Println("Error reading file:", err)
					continue
				}

				newContent := string(content)
				if oldContent != "" {
					// Compare old and new content
					selector, fragment := utils.DetectChanges(oldContent, newContent)
					fmt.Printf("Selector: %s\nFragment: %s\n", selector, fragment)
					if selector != "" && fragment != "" {
						callback(file, selector, fragment) // Notify clients of the change
					}
				}

				oldContent = newContent // Update old content
			}
		}
		time.Sleep(2 * time.Second) // Check every 2 seconds
	}
}
