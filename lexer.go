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
	// Pula espaços e tabs, mas NÃO pula '\n' (que é um token END significativo)
	for le.position < len(le.source) && (le.source[le.position] == ' ' || le.source[le.position] == '\t' || le.source[le.position] == '\r') {
		le.position++
	}

	if le.position >= len(le.source) {
		le.next = Token{tipo: "EOF", value: ""}
		return
	}

	char := le.source[le.position]

	if char == '+' {
		le.next = Token{tipo: "PLUS", value: "+"}
		le.position++
		return
	}

	if char == '-' {
		le.next = Token{tipo: "MINUS", value: "-"}
		le.position++
		return
	}

	if char == '*' {
		le.next = Token{tipo: "MULT", value: "*"}
		le.position++
		return
	}

	if char == '/' {
		le.next = Token{tipo: "DIV", value: "/"}
		le.position++
		return
	}

	if char == '(' {
		le.next = Token{tipo: "OPEN_PAR", value: "("}
		le.position++
		return
	}

	if char == ')' {
		le.next = Token{tipo: "CLOSE_PAR", value: ")"}
		le.position++
		return
	}

	if char == '{' {
		le.next = Token{tipo: "OPEN_BRA", value: "{"}
		le.position++
		return
	}

	if char == '}' {
		le.next = Token{tipo: "CLOSE_BRA", value: "}"}
		le.position++
		return
	}

	if char == '&' {
		if le.position+1 < len(le.source) && le.source[le.position+1] == '&' {
			le.next = Token{tipo: "AND", value: "&&"}
			le.position += 2
			return
		}
		panic("[Lexer] caractere '&' inesperado")
	}

	if char == '|' {
		if le.position+1 < len(le.source) && le.source[le.position+1] == '|' {
			le.next = Token{tipo: "OR", value: "||"}
			le.position += 2
			return
		}
		panic("[Lexer] caractere '|' inesperado")
	}

	if char == '!' {
		le.next = Token{tipo: "NOT", value: "!"}
		le.position++
		return
	}

	if char == '=' {
		if le.position+1 < len(le.source) && le.source[le.position+1] == '=' {
			le.next = Token{tipo: "EQ", value: "=="}
			le.position += 2
			return
		}
		le.next = Token{tipo: "ASSIGN", value: "="}
		le.position++
		return
	}

	if char == '>' {
		le.next = Token{tipo: "GT", value: ">"}
		le.position++
		return
	}

	if char == '<' {
		le.next = Token{tipo: "LT", value: "<"}
		le.position++
		return
	}

	if char == '\n' {
		le.next = Token{tipo: "END", value: "\n"}
		le.position++
		return
	}

	// Literal string: tudo entre aspas duplas
	if char == '"' {
		le.position++ // pula a aspas de abertura
		startPos := le.position
		for le.position < len(le.source) && le.source[le.position] != '"' && le.source[le.position] != '\n' {
			le.position++
		}
		if le.position >= len(le.source) || le.source[le.position] != '"' {
			panic("[Lexer] string não fechada")
		}
		strContent := le.source[startPos:le.position]
		le.position++ // pula a aspas de fechamento
		le.next = Token{tipo: "STR", value: strContent}
		return
	}

	if unicode.IsDigit(rune(char)) || (char == '.' && le.position+1 < len(le.source) && unicode.IsDigit(rune(le.source[le.position+1]))) {
		startPos := le.position
		dotCount := 0

		for le.position < len(le.source) {
			c := le.source[le.position]
			if unicode.IsDigit(rune(c)) {
				le.position++
			} else if c == '.' {
				dotCount++
				le.position++
			} else if unicode.IsLetter(rune(c)) {
				panic("[Lexer] formato de número inválido: letra em número")
			} else {
				break
			}
		}

		if dotCount > 1 {
			panic("[Lexer] formato de número inválido: múltiplos pontos decimais")
		}

		numStr := le.source[startPos:le.position]
		if dotCount == 1 {
			le.next = Token{tipo: "FLOAT", value: numStr}
		} else {
			le.next = Token{tipo: "INT", value: numStr}
		}
		return
	}

	if unicode.IsLetter(rune(char)) {
		startPos := le.position
		for le.position < len(le.source) && (unicode.IsDigit(rune(le.source[le.position])) || unicode.IsLetter(rune(le.source[le.position])) || le.source[le.position] == '_') {
			le.position++
		}
		idenStr := le.source[startPos:le.position]

		switch idenStr {
		case "Println":
			le.next = Token{tipo: "PRINT", value: idenStr}
		case "if":
			le.next = Token{tipo: "IF", value: idenStr}
		case "else":
			le.next = Token{tipo: "ELSE", value: idenStr}
		case "for":
			le.next = Token{tipo: "WHILE", value: idenStr}
		case "Scanln":
			le.next = Token{tipo: "READ", value: idenStr}
		case "ScanlnStr":
			le.next = Token{tipo: "READ_STR", value: idenStr}
		case "var":
			le.next = Token{tipo: "VAR", value: idenStr}
		case "true", "false":
			le.next = Token{tipo: "BOOL", value: idenStr}
		case "string", "int", "bool", "float":
			le.next = Token{tipo: "TYPE", value: idenStr}
		default:
			le.next = Token{tipo: "IDEN", value: idenStr}
		}
		return
	}

	panic("[Lexer] caractere inválido")
}