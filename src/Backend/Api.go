package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "sync"
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

func validateURL(url string) bool {
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Failed to fetch URL %s: %v", url, err)
        return false
    }
    defer resp.Body.Close()
    return resp.StatusCode == 200
}

func BFSHandler(w http.ResponseWriter, r *http.Request) {
    
    startArticle := r.URL.Query().Get("start")
    targetArticle := r.URL.Query().Get("target")
    
    // fmt.Println(startArticle)
    // fmt.Println(targetArticle)
    
    startArticleName := extractArticleName(startArticle)
    targetArticleName := extractArticleName(targetArticle)

    fullStartURL := "https://en.wikipedia.org/wiki/" + startArticleName
    fullTargetURL := "https://en.wikipedia.org/wiki/" + targetArticleName

    // Validate URLs
    if !validateURL(fullStartURL) || !validateURL(fullTargetURL) {
        http.Error(w, "Start or target articles do not exist", http.StatusBadRequest)
        return
    }

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

    // Validate URLs
    if !validateURL(fullStartURL) || !validateURL(fullTargetURL) {
        http.Error(w, "Start or target articles do not exist", http.StatusBadRequest)
        return
    }

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

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
	
func BFS(startURL, endURL string) ([]string, int, int, time.Duration) {
    var articlesChecked int

    visited := make(map[string]bool)
    queue := []*Node{{URL: startURL}}
    visited[startURL] = true
    batchSize := 15

    file, err := os.Create("log-bfs.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    start := time.Now()

    var mutex sync.Mutex

    for len(queue) > 0 {
        batch := queue[:min(len(queue), batchSize)]
        queue = queue[min(len(queue), batchSize):]

        var wg sync.WaitGroup
        found := false
        var foundNode *Node

        for _, current := range batch {
            wg.Add(1)
            go func(current *Node) {
                defer wg.Done()
                links := getLinks(current.URL)

                mutex.Lock()
                for _, link := range links {
                    if !visited[link] {
                        articlesChecked++
                        visited[link] = true
                        child := &Node{URL: link, Parent: current}
                        current.Children = append(current.Children, child)
                        queue = append(queue, child)

                        if link == endURL {
                            found = true
                            foundNode = child
                            mutex.Unlock()
                            return
                        }
                    }
                }
                mutex.Unlock()
            }(current)
            if found {
                break
            }
        }
        wg.Wait()
        if found {
            end := time.Since(start)
            path := getPath(foundNode)
            articlesVisited := len(path)
            return path, articlesVisited, articlesChecked, end
        }
    }

    return nil, 0, articlesChecked, time.Since(start)
}





func IDS(startURL, endURL string,  file *os.File) ([]string, int, int, time.Duration) {
	startTime := time.Now()
	visited := make(map[string]bool)
	stack := []*Node{{URL: startURL, Depth: 0}}
	var depthLimit, visits, checks int
	var result []string

	var wg sync.WaitGroup
	var mutex sync.Mutex

    checks = 0

	runSearch := func(stack []*Node, endURL string, depthLimit int, file *os.File, visited map[string]bool) ([]string, int, int, bool) {
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		return DLS(stack, endURL, depthLimit, file, visited)
	}

	for {
		wg.Add(1)
		path, localVisits, localChecks, found := runSearch(stack, endURL, depthLimit, file, visited)
		wg.Wait()
		if found {
			result = path
			visits = localVisits
			checks = localChecks
			break
		}
		depthLimit++
	}

	return result, visits, checks, time.Since(startTime)
}

var checks int

func DLS(stack []*Node, endURL string, depthLimit int, f *os.File, visited map[string]bool) ([]string, int, int, bool) {
	current := stack[len(stack)-1]
	f.WriteString(fmt.Sprintf("Scraping: %s\n", current.URL))

	if !visited[current.URL] {
		if current.URL != stack[0].URL {
			checks++
		}
        visited[current.URL] = true
	}

	if current.URL == endURL {
		path := getPath(current)
		return path, len(path) - 1, checks, true
	}

	if depthLimit <= 0 {
		return nil, 0, checks, false
	}

	links := getLinks(current.URL)
	for _, link := range links {
		child := &Node{URL: link, Parent: current}
		current.Children = append(current.Children, child)
		stack = append(stack, child)
		path, vis, chk, found := DLS(stack, endURL, depthLimit-1, f, visited)
		if found {
			return path, vis, chk, true
		}
	}

	return nil, 0, checks, false
}


func extractArticleName(url string) string {
    parts := strings.Split(url, "/")
    return parts[len(parts)-1]
}

var linkCache sync.Map

func getLinks(URL string) []string {
    if value, ok := linkCache.Load(URL); ok {
        return value.([]string)
    }

    resp, err := http.Get(URL)
    if err != nil {
        log.Printf("Error fetching URL %s: %v", URL, err)
        return nil // Return nil or an empty slice depending on your use case
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        log.Printf("Status code error: %d %s", resp.StatusCode, resp.Status)
        return nil
    }

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    var prefixes = []string{
        "/wiki/Main_Page",
        "/wiki/Main_Page:",
        "/wiki/Special:",
        "/wiki/Talk:",
        "/wiki/User:",
        "/wiki/Portal:",
        "/wiki/Wikipedia:",
        "/wiki/File:",
        "/wiki/Category:",
        "/wiki/Help:",
        "/wiki/Template:",
        "/wiki/Draft:",
        "/wiki/Module:",
        "/wiki/MediaWiki:",
        "/wiki/Index:",
        "/wiki/Education_Program:",
        "/wiki/TimedText:",
        "/wiki/Gadget:",
        "/wiki/Gadget_Definition:",
    }

    links := []string{}
    doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
        link, _ := s.Attr("href")
        if strings.HasPrefix(link, "/wiki/") {
            skip := false
            for _, prefix := range prefixes {
                if strings.HasPrefix(link, prefix) {
                    skip = true
                    break
                }
            }
            if !skip {
                links = append(links, "https://en.wikipedia.org"+link)
            }
        }
    })
    linkCache.Store(URL, links)
    return links
}


func getPath(endNode *Node) []string {
    path := []string{}
    current := endNode
    for current != nil {
        path = append([]string{current.URL}, path...) // Tambahkan URL di depan slice
        current = current.Parent
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

func recoverMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Recovered from a panic: %v", err)
                http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func main() {
    corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowCredentials: true,
    })

    mux := http.NewServeMux()
    mux.HandleFunc("/shortestpath", ShortestPathHandler)

    handler := corsHandler.Handler(recoverMiddleware(mux))

    fmt.Println("Server listening on port 8080...")
    if err := http.ListenAndServe(":8080", handler); err != nil {
        log.Printf("Error starting server: %s", err)
    }
}
