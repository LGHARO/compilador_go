package main

import (
	"fmt"
	"strconv"
)

// ============================================================================
// Versão 2.3 - Funções e Escopo de Variáveis
// ============================================================================
// Mudanças realizadas:
// 1. Adicionados os nós AST: FuncDec, FuncCall e Return.
// 2. FuncDec armazena n+2 filhos: IdentifierOp (nome), n nós VarDec (parâmetros)
//    e Block (corpo). O tipo de retorno é armazenado em value.
// 3. FuncCall armazena o nome da função em value e os nós de argumentos
//    (expressões) em children.
// 4. Return avalia a expressão filha e sinaliza isReturn: true.
// 5. Block, If e While foram adaptados para propagar o valor de retorno caso
//    um nó Return seja executado (interrompendo a execução do bloco).
// 6. Block, If e While criam uma nova SymbolTable encadeada para seus blocos,
//    garantindo isolamento de escopo (escopo léxico).
// 7. parseFuncDeclaration() adicionado para realizar o parsing de 'func'.
// 8. parseProgram() atualizado para suportar declarações de funções e variáveis globais.
// 9. parseStatement() e ParseFactor() atualizados para suportar chamadas de função (FuncCall)
//    e comando de retorno (Return).
// 10. No Parser.run(), adicionada a chamada à função 'main' (FuncCall) na raiz.
// 11. Removida a restrição que impedia '}' na mesma linha que '{' em parseBlock,
//     permitindo blocos vazios inline como '{{ }}'.
// ============================================================================

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
	return Variable{value: n.value, varType: "int", isFunc: false, isReturn: false}
}

// ─────────────────────────────────────────────
// BoolVal
// ─────────────────────────────────────────────

type BoolVal struct {
	value int // 1 = true, 0 = false
}

func (n *BoolVal) Evaluate(st SymbolTable) Variable {
	return Variable{value: n.value, varType: "bool", isFunc: false, isReturn: false}
}

// ─────────────────────────────────────────────
// StringVal
// ─────────────────────────────────────────────

type StringVal struct {
	value string
}

