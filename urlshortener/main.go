//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

var (
	mu      sync.RWMutex
	urls    = make(map[string]string)
	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

func keyGeneration() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[randGen.Intn(len(letters))]
	}
	return string(b)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req shortenRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	mu.RLock()
	for key, storedURL := range urls {
		if storedURL == req.URL {
			mu.RUnlock()
			resp := shortenResponse{
				URL: req.URL,
				Key: key,
			}
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				return
			}
			return
		}
	}
	mu.RUnlock()
	key := keyGeneration()
	mu.Lock()
	urls[key] = req.URL
	mu.Unlock()

	resp := shortenResponse{
		URL: req.URL,
		Key: key,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/go/")
	mu.RLock()
	url, exists := urls[key]
	mu.RUnlock()

	if !exists {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	port := flag.String("port", "6029", "port")
	flag.Parse()
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/go/", handleRedirect)

	addr := fmt.Sprintf(":%s", *port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return
	}
}
