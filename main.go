//colocamos o nome do package no topo para indicar que esse é um arquivo executavel

package main

//podemos importar as bibliotecas dessa forma
import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type Token struct {
	tipo  string
	value string
}

type Lexer struct {
	source   string
	position int
	next     Token
}

type Parser struct {
	lexer Lexer
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

	panic("[Lexer] digito nn valido")

}

func ParseFactor(lexer *Lexer) int {
	// 1. Se for um número, consome e retorna o valor
	if lexer.next.tipo == "INT" {
		num, err := strconv.Atoi(lexer.next.value)
		if err != nil {
			panic("[PARSER] número inválido")
		}
		lexer.selectNext() // Consome o INT
		return num
	}

	if lexer.next.tipo == "MINUS" ||  lexer.next.tipo == "PLUS"{
		if  lexer.next.tipo == "MINUS" {
			lexer.selectNext()
			return - ParseFactor(lexer)
		}else{
			lexer.selectNext()
			return + ParseFactor(lexer)
		}
		
	} 

	// 2. Se for um parênteses, resolve a expressão interna
	if lexer.next.tipo == "OPEN_PAR" {
		lexer.selectNext() // Consome o '('

		// Chama a expressão principal novamente para resolver o miolo
		resultadoInterno := ParseExpression(lexer)

		if lexer.next.tipo != "CLOSE_PAR" {
			panic("[PARSER] faltou fechar o parênteses")
		}
		lexer.selectNext() // Consome o ')'

		return resultadoInterno
	}

	panic("[PARSER] token inesperado, esperava número ou '(': " + lexer.next.value)
}

func ParseTerm(lexer *Lexer) int {

	resultado := ParseFactor(lexer)

	for lexer.next.tipo == "MULT" || lexer.next.tipo == "DIV" {
		operador := lexer.next.tipo
		lexer.selectNext()

		// puxa o factor
		proximoFator := ParseFactor(lexer)

		if operador == "MULT" {
			resultado *= proximoFator
		} else {
			resultado /= proximoFator
		}
	}

	return resultado
}



func ParseExpression(lexer *Lexer) int {

	// calcula o resultado no term
	resultado := ParseTerm(lexer)

	for lexer.next.tipo == "PLUS" || lexer.next.tipo == "MINUS" {
		operador := lexer.next.tipo
		lexer.selectNext()

		proximo := ParseTerm(lexer)

		if operador == "PLUS" {
			resultado += proximo
		} else {
			resultado -= proximo
		}

	}

	return resultado
}

func (p *Parser) run() int {
	// add o primeiro next
	p.lexer.selectNext()

	resultado := ParseExpression(&p.lexer)

	if p.lexer.next.tipo != "EOF" {
		panic("[PARSER] token inesperado: " + p.lexer.next.value)
	}

	return resultado
}

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
