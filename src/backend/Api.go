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

// validasi input URL
// -> cek startURL sama targetURL beneran valid apa ngga, kalo ngga gausah lakukan searching 
func validateURL(url string) bool {
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Failed to fetch URL %s: %v", url, err)
        return false
    }
    defer resp.Body.Close()
    return resp.StatusCode == 200
}

// handler untuk bfs search
// -> ambil startURL dan targetURL dari frontend, return hasil ke frontend juga
func BFSHandler(w http.ResponseWriter, r *http.Request) {
    
    startArticle := r.URL.Query().Get("start")
    targetArticle := r.URL.Query().Get("target")
    
    // fmt.Println(startArticle)
    // fmt.Println(targetArticle)
    
    startArticleName := extractArticleName(startArticle)
    targetArticleName := extractArticleName(targetArticle)

    fullStartURL := "https://en.wikipedia.org/wiki/" + startArticleName
    fullTargetURL := "https://en.wikipedia.org/wiki/" + targetArticleName

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

// handler untuk ids search
// -> ambil startURL dan targetURL dari frontend, return hasil ke frontend juga
func IDSHandler(w http.ResponseWriter, r *http.Request, f *os.File) {
    startArticle := r.URL.Query().Get("start")
    targetArticle := r.URL.Query().Get("target")

    startArticleName := extractArticleName(startArticle)
    targetArticleName := extractArticleName(targetArticle)

    fullStartURL := "https://en.wikipedia.org/wiki/" + startArticleName
    fullTargetURL := "https://en.wikipedia.org/wiki/" + targetArticleName

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

// fungsi bfs
func BFS(startURL, endURL string) ([]string, int, int, time.Duration) {
    var articlesChecked int

    visited := make(map[string]bool)
    queue := []*Node{{URL: startURL}}
    visited[startURL] = true
    batchSize := 15

    start := time.Now()

    var mutex sync.Mutex

    // Jika URL awal dan akhir sama, kembalikan jalur yang hanya berisi URL awal
    if startURL == endURL {
        return []string{startURL}, 0, 0, time.Since(start)
    }

    for len(queue) > 0 {
        // Ekstrak sebuah batch URL dari antrian berdasarkan ukuran batch
        batch := queue[:min(len(queue), batchSize)]
        queue = queue[min(len(queue), batchSize):]


        var wg sync.WaitGroup // Gunakan WaitGroup untuk menyinkronkan goroutine
        found := false // Flag untuk menandakan jika URL akhir ditemukan dalam batch saat ini
        var foundNode *Node

        // Iterasi setiap URL dalam batch dan proses secara bersamaan
        for _, current := range batch {
            time.Sleep(1 * time.Millisecond) // delay biar ga ke block wikipedianya
            fmt.Println("URL: ", current.URL)
            fmt.Println("Articles Checked: ", articlesChecked)
            wg.Add(1)
            go func(current *Node) {
                defer wg.Done()
                links := getLinks(current.URL)

                mutex.Lock()

                // Iterasi setiap tautan dan prosesnya
                for _, link := range links {
                    if !visited[link] {
                        articlesChecked++ // Tambahkan jumlah artikel yang diperiksa
                        visited[link] = true

                        // Buat node anak untuk tautan saat ini
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

        wg.Wait()  // Tunggu semua goroutine dalam batch selesai diproses

        if found {
            end := time.Since(start)
            path := getPath(foundNode)  // Dapatkan jalur terpendek dari foundNode
            articlesVisited := len(path)
            return path, articlesVisited-1, articlesChecked, end
        }
    }

    // Jika URL akhir tidak ditemukan, kembalikan nilai-nilai nol untuk jalur dan metrik lainnya
    return nil, len(visited)-1, articlesChecked, time.Since(start)
}



//fungsi ids
func IDS(startURL, endURL string, file *os.File) ([]string, int, int, time.Duration) {
    startTime := time.Now()
    visited := make(map[string]bool)
    stack := []*Node{{URL: startURL, Depth: 0}}
    var visits, checks int
    var result []string

    var wg sync.WaitGroup

    // Fungsi untuk menjalankan goroutine dalam pencarian dengan batasan kedalaman tertentu
    runSearch := func(stack []*Node, endURL string, depthLimit int, file *os.File, visited map[string]bool) ([]string, int, int, bool) {
        defer wg.Done()
        return DLS(stack, endURL, depthLimit, file, visited)
    }

    localfound := false // Flag untuk menandakan jika URL akhir ditemukan secara lokal

    // Lakukan IDS dengan batasan kedalaman dari 0 hingga 5, sehingga hanya ada 5 goroutine dulu berjalan bersamaan
    for depthLimit := 0; depthLimit <= 5; depthLimit++ {
        // time.Sleep(5 * time.Millisecond)
        wg.Add(1)
        path, localVisits, localChecks, found := runSearch(stack, endURL, depthLimit, file, visited)
        if found {
            localfound = true
            result = path
            visits = localVisits
            checks = localChecks
            break
        }
    }

    // Jika URL akhir tidak ditemukan dalam batasan kedalaman 0 hingga 5, lanjutkan dengan batasan 5 hingga 9
    if !localfound{
        for depthLimit := 5; depthLimit <= 9; depthLimit++ {
            // time.Sleep(5 * time.Millisecond)
            wg.Add(1)
            path, localVisits, localChecks, found := runSearch(stack, endURL, depthLimit, file, visited)
            if found {
                localfound = true
                result = path
                visits = localVisits
                checks = localChecks
                break
            }
        }
    }


    wg.Wait() // Tunggu semua goroutine selesai

    return result, visits, checks, time.Since(startTime)
}

var checks int // Variabel global untuk menyimpan jumlah pemeriksaan artikel di IDS

// fungsi dls 
// -> dipakai di ids
func DLS(stack []*Node, endURL string, depthLimit int, f *os.File, visited map[string]bool) ([]string, int, int, bool) {

    var mutex sync.Mutex
    var articlesMutex sync.Mutex // Mutex untuk sinkronisasi akses ke variabel jumlah pemeriksaan artikel

    mutex.Lock()
    current := stack[len(stack)-1]
    fmt.Println("URL: ", current.URL)
    fmt.Println("Articles Checked: ", checks)

	if !visited[current.URL] {
		if current.URL != stack[0].URL {
            articlesMutex.Lock()
			checks++
            articlesMutex.Unlock()
		}
        visited[current.URL] = true
	}

	if current.URL == endURL {
		path := getPath(current)
        mutex.Unlock()
		return path, len(path) - 1, checks, true
	}

    // Jika kedalaman batasan telah tercapai, kembalikan tanpa menemukan jalur
	if depthLimit <= 0 {
        mutex.Unlock()
		return nil, 0, checks, false
	}

	links := getLinks(current.URL)
    // Iterasi setiap tautan dan tambahkan anak ke stack untuk dilakukan pencarian lebih lanjut
	for _, link := range links {
		child := &Node{URL: link, Parent: current}
		current.Children = append(current.Children, child)
		stack = append(stack, child)
        // Lakukan DLS rekursif pada anak ini
		path, vis, chk, found := DLS(stack, endURL, depthLimit-1, f, visited)
		
        // Jika jalur ditemukan, kembalikan jalur tersebut
        if found {
            mutex.Unlock()
			return path, vis, chk, true
		}
	}
    mutex.Unlock()
	return nil, 0, checks, false
}

// ambil nama artikel doang
func extractArticleName(url string) string {
    parts := strings.Split(url, "/")
    return parts[len(parts)-1]
}

var (
    linkCache  sync.Map
    htmlCache  sync.Map
)

// ambil HTML content dari URL (parse pakai goquery)
func getHTMLContent(URL string) (*goquery.Document, error) {
    // check kalo HTML content nya udah ada di cache
    if value, ok := htmlCache.Load(URL); ok {
        return value.(*goquery.Document), nil
    }

    resp, err := http.Get(URL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
    }

    // parse html content dengan goquery
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, err
    }

    // cache html nya buat nanti dipake (kalo dibutuhkan pas search lain)
    htmlCache.Store(URL, doc)
    return doc, nil
}

// ambil links dari wikipedia
func getLinks(URL string) []string {
    // cek apakah link nya udah ada di cache
    if value, ok := linkCache.Load(URL); ok {
        return value.([]string)
    }

    // fetch html content terus parse utk dpt link nya
    doc, err := getHTMLContent(URL)
    if err != nil {
        log.Printf("Failed to fetch HTML content for %s: %v", URL, err)
        // handle 404 error -> langsung mark visited aja
        linkCache.Store(URL, []string{})
        return []string{}
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

    // cache linknya utk dipake lagi
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

// handler utk usernya pilih bfs/ ids
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