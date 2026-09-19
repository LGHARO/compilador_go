package main

// só serve como um wrapper
type  Variable struct {
	value  int
}


//vai guardar as variaveis e seus valores aqui
type SymbolTable struct {
	table   map[string]Variable
}

//getter
func (symbletable *SymbolTable) getter(variavel string) int{
	info, existe := symbletable.table[variavel]

	if !existe{
		panic("[Semantic] variavel nao existe no dicionario")
	}else{
		return info.value
	}

	
}

//setter
func (symbletable *SymbolTable) setter(variavel string, valor int){
	// cria o wrapper e adiciona o valor
	varCriada := &Variable{
		value: valor,
	}
	
	//adiciona na table
	symbletable.table[variavel] = *varCriada
	
}