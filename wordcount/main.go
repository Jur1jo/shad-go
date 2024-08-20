package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	counter := make(map[string]int)
	for _, s := range args[:] {
		f, err := os.Open(fmt.Sprintf("%s", s))
		if err != nil {
			panic("Can't open the file")
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			counter[line]++
		}
	}
	for key, value := range counter {
		if value > 1 {
			fmt.Printf("%v\t%v\n", value, key)
		}
	}
}
