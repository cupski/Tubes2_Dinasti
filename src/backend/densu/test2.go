package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "math"
    "net/http"
    "strings"
    "time"

    "golang.org/x/net/html"
)

// Struktur untuk menyimpan informasi artikel
type Artikel struct {
    Judul      string
    Tautan     []string
    Diperiksa  bool
    Jarak      int
    Sebelumnya *Artikel
}

// Fungsi untuk mengunduh isi artikel dari Wikipedia
func getArtikel(judul string) (*Artikel, error) {
    url := fmt.Sprintf("https://id.wikipedia.org/wiki/%s", judul)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("gagal mengunduh artikel: %s", resp.Status)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return parseArtikel(body), nil
}

// Fungsi untuk mengurai isi artikel dan mengekstrak tautan
func parseArtikel(body []byte) *Artikel {
    judul := extractJudul(body)
    tautan := extractTautan(body)

    return &Artikel{
        Judul:      judul,
        Tautan:     tautan,
        Diperiksa:  false,
        Jarak:      math.MaxInt,
        Sebelumnya: nil,
    }
}

// Fungsi untuk mengekstrak judul artikel dari HTML
func extractJudul(body []byte) string {
    // Gunakan pustaka `html` untuk mengurai HTML
    doc, err := html.Parse(bytes.NewReader(body))
    if err != nil {
        return ""
    }

    var judul string

    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "h1" {
            judul = n.FirstChild.Data
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)

    return judul
}

// Fungsi untuk mengekstrak tautan artikel dari HTML
func extractTautan(body []byte) []string {
    // Gunakan pustaka `html` untuk mengurai HTML
    doc, err := html.Parse(bytes.NewReader(body))
    if err != nil {
        return nil
    }

    var tautan []string

    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, attr := range n.Attr {
                if attr.Key == "href" {
                    href := attr.Val
                    if strings.HasPrefix(href, "/wiki/") && !strings.HasPrefix(href, "/wiki/Special:") {
                        tautan = append(tautan, strings.TrimPrefix(href, "/wiki/"))
                    }
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)

    return tautan
}

// Fungsi untuk melakukan pencarian dengan algoritma BFS
func bfs(awal *Artikel, tujuan string) ([]*Artikel, int, time.Duration) {
    antrian := []*Artikel{awal}
    diperiksa := make(map[string]bool)
    diperiksa[awal.Judul] = true

    waktuMulai := time.Now()

    for len(antrian) > 0 {
        artikel := antrian[0]
        antrian = antrian[1:]

        if artikel.Judul == tujuan {
            return reconstructPath(artikel), len(diperiksa), time.Since(waktuMulai)
        }

        for _, tautan := range artikel.Tautan {
            nextArtikel, err := getArtikel(tautan)
            if err != nil {
                continue
            }

            if !diperiksa[nextArtikel.Judul] {
                nextArtikel.Diperiksa = true
                nextArtikel.Jarak = artikel.Jarak + 1
                nextArtikel.Sebelumnya = artikel
                diperiksa[nextArtikel.Judul] = true

                antrian = append(antrian, nextArtikel)
            }
        }
    }

    return nil, len(diperiksa), time.Since(waktuMulai)
}

// Fungsi untuk melakukan pencarian dengan algoritma IDS
func ids(awal *Artikel, tujuan string, batasKedalaman int) ([]*Artikel, int, time.Duration) {
    for kedalaman := 1; kedalaman <= batasKedalaman; kedalaman++ {
        hasil, diperiksa, waktu := bfs(awal, tujuan)
        if hasil != nil {
            return hasil, diperiksa, waktu
        }
    }

    return nil, 0, 0
}

// Fungsi untuk merekonstruksi jalur dari artikel tujuan ke artikel awal
func reconstructPath(artikel *Artikel) []*Artikel {
    jalur := []*Artikel{}
    for artikel != nil {
        jalur = append(jalur, artikel)
        artikel = artikel.Sebelumnya
    }

    return reverseSlice(jalur)
}

// Fungsi untuk membalikkan urutan elemen dalam slice
func reverseSlice(s []*Artikel) []*Artikel {
    reversed := make([]*Artikel, len(s))
    for i, j := len(s)-1, 0; i >= 0; i, j = i-1, j+1 {
        reversed[j] = s[i]
    }
    return reversed
}

// Fungsi utama
func main() {
    // Baca masukan dari pengguna
    var algoritma string
    var judulAwal string
    var judulTujuan string

    fmt.Print("Algoritma (bfs/ids): ")
    fmt.Scanln(&algoritma)

    fmt.Print("Judul artikel awal: ")
    fmt.Scanln(&judulAwal)

    fmt.Print("Judul artikel tujuan: ")
    fmt.Scanln(&judulTujuan)

    // Unduh artikel awal dan tujuan
    awal, errAwal := getArtikel(judulAwal)
    if errAwal != nil {
        fmt.Println("Gagal mengunduh artikel awal:", errAwal)
        return
    }

    tujuan, errTujuan := getArtikel(judulTujuan)
    if errTujuan != nil {
        fmt.Println("Gagal mengunduh artikel tujuan:", errTujuan)
        return
    }

    // Jalankan pencarian berdasarkan algoritma yang dipilih
    var hasil []*Artikel
    var diperiksa int
    var waktu time.Duration

    if algoritma == "bfs" {
        hasil, diperiksa, waktu = bfs(awal, tujuan.Judul)
    } else if algoritma == "ids" {
        fmt.Print("Batas kedalaman (integer): ")
        var batasKedalaman int
        fmt.Scanln(&batasKedalaman)

        hasil, diperiksa, waktu = ids(awal, tujuan.Judul, batasKedalaman)
    } else {
        fmt.Println("Algoritma tidak valid (pilih bfs atau ids)")
        return
    }

    // Cetak hasil pencarian
    if hasil != nil {
        fmt.Println("\n**Hasil Pencarian:**")
        fmt.Printf("Jumlah artikel yang diperiksa: %d\n", diperiksa)
        fmt.Printf("Jumlah artikel yang dilalui: %d\n", len(hasil))
        fmt.Printf("Waktu pencarian: %v\n", waktu)

        fmt.Println("\n**Rute Penjelajahan:**")
        for _, artikel := range hasil {
            fmt.Println("-", artikel.Judul)
        }
    } else {
        fmt.Println("Gagal menemukan rute antara artikel awal dan tujuan.")
    }
}
