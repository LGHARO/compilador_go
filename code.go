package main

import (
	"fmt"
	"os"
	"strings"
)

// Code gerencia o armazenamento e gravação do código assembly
type Code struct{}

// codeInstructions armazena as instruções geradas sequencialmente no _start
var codeInstructions []string

// funcInstructions armazena os procedimentos de funções gerados fora de _start
var funcInstructions []string

// globalVars armazena os rótulos de variáveis globais declaradas na seção .data
var globalVars []string

// isGeneratingFunc indica se as instruções emitidas são de procedimento de função
var isGeneratingFunc bool

// hasProgramFunctions indica se o programa compilado contém funções
var hasProgramFunctions bool

// Append adiciona uma instrução ao buffer de código
func (c *Code) Append(code string) {
	CodeAppend(code)
}

// append versão lowercase para compatibilidade com a especificação
func (c *Code) append(code string) {
	CodeAppend(code)
}

// Dump escreve o código final no arquivo especificado
func (c *Code) Dump(filename string) {
	CodeDump(filename)
}

// dump versão lowercase para compatibilidade com a especificação
func (c *Code) dump(filename string) {
	CodeDump(filename)
}

// CodeAppend adiciona sequencialmente uma instrução assembly ao buffer ativo
func CodeAppend(instruction string) {
	if isGeneratingFunc {
		funcInstructions = append(funcInstructions, instruction)
	} else {
		codeInstructions = append(codeInstructions, instruction)
	}
}

// CodeAddGlobalVar registra uma variável global para a seção .data
func CodeAddGlobalVar(name string) {
	label := "_var_" + name
	for _, v := range globalVars {
		if v == label {
			return
		}
	}
	globalVars = append(globalVars, label)
}

// CodeReset limpa a lista de instruções (útil entre compilações e testes)
func CodeReset() {
	codeInstructions = []string{}
	funcInstructions = []string{}
	globalVars = []string{}
	isGeneratingFunc = false
	hasProgramFunctions = false
}

// CodeDump grava o código Assembly NASM completo no arquivo indicado
func CodeDump(filename string) {
	var sb strings.Builder

	// Cabeçalho da seção .data
	sb.WriteString("section .data\n")
	sb.WriteString("format_out: db \"%d\", 10, 0 ; format do printf\n")
	sb.WriteString("format_in: db \"%d\", 0 ; format do scanf\n")
	sb.WriteString("scan_int: dd 0; 32-bits integer\n")
	for _, gv := range globalVars {
		sb.WriteString(fmt.Sprintf("%s: dd 0\n", gv))
	}
	sb.WriteString("\n")

	// Seção .text
	sb.WriteString("section .text\n")
	sb.WriteString("extern printf\n")
	sb.WriteString("extern scanf\n")
	sb.WriteString("global _start ; início do programa\n\n")

	// Ponto de entrada _start
	sb.WriteString("_start:\n")
	sb.WriteString("push ebp ; guarda o EBP\n")
	sb.WriteString("mov ebp, esp ; zera a pilha\n")
	sb.WriteString("; aqui começa o codigo gerado:\n")
	for _, inst := range codeInstructions {
		sb.WriteString(inst)
		sb.WriteString("\n")
	}
	sb.WriteString("; aqui termina o código gerado\n")
	sb.WriteString("mov esp, ebp ; reestabelece a pilha\n")
	sb.WriteString("pop ebp ; chamada da interrupcao de saida\n")
	sb.WriteString("mov eax, 1\n")
	sb.WriteString("xor ebx, ebx\n")
	sb.WriteString("int 0x80\n\n")

	// Funções / procedimentos gerados fora de _start
	for _, inst := range funcInstructions {
		sb.WriteString(inst)
		sb.WriteString("\n")
	}

	err := os.WriteFile(filename, []byte(sb.String()), 0644)
	if err != nil {
		panic(err)
	}
}

