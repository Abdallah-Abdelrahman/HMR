package utils

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Compare the old HTML tree with the new one
// This version handles insertions, deletions, and changes, and returns the selector and the new fragment
func DetectChanges(oldContent, newContent string) (string, string) {
	// Parse the old and new HTML documents
	oldDoc, err := goquery.NewDocumentFromReader(strings.NewReader(oldContent))
	if err != nil {
		fmt.Println("Error parsing old HTML:", err)
		return "", ""
	}

	newDoc, err := goquery.NewDocumentFromReader(strings.NewReader(newContent))
	if err != nil {
		fmt.Println("Error parsing new HTML:", err)
		return "", ""
	}

	// Start DFS comparison from the root element
	return compareTrees(oldDoc.Selection, newDoc.Selection)
}

// Recursive DFS-based comparison of two elements
func compareTrees(oldElem, newElem *goquery.Selection) (string, string) {
	// Compare tag names
	if goquery.NodeName(oldElem) != goquery.NodeName(newElem) {
		selector := generateSelector(oldElem)
		fragment, _ := goquery.OuterHtml(newElem)
		fmt.Println("Tag name difference detected")
		return selector, fragment
	}

	// Compare attributes
	oldAttrs := getAttributesMap(oldElem)
	newAttrs := getAttributesMap(newElem)

	if !compareAttributes(oldAttrs, newAttrs) {
		selector := generateSelector(oldElem)
		fragment, _ := goquery.OuterHtml(newElem)
		fmt.Println("Attribute difference detected")
		return selector, fragment
	}

	// Compare children recursively
	oldChildren := oldElem.Children()
	newChildren := newElem.Children()

	oldLen := oldChildren.Length()
	newLen := newChildren.Length()

	// Detect insertion or deletion
	if oldLen != newLen {
		fmt.Printf("Difference in number of children: old=%d, new=%d\n", oldLen, newLen)
		if newLen > oldLen {
			// Insertion detected
			fmt.Println("Insertion detected!")
			selector := generateSelector(oldElem)
			fragment, _ := goquery.OuterHtml(newChildren.Parent()) // New child added
			return selector, fragment
		} else {
			// Deletion detected
			fmt.Println("Deletion detected!")
			selector := generateSelector(oldElem)
			fragment, _ := goquery.OuterHtml(newChildren.Parent()) // Old child deleted
			return selector, fragment
		}
	}

	// Continue DFS on each child
	for i := 0; i < oldLen; i++ {
		selector, fragment := compareTrees(oldChildren.Eq(i), newChildren.Eq(i))
		if selector != "" && fragment != "" {
			// Return the first detected change
			return selector, fragment
		}
	}

	// Compare text content
	if oldElem.Text() != newElem.Text() {
		selector := generateSelector(oldElem)
		fragment, _ := goquery.OuterHtml(newElem)
		fmt.Println("Text content difference detected")
		return selector, fragment
	}
	return "", ""
}

// Helper function to convert element attributes to a map
func getAttributesMap(elem *goquery.Selection) map[string]string {
	attrs := make(map[string]string)
	for _, attr := range elem.Nodes[0].Attr {
		attrs[attr.Key] = attr.Val
	}
	return attrs
}

// Compare attributes of old and new elements
func compareAttributes(oldAttrs, newAttrs map[string]string) bool {
	if len(oldAttrs) != len(newAttrs) {
		return false
	}
	for key, val := range oldAttrs {
		if newAttrs[key] != val {
			return false
		}
	}
	return true
}

// Helper function to generate a fine-grained CSS selector for an element
func generateSelector(elem *goquery.Selection) string {
	selector := goquery.NodeName(elem)

	// Add ID if it exists
	if id, exists := elem.Attr("id"); exists {
		selector += fmt.Sprintf("#%s", id)
	}

	// Add classes if they exist
	if class, exists := elem.Attr("class"); exists {
		classes := strings.Split(class, " ")
		for _, className := range classes {
			selector += fmt.Sprintf(".%s", className)
		}
	}

	// Special handling for list items (li) or similar elements
	// Add nth-child for finer granularity
	if selector == "li" || selector == "div" || selector == "span" {
		// Find the position of the element among its siblings
		parent := elem.Parent()
		allSiblings := parent.ChildrenFiltered(goquery.NodeName(elem)) // Find all siblings with the same tag
		for i := 0; i < allSiblings.Length(); i++ {
			if allSiblings.Eq(i).Get(0) == elem.Get(0) {
				// We use nth-child to represent its index in the parent
				selector += fmt.Sprintf(":nth-child(%d)", i+1)
				break
			}
		}
	}

	// Add parent selectors for more specificity
	for parent := elem.Parent(); parent.Length() > 0 && goquery.NodeName(parent) != "html" && goquery.NodeName(parent) != "body"; parent = parent.Parent() {
		parentSelector := goquery.NodeName(parent)
		if id, exists := parent.Attr("id"); exists {
			parentSelector += fmt.Sprintf("#%s", id)
		}

		if class, exists := parent.Attr("class"); exists {
			classes := strings.Split(class, " ")
			for _, className := range classes {
				parentSelector += fmt.Sprintf(".%s", className)
			}
		}

		// Prepend parent selector to the current selector
		selector = parentSelector + " > " + selector
	}

	return selector
}
