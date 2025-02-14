package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	oldContent string
	newContent string
}

var CONTENT = make(map[string]*File)

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

				if CONTENT[file] == nil {
					CONTENT[file] = &File{}
				}

				CONTENT[file].newContent = string(content)
				if CONTENT[file].oldContent != "" {
					// Compare old and new content
					selector, fragment := DetectChanges(
						CONTENT[file].oldContent,
						CONTENT[file].newContent,
					)
					fmt.Printf("Selector: %s\nFragment: %s\n", selector, fragment)
					if selector != "" && fragment != "" {
						callback(file, selector, fragment) // Notify clients of the change
					}
				}

				CONTENT[file].oldContent = CONTENT[file].newContent // Update old content
			}
		}
		time.Sleep(2 * time.Second) // Check every 2 seconds
	}
}

// Extract all .html files recursively for the given directory
func ExtractHTMLFiles(rootDir string) []string {
	var htmlFiles []string

	err := filepath.WalkDir(rootDir, func(path string, info os.DirEntry, err error) error {
		// cb function called for each node in the tree
		if err != nil {
			return err // Handle errors (e.g., permission issues)
		}

		// Check if the file has a .html extension
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			fmt.Println("Found HTML file:", path)
			htmlFiles = append(htmlFiles, path)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}

	fmt.Printf("%d files found.\n", len(htmlFiles))

	return htmlFiles
}
