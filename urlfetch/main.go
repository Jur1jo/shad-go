//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	for _, s := range args {
		res, err := http.Get(s)
		if err != nil {
			os.Exit(1)
		}
		body, _ := io.ReadAll(res.Body)
		fmt.Println(string(body))
		res.Body.Close()
	}
}
