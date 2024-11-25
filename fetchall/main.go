//go:build !solution

package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func logURL(url string, ch chan<- string) {
	t := time.Now()
	_, err := http.Get(url)
	if err != nil {
		ch <- "wrong url: " + url
		return
	}

	t2 := time.Since(t)

	ch <- strconv.FormatFloat(t2.Seconds(), 'f', 2, 64) + " " + url + "\n"
}

func main() {
	urls := os.Args[1:]
	messages := make(chan string)
	for _, url := range urls {
		go logURL(url, messages)
	}
	for range urls {
		msg := <-messages
		fmt.Println(msg)
	}
}
