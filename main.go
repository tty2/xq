package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	q := getQuery()
	q.parse()
	proc, err := getProcessor(q)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(os.Stdin)

	c := proc.Process(r)

	for {
		line, ok := <-c
		if !ok {
			break
		}

		fmt.Println(line) // nolint forbidigo: print is executed on purpose here
	}
}
