package main

import (
	"unicode"
)

type Lexer struct {
	source   string
	position int
	next     Token
}

func (le *Lexer) selectNext() {

	// percorre ate nn ser mais espaco em branco ou acabar as informacoes
	for le.position < len(le.source) && unicode.IsSpace(rune(le.source[le.position])) {
		le.position++
	}

	// se for fim adiciona o token
	if le.position >= len(le.source) {
		le.next = Token{tipo: "EOF", value: ""}
		return
	}

	// Caractere atual
	char := le.source[le.position]

	// Se o caractere é igual a '+'
	if char == '+' {
		le.next = Token{tipo: "PLUS", value: "+"}
		le.position++
		return
	}

	// Se o caractere é igual a '-'
	if char == '-' {
		le.next = Token{tipo: "MINUS", value: "-"}
		le.position++
		return
	}

	// Se o caractere é igual a '*'
	if char == '*' {
		le.next = Token{tipo: "MULT", value: "*"}
		le.position++
		return
	}

	// Se o caractere é igual a '/'
	if char == '/' {
		le.next = Token{tipo: "DIV", value: "/"}
		le.position++
		return
	}

	// Se o caractere é igual a '('
	if char == '(' {
		le.next = Token{tipo: "OPEN_PAR", value: "("}
		le.position++
		return
	}

	// Se o caractere é igual a ')'
	if char == ')' {
		le.next = Token{tipo: "CLOSE_PAR", value: ")"}
		le.position++
		return
	}

	// Se o caractere é igual a '='
	if char == '=' {
		le.next = Token{tipo: "ASSIGN", value: "="}
		le.position++
		return
	}

	// Se o caractere é igual a '\n'
	if char == '\n' {
		le.next = Token{tipo: "END", value: "\n"}
		le.position++
		return
	}

	

	// se for digito
	if unicode.IsDigit(rune(char)) {

		startPos := le.position
		for le.position < len(le.source) && unicode.IsDigit(rune(le.source[le.position])) {
			le.position++
		}

		numStr := le.source[startPos:le.position]

		le.next = Token{tipo: "INT", value: numStr}

		return
	}

	if unicode.IsLetter(rune(char)) {
		
		startPos := le.position

		for le.position < len(le.source) && (unicode.IsDigit(rune(le.source[le.position])) || unicode.IsLetter(rune(le.source[le.position])) || le.source[le.position] == '_') {
			le.position++
		}

		idenStr := le.source[startPos:le.position]

		if idenStr == "Println" {
			le.next = Token{tipo: "PRINT", value: idenStr}
		}else{
			le.next = Token{tipo: "IDEN", value: idenStr}
		}
		return
	}

	panic("[Lexer] digito nn valido")
}