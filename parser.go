package main

import (
	"fmt"
	"strconv"
)

// ============================================================================
// Roteiro 8 - Geração de Código x86 (NASM)
// ============================================================================
// Mudanças realizadas:
// 1. Adicionado o método Generate(st *SymbolTable) na interface Node.
// 2. Adicionado o atributo id em cada struct de nó e implementada a função
//    newId() para geração de identificadores únicos sequenciais (usados em rótulos).
// 3. Atualizados todos os pontos de instanciação de nós no Parser para inicializar
//    o campo id com newId().
// 4. Implementado o método Generate(st *SymbolTable) para todos os nós:
//    - IntVal / BoolVal: mov eax, <valor>
//    - IdentifierOp: mov eax, [ebp-<shift>]
//    - VarDec: sub esp, 4 e inicialização com mov [ebp-<shift>], eax se houver expr
//    - Assignment: calcula expr da direita e faz mov [ebp-<shift>], eax
//    - BinOp: avalia operando da direita (push eax), depois da esquerda (pop ecx),
//      e executa a instrução correspondente:
//      add, sub, imul, cdq/idiv, cmp/cmov(e/l/g), and, or
//    - UnOp: avalia filho em eax e executa neg eax (-) ou cmove (!)
//    - PrintOp: avalia filho, empilha eax, empilha format_out, chama printf, add esp, 8
//    - ReadOp (Scanln): empilha scan_int, format_in, chama scanf, mov eax, dword [scan_int]
//    - If: rótulos únicos else_<id> e exit_<id>, cmp eax, 0 e je
//    - While: rótulos únicos loop_<id> e exit_<id>, cmp eax, 0, je exit_<id>, jmp loop_<id>
//    - Block / NoOp: iteração sobre filhos ou sem operação
//    - StringVal / ReadStringOp: sem geração de código de string conforme especificação
// 5. Métodos Evaluate() de todos os nós foram mantidos intactos.
// 6. Removida a restrição que impedia '}' na mesma linha que '{' em parseBlock,
//    permitindo blocos vazios inline como '{{ }}'.
// ============================================================================

// Identificador único para rótulos em nós AST
var nextNodeId int = 0

// newId incrementa e retorna o próximo identificador único
func newId() int {
	nextNodeId++
	return nextNodeId
}

// resetNodeId reinicia o contador de identificadores únicos
func resetNodeId() {
	nextNodeId = 0
}

// ─────────────────────────────────────────────
// Interface Node
// ─────────────────────────────────────────────

type Node interface {
	Evaluate(st SymbolTable) Variable
	Generate(st *SymbolTable)
}

// ─────────────────────────────────────────────
// IntVal
// ─────────────────────────────────────────────

type IntVal struct {
	id    int
	value int
}

func (n *IntVal) Evaluate(st SymbolTable) Variable {
	return Variable{value: n.value, varType: "int"}
}

func (n *IntVal) Generate(st *SymbolTable) {
	CodeAppend(fmt.Sprintf("mov eax, %d", n.value))
}

// ─────────────────────────────────────────────
// BoolVal
// ─────────────────────────────────────────────

type BoolVal struct {
	id    int
	value int // 1 = true, 0 = false
}

func (n *BoolVal) Evaluate(st SymbolTable) Variable {
	return Variable{value: n.value, varType: "bool"}
}

func (n *BoolVal) Generate(st *SymbolTable) {
	CodeAppend(fmt.Sprintf("mov eax, %d", n.value))
}

// ─────────────────────────────────────────────
// StringVal
// ─────────────────────────────────────────────

type StringVal struct {
	id    int
	value string
}

func (n *StringVal) Evaluate(st SymbolTable) Variable {
	return Variable{value: n.value, varType: "string"}
}

func (n *StringVal) Generate(st *SymbolTable) {
	// Não será gerado código para nós de strings
}

// ─────────────────────────────────────────────
// UnOp
// ─────────────────────────────────────────────

type UnOp struct {
	id       int
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
		if child.varType != "bool" {
			panic("[Semantic] operador '!' requer bool, recebeu " + child.varType)
		}
		if child.value.(int) == 0 {
			return Variable{value: 1, varType: "bool"}
		}
		return Variable{value: 0, varType: "bool"}
	}
	panic("[Semantic] operador unário inválido: " + n.value)
}

func (n *UnOp) Generate(st *SymbolTable) {
	n.children.Generate(st)
	switch n.value {
	case "+":
		// Valor já se encontra em EAX
	case "-":
		CodeAppend("neg eax")
	case "!":
		CodeAppend("cmp eax, 0")
		CodeAppend("mov eax, 0")
		CodeAppend("mov ecx, 1")
		CodeAppend("cmove eax, ecx")
	}
}

