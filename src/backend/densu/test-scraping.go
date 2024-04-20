package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	startURL := "https://en.wikipedia.org/wiki/Contract_bridge"

	fmt.Println("Scraping contents of the start URL:", startURL)
	startLinks := getLinks(startURL)

	fmt.Println("Links found on the start URL:")
	for _, link := range startLinks {
		fmt.Println(link)
	}

	err := writeToFile("scraped_links.txt", startLinks)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Scraped links saved to scraped_links.txt")
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

func writeToFile(filename string, links []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, link := range links {
		_, err := file.WriteString(link + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
