package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func convertToURL(title string) string {
	return fmt.Sprintf("https://en.wikipedia.org/wiki/%s", strings.ReplaceAll(title, " ", "_"))
}

func linkScraper(url string, visited map[string]bool) []string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	var uniqueLinks []string

	doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("href")
		if isValidArticleLink(link) && !visited[link] {
			visited[link] = true
			uniqueLinks = append(uniqueLinks, "https://en.wikipedia.org"+link)
		}
	})

	return uniqueLinks
}

func isValidArticleLink(link string) bool {
	prefixes := []string{
		"/wiki/Special:",
		"/wiki/Talk:",
		"/wiki/User:",
		"/wiki/Portal:",
		"/wiki/Wikipedia:",
		"/wiki/File:",
		"/wiki/Category:",
		"/wiki/Help:",
		"/wiki/Template:",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(link, prefix) {
			return false
		}
	}
	return strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":")
}

func IDS(startURL, goalURL string, maxDepth int, resultChan chan<- []string) {
	visited := make(map[string]bool)
	visited[startURL] = true

	for depth := 0; depth <= maxDepth; depth++ {
		path := DLS(startURL, goalURL, depth, visited)
		if len(path) > 0 {
			resultChan <- path
			return
		}
	}
	resultChan <- nil
}

func DLS(currentURL, goalURL string, depth int, visited map[string]bool) []string {
	if depth == 0 && currentURL == goalURL {
		return []string{currentURL}
	}
	if depth <= 0 || visited[currentURL] {
		return nil
	}

	visited[currentURL] = true
	links := linkScraper(currentURL, visited)
	for _, link := range links {
		if link == goalURL {
			return []string{currentURL, link}
		}
		if path := DLS(link, goalURL, depth-1, visited); path != nil {
			return append([]string{currentURL}, path...)
		}
	}
	return nil
}

func main() {
	var startTitle, goalTitle string
	fmt.Printf("Enter the start title: ")
	fmt.Scanln(&startTitle)
	fmt.Printf("Enter the goal title: ")
	fmt.Scanln(&goalTitle)

	startURL := convertToURL(startTitle)
	goalURL := convertToURL(goalTitle)

	start := time.Now()

	maxDepth := 0

	resultChan := make(chan []string)
	go IDS(startURL, goalURL, maxDepth, resultChan)

	select {
	case path := <-resultChan:
		if path != nil {
			fmt.Println("IDS Shortest Path:")
			for _, node := range path {
				title := getTitle(node)
				fmt.Println(strings.ReplaceAll(title, "_", " "))
			}
			fmt.Println("Length of Path:", len(path)-1)
			fmt.Println("IDS Time:", time.Since(start))
		} else {
			fmt.Println("No path found.")
		}
	case <-time.After(60 * time.Second):
		fmt.Println("Search time exceeded 60 seconds.")
	}
}

func getTitle(urlString string) string {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return ""
	}

	pathParts := strings.Split(parsedURL.Path, "/")
	lastPart := pathParts[len(pathParts)-1]

	title, err := url.PathUnescape(lastPart)
	if err != nil {
		return ""
	}

	return title
}
