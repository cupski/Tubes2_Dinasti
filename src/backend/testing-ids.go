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

// Node merupakan representasi dari sebuah simpul dalam graf jalur.
type Node struct {
	URL      string
	Parent   *Node
	Children []*Node
	Depth    int
}

// IDS adalah fungsi utama yang mengimplementasikan algoritma Iterative Deepening Search (IDS).
func IDS(startURL, endURL string) ([]string, int, int, time.Duration) {
	startTime := time.Now()
	var path []string
	var articlesVisited, articlesChecked int
	var depthLimit int

	// Iterasi IDS dengan peningkatan batas kedalaman hingga solusi ditemukan atau batas waktu tercapai.
	for {
		path, articlesVisited, articlesChecked = DLS(startURL, endURL, depthLimit)
		// Jika solusi ditemukan atau waktu maksimum tercapai, keluar dari iterasi.
		if path != nil || time.Since(startTime).Seconds() > 60 {
			break
		}
		depthLimit++
	}

	// Jika tidak ditemukan jalur dari startURL ke endURL, kembalikan nilai-nilai default.
	if path == nil {
		return nil, articlesVisited, articlesChecked, 0
	}

	execTime := time.Since(startTime)
	return path, articlesVisited, articlesChecked, execTime
}

// DLS (Depth-Limited Search) adalah fungsi pencarian dengan batasan kedalaman.
func DLS(startURL, endURL string, depthLimit int) ([]string, int, int) {
	visited := make(map[string]bool)
	stack := []*Node{{URL: startURL, Depth: 0}}
	file, err := os.Create("log-ids-test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var articlesVisited, articlesChecked int

	// Melakukan pencarian selama masih ada simpul di tumpukan (stack).
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Menyimpan informasi tentang simpul yang sedang diperiksa ke dalam file log.
		msg := fmt.Sprintf("Checking stack for: %s\n", current.URL)
		fmt.Print(msg)
		file.WriteString(msg)
		articlesVisited++

		// Jika simpul tujuan ditemukan, kembalikan jalur yang ditemukan.
		if current.URL == endURL {
			return getPath(current), articlesVisited, articlesChecked
		}

		// Jika simpul telah dikunjungi sebelumnya atau kedalaman maksimum telah tercapai, lanjutkan ke simpul berikutnya.
		if visited[current.URL] || current.Depth >= depthLimit {
			continue
		}
		visited[current.URL] = true

		// Dapatkan tautan dari halaman URL saat ini dan tambahkan ke dalam tumpukan (stack) untuk diperiksa lebih lanjut.
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

	// Jika tidak ditemukan jalur dari startURL ke endURL dengan batasan kedalaman tertentu, kembalikan nilai-nilai default.
	return nil, articlesVisited, articlesChecked
}

// getLinks digunakan untuk mendapatkan daftar tautan dari sebuah halaman web.
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

// getPath digunakan untuk mendapatkan jalur dari simpul akhir ke simpul awal.
func getPath(endNode *Node) []string {
	path := []string{}
	current := endNode
	for current != nil {
		path = append(path, current.URL)
		current = current.Parent
	}
	// Membalikkan jalur agar menjadi urutan dari simpul awal ke simpul akhir.
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func main() {
	// startURL := "https://en.wikipedia.org/wiki/Knowledge"
	// endURL := "https://en.wikipedia.org/wiki/Fortune-telling"

	startURL := "https://en.wikipedia.org/wiki/French-suited_playing_cards"
	endURL := "https://en.wikipedia.org/wiki/Indian_Premier_League"

	fmt.Println("Mencari rute dari", startURL, "ke", endURL)

	// Memulai pencarian menggunakan IDS.
	path, articlesVisited, articlesChecked, execTime := IDS(startURL, endURL)
	if path == nil {
		fmt.Println("Rute tidak ditemukan!")
		return
	}

	// Menampilkan jalur yang ditemukan.
	fmt.Println("Rute yang ditemukan:")
	for _, link := range path {
		fmt.Println(link)
	}

	// Menampilkan statistik pencarian.
	fmt.Printf("Waktu pencarian: %v ms\n", execTime.Milliseconds())
	fmt.Printf("Jumlah artikel yang dilalui (visited article(s)): %d\n", articlesVisited)
	fmt.Printf("Jumlah artikel yang diperiksa (checked article(s)): %d\n", articlesChecked)
}
