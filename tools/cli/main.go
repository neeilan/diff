package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		if !sc.Scan() {
			return
		}

		io.WriteString(out, sc.Text())
		io.WriteString(out, "\n")
	}
}

func main() {
	Start(os.Stdin, os.Stdout)
}
