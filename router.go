package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var version string = "?"

var (
	MapLock sync.Mutex
	Map     map[string]int = map[string]int{}
)

// PortMap returns the map
func PortMap(url string) (port int) {

	// Maps aren't concurrency safe in go, but I don't they get concurrent read errors,
	// If they're n ot updated
	// MapLock.Lock()
	// defer MapLock.Unlock()

	var ok bool
	if port, ok = Map[url]; !ok {
		port = 9999
	}
	return
}

// ProxyRequest proxies a get request
func ProxyRequest(r *http.Request) (resp *http.Response, err error) {
	port := PortMap(r.Host)

	// Clones the request, but shaves the host to something local:PORT
	context := r.Context()
	r2 := r.Clone(context)
	r2.URL.Scheme = "http"
	r2.URL.Host = fmt.Sprintf("127.0.0.1:%v", port)
	fmt.Printf("Proxy %v%v to \"%v\"\n", r.Host, r.URL, port)

	return http.DefaultTransport.RoundTrip(r2)
}

// This just parses lines here
func LoadConfig(file_name string) {
	fi, err := os.Open(file_name)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(2)
		return
	}

	scanner := bufio.NewScanner(fi)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			fmt.Println("Skipping empty line")
			continue
		}

		pieces := strings.Split(line, ":")

		if len(pieces) != 2 {
			fmt.Println("Skipping invalid line:", line)
			continue
		}

		pieces[0] = strings.TrimSpace(pieces[0])
		pieces[1] = strings.TrimSpace(pieces[1])

		if pieces[0] == "" {
			fmt.Println("Skipping empty hostname")
			continue
		}

		if n, e := strconv.Atoi(pieces[1]); e == nil {
			Map[pieces[0]] = n
		}
	}
}

func main() {

	port := flag.Int("port", 8080, "listen on port")
	config := flag.String("config", "", "")
	flag.Parse()

	// This mutates the map!
	LoadConfig(*config)

	// Single proxy handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if resp, err := ProxyRequest(r); err != nil {
			fmt.Println(err)
		} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
			fmt.Println(err)
		} else {
			if content_type := resp.Header.Get("Content-Type"); content_type != "" {
				w.Header().Add("Content-Type", content_type)
			}
			if n, e := w.Write(body); e != nil {
				fmt.Println(e)
			} else {
				fmt.Printf("Received %d bytes, Sent %d bytes\n", len(body), n)
			}
		}
	})

	fmt.Printf("Router v%v\n", version)
	fmt.Printf("Port ....... %v\n", *port)
	for host, port := range Map {
		fmt.Printf("%v -> %v\n", host, port)
	}
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}
