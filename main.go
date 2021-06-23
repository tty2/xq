package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	err := readStdin()
	if err != nil {
		log.Fatal(err)
	}
}

func readStdin() error {
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

	fmt.Println("")

	return nil
}
