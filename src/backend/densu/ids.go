package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Node struct {
	URL      string
	Parent   *Node
	Children []*Node
	Depth    int
}

func IDS(startURL, endURL string) ([]string, int, int, time.Duration) {
	startTime := time.Now()
	var path []string
	var articlesVisited, articlesChecked int
	var depthLimit int

	for {
		path, articlesVisited, articlesChecked = DLS(startURL, endURL, depthLimit)
		if path != nil || time.Since(startTime).Seconds() > 60 {
			break
		}
		depthLimit++
	}

	if path == nil {
		return nil, articlesVisited, articlesChecked, 0
	}

	execTime := time.Since(startTime)
	return path, articlesVisited, articlesChecked, execTime
}

func DLS(startURL, endURL string, depthLimit int) ([]string, int, int) {
	visited := make(map[string]bool)
	stack := []*Node{{URL: startURL, Depth: 0}}
	file, err := os.Create("log-ids.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var articlesVisited, articlesChecked int

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		msg := fmt.Sprintf("Checking stack for: %s\n", current.URL)
		fmt.Print(msg)
		file.WriteString(msg)
		articlesVisited++

		if current.URL == endURL {
			return getPath(current), articlesVisited, articlesChecked
		}

		if visited[current.URL] || current.Depth >= depthLimit {
			continue
		}
		visited[current.URL] = true

		links := getLinks(current.URL)
		articlesChecked++
		for _, link := range links {
			msg := fmt.Sprintf("Scraping: %s\n", link)
			fmt.Print(msg)
			file.WriteString(msg)

			child := &Node{URL: link, Parent: current, Depth: current.Depth + 1}
			current.Children = append(current.Children, child)
			stack = append(stack, child)
		}
	}

	return nil, articlesVisited, articlesChecked
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
	// reverse the path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func main() {
	startURL := "https://en.wikipedia.org/wiki/Artificial_intelligence"
	endURL := "https://en.wikipedia.org/wiki/Power_(physics)"

	fmt.Println("Mencari rute dari", startURL, "ke", endURL)

	path, articlesVisited, articlesChecked, execTime := IDS(startURL, endURL)
	if path == nil {
		fmt.Println("Rute tidak ditemukan!")
		return
	}

	fmt.Println("Rute yang ditemukan:")
	for _, link := range path {
		fmt.Println(link)
	}

	fmt.Printf("Waktu pencarian: %v ms\n", execTime.Milliseconds())
	fmt.Printf("Jumlah artikel yang dilalui (visited article(s)): %d\n", articlesVisited)
	fmt.Printf("Jumlah artikel yang diperiksa (checked article(s)): %d\n", articlesChecked)
}
