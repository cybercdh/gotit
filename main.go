package main

import (
	"golang.org/x/net/html"
	"io"
	"net/http"
	"fmt"
	"flag"
	"strings"
	"os"
	"bufio"
	"sync"
)

var wg sync.WaitGroup

func main() {

	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "set the concurrency level")

	flag.Parse()

	urls := make(chan string)

	for i := 0; i < concurrency; i++ {
	    wg.Add(1)

	    go func() {
	      for url := range urls {

	        resp, err := http.Get(url)
	        if err != nil {
	        	continue
	        }
	        defer resp.Body.Close()

	        if title, ok := GetHtmlTitle(resp.Body); ok {
	        	fmt.Printf("%s : %s\n", url, title)
	        } else {
	        	fmt.Printf("%s : No title\n", url)
	        }
	      }
	      wg.Done()
	    }()
	}

	var input_urls io.Reader
	input_urls = os.Stdin

	arg_url := flag.Arg(0)
	if arg_url != "" {
	    input_urls = strings.NewReader(arg_url)
	}

	sc := bufio.NewScanner(input_urls)
	
	for sc.Scan() {
			url := sc.Text()
			if strings.HasPrefix(url, "http") {
				urls <- url
			} else {
				urls <- "http://" + url
				urls <- "https://" + url
			}
	    
	}

	close(urls)
	wg.Wait()

}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func GetHtmlTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		panic("Failed to parse html")
	}

	return traverse(doc)
}