package main

import (
	"strconv"
)

// interface node

type Node interface {
	Evaluate() int
}

// IntVal

type IntVal struct {
	value int
}

func (intVal *IntVal) Evaluate() int {
	return intVal.value
}

// UnOp
type UnOp struct {
	value    string
	children Node
}

func (unOp *UnOp) Evaluate() int {
	if unOp.value == "+" {
		return +unOp.children.Evaluate()
	} else {
		return -unOp.children.Evaluate()
	}
}

// FactOp
type FactOp struct {
	children Node
}

func (factOp *FactOp) Evaluate() int {
	val := factOp.children.Evaluate()
	if val < 0 {
		panic("[Semantic] fatorial de número negativo")
	}
	res := 1
	for i := 2; i <= val; i++ {
		res *= i
	}
	return res
}

// BinOp
type BinOp struct {
	value         string
	childrenLeft  Node
	childrenRight Node
}

func (binOp *BinOp) Evaluate() int {
	if binOp.value == "+" {
		return binOp.childrenLeft.Evaluate() + binOp.childrenRight.Evaluate()
	} else if binOp.value == "-" {
		return binOp.childrenLeft.Evaluate() - binOp.childrenRight.Evaluate()
	} else if binOp.value == "*" {
		return binOp.childrenLeft.Evaluate() * binOp.childrenRight.Evaluate()
	} else if binOp.value == "/" {
		divisor := binOp.childrenRight.Evaluate()
        if divisor == 0 {
            panic("[Semantic] divisão por zero")
        }
        return binOp.childrenLeft.Evaluate() / divisor
	} else {
		panic("[Semantic] erro em bin op")
	}

}

type Parser struct {
	lexer Lexer
}

func ParseBase(lexer *Lexer) Node {
	// 1. Se for um número, consome e retorna o valor
	if lexer.next.tipo == "INT" {
		num, err := strconv.Atoi(lexer.next.value)
		if err != nil {
			panic("[PARSER] número inválido")
		}
		lexer.selectNext() // Consome o INT

		intVal := &IntVal{
			value: num,
		}
		return intVal
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

func ParseFact(lexer *Lexer) Node {
	node := ParseBase(lexer)

	for lexer.next.tipo == "FACT" {
		lexer.selectNext()
		node = &FactOp{
			children: node,
		}
	}

	return node
}

func ParseFactor(lexer *Lexer) Node {
	if lexer.next.tipo == "MINUS" || lexer.next.tipo == "PLUS" {
		op := lexer.next.value
		lexer.selectNext()
		return &UnOp{
			value:    op,
			children: ParseFactor(lexer),
		}
	}

	return ParseFact(lexer)
}

func ParseTerm(lexer *Lexer) Node {
    resultado := ParseFactor(lexer)

    for lexer.next.tipo == "MULT" || lexer.next.tipo == "DIV" {
        opTipo := lexer.next.tipo
        lexer.selectNext()

        proximoFator := ParseFactor(lexer)

        simbolo := "*"
        if opTipo == "DIV" {
            simbolo = "/"
        }

        // Encadeia o nó atual como filho à esquerda do novo operador
        resultado = &BinOp{
            value:         simbolo,
            childrenLeft:  resultado,
            childrenRight: proximoFator,
        }
    }

    return resultado
}

func ParseExpression(lexer *Lexer) Node {
    resultado := ParseTerm(lexer)

    for lexer.next.tipo == "PLUS" || lexer.next.tipo == "MINUS" {
        opTipo := lexer.next.tipo
        lexer.selectNext()

        proximoTermo := ParseTerm(lexer)

        simbolo := "+"
        if opTipo == "MINUS" {
            simbolo = "-"
        }

        // Atualiza a raiz parcial da árvore
        resultado = &BinOp{
            value:         simbolo,
            childrenLeft:  resultado,
            childrenRight: proximoTermo,
        }
    }

    return resultado
}

func (p *Parser) run() int {
	// add o primeiro next
	p.lexer.selectNext()

	//vou ter a AST montada
	AST := ParseExpression(&p.lexer)

	if p.lexer.next.tipo != "EOF" {
		panic("[PARSER] token inesperado: " + p.lexer.next.value)
	}

	return AST.Evaluate()
}
