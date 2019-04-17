package main

import (
	"fmt"
	"os"
)

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println("Hello world!")
	}
	os.Exit(1)
}
