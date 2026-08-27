//colocamos o nome do package no topo para indicar que esse é um arquivo executavel

package main

//podemos importar as bibliotecas dessa forma
import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type Token struct{
	tipo string
	value string
}

type Lexer struct{
	source string
	position int
	next Token
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

		// Se o caractere é igual a '^'
		if char == '^' {
			le.next = Token{tipo: "XOR", value: "^"}
			le.position++
			return
		}
		

		// Se o caractere é igual a '-'
		if char == '-' {
			le.next = Token{tipo: "MINUS", value: "-"}
			le.position++
			return
		}

		// se for digito
		if unicode.IsDigit(rune(char)){

			startPos := le.position
			for le.position < len(le.source) && unicode.IsDigit(rune(le.source[le.position])){
				le.position++
			}

			numStr := le.source[startPos:le.position]

			le.next = Token{tipo: "INT", value: numStr}

			return

		}

		panic("[Lexer] digito nn valido")

}

type Parser struct{
	lexer Lexer
}

func ParseExpression(lexer *Lexer) int {
	if lexer.next.tipo != "INT" {
		panic("[PARSER] esperava um número")
	}

	num, err := strconv.Atoi(lexer.next.value)
	if err != nil {
		panic("[PARSER] número inválido")
	}

	resultado := num
	lexer.selectNext()

	for lexer.next.tipo == "PLUS" || lexer.next.tipo == "MINUS" || lexer.next.tipo == "XOR"{
		operador := lexer.next.tipo
		lexer.selectNext()

		if lexer.next.tipo != "INT" {
			panic("[PARSER] esperava um número depois do operador")
		}

		numero, err := strconv.Atoi(lexer.next.value)
		if err != nil {
			panic("[PARSER] número inválido")
		}

		if operador == "PLUS" {
			resultado += numero
		}
		if operador == "MINUS" {
			resultado -= numero
		}
		if operador == "XOR"{
			resultado = resultado ^ numero
		}

		lexer.selectNext()
	}

	// se sobrou da erro
	if lexer.next.tipo != "EOF" {
		panic("[PARSER] token inesperado: " + lexer.next.value)
	}

	return resultado
}


func(p *Parser) run() int{
	// add o primeiro next
	p.lexer.selectNext()

	// le a expressao
	resultado := ParseExpression(&p.lexer)


	return resultado
}




func main() {

	source := os.Args[1]

	l := &Lexer{
		source: source,
		position: 0,
	}

	p := &Parser{
		lexer: *l,
	}


	fmt.Println(p.run())
}
