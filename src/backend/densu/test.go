package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Node struct {
	URL      string
	Parent   *Node
	Children []*Node
}

func BFS(startURL, endURL string) *Node {
	visited := make(map[string]bool)
	queue := []*Node{{URL: startURL}}

	for len(queue) > 0 {
		// fmt.Println("Queue contents:")
		// for _, node := range queue {
		// 	fmt.Println(node.URL)
		// }
		// fmt.Println("-------------------")

		current := queue[0]
		queue = queue[1:]

		if current.URL == endURL {
			return current
		}

		if visited[current.URL] {
			continue
		}
		visited[current.URL] = true

		links := getLinks(current.URL)
		for _, link := range links {
			child := &Node{URL: link, Parent: current}
			current.Children = append(current.Children, child)
			queue = append(queue, child)
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
	startURL := "https://en.wikipedia.org/wiki/Mathematics"
	endURL := "https://en.wikipedia.org/wiki/Suicide"

	fmt.Println("Finding path from", startURL, "to", endURL)

	endNode := BFS(startURL, endURL)
	if endNode == nil {
		fmt.Println("Path not found!")
		return
	}

	path := getPath(endNode)
	fmt.Println("Path found:")
	for _, link := range path {
		fmt.Println(link)
	}
}
