package main

import (
	"os"
	"path/filepath"
	"strings"
)

// ============================================================================
// Roteiro 8 - Geração de Código x86 (NASM)
// ============================================================================
// Mudanças realizadas:
// 1. O compilador agora gera código Assembly NASM chamando AST.Generate(st)
//    em vez de executar o programa pelo Evaluate().
// 2. O arquivo de saída possui o mesmo prefixo do arquivo de entrada e extensão
//    .asm (exemplo: teste.go -> teste.asm), conforme exigido pela especificação.
// 3. Chamada a CodeDump(asmFileName) para persistir o código Assembly com
//    cabeçalho e rodapé mínimos para ambiente Linux x86 (32 bits).
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

	// Reinicia geradores para compilação limpa
	resetNodeId()
	CodeReset()

	// Montamos a AST
	AST := p.run()

	st := &SymbolTable{
		table:        make(map[string]Variable),
		currentShift: 0,
	}

	// Geração de código Assembly NASM
	AST.Generate(st)

	// Nome do arquivo de saída com extensão .asm e mesmo prefixo
	ext := filepath.Ext(caminhoArquivo)
	asmFileName := strings.TrimSuffix(caminhoArquivo, ext) + ".asm"

	// Grava o código gerado no arquivo .asm
	CodeDump(asmFileName)
}