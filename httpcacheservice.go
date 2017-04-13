package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/lox/httpcache"
)

func main() {
	configFileName := flag.String("config", "", "The config file path")
	flag.Parse()

	if *configFileName == "" {
		fmt.Printf("Error starting service: The --config argument is required\n")
		os.Exit(1)
	}

	cfg, err := LoadConfig(*configFileName)
	if err != nil {
		fmt.Printf("Error starting service: loading config: %v\n")
		os.Exit(1)
	}

	cache, err := httpcache.NewLRUCache(uint64(cfg.CacheSizeBytes))
	if err != nil {
		fmt.Printf("Error starting service: creating cache: %v\n")
		os.Exit(1)
	}
	// cache := httpcache.NewMemoryCache()

	handler := httpcache.NewHandler(cache, http.HandlerFunc(handle))
	handler.Shared = true

	listen := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("proxy listening on http://%s\n", listen)
	if err := http.ListenAndServe(listen, handler); err != nil {
		fmt.Printf("Error serving: %v\n", err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	target := "https://example.net" // debug
	uri := target + r.RequestURI
	fmt.Printf("%v %v %v\n", time.Now(), r.Method, uri)

	rr, err := http.NewRequest(r.Method, uri, r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	copyHeader(r.Header, &rr.Header)
	// Create a client and query the target
	var transport http.Transport
	resp, err := transport.RoundTrip(rr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	dH := w.Header()
	copyHeader(resp.Header, &dH)
	dH.Add("Requested-Host", rr.Host)

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}
