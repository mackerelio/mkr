package main

import (
	"fmt"
	"os"
)

func main() {
	for i := range 100 {
		fmt.Fprintln(os.Stdout, "out", i)
		fmt.Fprintln(os.Stderr, "err", i)
	}
}
