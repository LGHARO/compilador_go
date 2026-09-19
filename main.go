package main

import (
	"os"
)

func main() {

	caminho_arquivo := os.Args[1]

	dados,_ := os.ReadFile(caminho_arquivo)
	filter := &Prepro{

	}
	dados_filtrados := filter.filter(string(dados))

	l := &Lexer{
		source:   dados_filtrados,
		position: 0,
	}

	p := &Parser{
		lexer: *l,
	}

	//montamos a ast
	AST := p.run()

	st := &SymbolTable{
		table: make(map[string]Variable),
	}
	AST.Evaluate(*st)
}