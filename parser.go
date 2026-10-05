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
    } else if unOp.value == "-" {
        return -unOp.children.Evaluate(st)
    } else if unOp.value == "!" {
        if unOp.children.Evaluate(st) == 0 {
            return 1
        }
        return 0
    }
    panic("[Semantic] operador unário inválido")
}

// BinOp
type BinOp struct {
    value         string
    childrenLeft  Node
    childrenRight Node
}

func (binOp *BinOp) Evaluate(st SymbolTable) int {
    left := binOp.childrenLeft.Evaluate(st)

    if binOp.value == "&&" {
        if left != 0 && binOp.childrenRight.Evaluate(st) != 0 {
            return 1
        }
        return 0
    } else if binOp.value == "||" {
        if left != 0 || binOp.childrenRight.Evaluate(st) != 0 {
            return 1
        }
        return 0
    }

    right := binOp.childrenRight.Evaluate(st)

    switch binOp.value {
    case "+":
        return left + right
    case "-":
        return left - right
    case "*":
        return left * right
    case "/":
        if right == 0 {
            panic("[Semantic] divisão por zero")
        }
        return left / right
    case "==":
        if left == right {
            return 1
        }
        return 0
    case "<":
        if left < right {
            return 1
        }
        return 0
    case ">":
        if left > right {
            return 1
        }
        return 0
    default:
        panic("[Semantic] erro em bin op: " + binOp.value)
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
        child.Evaluate(st) // Executa a instrução
    }
    return 0
}

type NoOp struct{
    
}
func (noOp *NoOp) Evaluate(st SymbolTable) int {
    return 0
}

type If struct {
    condition   Node
    trueBranch  Node
    falseBranch Node
}

func (ifNode *If) Evaluate(st SymbolTable) int {
    if ifNode.condition.Evaluate(st) != 0 {
        return ifNode.trueBranch.Evaluate(st)
    } else if ifNode.falseBranch != nil {
        return ifNode.falseBranch.Evaluate(st)
    }
    return 0
}

type While struct {
    condition Node
    block     Node
}

func (whileNode *While) Evaluate(st SymbolTable) int {
    for whileNode.condition.Evaluate(st) != 0 {
        whileNode.block.Evaluate(st)
    }
    return 0
}

type ReadOp struct{}

func (readOp *ReadOp) Evaluate(st SymbolTable) int {
    var val int
    fmt.Scan(&val)
    return val
}

func ParseFactor(lexer *Lexer) Node {
    if lexer.next.tipo == "INT" {
        num, err := strconv.Atoi(lexer.next.value)
        if err != nil {
            panic("[PARSER] número inválido")
        }
        lexer.selectNext()
        return &IntVal{value: num}
    }

    if lexer.next.tipo == "MINUS" || lexer.next.tipo == "PLUS" || lexer.next.tipo == "NOT" {
        op := lexer.next.value
        lexer.selectNext()
        return &UnOp{
            value:    op,
            children: ParseFactor(lexer),
        }
    }

    if lexer.next.tipo == "READ" {
        lexer.selectNext()
        if lexer.next.tipo != "OPEN_PAR" {
            panic("[PARSER] Esperava '(' após Scanln")
        }
        lexer.selectNext()
        if lexer.next.tipo != "CLOSE_PAR" {
            panic("[PARSER] Esperava ')' após Scanln")
        }
        lexer.selectNext()
        return &ReadOp{}
    }

    if lexer.next.tipo == "OPEN_PAR" {
        lexer.selectNext()
        resultadoInterno := ParseBoolExpression(lexer)
        if lexer.next.tipo != "CLOSE_PAR" {
            panic("[PARSER] faltou fechar o parênteses")
        }
        lexer.selectNext()
        return resultadoInterno
    }

    if lexer.next.tipo == "IDEN" {
        idVal := &IdentifierOp{value: lexer.next.value}
        lexer.selectNext()
        return idVal
    }

    panic("[PARSER] token inesperado: " + lexer.next.value)
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

        resultado = &BinOp{
            value:         simbolo,
            childrenLeft:  resultado,
            childrenRight: proximoTermo,
        }
    }

    return resultado
}

func ParseRelExpression(lexer *Lexer) Node {
    resultado := ParseExpression(lexer)
    if lexer.next.tipo == "EQ" || lexer.next.tipo == "LT" || lexer.next.tipo == "GT" {
        opTipo := lexer.next.tipo
        lexer.selectNext()
        proximo := ParseExpression(lexer)
        simbolo := "=="
        if opTipo == "LT" {
            simbolo = "<"
        } else if opTipo == "GT" {
            simbolo = ">"
        }
        resultado = &BinOp{
            value:         simbolo,
            childrenLeft:  resultado,
            childrenRight: proximo,
        }
    }
    return resultado
}

