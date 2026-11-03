package main

import (
	"os"
)

// ============================================================================
// Versão 2.3 - Funções e Escopo de Variáveis
// ============================================================================
// Mudanças realizadas:
// 1. O pipeline lê o arquivo de entrada, aplica o pré-processamento de comentários,
//    executa a análise léxica e constrói a AST através do parser.
// 2. A raiz da AST passa a conter declarações de variáveis globais, declarações
//    de funções (FuncDec) e a chamada automática para a função main() (FuncCall).
// 3. Inicialização da SymbolTable global (parent = nil), onde variáveis globais
//    e funções são registradas durante a avaliação (Evaluate).
// ============================================================================

func main() {
	if len(os.Args) < 2 {
		panic("Uso: compilador_go <arquivo_fonte>")
	}

	caminhoArquivo := os.Args[1]

	dados, err := os.ReadFile(caminhoArquivo)
	if err != nil {
		panic(err)
	}

	filter := &Prepro{}
	dadosFiltrados := filter.filter(string(dados))

	l := &Lexer{
		source:   dadosFiltrados,
		position: 0,
	}

	p := &Parser{
		lexer: *l,
	}

	// Montamos a AST com a chamada para main() incluída
	AST := p.run()

	st := &SymbolTable{
		table:  make(map[string]Variable),
		parent: nil,
	}
	AST.Evaluate(*st)
}