// ─────────────────────────────────────────────
// BinOp
// ─────────────────────────────────────────────

type BinOp struct {
	id            int
	value         string
	childrenLeft  Node
	childrenRight Node
}

func (n *BinOp) Evaluate(st SymbolTable) Variable {
	left := n.childrenLeft.Evaluate(st)

	// Short-circuit para && e ||
	if n.value == "&&" {
		if left.varType != "bool" {
			panic("[Semantic] '&&' requer bool nos operandos; recebeu " + left.varType)
		}
		if left.value.(int) == 0 {
			return Variable{value: 0, varType: "bool"}
		}
		right := n.childrenRight.Evaluate(st)
		if right.varType != "bool" {
			panic("[Semantic] '&&' requer bool nos operandos; recebeu " + right.varType)
		}
		if right.value.(int) != 0 {
			return Variable{value: 1, varType: "bool"}
		}
		return Variable{value: 0, varType: "bool"}
	}

	if n.value == "||" {
		if left.varType != "bool" {
			panic("[Semantic] '||' requer bool nos operandos; recebeu " + left.varType)
		}
		if left.value.(int) != 0 {
			return Variable{value: 1, varType: "bool"}
		}
		right := n.childrenRight.Evaluate(st)
		if right.varType != "bool" {
			panic("[Semantic] '||' requer bool nos operandos; recebeu " + right.varType)
		}
		if right.value.(int) != 0 {
			return Variable{value: 1, varType: "bool"}
		}
		return Variable{value: 0, varType: "bool"}
	}

	right := n.childrenRight.Evaluate(st)

	switch n.value {
	case "+":
		// Se pelo menos um operando for string, converte o outro e concatena
		if left.varType == "string" || right.varType == "string" {
			return Variable{value: varToString(left) + varToString(right), varType: "string"}
		}
		// Soma de dois ints
		if left.varType == "int" && right.varType == "int" {
			return Variable{value: left.value.(int) + right.value.(int), varType: "int"}
		}
		panic("[Semantic] '+' requer dois int, duas strings, ou string com qualquer tipo; recebeu " + left.varType + " e " + right.varType)
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
		// Comparação de int ou string (lexicográfica)
		if left.varType == "string" && right.varType == "string" {
			result := 0
			if left.value.(string) < right.value.(string) {
				result = 1
			}
			return Variable{value: result, varType: "bool"}
		}
		requireBothInt(left, right, "<")
		result := 0
		if left.value.(int) < right.value.(int) {
			result = 1
		}
		return Variable{value: result, varType: "bool"}
	case ">":
		// Comparação de int ou string (lexicográfica)
		if left.varType == "string" && right.varType == "string" {
			result := 0
			if left.value.(string) > right.value.(string) {
				result = 1
			}
			return Variable{value: result, varType: "bool"}
		}
		requireBothInt(left, right, ">")
		result := 0
		if left.value.(int) > right.value.(int) {
			result = 1
		}
		return Variable{value: result, varType: "bool"}
	}
	panic("[Semantic] operador binário inválido: " + n.value)
}

