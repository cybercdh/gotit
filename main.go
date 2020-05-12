package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gookit/color"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "set the concurrency level")

	var ignoreblanks bool
	flag.BoolVar(&ignoreblanks, "b", false, "ignore sites which don't return a title")

	flag.Parse()

	urls := make(chan string)

	var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   time.Second,
			KeepAlive: time.Second,
		}).DialContext,
	}

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	client := &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       time.Second,
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			for url := range urls {

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					continue
				}

				req.Header.Add("Connection", "close")
				req.Close = true

				resp, err := client.Do(req)

				if err != nil {
					continue
				}
				defer resp.Body.Close()

				doc, err := goquery.NewDocumentFromReader(resp.Body)
				if err != nil {
					continue
				}

				title := doc.Find("title").Text()
				if title != "" {
					if strings.Contains(title, "403") {
						color.Yellow.Printf("%s : %s\n", url, title)
					} else if strings.Contains(title, "30") {
						color.Cyan.Printf("%s : %s\n", url, title)
					} else {
						fmt.Printf("%s : %s\n", url, title)
					}
				} else {
					if !ignoreblanks {
						fmt.Printf("%s : No title\n", url)
					}
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
