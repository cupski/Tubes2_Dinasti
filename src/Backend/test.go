package main

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"net/http"
	"strings"
	"time"
	"container/list"
   ) 
   
type Graph map[string][]string


func getLinks(title string) ([]string, error) {
// Ubah spasi menjadi underscore
title = strings.Replace(title, " ", "_", -1)
url := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", title)

// Buat request HTTP
resp, err := http.Get(url)
if err != nil {
	return nil, fmt.Errorf("error fetching page: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != 200 {
	return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
}

// Parse HTML
doc, err := goquery.NewDocumentFromReader(resp.Body)
if err != nil {
	return nil, fmt.Errorf("error parsing page: %w", err)
}

var links []string
doc.Find("#mw-content-text a").Each(func(i int, s *goquery.Selection) {
	if href, exists := s.Attr("href"); exists {
	// Check if the link is an internal wiki link and not an external link or a reference link
	if strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
	linkTitle := strings.TrimPrefix(href, "/wiki/")
	links = append(links, linkTitle)
	}
	}
})

return links, nil
}

func bfs(graph Graph, start, end string) (int, []string, bool) {
visited := make(map[string]bool)
parent := make(map[string]string)
queue := list.New()

queue.PushBack(start)
visited[start] = true

for queue.Len() > 0 {
	current := queue.Remove(queue.Front()).(string)

	if current == end {
		path := make([]string, 0)
		step := current
		for step != "" {
			path = append([]string{step}, path...)
			step = parent[step]
		}

		// Cetak rute setelah ditemukan
		for _, p := range path {
			url := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", strings.Replace(p, " ", "_", -1))
			fmt.Println("Rute:", url)
		}

		return len(visited), path, true
	}

	for _, neighbor := range graph[current] {
		if !visited[neighbor] {
			visited[neighbor] = true
			parent[neighbor] = current
			queue.PushBack(neighbor)
		}
	}
}

return len(visited), nil, false
}



func ids(graph Graph, start, end string, maxDepth int) (int, []string, bool) {
for depth := 0; depth <= maxDepth; depth++ {
	visited := make(map[string]bool)
	path, found := dls(graph, start, end, depth, visited)
	if found {
		// Cetak rute setelah ditemukan
		for _, p := range path {
			url := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", strings.Replace(p, " ", "_", -1))
			fmt.Println("Rute:", url)
		}
		return len(visited), path, true
	}
}
return 0, nil, false
}


func dls(graph Graph, node, end string, depth int, visited map[string]bool) ([]string, bool) {
if node == end {
	return []string{end}, true
}
if depth == 0 {
	return nil, false
}
visited[node] = true


for _, neighbor := range graph[node] {
	if _, seen := visited[neighbor]; !seen {
		path, found := dls(graph, neighbor, end, depth-1, visited)
		if found {
			return append([]string{node}, path...), true
		}
	}
}
return nil, false
}


func main() {
var algorithm, startArticle, endArticle string
fmt.Print("Enter algorithm (IDS/BFS): ")
fmt.Scanln(&algorithm)
fmt.Print("Enter start article title: ")
fmt.Scanln(&startArticle)
fmt.Print("Enter target article title: ")
fmt.Scanln(&endArticle)

graph := make(Graph)


links, err := getLinks(startArticle)
if err != nil {
	fmt.Printf("Failed to get links for start article: %v\n", err)
	return
}
graph[startArticle] = links


startTime := time.Now()

var checked int
var path []string
var found bool

switch algorithm {
case "BFS":
	checked, path, found = bfs(graph, startArticle, endArticle)
case "IDS":
	for depth := 0; !found; depth++ {
		checked, path, found = ids(graph, startArticle, endArticle, depth)
	}
default:
	fmt.Println("Invalid algorithm specified")
	return
}

elapsed := time.Since(startTime)

if found {
	fmt.Println("Checked articles:", checked)
	fmt.Println("Path to target article:", path)
	fmt.Printf("Search time: %v ms\n", elapsed.Milliseconds())
} else {
	fmt.Println("Target article not found.")
}
}