func ParseBoolTerm(lexer *Lexer) Node {
    resultado := ParseRelExpression(lexer)
    for lexer.next.tipo == "AND" {
        lexer.selectNext()
        proximo := ParseRelExpression(lexer)
        resultado = &BinOp{
            value:         "&&",
            childrenLeft:  resultado,
            childrenRight: proximo,
        }
    }
    return resultado
}

func ParseBoolExpression(lexer *Lexer) Node {
    resultado := ParseBoolTerm(lexer)
    for lexer.next.tipo == "OR" {
        lexer.selectNext()
        proximo := ParseBoolTerm(lexer)
        resultado = &BinOp{
            value:         "||",
            childrenLeft:  resultado,
            childrenRight: proximo,
        }
    }
    return resultado
}

func parseBlock(lexer *Lexer) Node {
    if lexer.next.tipo != "OPEN_BRA" {
        panic("[PARSER] Esperava '{'")
    }
    lexer.selectNext()

    block := &Block{children: []Node{}}
    for lexer.next.tipo != "CLOSE_BRA" && lexer.next.tipo != "EOF" {
        if lexer.next.tipo == "END" {
            lexer.selectNext()
            continue
        }
        stmt := parseStatement(lexer)
        block.children = append(block.children, stmt)
    }

    if lexer.next.tipo != "CLOSE_BRA" {
        panic("[PARSER] Esperava '}'")
    }
    lexer.selectNext()
    return block
}

func parseStatement(lexer *Lexer) Node {
    if lexer.next.tipo == "END" {
        lexer.selectNext()
        return &NoOp{}
    }

    // Se encontrar um bloco isolado `{ ... }` no meio das instruções
    if lexer.next.tipo == "OPEN_BRA" {
        return parseBlock(lexer)
    }

    if lexer.next.tipo == "PRINT" {
        lexer.selectNext()
        if lexer.next.tipo != "OPEN_PAR" {
            panic("[PARSER] Esperava '(' após Println")
        }
        lexer.selectNext()
        expr := ParseBoolExpression(lexer)
        if lexer.next.tipo != "CLOSE_PAR" {
            panic("[PARSER] Esperava ')'")
        }
        lexer.selectNext()
        if lexer.next.tipo == "END" {
            lexer.selectNext()
        }
        return &PrintOp{children: expr}
    }

    if lexer.next.tipo == "IF" {
        lexer.selectNext()
        cond := ParseBoolExpression(lexer)
        
        // Garante que o bloco do IF abre corretamente (evita quebra de linha indevida)
        if lexer.next.tipo != "OPEN_BRA" {
            panic("[PARSER] Esperava '{' após condição do if")
        }
        
        trueBr := parseBlock(lexer)
        var falseBr Node = nil
        
        if lexer.next.tipo == "ELSE" {
            lexer.selectNext()
            if lexer.next.tipo == "IF" {
                falseBr = parseStatement(lexer)
            } else {
                if lexer.next.tipo != "OPEN_BRA" {
                    panic("[PARSER] Esperava '{' após else")
                }
                falseBr = parseBlock(lexer)
            }
        }
        return &If{condition: cond, trueBranch: trueBr, falseBranch: falseBr}
    }

    if lexer.next.tipo == "WHILE" {
        lexer.selectNext()
        cond := ParseBoolExpression(lexer)
        
        if lexer.next.tipo != "OPEN_BRA" {
            panic("[PARSER] Esperava '{' após condição do while")
        }
        
        blk := parseBlock(lexer)
        return &While{condition: cond, block: blk}
    }

    if lexer.next.tipo == "IDEN" {
        idToken := lexer.next
        lexer.selectNext()

        if lexer.next.tipo == "ASSIGN" {
            lexer.selectNext()
            exprRight := ParseBoolExpression(lexer)
            identifier := &IdentifierOp{value: idToken.value}
            assign := &Assignment{
                childrenLeft:  identifier,
                childrenRight: exprRight,
            }
            if lexer.next.tipo == "END" {
                lexer.selectNext()
            }
            return assign
        } else {
            panic("[PARSER] Esperava '=' após identificador")
        }
    }

    panic("[PARSER] Instrução inválida: " + lexer.next.value)
}

func parseProgram(lexer *Lexer) *Block {
    block := &Block{
        children: []Node{},
    }
    for lexer.next.tipo != "EOF" {
        if lexer.next.tipo == "END" {
            lexer.selectNext()
            continue
        }
        // Se o programa começar direto com um bloco `{` (como no Test 1)
        if lexer.next.tipo == "OPEN_BRA" {
            block.children = append(block.children, parseBlock(lexer))
            continue
        }
        stmt := parseStatement(lexer)
        block.children = append(block.children, stmt)
    }

    return block
}

func (p *Parser) run() Node {
    p.lexer.selectNext()
    AST := parseProgram(&p.lexer)

    if p.lexer.next.tipo != "EOF" {
        panic("[PARSER] token inesperado: " + p.lexer.next.value)
    }

    return AST
}