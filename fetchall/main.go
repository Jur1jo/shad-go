//package main
//
//import (
//	"fmt"
//	"io"
//	"net/http"
//	"os"
//)
//
//func fetchPage(url string)
//
//func main() {
//	args := os.Args[1:]
//	for _, s := range args {
//		go func() {
//			res, err := http.Get(s)
//			if err != nil {
//				fmt.Printf("Error %d while download page", err)
//				return
//			}
//			body, err := io.ReadAll(res.Body)
//			if err != nil {
//				fmt.Printf("Error %d while parsing body", err)
//			}
//			fmt.Printf("Download page pass succesed, lenght of page %d", len(string(body)))
//			res.Body.Close()
//		}()
//	}
//}

package main

func main() {
}
