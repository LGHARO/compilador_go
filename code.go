package main

import (
	"os"
	"strings"
)

// ============================================================================
// Roteiro 8 - Geração de Código x86 (NASM)
// ============================================================================
// Mudanças realizadas:
// 1. Criação da classe/struct Code para gerenciar as instruções Assembly geradas.
// 2. Implementação dos métodos append() e dump() conforme especificação:
//    - append: adiciona sequencialmente instruções Assembly na lista.
//    - dump: escreve o código final no arquivo de saída, adicionando o cabeçalho
//            mínimo (.data, .text, extern printf, extern scanf, global _start,
//            push/mov ebp) e o rodapé mínimo (restaura pilha e int 0x80 / sys_exit).
// ============================================================================

// Code gerencia o armazenamento e gravação do código assembly
type Code struct{}

// codeInstructions armazena as instruções geradas sequencialmente
var codeInstructions []string

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

// CodeAppend adiciona sequencialmente uma instrução assembly ao arquivo final
func CodeAppend(instruction string) {
	codeInstructions = append(codeInstructions, instruction)
}

// CodeReset limpa a lista de instruções (útil entre compilações e testes)
func CodeReset() {
	codeInstructions = []string{}
}

// CodeDump grava o código Assembly NASM completo no arquivo indicado
func CodeDump(filename string) {
	header := `section .data
format_out: db "%d", 10, 0 ; format do printf
format_in: db "%d", 0 ; format do scanf
scan_int: dd 0; 32-bits integer

section .text
extern printf
extern scanf
global _start ; início do programa

_start:
push ebp ; guarda o EBP
mov ebp, esp ; zera a pilha
; aqui começa o codigo gerado:
`

	footer := `; aqui termina o código gerado
mov esp, ebp ; reestabelece a pilha
pop ebp ; chamada da interrupcao de saida
mov eax, 1
xor ebx, ebx
int 0x80
`

	var sb strings.Builder
	sb.WriteString(header)
	for _, inst := range codeInstructions {
		sb.WriteString(inst)
		sb.WriteString("\n")
	}
	sb.WriteString(footer)

	err := os.WriteFile(filename, []byte(sb.String()), 0644)
	if err != nil {
		panic(err)
	}
}
