package main

import (
	"fmt"
	"strconv"
)

// ─────────────────────────────────────────────
// Interface Node
// ─────────────────────────────────────────────

type Node interface {
	Evaluate(st SymbolTable) Variable
}

// ─────────────────────────────────────────────
// IntVal
// ─────────────────────────────────────────────

type IntVal struct {
	value int
}

func (n *IntVal) Evaluate(st SymbolTable) Variable {
	return Variable{value: n.value, varType: "int"}
}

// ─────────────────────────────────────────────
// BoolVal
// ─────────────────────────────────────────────

type BoolVal struct {
	value int // 1 = true, 0 = false
}

func (n *BoolVal) Evaluate(st SymbolTable) Variable {
	return Variable{value: n.value, varType: "bool"}
}

// ─────────────────────────────────────────────
// StringVal
// ─────────────────────────────────────────────

type StringVal struct {
	value string
}

func (n *StringVal) Evaluate(st SymbolTable) Variable {
	return Variable{value: n.value, varType: "string"}
}

// ─────────────────────────────────────────────
// UnOp
// ─────────────────────────────────────────────

type UnOp struct {
	value    string
	children Node
}

func (n *UnOp) Evaluate(st SymbolTable) Variable {
	child := n.children.Evaluate(st)

	switch n.value {
	case "+":
		if child.varType != "int" {
			panic("[Semantic] operador '+' unário requer int, recebeu " + child.varType)
		}
		return Variable{value: child.value.(int), varType: "int"}
	case "-":
		if child.varType != "int" {
			panic("[Semantic] operador '-' unário requer int, recebeu " + child.varType)
		}
		return Variable{value: -child.value.(int), varType: "int"}
	case "!":
		if child.varType != "bool" && child.varType != "int" {
			panic("[Semantic] operador '!' requer bool ou int, recebeu " + child.varType)
		}
		if child.varType == "bool" {
			if child.value.(int) == 0 {
				return Variable{value: 1, varType: "bool"}
			}
			return Variable{value: 0, varType: "bool"}
		}
		// int
		if child.value.(int) == 0 {
			return Variable{value: 1, varType: "int"}
		}
		return Variable{value: 0, varType: "int"}
	}
	panic("[Semantic] operador unário inválido: " + n.value)
}

// ─────────────────────────────────────────────
// BinOp
// ─────────────────────────────────────────────

type BinOp struct {
	value         string
	childrenLeft  Node
	childrenRight Node
}

func (n *BinOp) Evaluate(st SymbolTable) Variable {
	left := n.childrenLeft.Evaluate(st)

	// Short-circuit para && e ||
	if n.value == "&&" {
		if left.varType != "bool" && left.varType != "int" {
			panic("[Semantic] '&&' requer bool ou int nos operandos")
		}
		leftInt := toInt(left)
		if leftInt == 0 {
			return Variable{value: 0, varType: "bool"}
		}
		right := n.childrenRight.Evaluate(st)
		if right.varType != "bool" && right.varType != "int" {
			panic("[Semantic] '&&' requer bool ou int nos operandos")
		}
		if toInt(right) != 0 {
			return Variable{value: 1, varType: "bool"}
		}
		return Variable{value: 0, varType: "bool"}
	}

	if n.value == "||" {
		if left.varType != "bool" && left.varType != "int" {
			panic("[Semantic] '||' requer bool ou int nos operandos")
		}
		leftInt := toInt(left)
		if leftInt != 0 {
			return Variable{value: 1, varType: "bool"}
		}
		right := n.childrenRight.Evaluate(st)
		if right.varType != "bool" && right.varType != "int" {
			panic("[Semantic] '||' requer bool ou int nos operandos")
		}
		if toInt(right) != 0 {
			return Variable{value: 1, varType: "bool"}
		}
		return Variable{value: 0, varType: "bool"}
	}

	right := n.childrenRight.Evaluate(st)

	switch n.value {
	case "+":
		// Concatenação de strings ou soma de ints
		if left.varType == "string" && right.varType == "string" {
			return Variable{value: left.value.(string) + right.value.(string), varType: "string"}
		}
		if left.varType == "int" && right.varType == "int" {
			return Variable{value: left.value.(int) + right.value.(int), varType: "int"}
		}
		panic("[Semantic] '+' requer dois int ou duas strings; recebeu " + left.varType + " e " + right.varType)
	case "-":
		requireBothInt(left, right, "-")
		return Variable{value: left.value.(int) - right.value.(int), varType: "int"}
	case "*":
		requireBothInt(left, right, "*")
		return Variable{value: left.value.(int) * right.value.(int), varType: "int"}
	case "/":
		requireBothInt(left, right, "/")
		if right.value.(int) == 0 {
			panic("[Semantic] divisão por zero")
		}
		return Variable{value: left.value.(int) / right.value.(int), varType: "int"}
	case "==":
		if left.varType != right.varType {
			panic("[Semantic] '==' requer operandos do mesmo tipo; recebeu " + left.varType + " e " + right.varType)
		}
		result := 0
		switch left.varType {
		case "int", "bool":
			if left.value.(int) == right.value.(int) {
				result = 1
			}
		case "string":
			if left.value.(string) == right.value.(string) {
				result = 1
			}
		}
		return Variable{value: result, varType: "bool"}
	case "<":
		requireBothInt(left, right, "<")
		result := 0
		if left.value.(int) < right.value.(int) {
			result = 1
		}
		return Variable{value: result, varType: "bool"}
	case ">":
		requireBothInt(left, right, ">")
		result := 0
		if left.value.(int) > right.value.(int) {
			result = 1
		}
		return Variable{value: result, varType: "bool"}
	}
	panic("[Semantic] operador binário inválido: " + n.value)
}

