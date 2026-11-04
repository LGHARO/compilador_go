package main

import "fmt"

// Variable armazena um valor com seu tipo correspondente, deslocamento na pilha e flags de função/parâmetro
type Variable struct {
	value    interface{} // int, string, bool ou referência ao nó *FuncDec
	varType  string      // "int", "string", "bool" ou tipo de retorno da função
	shift    int         // deslocamento relativo ao EBP em bytes (ex: 4 para [ebp-4] ou 8 para [ebp+8])
	isParam  bool        // true se for parâmetro de função ([ebp+shift])
	isGlobal bool        // true se for variável global ([_var_name])
	isFunc   bool        // true se for função
	isReturn bool        // true se for valor originado de um return statement
}

// Operand retorna a string do operando Assembly NASM correspondente à variável
func (v Variable) Operand(name string) string {
	if v.isGlobal {
		return fmt.Sprintf("[_var_%s]", name)
	}
	if v.isParam {
		return fmt.Sprintf("[ebp+%d]", v.shift)
	}
	return fmt.Sprintf("[ebp-%d]", v.shift)
}

// SymbolTable guarda as variáveis, funções e deslocamentos na pilha em escopos encadeados
type SymbolTable struct {
	table        map[string]Variable
	parent       *SymbolTable
	currentShift int // deslocamento acumulado na pilha para variáveis locais
}

// getRoot retorna o escopo raiz (global)
func (st *SymbolTable) getRoot() *SymbolTable {
	curr := st
	for curr.parent != nil {
		curr = curr.parent
	}
	return curr
}

// getter retorna a Variable buscando no escopo atual e recursivamente nos pais
func (st *SymbolTable) getter(variavel string) Variable {
	info, existe := st.table[variavel]
	if existe {
		return info
	}
	if st.parent != nil {
		return st.parent.getter(variavel)
	}
	panic("[Semantic] variavel nao declarada: " + variavel)
}

// createVariable declara uma nova variável no escopo atual; panics se já existir no escopo atual
func (st *SymbolTable) createVariable(nome string, tipo string) {
	if _, existe := st.table[nome]; existe {
		panic("[Semantic] variavel ja declarada: " + nome)
	}
	var zeroVal interface{}
	switch tipo {
	case "int":
		zeroVal = 0
	case "string":
		zeroVal = ""
	case "bool":
		zeroVal = 0 // false
	default:
		panic("[Semantic] tipo desconhecido: " + tipo)
	}

	isGlobal := (st.parent == nil && hasProgramFunctions)
	shift := 0
	if !isGlobal {
		st.currentShift += 4
		shift = st.currentShift
	} else {
		CodeAddGlobalVar(nome)
	}

	st.table[nome] = Variable{
		value:    zeroVal,
		varType:  tipo,
		shift:    shift,
		isParam:  false,
		isGlobal: isGlobal,
		isFunc:   false,
		isReturn: false,
	}
}

// createParam declara um parâmetro com deslocamento positivo na pilha ([ebp+offset])
func (st *SymbolTable) createParam(nome string, tipo string, offset int) {
	if _, existe := st.table[nome]; existe {
		panic("[Semantic] parametro ja declarado: " + nome)
	}
	st.table[nome] = Variable{
		value:    0,
		varType:  tipo,
		shift:    offset,
		isParam:  true,
		isGlobal: false,
		isFunc:   false,
		isReturn: false,
	}
}

// createFunction declara uma função na tabela raiz (global)
func (st *SymbolTable) createFunction(nome string, returnType string, funcNode interface{}) {
	root := st.getRoot()
	if _, existe := root.table[nome]; existe {
		panic("[Semantic] funcao ja declarada: " + nome)
	}
	root.table[nome] = Variable{
		value:    funcNode,
		varType:  returnType,
		shift:    0,
		isParam:  false,
		isGlobal: false,
		isFunc:   true,
		isReturn: false,
	}
}

// setter atribui um valor a uma variável já declarada, buscando no escopo onde foi definida
func (st *SymbolTable) setter(nome string, val Variable) {
	existente, existe := st.table[nome]
	if existe {
		if existente.isFunc {
			panic("[Semantic] nao e possivel reatribuir valor a uma funcao: " + nome)
		}
		if existente.varType != val.varType {
			panic("[Semantic] tipo incompativel: esperava '" + existente.varType + "' mas recebeu '" + val.varType + "'")
		}
		st.table[nome] = Variable{
			value:    val.value,
			varType:  existente.varType,
			shift:    existente.shift,
			isParam:  existente.isParam,
			isGlobal: existente.isGlobal,
			isFunc:   false,
			isReturn: false,
		}
		return
	}
	if st.parent != nil {
		st.parent.setter(nome, val)
		return
	}
	panic("[Semantic] variavel nao declarada: " + nome)
}