package main

import (
	tea "charm.land/bubbletea/v2"
	"flag"
	"fmt"
	"os"
)

func main() {

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: please provide a file path")
		os.Exit(1)
	}

	fileName := args[0]

	tools, err := parse(fileName)
	if err != nil {
		fmt.Println("Error occured while reading from file:", fileName)
		os.Exit(1)
	}

	m := &model{
		tools:  tools,
		styles: NewTheme(),
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error launching CLI: %v", err)
		os.Exit(1)
	}

}
