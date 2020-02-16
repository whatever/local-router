package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var version string = "?"

var Map map[string]int = map[string]int{
	"whatever.org":       9000,
	"nevermind.business": 9001,
	"anyway.net":         9002,
}

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

func main() {

	port := flag.Int("port", 8080, "listen on port")
	flag.Parse()

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
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}
