package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/PuerkitoBio/goquery"
    "github.com/rs/cors"
)

type Node struct {
    URL      string
    Parent   *Node
    Children []*Node
    Depth    int
}

type ShortestPathResult struct {
    Path            []string      `json:"path"`
    ArticlesVisited int           `json:"articlesVisited"`
    ArticlesChecked int           `json:"articlesChecked"`
    ExecutionTime   time.Duration `json:"executionTime"`
}

func BFSHandler(w http.ResponseWriter, r *http.Request) {
    
    startArticle := r.URL.Query().Get("start")
    targetArticle := r.URL.Query().Get("target")
    
    fmt.Println(startArticle)
    fmt.Println(targetArticle)
    
    startArticleName := extractArticleName(startArticle)
    targetArticleName := extractArticleName(targetArticle)

    fullStartURL := "https://en.wikipedia.org/wiki/" + startArticleName
    fullTargetURL := "https://en.wikipedia.org/wiki/" + targetArticleName

    path, articlesVisited, articlesChecked, execTime := BFS(fullStartURL, fullTargetURL)

    if path == nil {
        http.Error(w, "Route not found", http.StatusNotFound)
        return
    }

    execTimeDuration := time.Duration(execTime.Nanoseconds() / int64(time.Millisecond))

    result := ShortestPathResult{
        Path:            path,
        ArticlesVisited: articlesVisited,
        ArticlesChecked: articlesChecked,
        ExecutionTime:   execTimeDuration,
    }

    jsonResponse, err := json.Marshal(result)
    if err != nil {
        http.Error(w, "Unable to marshal JSON response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}

func IDSHandler(w http.ResponseWriter, r *http.Request, f *os.File) {
    startArticle := r.URL.Query().Get("start")
    targetArticle := r.URL.Query().Get("target")

    startArticleName := extractArticleName(startArticle)
    targetArticleName := extractArticleName(targetArticle)

    fullStartURL := "https://en.wikipedia.org/wiki/" + startArticleName
    fullTargetURL := "https://en.wikipedia.org/wiki/" + targetArticleName

    path, articlesVisited, articlesChecked, execTime := IDS(fullStartURL, fullTargetURL, f)

    if path == nil {
        http.Error(w, "Route not found", http.StatusNotFound)
        return
    }

    execTimeDuration := time.Duration(execTime.Nanoseconds() / int64(time.Millisecond))

    result := ShortestPathResult{
        Path:            path,
        ArticlesVisited: articlesVisited,
        ArticlesChecked: articlesChecked,
        ExecutionTime:   execTimeDuration,
    }

    jsonResponse, err := json.Marshal(result)
    if err != nil {
        http.Error(w, "Unable to marshal JSON response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}

func BFS(startURL, endURL string) ([]string, int, int, time.Duration) {
    visited := make(map[string]bool)
    queue := []*Node{{URL: startURL}}

    file, err := os.Create("log-bfs.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

    start := time.Now()
    var articlesVisited, articlesChecked int

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        if(current.URL == startURL){
            articlesVisited++
        }

        if current.URL == endURL {
            end := time.Since(start)
            return getPath(current), articlesVisited, articlesChecked, end
        }

        if visited[current.URL] {
            continue
        }
        visited[current.URL] = true

        links := getLinks(current.URL)
        if(current.URL == startURL){
            articlesChecked++
        }
        endFound := false
        for _, link := range links {
            articlesVisited++

            msg := fmt.Sprintf("Scraping: %s\n", link)
			fmt.Print(msg)
			file.WriteString(msg)

            articlesChecked++
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


func IDS(startURL, endURL string, f *os.File) ([]string, int, int, time.Duration) {
    startTime := time.Now()
    stack := []*Node{{URL: startURL, Depth: 0}}

    var path []string
    var articlesVisited, articlesChecked, depthLimit int
    var found bool

    for {
        path, articlesVisited, articlesChecked, found = DLS(stack, endURL, depthLimit, f)
        if found {
            break
        }
        depthLimit++
    }

    // if path == nil {
    //     return nil, articlesVisited, articlesChecked, 0
    // }

    execTime := time.Since(startTime)
    return path, articlesVisited, articlesChecked, execTime
}

var articlesVisited, articlesChecked int

func DLS(stack []*Node , endURL string, depthLimit int, f *os.File) ([]string, int, int, bool) {

    var found bool
    var path []string

    current := stack[len(stack)-1]

    msg := fmt.Sprintf("Scraping: %s\n", current.URL)
    fmt.Print(msg)
    f.WriteString(msg)
    
    articlesVisited++
    articlesChecked++

    if current.URL == endURL {
        return getPath(current), articlesVisited, articlesChecked, true
    }

    if depthLimit <= 0 {
        return getPath(current), articlesVisited, articlesChecked, false
    }

    links := getLinks(current.URL)

    for _, link := range links {
        child := &Node{URL: link, Parent: current}
        current.Children = append(current.Children, child)
        stack = append(stack, child)

        path, articlesVisited, articlesChecked, found = DLS(stack, endURL, depthLimit-1, f)
        if found {
            return  path, articlesVisited, articlesChecked, true
        }
    }


    return nil, articlesVisited, articlesChecked, false
}

func extractArticleName(url string) string {
    parts := strings.Split(url, "/")
    return parts[len(parts)-1]
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

func ShortestPathHandler(w http.ResponseWriter, r *http.Request) {
    algorithm := r.URL.Query().Get("algorithm")
    switch algorithm {
    case "bfs":
        BFSHandler(w, r)
    case "ids":
        file, err := os.Create("log-ids.txt")
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()
        IDSHandler(w, r, file)
    default:
        http.Error(w, "Invalid algorithm", http.StatusBadRequest)
    }
}

func main() {
    corsHandler := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowCredentials: true,
    })

    handler := corsHandler.Handler(http.DefaultServeMux)

    http.HandleFunc("/shortestpath", ShortestPathHandler)
    fmt.Println("Server listening on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
