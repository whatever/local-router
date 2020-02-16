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
)

var version string = "?"

var Map map[string]int = map[string]int{}

// PortMap returns the map
func PortMap(url string) (port int) {
	var ok bool
	if port, ok = Map[url]; !ok {
		port = 9999
	}
	return
}

// ProxyRequest proxies a get request
func ProxyRequest(r *http.Request) (resp *http.Response, err error) {

	port := PortMap(r.Host)

	context := r.Context()
	r2 := r.Clone(context)
	r2.URL.Scheme = "http"
	r2.URL.Host = fmt.Sprintf("0.0.0.0:%v", port)
	resp, err = http.DefaultTransport.RoundTrip(r2)

	return resp, err
}

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
			fmt.Fprintf(w, string(body))
		}
	})

	fmt.Printf("Router v%v\n", version)
	fmt.Printf("Port ....... %v\n", *port)
	for host, port := range Map {
		fmt.Printf("%v -> %v\n", host, port)
	}
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}
