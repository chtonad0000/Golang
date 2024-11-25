//go:build !solution

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	files := os.Args[1:]
	stringMap := make(map[string]int)
	for _, file := range files {
		dat, err := os.ReadFile(file)
		check(err)
		filesStrings := strings.Split(string(dat), "\n")
		for _, fileString := range filesStrings {
			stringMap[fileString]++
		}
	}
	for k, v := range stringMap {
		if v >= 2 {
			fmt.Println(strconv.Itoa(v) + "\t" + k)
		}
	}
}
