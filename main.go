package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	q := getQuery()
	q.parse()
	proc := getProcessor(q)

	r := bufio.NewReader(os.Stdin)

	err := proc.Process(r)
	if err != nil {
		log.Fatal(err)
	}
}
