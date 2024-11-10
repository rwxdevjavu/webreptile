package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

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

func main() {
	response, err := http.Get("https://gobyexample.com")
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
				fmt.Println(href)
			}
		case "h1", "h2", "h3":
			s.NextAllFiltered("p").First().Each(func(i int, p *goquery.Selection) {
				fmt.Println(goquery.NodeName(p), noExtraSpacesExp.ReplaceAllLiteralString(p.Text(), ""))
			})
		}
	})
}
