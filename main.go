package main

import (
	"fmt"
	"os"
)

func main() {

	source := os.Args[1]

	l := &Lexer{
		source:   source,
		position: 0,
	}

	p := &Parser{
		lexer: *l,
	}

	fmt.Println(p.run())
}