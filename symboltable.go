package main

// ============================================================================
// Versão 2.3 - Funções e Escopo de Variáveis
// ============================================================================
// Mudanças realizadas:
// 1. Variable: adicionados os campos isFunc (sinaliza declaração de função)
//    e isReturn (usado para propagação e captura de comandos return).
// 2. SymbolTable: adicionado o atributo parent (*SymbolTable) para permitir
//    encadeamento de escopos em lista ligada reversa.
// 3. getter(variavel string): busca recursivamente a variável no escopo atual
//    e nos escopos pais (parent).
// 4. setter(nome string, val Variable): busca recursivamente e atualiza o valor
//    no escopo onde a variável foi declarada, verificando compatibilidade de tipos.
// 5. createVariable(nome string, tipo string): cria a variável no escopo atual
//    sinalizando isFunc: false.
// 6. createFunction(nome string, returnType string, funcNode interface{}):
//    registra a função no escopo global (raiz da cadeia de SymbolTable) com
//    isFunc: true.
// ============================================================================

// Variable armazena um valor com seu tipo correspondente, flag de função e flag de return
type Variable struct {
	value    interface{} // int, string, bool ou referência ao nó *FuncDec
	varType  string      // "int", "string", "bool" ou tipo de retorno da função
	isFunc   bool        // true se for função
	isReturn bool        // true se for valor originado de um return statement
}

// SymbolTable guarda as variáveis em escopos encadeados
type SymbolTable struct {
	table  map[string]Variable
	parent *SymbolTable
}

// NewSymbolTable cria uma nova tabela de símbolos encadeada ao pai fornecido
func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		table:  make(map[string]Variable),
		parent: parent,
	}
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
	st.table[nome] = Variable{value: zeroVal, varType: tipo, isFunc: false, isReturn: false}
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
		st.table[nome] = Variable{value: val.value, varType: existente.varType, isFunc: false, isReturn: false}
		return
	}
	if st.parent != nil {
		st.parent.setter(nome, val)
		return
	}
	panic("[Semantic] variavel nao declarada: " + nome)
}