// requireBothInt panics se os dois operandos não forem int
func requireBothInt(l, r Variable, op string) {
	if l.varType != "int" || r.varType != "int" {
		panic("[Semantic] '" + op + "' requer dois int; recebeu " + l.varType + " e " + r.varType)
	}
}

// toInt converte o valor interno de bool/int para int (0 ou 1)
func toInt(v Variable) int {
	return v.value.(int)
}

// ─────────────────────────────────────────────
// IdentifierOp
// ─────────────────────────────────────────────

type IdentifierOp struct {
	value string // nome da variável
}

func (n *IdentifierOp) Evaluate(st SymbolTable) Variable {
	return st.getter(n.value)
}

// ─────────────────────────────────────────────
// PrintOp
// ─────────────────────────────────────────────

type PrintOp struct {
	children Node
}

func (n *PrintOp) Evaluate(st SymbolTable) Variable {
	val := n.children.Evaluate(st)
	switch val.varType {
	case "int":
		fmt.Println(val.value.(int))
	case "bool":
		if val.value.(int) != 0 {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
	case "string":
		fmt.Println(val.value.(string))
	}
	return Variable{value: 0, varType: "int"}
}

// ─────────────────────────────────────────────
// Assignment
// ─────────────────────────────────────────────

type Assignment struct {
	childrenLeft  Node // IdentifierOp
	childrenRight Node // expressão
}

func (n *Assignment) Evaluate(st SymbolTable) Variable {
	name := n.childrenLeft.(*IdentifierOp).value
	val := n.childrenRight.Evaluate(st)
	st.setter(name, val)
	return Variable{value: 0, varType: "int"}
}

// ─────────────────────────────────────────────
// VarDec
// ─────────────────────────────────────────────

type VarDec struct {
	value    string // tipo da variável ("int", "string", "bool")
	children []Node // [0]=IdentifierOp, [1]=expressão de inicialização (opcional)
}

func (n *VarDec) Evaluate(st SymbolTable) Variable {
	name := n.children[0].(*IdentifierOp).value
	st.createVariable(name, n.value)
	if len(n.children) > 1 {
		val := n.children[1].Evaluate(st)
		st.setter(name, val)
	}
	return Variable{value: 0, varType: "int"}
}

// ─────────────────────────────────────────────
// Block
// ─────────────────────────────────────────────

type Block struct {
	children []Node
}

func (n *Block) Evaluate(st SymbolTable) Variable {
	for _, child := range n.children {
		child.Evaluate(st)
	}
	return Variable{value: 0, varType: "int"}
}

// ─────────────────────────────────────────────
// NoOp
// ─────────────────────────────────────────────

type NoOp struct{}

func (n *NoOp) Evaluate(st SymbolTable) Variable {
	return Variable{value: 0, varType: "int"}
}

// ─────────────────────────────────────────────
// If
// ─────────────────────────────────────────────

type If struct {
	condition   Node
	trueBranch  Node
	falseBranch Node
}

func (n *If) Evaluate(st SymbolTable) Variable {
	cond := n.condition.Evaluate(st)
	if cond.varType != "bool" && cond.varType != "int" {
		panic("[Semantic] condição do if deve ser bool ou int")
	}
	if toInt(cond) != 0 {
		return n.trueBranch.Evaluate(st)
	} else if n.falseBranch != nil {
		return n.falseBranch.Evaluate(st)
	}
	return Variable{value: 0, varType: "int"}
}

// ─────────────────────────────────────────────
// While
// ─────────────────────────────────────────────

type While struct {
	condition Node
	block     Node
}

func (n *While) Evaluate(st SymbolTable) Variable {
	for {
		cond := n.condition.Evaluate(st)
		if cond.varType != "bool" && cond.varType != "int" {
			panic("[Semantic] condição do for deve ser bool ou int")
		}
		if toInt(cond) == 0 {
			break
		}
		n.block.Evaluate(st)
	}
	return Variable{value: 0, varType: "int"}
}

// ─────────────────────────────────────────────
// ReadOp  (Scanln — lê um inteiro)
// ─────────────────────────────────────────────

type ReadOp struct{}

func (n *ReadOp) Evaluate(st SymbolTable) Variable {
	var val int
	fmt.Scan(&val)
	return Variable{value: val, varType: "int"}
}

// ─────────────────────────────────────────────
// ReadStringOp  (ScanlnStr — lê uma string)
// ─────────────────────────────────────────────

type ReadStringOp struct{}

func (n *ReadStringOp) Evaluate(st SymbolTable) Variable {
	var val string
	fmt.Scan(&val)
	return Variable{value: val, varType: "string"}
}

// ─────────────────────────────────────────────
// Parser struct
// ─────────────────────────────────────────────

type Parser struct {
	lexer Lexer
}

// ─────────────────────────────────────────────
// ParseFactor
// ─────────────────────────────────────────────

func ParseFactor(lexer *Lexer) Node {
	// Literal inteiro
	if lexer.next.tipo == "INT" {
		num, err := strconv.Atoi(lexer.next.value)
		if err != nil {
			panic("[PARSER] número inválido")
		}
		lexer.selectNext()
		return &IntVal{value: num}
	}

	// Literal booleano
	if lexer.next.tipo == "BOOL" {
		val := 0
		if lexer.next.value == "true" {
			val = 1
		}
		lexer.selectNext()
		return &BoolVal{value: val}
	}

	// Literal string
	if lexer.next.tipo == "STR" {
		s := lexer.next.value
		lexer.selectNext()
		return &StringVal{value: s}
	}

	// Operadores unários
	if lexer.next.tipo == "MINUS" || lexer.next.tipo == "PLUS" || lexer.next.tipo == "NOT" {
		op := lexer.next.value
		lexer.selectNext()
		return &UnOp{
			value:    op,
			children: ParseFactor(lexer),
		}
	}

	// Scanln()
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

	// ScanlnStr()
	if lexer.next.tipo == "READ_STR" {
		lexer.selectNext()
		if lexer.next.tipo != "OPEN_PAR" {
			panic("[PARSER] Esperava '(' após ScanlnStr")
		}
		lexer.selectNext()
		if lexer.next.tipo != "CLOSE_PAR" {
			panic("[PARSER] Esperava ')' após ScanlnStr")
		}
		lexer.selectNext()
		return &ReadStringOp{}
	}

	// Parênteses
	if lexer.next.tipo == "OPEN_PAR" {
		lexer.selectNext()
		inner := ParseBoolExpression(lexer)
		if lexer.next.tipo != "CLOSE_PAR" {
			panic("[PARSER] faltou fechar o parênteses")
		}
		lexer.selectNext()
		return inner
	}

	// Identificador
	if lexer.next.tipo == "IDEN" {
		idVal := &IdentifierOp{value: lexer.next.value}
		lexer.selectNext()
		return idVal
	}

	panic("[PARSER] token inesperado em factor: " + lexer.next.value)
}

// ─────────────────────────────────────────────
// ParseTerm
// ─────────────────────────────────────────────

func ParseTerm(lexer *Lexer) Node {
	resultado := ParseFactor(lexer)

	for lexer.next.tipo == "MULT" || lexer.next.tipo == "DIV" {
		opTipo := lexer.next.tipo
		lexer.selectNext()
		proximo := ParseFactor(lexer)
		simbolo := "*"
		if opTipo == "DIV" {
			simbolo = "/"
		}
		resultado = &BinOp{
			value:         simbolo,
			childrenLeft:  resultado,
			childrenRight: proximo,
		}
	}

	return resultado
}

// ─────────────────────────────────────────────
// ParseExpression
// ─────────────────────────────────────────────

func ParseExpression(lexer *Lexer) Node {
	resultado := ParseTerm(lexer)

	for lexer.next.tipo == "PLUS" || lexer.next.tipo == "MINUS" {
		opTipo := lexer.next.tipo
		lexer.selectNext()
		proximo := ParseTerm(lexer)
		simbolo := "+"
		if opTipo == "MINUS" {
			simbolo = "-"
		}
		resultado = &BinOp{
			value:         simbolo,
			childrenLeft:  resultado,
			childrenRight: proximo,
		}
	}

	return resultado
}

// ─────────────────────────────────────────────
// ParseRelExpression
// ─────────────────────────────────────────────

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

// ─────────────────────────────────────────────
// ParseBoolTerm
// ─────────────────────────────────────────────

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

// ─────────────────────────────────────────────
// ParseBoolExpression
// ─────────────────────────────────────────────

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

// ─────────────────────────────────────────────
// parseBlock
// ─────────────────────────────────────────────

func parseBlock(lexer *Lexer) Node {
	if lexer.next.tipo != "OPEN_BRA" {
		panic("[PARSER] Esperava '{'")
	}
	lexer.selectNext()

	// '}' na mesma linha que '{' é inválido
	if lexer.next.tipo == "CLOSE_BRA" {
		panic("[PARSER] Unexpected token CLOSE_BRA: '}' não pode estar na mesma linha que '{'")
	}

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

// ─────────────────────────────────────────────
// parseStatement
// ─────────────────────────────────────────────

func parseStatement(lexer *Lexer) Node {
	if lexer.next.tipo == "END" {
		lexer.selectNext()
		return &NoOp{}
	}

	if lexer.next.tipo == "OPEN_BRA" {
		return parseBlock(lexer)
	}

	// Declaração de variável: var <iden> <type> [= <expr>]
	if lexer.next.tipo == "VAR" {
		lexer.selectNext()

		if lexer.next.tipo != "IDEN" {
			panic("[PARSER] Esperava identificador após 'var'")
		}
		iden := &IdentifierOp{value: lexer.next.value}
		lexer.selectNext()

		if lexer.next.tipo != "TYPE" {
			panic("[PARSER] Esperava tipo após o identificador em 'var'")
		}
		varType := lexer.next.value
		lexer.selectNext()

		node := &VarDec{
			value:    varType,
			children: []Node{iden},
		}

		// Inicialização opcional: = <expr>
		if lexer.next.tipo == "ASSIGN" {
			lexer.selectNext()
			expr := ParseBoolExpression(lexer)
			node.children = append(node.children, expr)
		}

		if lexer.next.tipo == "END" {
			lexer.selectNext()
		}
		return node
	}

	// Println(<expr>)
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

	// if <expr> { ... } [else { ... }]
	if lexer.next.tipo == "IF" {
		lexer.selectNext()
		cond := ParseBoolExpression(lexer)

		if lexer.next.tipo == "END" {
			panic("[PARSER] Esperava '{' na mesma linha após a condição do if")
		}
		if lexer.next.tipo != "OPEN_BRA" {
			panic("[PARSER] Esperava '{' após condição do if")
		}

		trueBr := parseBlock(lexer)
		var falseBr Node = nil

		// 'else' deve estar na mesma linha que '}'
		if lexer.next.tipo == "ELSE" {
			lexer.selectNext()

			if lexer.next.tipo == "END" {
				panic("[PARSER] Unexpected token NEWLINE: esperava '{' na mesma linha após else")
			}

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

	// for <expr> { ... }
	if lexer.next.tipo == "WHILE" {
		lexer.selectNext()
		cond := ParseBoolExpression(lexer)

		if lexer.next.tipo == "END" || lexer.next.tipo != "OPEN_BRA" {
			panic("[PARSER] Esperava '{' após condição do for")
		}

		blk := parseBlock(lexer)
		return &While{condition: cond, block: blk}
	}

	// Atribuição: <iden> = <expr>
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
		}
		panic("[PARSER] Esperava '=' após identificador")
	}

	panic("[PARSER] Instrução inválida: " + lexer.next.value)
}

// ─────────────────────────────────────────────
// parseProgram
// ─────────────────────────────────────────────

func parseProgram(lexer *Lexer) *Block {
	block := &Block{
		children: []Node{},
	}
	for lexer.next.tipo != "EOF" {
		if lexer.next.tipo == "END" {
			lexer.selectNext()
			continue
		}
		if lexer.next.tipo == "OPEN_BRA" {
			block.children = append(block.children, parseBlock(lexer))
			continue
		}
		stmt := parseStatement(lexer)
		block.children = append(block.children, stmt)
	}

	return block
}

// ─────────────────────────────────────────────
// Parser.run
// ─────────────────────────────────────────────

func (p *Parser) run() Node {
	p.lexer.selectNext()
	AST := parseProgram(&p.lexer)

	if p.lexer.next.tipo != "EOF" {
		panic("[PARSER] token inesperado: " + p.lexer.next.value)
	}

	return AST
}