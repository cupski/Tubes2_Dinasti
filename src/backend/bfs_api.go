package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/PuerkitoBio/goquery"
)

type Node struct {
    URL      string
    Parent   *Node
    Children []*Node
}

type BFSResult struct {
    Path            []string      `json:"path"`
    ArticlesVisited int           `json:"articlesVisited"`
    ArticlesChecked int           `json:"articlesChecked"`
    ExecutionTime   time.Duration `json:"executionTime"`
}

func BFS(startURL, endURL string) ([]string, int, int, time.Duration) {
    visited := make(map[string]bool)
    queue := []*Node{{URL: startURL}}
    start := time.Now()
    var articlesVisited, articlesChecked int

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        articlesVisited++

        if current.URL == endURL {
            end := time.Since(start)
            return getPath(current), articlesVisited, articlesChecked, end
        }

        if visited[current.URL] {
            continue
        }
        visited[current.URL] = true

        links := getLinks(current.URL)
        articlesChecked++
        endFound := false
        for _, link := range links {
            if link == endURL {
                endFound = true
                break
            }

            child := &Node{URL: link, Parent: current}
            current.Children = append(current.Children, child)
            queue = append(queue, child)
        }

        if endFound {
            end := time.Since(start)
            return getPath(&Node{URL: endURL, Parent: current}), articlesVisited, articlesChecked, end
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

func BFSHandler(w http.ResponseWriter, r *http.Request) {
    startArticle := r.URL.Query().Get("start")
    targetArticle := r.URL.Query().Get("target")

    log.Printf("Received request: startArticle=%s, targetArticle=%s\n", startArticle, targetArticle)

    path, articlesVisited, articlesChecked, execTime := BFS(startArticle, targetArticle)
    if path == nil {
        http.Error(w, "Route not found", http.StatusNotFound)
        return
    }

    log.Printf("Route found! Path: %v\n", path)
    log.Printf("Articles visited: %d\n", articlesVisited)
    log.Printf("Articles checked: %d\n", articlesChecked)
    log.Printf("Execution time: %v\n", execTime)
    
    result := BFSResult{
        Path:            path,
        ArticlesVisited: articlesVisited,
        ArticlesChecked: articlesChecked,
        ExecutionTime:   execTime,
    }

    jsonResponse, err := json.Marshal(result)
    if err != nil {
        http.Error(w, "Unable to marshal JSON response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}



func main() {
    http.HandleFunc("/shortestpath", BFSHandler)
    fmt.Println("Server listening on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
