package main

// só serve como um wrapper
type Variable struct {
	value   int
	isConst bool
}

// vai guardar as variaveis e seus valores aqui
type SymbolTable struct {
	table map[string]Variable
}

// getter
func (symbletable *SymbolTable) getter(variavel string) int {
	info, existe := symbletable.table[variavel]

	if !existe {
		panic("[Semantic] variavel nao existe no dicionario")
	} else {
		return info.value
	}
}

// setter
func (symbletable *SymbolTable) setter(variavel string, valor int) {
	info, existe := symbletable.table[variavel]
	if existe && info.isConst {
		panic("[Semantic] cannot change the value of " + variavel)
	}

	symbletable.table[variavel] = Variable{
		value:   valor,
		isConst: false,
	}
}

// createConst
func (symbletable *SymbolTable) createConst(variavel string, valor int) {
	info, existe := symbletable.table[variavel]
	if existe && info.isConst {
		panic("[Semantic] cannot change the value of " + variavel)
	}

	symbletable.table[variavel] = Variable{
		value:   valor,
		isConst: true,
	}
}