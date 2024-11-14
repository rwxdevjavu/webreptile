package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Queue map[string]bool

func (q Queue) visitedURL(url string) bool {
	_, ok := q[url]
	return ok
}

func (q Queue) AddURL(url string) {
	if q.visitedURL(url) {
		q[url] = true
	} else {
		q[url] = false
	}
}

func Crawler(targetURL string, visitedURLQueue Queue) {
	// defer wg.Done()
	mux.Lock()
	base, err := url.Parse(targetURL)
	check(err)
	response, err := http.Get(targetURL)
	check(err)
	defer response.Body.Close()
	if response.StatusCode != 200 {
		fmt.Printf("response.StatusCode %d", response.StatusCode)
		os.Exit(1)
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	check(err)
	noExtraSpacesExp, err := regexp.Compile(`(?m)\s{2,}`)
	check(err)
	doc.Find(`meta[name="description"], meta[name="keywords"],meta[property="og:description"], title, h1, h2, h3, a`).Each(func(i int, s *goquery.Selection) {
		switch goquery.NodeName(s) {
		case "meta":
			if name, _ := s.Attr("name"); name == "description" {
				fmt.Println(s.Attr("content"))
			}
			if property, _ := s.Attr("property"); property == "og:description" {
				fmt.Println(s.Attr("content"))
			}

		case "a":
			href, ok := s.Attr("href")
			if ok {
				parsedURL, err := url.Parse(href)
				if err != nil {
					fmt.Println("err parsing url:", err)
				}
				absoluteURL := base.ResolveReference(parsedURL).String()
				visitedURLQueue.AddURL(absoluteURL)
			}
		case "h1", "h2", "h3":
			s.NextAllFiltered("p").First().Each(func(i int, p *goquery.Selection) {
				fmt.Println(goquery.NodeName(p), noExtraSpacesExp.ReplaceAllLiteralString(p.Text(), ""))
			})
		}
	})

	visitedURLQueue.AddURL(targetURL)
	fmt.Printf("%s Completed", targetURL)
	mux.Unlock()
}

// urlQueue := make(chan string, 100)
var mux sync.Mutex

var wg sync.WaitGroup

func main() {
	visitedURLQueue := Queue{}
	targetURL := `https://gobyexample.com`
	visitedURLQueue.AddURL(targetURL)
	visitedURLQueue.AddURL(`https://qikoffice.com`)
	for url, visited := range visitedURLQueue {
		if !visited {
			wg.Add(1)
			go Crawler(url, visitedURLQueue)
		}
	}
	wg.Wait()
	//fmt.Println(visitedURLQueue)
}
