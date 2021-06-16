package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var tokens = []string{}
var openedBrackets bool
var startToken int
var interruptedToken string

func main() {

	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)

	for {
		n, err := r.Read(buf[:cap(buf)])
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		buf = buf[:n]

		processChunk(string(buf))
	}
}

func processChunk(chunk string) {

	var printString int
	ln := len(chunk)

	for i, rn := range chunk {
		if openedBrackets {
			if rn == '>' {
				openedBrackets = false
				if ln-1 == i {
					printString = -1
				} else if ln-1 > i {
					printString = i + 1
				}

				var token string
				if len(interruptedToken) > 0 {
					token = interruptedToken + chunk[:i+1]
					interruptedToken = ""
				} else {
					token = chunk[startToken : i+1]
				}

				if token[1] == '/' {
					fmt.Printf("\n%s%s", strings.Repeat("  ", len(tokens)-1), token)
					tokens = tokens[:len(tokens)-1]
				} else if token[1] == '!' {
					fmt.Printf("%s", token)
				} else {
					tokens = append(tokens, token)
					fmt.Printf("\n%s%s", strings.Repeat("  ", len(tokens)-1), token)
				}

			}
			continue
		}

		if rn == '<' {
			startToken = i
			openedBrackets = true

			if printString >= 0 {
				fmt.Printf("%s", chunk[printString:i])
			}
		}
	}

	if openedBrackets {
		interruptedToken = chunk[startToken:ln]
	} else if printString >= 0 {
		fmt.Printf("%s", chunk[printString:ln])
	}

	startToken = 0
}