func (n *StringVal) Evaluate(st SymbolTable) Variable {
	return Variable{value: n.value, varType: "string", isFunc: false, isReturn: false}
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
		return Variable{value: child.value.(int), varType: "int", isFunc: false, isReturn: false}
	case "-":
		if child.varType != "int" {
			panic("[Semantic] operador '-' unário requer int, recebeu " + child.varType)
		}
		return Variable{value: -child.value.(int), varType: "int", isFunc: false, isReturn: false}
	case "!":
		if child.varType != "bool" {
			panic("[Semantic] operador '!' requer bool, recebeu " + child.varType)
		}
		if child.value.(int) == 0 {
			return Variable{value: 1, varType: "bool", isFunc: false, isReturn: false}
		}
		return Variable{value: 0, varType: "bool", isFunc: false, isReturn: false}
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
		if left.varType != "bool" {
			panic("[Semantic] '&&' requer bool nos operandos; recebeu " + left.varType)
		}
		if left.value.(int) == 0 {
			return Variable{value: 0, varType: "bool", isFunc: false, isReturn: false}
		}
		right := n.childrenRight.Evaluate(st)
		if right.varType != "bool" {
			panic("[Semantic] '&&' requer bool nos operandos; recebeu " + right.varType)
		}
		if right.value.(int) != 0 {
			return Variable{value: 1, varType: "bool", isFunc: false, isReturn: false}
		}
		return Variable{value: 0, varType: "bool", isFunc: false, isReturn: false}
	}

	if n.value == "||" {
		if left.varType != "bool" {
			panic("[Semantic] '||' requer bool nos operandos; recebeu " + left.varType)
		}
		if left.value.(int) != 0 {
			return Variable{value: 1, varType: "bool", isFunc: false, isReturn: false}
		}
		right := n.childrenRight.Evaluate(st)
		if right.varType != "bool" {
			panic("[Semantic] '||' requer bool nos operandos; recebeu " + right.varType)
		}
		if right.value.(int) != 0 {
			return Variable{value: 1, varType: "bool", isFunc: false, isReturn: false}
		}
		return Variable{value: 0, varType: "bool", isFunc: false, isReturn: false}
	}

	right := n.childrenRight.Evaluate(st)

	switch n.value {
	case "+":
		// Se pelo menos um operando for string, converte o outro e concatena
		if left.varType == "string" || right.varType == "string" {
			return Variable{value: varToString(left) + varToString(right), varType: "string", isFunc: false, isReturn: false}
		}
		// Soma de dois ints
		if left.varType == "int" && right.varType == "int" {
			return Variable{value: left.value.(int) + right.value.(int), varType: "int", isFunc: false, isReturn: false}
		}
		panic("[Semantic] '+' requer dois int, duas strings, ou string com qualquer tipo; recebeu " + left.varType + " e " + right.varType)
	case "-":
		requireBothInt(left, right, "-")
		return Variable{value: left.value.(int) - right.value.(int), varType: "int", isFunc: false, isReturn: false}
	case "*":
		requireBothInt(left, right, "*")
		return Variable{value: left.value.(int) * right.value.(int), varType: "int", isFunc: false, isReturn: false}
	case "/":
		requireBothInt(left, right, "/")
		if right.value.(int) == 0 {
			panic("[Semantic] divisão por zero")
		}
		return Variable{value: left.value.(int) / right.value.(int), varType: "int", isFunc: false, isReturn: false}
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
		return Variable{value: result, varType: "bool", isFunc: false, isReturn: false}
	case "<":
		// Comparação de int ou string (lexicográfica)
		if left.varType == "string" && right.varType == "string" {
			result := 0
			if left.value.(string) < right.value.(string) {
				result = 1
			}
			return Variable{value: result, varType: "bool", isFunc: false, isReturn: false}
		}
		requireBothInt(left, right, "<")
		result := 0
		if left.value.(int) < right.value.(int) {
			result = 1
		}
		return Variable{value: result, varType: "bool", isFunc: false, isReturn: false}
	case ">":
		// Comparação de int ou string (lexicográfica)
		if left.varType == "string" && right.varType == "string" {
			result := 0
			if left.value.(string) > right.value.(string) {
				result = 1
			}
			return Variable{value: result, varType: "bool", isFunc: false, isReturn: false}
		}
		requireBothInt(left, right, ">")
		result := 0
		if left.value.(int) > right.value.(int) {
			result = 1
		}
		return Variable{value: result, varType: "bool", isFunc: false, isReturn: false}
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
	case "void":
		panic("[Semantic] nao e possivel imprimir expressao do tipo void")
	}
	return Variable{value: 0, varType: "int", isFunc: false, isReturn: false}
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
	if val.varType == "void" {
		panic("[Semantic] nao e possivel atribuir expressao do tipo void para: " + name)
	}
	st.setter(name, val)
	return Variable{value: 0, varType: "int", isFunc: false, isReturn: false}
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
		if val.varType == "void" {
			panic("[Semantic] nao e possivel inicializar variavel com tipo void: " + name)
		}
		st.setter(name, val)
	}
	return Variable{value: 0, varType: "int", isFunc: false, isReturn: false}
}

// ─────────────────────────────────────────────
// Block
// ─────────────────────────────────────────────

type Block struct {
	children []Node
}

func (n *Block) Evaluate(st SymbolTable) Variable {
	for _, child := range n.children {
		var res Variable
		if b, ok := child.(*Block); ok {
			// Cria uma nova SymbolTable encadeada para o bloco interno
			newSt := SymbolTable{
				table:  make(map[string]Variable),
				parent: &st,
			}
			res = b.Evaluate(newSt)
		} else {
			res = child.Evaluate(st)
		}
		if res.isReturn {
			return res
		}
	}
	return Variable{value: 0, varType: "int", isFunc: false, isReturn: false}
}

