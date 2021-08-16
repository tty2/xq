package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	q := parseQuery()

	p := newParser(q)

	r := bufio.NewReader(os.Stdin)

	err := p.process(r)
	if err != nil {
		log.Fatal(err)
	}
}
