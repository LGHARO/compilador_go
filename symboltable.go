package main

// ============================================================================
// Roteiro 8 - Geração de Código x86 (NASM)
// ============================================================================
// Mudanças realizadas:
// 1. Adicionado o atributo shift na struct Variable para armazenar o deslocamento
//    relativo a EBP na pilha (ex: 4, 8, 12...).
// 2. Adicionado o atributo currentShift na struct SymbolTable para acumular o
//    deslocamento das variáveis declaradas.
// 3. Atualizado o método createVariable() para incrementar currentShift em 4 bytes
//    para cada nova variável inteira/booleana declarada e associar esse shift.
// 4. Atualizado o método setter() para preservar o shift original da variável
//    ao atualizar seu valor.
// ============================================================================

// Variable armazena um valor com seu tipo correspondente e seu deslocamento na pilha
type Variable struct {
	value   interface{} // int, string ou bool (representado como int: 0=false, 1=true)
	varType string      // "int", "string" ou "bool"
	shift   int         // deslocamento relativo ao EBP em bytes (ex: 4 para [EBP-4])
}

// SymbolTable guarda as variáveis, seus valores e deslocamentos na pilha
type SymbolTable struct {
	table        map[string]Variable
	currentShift int // deslocamento acumulado na pilha
}

// getter retorna a Variable completa (com valor, tipo e shift)
func (st *SymbolTable) getter(variavel string) Variable {
	info, existe := st.table[variavel]
	if !existe {
		panic("[Semantic] variavel nao declarada: " + variavel)
	}
	return info
}

// createVariable declara uma nova variável com o tipo dado e aloca 4 bytes de deslocamento na pilha
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
	st.currentShift += 4
	st.table[nome] = Variable{value: zeroVal, varType: tipo, shift: st.currentShift}
}

// setter atribui um valor a uma variável já declarada, verificando tipos e preservando o shift
func (st *SymbolTable) setter(nome string, val Variable) {
	existente, existe := st.table[nome]
	if !existe {
		panic("[Semantic] variavel nao declarada: " + nome)
	}
	if existente.varType != val.varType {
		panic("[Semantic] tipo incompativel: esperava '" + existente.varType + "' mas recebeu '" + val.varType + "'")
	}
	st.table[nome] = Variable{value: val.value, varType: existente.varType, shift: existente.shift}
}