func (n *BinOp) Generate(st *SymbolTable) {
	// Avalia a subárvore direita e empilha o resultado
	n.childrenRight.Generate(st)
	CodeAppend("push eax")

	// Avalia a subárvore esquerda (resultado fica em EAX)
	n.childrenLeft.Generate(st)

	// Desempilha o resultado da direita para ECX
	CodeAppend("pop ecx")

	// Executa a operação com EAX (esquerda) e ECX (direita)
	switch n.value {
	case "+":
		CodeAppend("add eax, ecx")
	case "-":
		CodeAppend("sub eax, ecx")
	case "*":
		CodeAppend("imul ecx")
	case "/":
		CodeAppend("cdq")
		CodeAppend("idiv ecx")
	case "==":
		CodeAppend("cmp eax, ecx")
		CodeAppend("mov eax, 0")
		CodeAppend("mov ecx, 1")
		CodeAppend("cmove eax, ecx")
	case "<":
		CodeAppend("cmp eax, ecx")
		CodeAppend("mov eax, 0")
		CodeAppend("mov ecx, 1")
		CodeAppend("cmovl eax, ecx")
	case ">":
		CodeAppend("cmp eax, ecx")
		CodeAppend("mov eax, 0")
		CodeAppend("mov ecx, 1")
		CodeAppend("cmovg eax, ecx")
	case "&&":
		CodeAppend("and eax, ecx")
	case "||":
		CodeAppend("or eax, ecx")
	}
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

// varToString converte qualquer Variable para sua representação em string
func varToString(v Variable) string {
	switch v.varType {
	case "string":
		return v.value.(string)
	case "int":
		return strconv.Itoa(v.value.(int))
	case "bool":
		if v.value.(int) != 0 {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v.value)
	}
}

// ─────────────────────────────────────────────
// IdentifierOp
// ─────────────────────────────────────────────

type IdentifierOp struct {
	id    int
	value string // nome da variável
}

func (n *IdentifierOp) Evaluate(st SymbolTable) Variable {
	return st.getter(n.value)
}

func (n *IdentifierOp) Generate(st *SymbolTable) {
	v := st.getter(n.value)
	CodeAppend(fmt.Sprintf("mov eax, [ebp-%d]", v.shift))
}

// ─────────────────────────────────────────────
// PrintOp
// ─────────────────────────────────────────────

type PrintOp struct {
	id       int
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

func (n *PrintOp) Generate(st *SymbolTable) {
	n.children.Generate(st)
	CodeAppend("push eax")
	CodeAppend("push format_out")
	CodeAppend("call printf")
	CodeAppend("add esp, 8")
}

// ─────────────────────────────────────────────
// Assignment
// ─────────────────────────────────────────────

type Assignment struct {
	id            int
	childrenLeft  Node // IdentifierOp
	childrenRight Node // expressão
}

func (n *Assignment) Evaluate(st SymbolTable) Variable {
	name := n.childrenLeft.(*IdentifierOp).value
	val := n.childrenRight.Evaluate(st)
	st.setter(name, val)
	return Variable{value: 0, varType: "int"}
}

func (n *Assignment) Generate(st *SymbolTable) {
	name := n.childrenLeft.(*IdentifierOp).value
	n.childrenRight.Generate(st)
	v := st.getter(name)
	CodeAppend(fmt.Sprintf("mov [ebp-%d], eax", v.shift))
}

// ─────────────────────────────────────────────
// VarDec
// ─────────────────────────────────────────────

type VarDec struct {
	id       int
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

func (n *VarDec) Generate(st *SymbolTable) {
	name := n.children[0].(*IdentifierOp).value
	st.createVariable(name, n.value)
	CodeAppend("sub esp, 4")
	if len(n.children) > 1 {
		n.children[1].Generate(st)
		v := st.getter(name)
		CodeAppend(fmt.Sprintf("mov [ebp-%d], eax", v.shift))
	}
}

// ─────────────────────────────────────────────
// Block
// ─────────────────────────────────────────────

type Block struct {
	id       int
	children []Node
}

func (n *Block) Evaluate(st SymbolTable) Variable {
	for _, child := range n.children {
		child.Evaluate(st)
	}
	return Variable{value: 0, varType: "int"}
}

func (n *Block) Generate(st *SymbolTable) {
	for _, child := range n.children {
		child.Generate(st)
	}
}

// ─────────────────────────────────────────────
// NoOp
// ─────────────────────────────────────────────

type NoOp struct {
	id int
}

func (n *NoOp) Evaluate(st SymbolTable) Variable {
	return Variable{value: 0, varType: "int"}
}

func (n *NoOp) Generate(st *SymbolTable) {
	// Nenhuma instrução necessária
}

// ─────────────────────────────────────────────
// If
// ─────────────────────────────────────────────

type If struct {
	id          int
	condition   Node
	trueBranch  Node
	falseBranch Node
}

func (n *If) Evaluate(st SymbolTable) Variable {
	cond := n.condition.Evaluate(st)
	if cond.varType != "bool" {
		panic("[Semantic] condição do if deve ser bool; recebeu " + cond.varType)
	}
	if toInt(cond) != 0 {
		return n.trueBranch.Evaluate(st)
	} else if n.falseBranch != nil {
		return n.falseBranch.Evaluate(st)
	}
	return Variable{value: 0, varType: "int"}
}

func (n *If) Generate(st *SymbolTable) {
	n.condition.Generate(st)
	CodeAppend("cmp eax, 0")
	if n.falseBranch != nil {
		CodeAppend(fmt.Sprintf("je else_%d", n.id))
		n.trueBranch.Generate(st)
		CodeAppend(fmt.Sprintf("jmp exit_%d", n.id))
		CodeAppend(fmt.Sprintf("else_%d:", n.id))
		n.falseBranch.Generate(st)
		CodeAppend(fmt.Sprintf("exit_%d:", n.id))
	} else {
		CodeAppend(fmt.Sprintf("je exit_%d", n.id))
		n.trueBranch.Generate(st)
		CodeAppend(fmt.Sprintf("exit_%d:", n.id))
	}
}

// ─────────────────────────────────────────────
// While
// ─────────────────────────────────────────────

type While struct {
	id        int
	condition Node
	block     Node
}

func (n *While) Evaluate(st SymbolTable) Variable {
	for {
		cond := n.condition.Evaluate(st)
		if cond.varType != "bool" {
			panic("[Semantic] condição do for deve ser bool; recebeu " + cond.varType)
		}
		if toInt(cond) == 0 {
			break
		}
		n.block.Evaluate(st)
	}
	return Variable{value: 0, varType: "int"}
}

func (n *While) Generate(st *SymbolTable) {
	CodeAppend(fmt.Sprintf("loop_%d:", n.id))
	n.condition.Generate(st)
	CodeAppend("cmp eax, 0")
	CodeAppend(fmt.Sprintf("je exit_%d", n.id))
	n.block.Generate(st)
	CodeAppend(fmt.Sprintf("jmp loop_%d", n.id))
	CodeAppend(fmt.Sprintf("exit_%d:", n.id))
}

// ─────────────────────────────────────────────
// ReadOp  (Scanln — lê um inteiro)
// ─────────────────────────────────────────────

type ReadOp struct {
	id int
}

func (n *ReadOp) Evaluate(st SymbolTable) Variable {
	var val int
	fmt.Scan(&val)
	return Variable{value: val, varType: "int"}
}

func (n *ReadOp) Generate(st *SymbolTable) {
	CodeAppend("push scan_int")
	CodeAppend("push format_in")
	CodeAppend("call scanf")
	CodeAppend("add esp, 8")
	CodeAppend("mov eax, dword [scan_int]")
}

// ─────────────────────────────────────────────
// ReadStringOp  (ScanlnStr — lê uma string)
// ─────────────────────────────────────────────

type ReadStringOp struct {
	id int
}

func (n *ReadStringOp) Evaluate(st SymbolTable) Variable {
	var val string
	fmt.Scan(&val)
	return Variable{value: val, varType: "string"}
}

func (n *ReadStringOp) Generate(st *SymbolTable) {
	// Não será gerado código para strings
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
		return &IntVal{id: newId(), value: num}
	}

	// Literal booleano
	if lexer.next.tipo == "BOOL" {
		val := 0
		if lexer.next.value == "true" {
			val = 1
		}
		lexer.selectNext()
		return &BoolVal{id: newId(), value: val}
	}

	// Literal string
	if lexer.next.tipo == "STR" {
		s := lexer.next.value
		lexer.selectNext()
		return &StringVal{id: newId(), value: s}
	}

	// Operadores unários
	if lexer.next.tipo == "MINUS" || lexer.next.tipo == "PLUS" || lexer.next.tipo == "NOT" {
		op := lexer.next.value
		lexer.selectNext()
		return &UnOp{
			id:       newId(),
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
		return &ReadOp{id: newId()}
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
		return &ReadStringOp{id: newId()}
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
		idVal := &IdentifierOp{id: newId(), value: lexer.next.value}
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
			id:            newId(),
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
			id:            newId(),
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
			id:            newId(),
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
			id:            newId(),
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
			id:            newId(),
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

	block := &Block{id: newId(), children: []Node{}}
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
		return &NoOp{id: newId()}
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
		iden := &IdentifierOp{id: newId(), value: lexer.next.value}
		lexer.selectNext()

		if lexer.next.tipo != "TYPE" {
			panic("[PARSER] Esperava tipo após o identificador em 'var'")
		}
		varType := lexer.next.value
		lexer.selectNext()

		node := &VarDec{
			id:       newId(),
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
		return &PrintOp{id: newId(), children: expr}
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
		return &If{id: newId(), condition: cond, trueBranch: trueBr, falseBranch: falseBr}
	}

	// for <expr> { ... }
	if lexer.next.tipo == "WHILE" {
		lexer.selectNext()
		cond := ParseBoolExpression(lexer)

		if lexer.next.tipo == "END" || lexer.next.tipo != "OPEN_BRA" {
			panic("[PARSER] Esperava '{' após condição do for")
		}

		blk := parseBlock(lexer)
		return &While{id: newId(), condition: cond, block: blk}
	}

	// Atribuição: <iden> = <expr>
	if lexer.next.tipo == "IDEN" {
		idToken := lexer.next
		lexer.selectNext()

		if lexer.next.tipo == "ASSIGN" {
			lexer.selectNext()
			exprRight := ParseBoolExpression(lexer)
			identifier := &IdentifierOp{id: newId(), value: idToken.value}
			assign := &Assignment{
				id:            newId(),
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
		id:       newId(),
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