package main

// Variable armazena um valor com seu tipo correspondente
type Variable struct {
	value   interface{} // int, string ou bool (representado como int: 0=false, 1=true)
	varType string      // "int", "string" ou "bool"
}

// SymbolTable guarda as variáveis e seus valores
type SymbolTable struct {
	table map[string]Variable
}

// getter retorna a Variable completa (com valor e tipo)
func (st *SymbolTable) getter(variavel string) Variable {
	info, existe := st.table[variavel]
	if !existe {
		panic("[Semantic] variavel nao declarada: " + variavel)
	}
	return info
}

// createVariable declara uma nova variável com o tipo dado; panics se já existir
func (st *SymbolTable) createVariable(nome string, tipo string) {
	if _, existe := st.table[nome]; existe {
		panic("[Semantic] variavel ja declarada: " + nome)
	}
	var zeroVal interface{}
	switch tipo {
	case "int":
		zeroVal = 0
	case "float":
		zeroVal = 0.0
	case "string":
		zeroVal = ""
	case "bool":
		zeroVal = 0 // false
	default:
		panic("[Semantic] tipo desconhecido: " + tipo)
	}
	st.table[nome] = Variable{value: zeroVal, varType: tipo}
}

// setter atribui um valor a uma variável já declarada; verifica tipos
func (st *SymbolTable) setter(nome string, val Variable) {
	existente, existe := st.table[nome]
	if !existe {
		panic("[Semantic] variavel nao declarada: " + nome)
	}
	if existente.varType != val.varType {
		panic("[Semantic] tipo incompativel: esperava '" + existente.varType + "' mas recebeu '" + val.varType + "'")
	}
	st.table[nome] = Variable{value: val.value, varType: existente.varType}
}