package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func main() {
	args := os.Args[1:]
	wg := sync.WaitGroup{}
	for _, s := range args {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			res, err := http.Get(s)
			if err != nil {
				fmt.Printf("Error %d while download page", err)
				return
			}
			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("Error %d while parsing body", err)
			}
			fmt.Printf("Download page pass succesed, lenght of page %d", len(string(body)))
			res.Body.Close()
		}(s)
	}
	wg.Wait()
}