// ─────────────────────────────────────────────
// NoOp
// ─────────────────────────────────────────────

type NoOp struct{}

func (n *NoOp) Evaluate(st SymbolTable) Variable {
	return Variable{value: 0, varType: "int", isFunc: false, isReturn: false}
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
	if cond.varType != "bool" {
		panic("[Semantic] condição do if deve ser bool; recebeu " + cond.varType)
	}
	if toInt(cond) != 0 {
		newSt := SymbolTable{
			table:  make(map[string]Variable),
			parent: &st,
		}
		res := n.trueBranch.Evaluate(newSt)
		if res.isReturn {
			return res
		}
	} else if n.falseBranch != nil {
		var res Variable
		if _, isIf := n.falseBranch.(*If); isIf {
			res = n.falseBranch.Evaluate(st)
		} else {
			newSt := SymbolTable{
				table:  make(map[string]Variable),
				parent: &st,
			}
			res = n.falseBranch.Evaluate(newSt)
		}
		if res.isReturn {
			return res
		}
	}
	return Variable{value: 0, varType: "int", isFunc: false, isReturn: false}
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
		if cond.varType != "bool" {
			panic("[Semantic] condição do for deve ser bool; recebeu " + cond.varType)
		}
		if toInt(cond) == 0 {
			break
		}
		newSt := SymbolTable{
			table:  make(map[string]Variable),
			parent: &st,
		}
		res := n.block.Evaluate(newSt)
		if res.isReturn {
			return res
		}
	}
	return Variable{value: 0, varType: "int", isFunc: false, isReturn: false}
}

// ─────────────────────────────────────────────
// ReadOp  (Scanln — lê um inteiro)
// ─────────────────────────────────────────────

type ReadOp struct{}

func (n *ReadOp) Evaluate(st SymbolTable) Variable {
	var val int
	fmt.Scan(&val)
	return Variable{value: val, varType: "int", isFunc: false, isReturn: false}
}

// ─────────────────────────────────────────────
// ReadStringOp  (ScanlnStr — lê uma string)
// ─────────────────────────────────────────────

type ReadStringOp struct{}

func (n *ReadStringOp) Evaluate(st SymbolTable) Variable {
	var val string
	fmt.Scan(&val)
	return Variable{value: val, varType: "string", isFunc: false, isReturn: false}
}

// ─────────────────────────────────────────────
// Return
// ─────────────────────────────────────────────

type Return struct {
	children Node // expressão a ser retornada
}

func (n *Return) Evaluate(st SymbolTable) Variable {
	res := n.children.Evaluate(st)
	res.isReturn = true
	return res
}

// ─────────────────────────────────────────────
// FuncDec
// ─────────────────────────────────────────────

type FuncDec struct {
	value    string // tipo de retorno ("int", "string", "bool" ou "")
	children []Node // [0]=IdentifierOp, [1..n]=VarDec, [n+1]=Block
}

func (n *FuncDec) Evaluate(st SymbolTable) Variable {
	funcName := n.children[0].(*IdentifierOp).value
	st.createFunction(funcName, n.value, n)
	return Variable{value: nil, varType: "void", isFunc: false, isReturn: false}
}

// ─────────────────────────────────────────────
// FuncCall
// ─────────────────────────────────────────────

type FuncCall struct {
	value    string // nome da função
	children []Node // expressões dos argumentos
}

