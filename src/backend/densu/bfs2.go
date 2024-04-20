package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Node struct {
	URL      string
	Parent   *Node
	Children []*Node
}

func BFS(startURL, endURL string) []string {
	visited := make(map[string]bool)
	queue := []*Node{{URL: startURL}}
	file, err := os.Create("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		msg := fmt.Sprintf("Checking queue for: %s\n", current.URL)
		fmt.Print(msg)
		file.WriteString(msg)

		if current.URL == endURL {
			return getPath(current)
		}

		if visited[current.URL] {
			continue
		}
		visited[current.URL] = true

		links := getLinks(current.URL)
		endFound := false
		for _, link := range links {
			msg := fmt.Sprintf("Scraping: %s\n", link)
			fmt.Print(msg)
			file.WriteString(msg)

			if link == endURL {
				endFound = true
				break
			}

			child := &Node{URL: link, Parent: current}
			current.Children = append(current.Children, child)
			queue = append(queue, child)
		}

		if endFound {
			return getPath(&Node{URL: endURL, Parent: current})
		}
	}
	return nil
}

func getLinks(URL string) []string {
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	links := []string{}
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		if strings.HasPrefix(link, "/wiki/") {
			links = append(links, "https://en.wikipedia.org"+link)
		}
	})
	return links
}

func getPath(endNode *Node) []string {
	path := []string{}
	current := endNode
	for current != nil {
		path = append(path, current.URL)
		current = current.Parent
	}
	// Reverse the path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func main() {
	// startURL := "https://en.wikipedia.org/wiki/Justification_(epistemology)"
	// endURL := "https://en.wikipedia.org/wiki/Donald_Davidson_(philosopher)"

	startURL := "https://en.wikipedia.org/wiki/Knowledge"
	endURL := "https://en.wikipedia.org/wiki/Fortune-telling"

	fmt.Println("Finding path from", startURL, "to", endURL)

	path := BFS(startURL, endURL)
	if path == nil {
		fmt.Println("Path not found!")
		return
	}

	fmt.Println("Path found:")
	for _, link := range path {
		fmt.Println(link)
	}
}