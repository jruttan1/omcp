package main

import (
	"fmt"
	"os"
)

func main() {

	cmd := os.Args

	if len(cmd) < 2 {
		fmt.Println("No arguments passed")
		os.Exit(1)
	}

	fileName := os.Args[1]

	tools, err := parse(fileName)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println(tools)

}
