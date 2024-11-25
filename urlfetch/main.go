//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	urls := os.Args[1:]
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			os.Exit(1)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Print(string(body))
	}
}
