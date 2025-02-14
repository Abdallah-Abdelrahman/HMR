package utils

import (
	"fmt"
	"os"
	"time"
)

var oldContent string

// Detect changes between old and new content in a slice of watched files
func WatchFiles(files []string, callback func(string, string, string)) {
	lastModTimes := make(map[string]time.Time)

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
					selector, fragment := DetectChanges(oldContent, newContent)
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
