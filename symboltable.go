package main

import (
	"fmt"
	"strings"
)

// Variable armazena um valor com seu tipo correspondente, flag de função e flag de return
type Variable struct {
	value    interface{} // int, string, bool, referência ao nó *FuncDec ou *StructInstance
	varType  string      // "int", "string", "bool", tipo de retorno da função ou nome da struct
	isFunc   bool        // true se for função
	isReturn bool        // true se for valor originado de um return statement
}

// StructTypeDef armazena a definição (molde) de uma struct
type StructTypeDef struct {
	Name   string
	Fields map[string]string // nomeCampo -> tipoCampo
	Order  []string          // ordem de declaração dos campos
}

// StructInstance representa um objeto concreto instanciado a partir do molde da struct
type StructInstance struct {
	TypeName string
	Fields   map[string]Variable
}

// Clone cria uma cópia profunda da instância da struct (garantindo semântica de valor em atribuições)
func (s *StructInstance) Clone() *StructInstance {
	if s == nil {
		return nil
	}
	clone := &StructInstance{
		TypeName: s.TypeName,
		Fields:   make(map[string]Variable),
	}
	for k, v := range s.Fields {
		if subInst, ok := v.value.(*StructInstance); ok && subInst != nil {
			clone.Fields[k] = Variable{
				value:    subInst.Clone(),
				varType:  v.varType,
				isFunc:   v.isFunc,
				isReturn: v.isReturn,
			}
		} else {
			clone.Fields[k] = v
		}
	}
	return clone
}

// String retorna uma representação legível da struct
func (s *StructInstance) String() string {
	if s == nil {
		return "{}"
	}
	var parts []string
	for _, v := range s.Fields {
		parts = append(parts, fmt.Sprintf("%v", v.value))
	}
	return "{" + strings.Join(parts, " ") + "}"
}

// SymbolTable guarda as variáveis em escopos encadeados e tipos de structs
type SymbolTable struct {
	table   map[string]Variable
	parent  *SymbolTable
	structs map[string]*StructTypeDef
}

// NewSymbolTable cria uma nova tabela de símbolos encadeada ao pai fornecido
func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		table:   make(map[string]Variable),
		parent:  parent,
		structs: make(map[string]*StructTypeDef),
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

// createStructType registra a definição de uma nova struct no escopo raiz
func (st *SymbolTable) createStructType(name string, def *StructTypeDef) {
	root := st.getRoot()
	if root.isStructType(name) {
		panic("[Semantic] struct ja declarada: " + name)
	}
	if root.structs != nil {
		root.structs[name] = def
	}
	if root.table != nil {
		root.table["::struct::"+name] = Variable{
			value:    def,
			varType:  "struct_def",
			isFunc:   false,
			isReturn: false,
		}
	}
}

// getStructType busca a definição da struct no escopo raiz
func (st *SymbolTable) getStructType(name string) *StructTypeDef {
	root := st.getRoot()
	if root.structs != nil {
		if def, exists := root.structs[name]; exists && def != nil {
			return def
		}
	}
	if root.table != nil {
		if v, exists := root.table["::struct::"+name]; exists {
			if def, ok := v.value.(*StructTypeDef); ok && def != nil {
				return def
			}
		}
	}
	panic("[Semantic] tipo struct desconhecido: " + name)
}

// isStructType verifica se o tipo fornecido corresponde a uma struct registrada
func (st *SymbolTable) isStructType(name string) bool {
	root := st.getRoot()
	if root.structs != nil {
		if _, exists := root.structs[name]; exists {
			return true
		}
	}
	if root.table != nil {
		if _, exists := root.table["::struct::"+name]; exists {
			return true
		}
	}
	return false
}

// instantiateStruct cria uma nova instância de struct com todos os atributos inicializados recursivamente com seus valores padrão
func (st *SymbolTable) instantiateStruct(typeName string, visiting map[string]bool) *StructInstance {
	if visiting == nil {
		visiting = make(map[string]bool)
	}
	if visiting[typeName] {
		panic("[Semantic] ciclo ou recursao direta invalida na struct: " + typeName)
	}
	visiting[typeName] = true

	def := st.getStructType(typeName)
	instance := &StructInstance{
		TypeName: typeName,
		Fields:   make(map[string]Variable),
	}

	for _, fieldName := range def.Order {
		fieldType := def.Fields[fieldName]
		var fieldZero interface{}
		switch fieldType {
		case "int":
			fieldZero = 0
		case "string":
			fieldZero = ""
		case "bool":
			fieldZero = 0
		default:
			if st.isStructType(fieldType) {
				fieldZero = st.instantiateStruct(fieldType, visiting)
			} else {
				panic("[Semantic] tipo desconhecido para campo '" + fieldName + "': " + fieldType)
			}
		}
		instance.Fields[fieldName] = Variable{
			value:    fieldZero,
			varType:  fieldType,
			isFunc:   false,
			isReturn: false,
		}
	}

	delete(visiting, typeName)
	return instance
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
		if st.isStructType(tipo) {
			zeroVal = st.instantiateStruct(tipo, nil)
		} else {
			panic("[Semantic] tipo desconhecido: " + tipo)
		}
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
		var newVal interface{} = val.value
		if inst, ok := val.value.(*StructInstance); ok && inst != nil {
			newVal = inst.Clone()
		}
		st.table[nome] = Variable{value: newVal, varType: existente.varType, isFunc: false, isReturn: false}
		return
	}
	if st.parent != nil {
		st.parent.setter(nome, val)
		return
	}
	panic("[Semantic] variavel nao declarada: " + nome)
}