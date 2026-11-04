package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutputStr executa uma função e captura sua saída padrão
func captureOutputStr(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return strings.TrimSpace(buf.String())
}

// runCodeStr executa uma string de código no compilador
func runCodeStr(code string) {
	filter := &Prepro{}
	dadosFiltrados := filter.filter(code)

	l := &Lexer{
		source:   dadosFiltrados,
		position: 0,
	}

	p := &Parser{
		lexer: *l,
	}

	AST := p.run()

	st := &SymbolTable{
		table:  make(map[string]Variable),
		parent: nil,
	}
	AST.Evaluate(*st)
}

// TestOfficialPromptExample verifica o exemplo fornecido no enunciado
func TestOfficialPromptExample(t *testing.T) {
	code := `
type A struct {
  var b int
}

func main() {
  var x A // Creates an instance of A
  x.b = 1 // Assigns 1 to x.b
  Println(x.b) // Prints 1
}
`
	out := captureOutputStr(func() {
		runCodeStr(code)
	})

	expected := "1"
	if out != expected {
		t.Fatalf("Esperava saída '%s', mas obteve '%s'", expected, out)
	}
}

// TestUserTest1ValidInput verifica o Test 1 fornecido pelo usuário (Student com string, int e bool)
func TestUserTest1ValidInput(t *testing.T) {
	code := `
// Struct implementation

type Student struct {
  var name string
  var age int
  var grade bool
}

func main() {
  var student1 Student
  student1.name = "John Doe"
  student1.age = 20
  student1.grade = true

  student1.age = student1.age + 1
  Println(student1.name)
  Println(student1.age)
  Println(student1.grade)

}
`
	out := captureOutputStr(func() {
		runCodeStr(code)
	})

	expected := "John Doe\n21\ntrue"
	if out != expected {
		t.Fatalf("Esperava saída '%s', mas obteve '%s'", expected, out)
	}
}

// TestNestedStructs verifica structs aninhadas (Line -> Point -> x, y)
func TestNestedStructs(t *testing.T) {
	code := `
type Point struct {
  var x int
  var y int
}

type Line struct {
  var p1 Point
  var p2 Point
}

func main() {
  var l Line
  l.p1.x = 10
  l.p1.y = 20
  l.p2.x = 30
  l.p2.y = 40
  Println(l.p1.x)
  Println(l.p1.y)
  Println(l.p2.x)
  Println(l.p2.y)
}
`
	out := captureOutputStr(func() {
		runCodeStr(code)
	})

	expected := "10\n20\n30\n40"
	if out != expected {
		t.Fatalf("Esperava saída:\n%s\nMas obteve:\n%s", expected, out)
	}
}

// TestDeeplyNestedStructs verifica múltiplos níveis de aninhamento
func TestDeeplyNestedStructs(t *testing.T) {
	code := `
type Inner struct {
  var val int
}

type Mid struct {
  var inn Inner
}

type Outer struct {
  var m Mid
}

func main() {
  var o Outer
  o.m.inn.val = 42
  Println(o.m.inn.val)
}
`
	out := captureOutputStr(func() {
		runCodeStr(code)
	})

	expected := "42"
	if out != expected {
		t.Fatalf("Esperava saída '%s', mas obteve '%s'", expected, out)
	}
}

// TestStructMultipleTypes testa struct com tipos int, string e bool
func TestStructMultipleTypes(t *testing.T) {
	code := `
type Person struct {
  var name string
  var age int
  var active bool
}

func main() {
  var p Person
  p.name = "Lucas"
  p.age = 25
  p.active = true
  Println(p.name)
  Println(p.age)
  Println(p.active)
}
`
	out := captureOutputStr(func() {
		runCodeStr(code)
	})

	expected := "Lucas\n25\ntrue"
	if out != expected {
		t.Fatalf("Esperava saída:\n%s\nMas obteve:\n%s", expected, out)
	}
}

