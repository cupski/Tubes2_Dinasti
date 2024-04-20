package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"
	"container/list"

    "github.com/PuerkitoBio/goquery"
)

// Graph represents the adjacency list of a graph
type Graph map[string][]string

// ShortestPathResponse represents the response structure for the shortest path
type ShortestPathResponse struct {
    Checked      int      `json:"checked"`
    Path         []string `json:"path"`
    SearchTimeMs int64    `json:"search_time_ms"`
    Error        string   `json:"error,omitempty"`
}

func addCorsHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
    mux.HandleFunc("/shortestpath", shortestPathHandler)

    // Apply CORS middleware
    corsMux := addCorsHeaders(mux)

    log.Fatal(http.ListenAndServe(":8080", corsMux))
}

func shortestPathHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    startArticle := r.URL.Query().Get("start")
    endArticle := r.URL.Query().Get("target")

    if startArticle == "" || endArticle == "" {
        http.Error(w, "Start and target articles are required", http.StatusBadRequest)
        return
    }

    graph := make(Graph)
    links, err := getLinks(startArticle)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to get links for start article: %v", err), http.StatusInternalServerError)
        return
    }
    graph[startArticle] = links

    startTime := time.Now()

    checked, path, found := bfs(graph, startArticle, endArticle)

    elapsed := time.Since(startTime)

    response := ShortestPathResponse{
        Checked:      checked,
        Path:         path,
        SearchTimeMs: elapsed.Milliseconds(),
    }

    if !found {
        response.Error = "Target article not found."
    }

    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func getLinks(title string) ([]string, error) {
	fmt.Println(title)
    title = strings.Replace(title, " ", "_", -1)
    url := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", title)

    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("error fetching page: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
    }

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error parsing page: %w", err)
    }

    var links []string
    doc.Find("#mw-content-text a").Each(func(i int, s *goquery.Selection) {
        if href, exists := s.Attr("href"); exists {
            if strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
                linkTitle := strings.TrimPrefix(href, "/wiki/")
                links = append(links, linkTitle)
            }
        }
    })

    return links, nil
}

func bfs(graph Graph, start string, end string) (int, []string, bool) {
    visited := make(map[string]bool)
    parent := make(map[string]string)
    queue := list.New()

    queue.PushBack(start)
    visited[start] = true

    for queue.Len() > 0 {
        current := queue.Remove(queue.Front()).(string)

        if current == end {
            path := make([]string, 0)
            for step := end; step != ""; step = parent[step] {
                path = append([]string{step}, path...)
            }
            return len(visited), path, true
        }

        if _, found := graph[current]; !found {
            newLinks, err := getLinks(current)
            if err != nil {
                fmt.Println("Error retrieving links for", current, ":", err)
                continue
            }
            graph[current] = newLinks
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

// import React, { useState } from 'react';

// function App() {
//   const [startArticle, setStartArticle] = useState('');
//   const [targetArticle, setTargetArticle] = useState('');
//   const [result, setResult] = useState(null);

//   const handleSubmit = async (e) => {
//     e.preventDefault();
//     try {
//       const response = await fetch(`http://localhost:8080/shortestpath?start=${startArticle}&target=${targetArticle}`);
//       const data = await response.json();
//       setResult(data);
//     } catch (error) {
//       console.error('Error:', error);
//     }
//   };

//   return (
//     <div>
//       <h1>Shortest Path Finder</h1>
//       <form onSubmit={handleSubmit}>
//         <label>
//           Start Article:
//           <input type="text" value={startArticle} onChange={(e) => setStartArticle(e.target.value)} />
//         </label>
//         <label>
//           Target Article:
//           <input type="text" value={targetArticle} onChange={(e) => setTargetArticle(e.target.value)} />
//         </label>
//         <button type="submit">Find Shortest Path</button>
//       </form>
//       {result && (
//         <div>
//           <h2>Result</h2>
//           <p>Checked articles: {result.checked}</p>
//           <p>Search time: {result.search_time_ms} ms</p>
//           <p>Path to target article: {result.path.join(' -> ')}</p>
//         </div>
//       )}
//     </div>
//   );
// }

// export default App;
