package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"container/list"
)

type Graph map[string][]string
type WikiResponse struct {
	Query struct {
		Pages map[string]struct {
			Links []struct {
				Title string `json:"title"`
			} `json:"links"`
		} `json:"pages"`
	} `json:"query"`
}

func getLinks(title string) ([]string, error) {
    title = strings.Replace(title, " ", "_", -1)
    url := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&titles=%s&prop=links&format=json&pllimit=max", title)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var wikiResp WikiResponse
    err = json.Unmarshal(body, &wikiResp)
    if err != nil {
        return nil, err
    }

    var links []string
    for _, page := range wikiResp.Query.Pages {
        for _, link := range page.Links {
            links = append(links, link.Title)
        }
    }

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

        // Print the current node being visited
        fmt.Println("Visiting node:", current)

        if current == end {
            path := make([]string, 0)
            step := current
            for step != "" {
                path = append([]string{step}, path...)
                step = parent[step]
            }

            // Print the path after it's found
            fmt.Println("Path found:")
            for _, p := range path {
                fmt.Println(p)
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
	visited := make(map[string]bool) // Initialize visited outside the loop

	for depth := 0; depth <= maxDepth; depth++ {
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

    // Print the current node being visited
    fmt.Println("Visiting node:", node)

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