// TestStructValueSemantics testa clonagem ao atribuir variável struct
func TestStructValueSemantics(t *testing.T) {
	code := `
type A struct {
  var b int
}

func main() {
  var x A
  x.b = 10
  var y A
  y = x
  x.b = 20
  Println(x.b)
  Println(y.b)
}
`
	out := captureOutputStr(func() {
		runCodeStr(code)
	})

	expected := "20\n10"
	if out != expected {
		t.Fatalf("Esperava saída:\n%s\nMas obteve:\n%s", expected, out)
	}
}

// TestStructFunctionArgAndReturn testa passar e retornar struct de funções
func TestStructFunctionArgAndReturn(t *testing.T) {
	code := `
type Point struct {
  var x int
  var y int
}

func makePoint(x int, y int) Point {
  var p Point
  p.x = x
  p.y = y
  return p
}

func sumCoords(p Point) int {
  return p.x + p.y
}

func main() {
  var pt Point
  pt = makePoint(3, 4)
  Println(sumCoords(pt))
}
`
	out := captureOutputStr(func() {
		runCodeStr(code)
	})

	expected := "7"
	if out != expected {
		t.Fatalf("Esperava saída '%s', mas obteve '%s'", expected, out)
	}
}

// TestStructInExpression testa uso de campos em expressões aritméticas
func TestStructInExpression(t *testing.T) {
	code := `
type Box struct {
  var w int
  var h int
}

func main() {
  var b Box
  b.w = 5
  b.h = 10
  var area int
  area = b.w * b.h
  Println(area)
}
`
	out := captureOutputStr(func() {
		runCodeStr(code)
	})

	expected := "50"
	if out != expected {
		t.Fatalf("Esperava saída '%s', mas obteve '%s'", expected, out)
	}
}

// TestErrorFieldTypeMismatch testa erro ao atribuir tipo incompatível a campo
func TestErrorFieldTypeMismatch(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por incompatibilidade de tipo no campo")
		}
	}()

	code := `
type A struct {
  var b int
}

func main() {
  var x A
  x.b = "texto"
}
`
	runCodeStr(code)
}

// TestErrorAccessNonExistentField testa erro ao acessar campo inexistente
func TestErrorAccessNonExistentField(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic ao acessar campo inexistente")
		}
	}()

	code := `
type A struct {
  var b int
}

func main() {
  var x A
  Println(x.c)
}
`
	runCodeStr(code)
}

// TestErrorAssignNonExistentField testa erro ao atribuir a campo inexistente
func TestErrorAssignNonExistentField(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic ao atribuir a campo inexistente")
		}
	}()

	code := `
type A struct {
  var b int
}

func main() {
  var x A
  x.c = 10
}
`
	runCodeStr(code)
}

// TestErrorFieldOnPrimitive testa erro ao tentar acessar campo em tipo primitivo
func TestErrorFieldOnPrimitive(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic ao tentar acessar campo em variável primitiva")
		}
	}()

	code := `
func main() {
  var n int
  n.x = 1
}
`
	runCodeStr(code)
}

// TestErrorFieldInitializationInDecl testa erro de inicialização dentro da struct
func TestErrorFieldInitializationInDecl(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por inicialização de campo na declaração da struct")
		}
	}()

	code := `
type A struct {
  var b int = 1
}

func main() {
}
`
	runCodeStr(code)
}

// TestErrorDuplicateField testa erro de campo duplicado na struct
func TestErrorDuplicateField(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por campo duplicado na struct")
		}
	}()

	code := `
type A struct {
  var b int
  var b string
}

func main() {
}
`
	runCodeStr(code)
}

// TestErrorDuplicateStruct testa erro de struct duplicada
func TestErrorDuplicateStruct(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por struct duplicada")
		}
	}()

	code := `
type A struct {
  var b int
}

type A struct {
  var c int
}

func main() {
}
`
	runCodeStr(code)
}

// TestErrorStructTypeMismatchInAssignment testa erro ao atribuir struct de tipo diferente
func TestErrorStructTypeMismatchInAssignment(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por incompatibilidade de tipos de struct")
		}
	}()

	code := `
type A struct {
  var b int
}

type B struct {
  var b int
}

func main() {
  var x A
  var y B
  x = y
}
`
	runCodeStr(code)
}
