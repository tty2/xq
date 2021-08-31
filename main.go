package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	q := getQuery()
	q.parse()

	p := newParser(q)

	r := bufio.NewReader(os.Stdin)

	err := p.process(r)
	if err != nil {
		log.Fatal(err)
	}
}