func (n *FuncCall) Evaluate(st SymbolTable) Variable {
	funcVar := st.getter(n.value)
	if !funcVar.isFunc {
		panic("[Semantic] identificador nao e uma funcao: " + n.value)
	}
	funcDec, ok := funcVar.value.(*FuncDec)
	if !ok || funcDec == nil {
		panic("[Semantic] definicao da funcao invalida: " + n.value)
	}

	numParams := len(funcDec.children) - 2
	if len(n.children) != numParams {
		panic(fmt.Sprintf("[Semantic] numero incorreto de argumentos para a funcao '%s': esperava %d, recebeu %d", n.value, numParams, len(n.children)))
	}

	// Avalia argumentos no escopo atual do chamador
	evaluatedArgs := make([]Variable, len(n.children))
	for i, argNode := range n.children {
		evaluatedArgs[i] = argNode.Evaluate(st)
	}

	// Cria a nova SymbolTable para a função apontando para o escopo global (raiz)
	rootSt := st.getRoot()
	funcSt := SymbolTable{
		table:  make(map[string]Variable),
		parent: rootSt,
	}

	// Declara os parâmetros e atribui os valores avaliados
	for i := 0; i < numParams; i++ {
		paramVarDec := funcDec.children[i+1].(*VarDec)
		paramName := paramVarDec.children[0].(*IdentifierOp).value
		paramType := paramVarDec.value
		argVal := evaluatedArgs[i]

		if argVal.varType != paramType {
			panic(fmt.Sprintf("[Semantic] tipo incompativel no argumento '%s': esperava '%s', recebeu '%s'", paramName, paramType, argVal.varType))
		}

		funcSt.createVariable(paramName, paramType)
		funcSt.setter(paramName, argVal)
	}

	// Executa o corpo da função
	body := funcDec.children[len(funcDec.children)-1].(*Block)
	res := body.Evaluate(funcSt)

	declaredReturnType := funcDec.value // "" ou "int", "string", "bool"

	if declaredReturnType == "" || declaredReturnType == "void" {
		if res.isReturn {
			panic("[Semantic] funcao sem tipo de retorno nao pode dar return: " + n.value)
		}
		return Variable{value: 0, varType: "void", isFunc: false, isReturn: false}
	} else {
		if !res.isReturn {
			panic(fmt.Sprintf("[Semantic] funcao '%s' com retorno '%s' nao retornou nenhum valor", n.value, declaredReturnType))
		}
		if res.varType != declaredReturnType {
			panic(fmt.Sprintf("[Semantic] tipo de retorno incompativel na funcao '%s': esperava '%s', recebeu '%s'", n.value, declaredReturnType, res.varType))
		}
		return Variable{value: res.value, varType: res.varType, isFunc: false, isReturn: false}
	}
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

	// Identificador ou Chamada de Função
	if lexer.next.tipo == "IDEN" {
		idName := lexer.next.value
		lexer.selectNext()

		if lexer.next.tipo == "OPEN_PAR" {
			lexer.selectNext()
			args := []Node{}
			if lexer.next.tipo != "CLOSE_PAR" {
				for {
					argExpr := ParseBoolExpression(lexer)
					args = append(args, argExpr)
					if lexer.next.tipo == "COMMA" {
						lexer.selectNext()
						if lexer.next.tipo == "CLOSE_PAR" {
							panic("[PARSER] Esperava expressão após ',' na chamada de função")
						}
					} else if lexer.next.tipo == "CLOSE_PAR" {
						break
					} else {
						panic("[PARSER] Esperava ',' ou ')' na chamada de função")
					}
				}
			}
			if lexer.next.tipo != "CLOSE_PAR" {
				panic("[PARSER] Esperava ')' após argumentos da chamada de função")
			}
			lexer.selectNext()
			return &FuncCall{
				value:    idName,
				children: args,
			}
		}

		return &IdentifierOp{value: idName}
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

	// '}' na mesma linha que '{' é inválido (BLOCK = "{", "\n", { STATEMENT }, "}")
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
// parseFuncDeclaration
// ─────────────────────────────────────────────

func parseFuncDeclaration(lexer *Lexer) *FuncDec {
	if lexer.next.tipo != "FUNC" {
		panic("[PARSER] Esperava 'func'")
	}
	lexer.selectNext()

	if lexer.next.tipo != "IDEN" {
		panic("[PARSER] Esperava identificador após 'func'")
	}
	funcName := &IdentifierOp{value: lexer.next.value}
	lexer.selectNext()

	if lexer.next.tipo != "OPEN_PAR" {
		panic("[PARSER] Esperava '(' após nome da função")
	}
	lexer.selectNext()

	params := []Node{}
	if lexer.next.tipo != "CLOSE_PAR" {
		for {
			if lexer.next.tipo != "IDEN" {
				panic("[PARSER] Esperava identificador de parâmetro")
			}
			paramId := &IdentifierOp{value: lexer.next.value}
			lexer.selectNext()

			if lexer.next.tipo != "TYPE" {
				panic("[PARSER] Esperava tipo após identificador do parâmetro")
			}
			paramType := lexer.next.value
			lexer.selectNext()

			params = append(params, &VarDec{
				value:    paramType,
				children: []Node{paramId},
			})

			if lexer.next.tipo == "COMMA" {
				lexer.selectNext()
				if lexer.next.tipo == "CLOSE_PAR" {
					panic("[PARSER] Esperava parâmetro após ',' na declaração da função")
				}
			} else if lexer.next.tipo == "CLOSE_PAR" {
				break
			} else {
				panic("[PARSER] Esperava ',' ou ')' na lista de parâmetros")
			}
		}
	}

	if lexer.next.tipo != "CLOSE_PAR" {
		panic("[PARSER] Esperava ')' após lista de parâmetros")
	}
	lexer.selectNext()

	returnType := ""
	if lexer.next.tipo == "TYPE" {
		returnType = lexer.next.value
		lexer.selectNext()
	}

	if lexer.next.tipo == "END" {
		panic("[PARSER] Esperava '{' na mesma linha após declaração da função")
	}
	if lexer.next.tipo != "OPEN_BRA" {
		panic("[PARSER] Esperava '{' para iniciar o corpo da função")
	}

	body := parseBlock(lexer)

	children := []Node{funcName}
	for _, p := range params {
		children = append(children, p)
	}
	children = append(children, body)

	return &FuncDec{
		value:    returnType,
		children: children,
	}
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

	// return <expr>
	if lexer.next.tipo == "RETURN" {
		lexer.selectNext()
		expr := ParseBoolExpression(lexer)
		if lexer.next.tipo == "END" {
			lexer.selectNext()
		}
		return &Return{children: expr}
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

	// Atribuição ou Chamada de Função
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

		if lexer.next.tipo == "OPEN_PAR" {
			lexer.selectNext()
			args := []Node{}
			if lexer.next.tipo != "CLOSE_PAR" {
				for {
					argExpr := ParseBoolExpression(lexer)
					args = append(args, argExpr)
					if lexer.next.tipo == "COMMA" {
						lexer.selectNext()
						if lexer.next.tipo == "CLOSE_PAR" {
							panic("[PARSER] Esperava expressão após ',' na chamada de função")
						}
					} else if lexer.next.tipo == "CLOSE_PAR" {
						break
					} else {
						panic("[PARSER] Esperava ',' ou ')' na chamada de função")
					}
				}
			}
			if lexer.next.tipo != "CLOSE_PAR" {
				panic("[PARSER] Esperava ')' após argumentos da chamada de função")
			}
			lexer.selectNext()
			if lexer.next.tipo == "END" {
				lexer.selectNext()
			}
			return &FuncCall{
				value:    idToken.value,
				children: args,
			}
		}

		panic("[PARSER] Esperava '=' ou '(' após identificador")
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
		if lexer.next.tipo == "FUNC" {
			fn := parseFuncDeclaration(lexer)
			block.children = append(block.children, fn)
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

	// Adiciona a chamada da função main() na raiz do programa
	mainCall := &FuncCall{
		value:    "main",
		children: []Node{},
	}
	AST.children = append(AST.children, mainCall)

	return AST
}