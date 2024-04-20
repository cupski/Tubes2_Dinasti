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
}

func BFS(startURL, endURL string) ([]string, int, int, time.Duration) {
	visited := make(map[string]bool)
	queue := []*Node{{URL: startURL}}
	file, err := os.Create("log-bfs.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	startTime := time.Now()
	var articlesVisited, articlesChecked int

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		msg := fmt.Sprintf("Checking queue for: %s\n", current.URL)
		fmt.Print(msg)
		file.WriteString(msg)
		articlesVisited++

		if current.URL == endURL {
			endTime := time.Since(startTime)
			return getPath(current), articlesVisited, articlesChecked, endTime
		}

		if visited[current.URL] {
			continue
		}
		visited[current.URL] = true

		links := getLinks(current.URL)
		articlesChecked++
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
			endTime := time.Since(startTime)
			return getPath(&Node{URL: endURL, Parent: current}), articlesVisited, articlesChecked, endTime
		}
	}
	return nil, articlesVisited, articlesChecked, 0
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
	// startURL := "https://en.wikipedia.org/wiki/French-suited_playing_cards"
	// endURL := "https://en.wikipedia.org/wiki/Indian_Premier_League"

	// startURL := "https://en.wikipedia.org/wiki/Physics"
	// endURL := "https://en.wikipedia.org/wiki/Indian_Premier_League"

	startURL := "https://en.wikipedia.org/wiki/Artificial_intelligence"
	endURL := "https://en.wikipedia.org/wiki/Physics"

	fmt.Println("Mencari rute dari", startURL, "ke", endURL)

	path, articlesVisited, articlesChecked, execTime := BFS(startURL, endURL)
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
