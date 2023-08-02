package main

import (
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

/*
Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func download(uri string, queue chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(uri)
	if err != nil {
		log.Printf("failed to get %s: %v", uri, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read body of %s: %v", uri, err)
		return
	}

	parsedURL, _ := url.Parse(uri)
	if strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path = parsedURL.Path[:len(parsedURL.Path)-1]
	}

	fileName := urlToFileName(uri)
	err = ioutil.WriteFile(fileName, body, 0666)
	if err != nil {
		log.Printf("failed to write to file %s: %v", fileName, err)
		return
	}

	newURLs := parseHTML(body)
	for _, u := range newURLs {
		absoluteURL, err := parsedURL.Parse(u)
		if err != nil {
			log.Printf("failed to parse url %s: %v", u, err)
			continue
		}

		if parsedURL.Hostname() == absoluteURL.Hostname() {
			wg.Add(1)
			queue <- absoluteURL.String()
		}
	}
}

func urlToFileName(url string) string {
	fileName := strings.Replace(url, "http://", "", 1)
	fileName = strings.Replace(fileName, "https://", "", 1)
	fileName = strings.Replace(fileName, "/", "-", -1)
	return fileName
}

func parseHTML(body []byte) []string {
	doc, _ := html.Parse(strings.NewReader(string(body)))
	var urls []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					urls = append(urls, a.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return urls
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("provide a starting url")
	}
	startURL := args[0]
	queue := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go download(startURL, queue, &wg)
	go func() {
		wg.Wait()
		close(queue)
	}()

	for uri := range queue {
		if _, err := os.Stat(urlToFileName(uri)); os.IsNotExist(err) {
			wg.Add(1)
			go download(uri, queue, &wg)
		}
	}
}
