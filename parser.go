package main

import (
    "fmt"
    "strconv"
)

// interface node

type Node interface {
    Evaluate(st SymbolTable) int
}

// IntVal

type IntVal struct {
    value int
}

func (intVal *IntVal) Evaluate(st SymbolTable) int {
    return intVal.value
}

// UnOp
type UnOp struct {
    value    string
    children Node
}

func (unOp *UnOp) Evaluate(st SymbolTable) int {
    if unOp.value == "+" {
        return +unOp.children.Evaluate(st)
    } else {
        return -unOp.children.Evaluate(st)
    }
}

// BinOp
type BinOp struct {
    value         string
    childrenLeft  Node
    childrenRight Node
}

func (binOp *BinOp) Evaluate(st SymbolTable) int {
    if binOp.value == "+" {
        return binOp.childrenLeft.Evaluate(st) + binOp.childrenRight.Evaluate(st)
    } else if binOp.value == "-" {
        return binOp.childrenLeft.Evaluate(st) - binOp.childrenRight.Evaluate(st)
    } else if binOp.value == "*" {
        return binOp.childrenLeft.Evaluate(st) * binOp.childrenRight.Evaluate(st)
    } else if binOp.value == "/" {
        divisor := binOp.childrenRight.Evaluate(st)
        if divisor == 0 {
            panic("[Semantic] divisão por zero")
        }
        return binOp.childrenLeft.Evaluate(st) / divisor
    } else {
        panic("[Semantic] erro em bin op")
    }

}

type IdentifierOp struct {
    value string // Guarda o nome da variável (ex: "x1")
}

//retornamos o valor da variavel na hora da execucao
func (idOp *IdentifierOp) Evaluate(st SymbolTable) int {
    return st.getter(idOp.value)
}

type PrintOp struct {
    children Node
}

//printamos a informacao do evaluate do filho
func (printOp *PrintOp) Evaluate(st SymbolTable) int {
    fmt.Println(printOp.children.Evaluate(st))
    return 0
}

type Assignment struct {
    childrenLeft  Node
    childrenRight Node
}

//atribuimos na tabela o valor da variavel
func (assignment *Assignment) Evaluate(st SymbolTable) int {
    st.setter(assignment.childrenLeft.(*IdentifierOp).value, assignment.childrenRight.Evaluate(st))
    return 0
}
type Parser struct {
    lexer Lexer
}

type Block struct {
    children []Node // Uma lista (slice) de nós (instruções)
}

func (block *Block) Evaluate(st SymbolTable) int {
    // Percorre cada instrução do bloco em ordem
    for _, child := range block.children {
        child.Evaluate(st) // Executa a instrução (
    }
    return 0
}

type NoOp struct{
    
}
func (noOp *NoOp) Evaluate(st SymbolTable) int {
    return 0
}

func ParseFactor(lexer *Lexer) Node {
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

	if lexer.next.tipo == "MINUS" || lexer.next.tipo == "PLUS" {

		// colocamos como filho o nodo resultado de parseFactor
		if lexer.next.tipo == "MINUS" {

			lexer.selectNext()
			unOp := &UnOp{
				value:    "-",
				children: ParseFactor(lexer),
			}
			return unOp

		} else {

			lexer.selectNext()
			unOp := &UnOp{
				value:    "+",
				children: ParseFactor(lexer),
			}
			return unOp

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

	// Se for um identificador (variável), cria o nó de busca na tabela de símbolos
    if lexer.next.tipo == "IDEN" {
        idVal := &IdentifierOp{
            value: lexer.next.value,
        }
        lexer.selectNext() // Consome o IDEN
        return idVal
    }

	panic("[PARSER] token inesperado, esperava número ou '(': " + lexer.next.value)
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

// aqui lemos o statement especifico da linha 
func parseStatement(lexer *Lexer) Node {
    
    // Se a linha estiver vazia (apenas uma quebra de linha)
    if lexer.next.tipo == "END" {
        lexer.selectNext() // Consome o END
        return &NoOp{}
    }

    // Se for uma instrução de Println
    if lexer.next.tipo == "PRINT" {
        lexer.selectNext() // Consome "Println"

        if lexer.next.tipo != "OPEN_PAR" {
            panic("[PARSER] Esperava '(' após Println")
        }
        lexer.selectNext() // Consome "("

        expr := ParseExpression(lexer) // Resolve a expressão interna

        if lexer.next.tipo != "CLOSE_PAR" {
            panic("[PARSER] Esperava ')' para fechar o Println")
        }
        lexer.selectNext() // Consome ")"

        // Consome o final da linha (END) ou EOF
        if lexer.next.tipo == "END" {
            lexer.selectNext()
        }

        return &PrintOp{children: expr}
    }

    // Se começar com um Identificador, pode ser uma Atribuição
    if lexer.next.tipo == "IDEN" {
        idToken := lexer.next
        lexer.selectNext() // Consome o IDEN

        // Verifica se o próximo token é o '=' (ASSIGN)
        if lexer.next.tipo == "ASSIGN" {
            lexer.selectNext() // Consome o "="

            exprRight := ParseExpression(lexer) // Resolve o valor à direita

            identifier := &IdentifierOp{value: idToken.value}

            assign := &Assignment{
                childrenLeft:  identifier,
                childrenRight: exprRight,
            }

            // Consome o final da linha (END) ou EOF
            if lexer.next.tipo == "END" {
                lexer.selectNext()
            }

            return assign
        } else {
            panic("[PARSER] Esperava '=' após identificador")
        }
    }

    panic("[PARSER] Instrução inválida ou token inesperado: " + lexer.next.value)
}


//parseProgram
func parseProgram(lexer *Lexer) *Block{
	block := &Block{
		children: []Node{},
	}
	for lexer.next.tipo != "EOF"{
		stmt := parseStatement(lexer)
		block.children = append(block.children, stmt)
	}

	return block
} 

func (p *Parser) run() Node{
	// add o primeiro next
	p.lexer.selectNext()
	//vou ter a AST montada
	AST := parseProgram(&p.lexer)

	if p.lexer.next.tipo != "EOF" {
		panic("[PARSER] token inesperado: " + p.lexer.next.value)
	}

	return AST
}
