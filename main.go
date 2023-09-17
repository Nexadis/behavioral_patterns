package main

import (
	"fmt"
	"log"
	"strings"
)

func NewIter(max int) func() (int, bool) {
	n := 0
	return func() (int, bool) {
		if n >= max {
			return 0, false
		}
		n++
		return n - 1, true
	}
}

func DemoIterator() {
	iter := NewIter(5)
	for {
		n, ok := iter()
		if !ok {
			break
		}
		log.Println(n)
	}
}

func main() {
	border("Iterator")
	DemoIterator()
}

func border(name string) {
	line := strings.Repeat("=", 80)
	out := fmt.Sprintf("%s\n\t\t\t%s\n%s\n", line, name, line)
	fmt.Println(out)
}
