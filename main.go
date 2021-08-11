package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	q := parseQuery()
	err := readStdin(q)
	if err != nil {
		log.Fatal(err)
	}
}

func readStdin(q query) error {
	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)

	pars := newParser()

	for {
		n, err := r.Read(buf[:cap(buf)])
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		buf = buf[:n]

		pars.process(buf)
	}

	return nil
}
