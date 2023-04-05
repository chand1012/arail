package pages

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func ExtractPageData(url string) (string, error) {
	// Download the webpage
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download webpage: %v", err)
	}
	defer resp.Body.Close()

	// Parse the HTML document
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML document: %v", err)
	}

	// Extract the main text content from the page
	var textContent string
	var traverseNode func(*html.Node)
	traverseNode = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "p" || n.Data == "div") {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					textContent += c.Data
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverseNode(c)
		}
	}
	traverseNode(doc)

	return textContent, nil
}
