package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
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

func runSource(code string) {
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
		table: make(map[string]Variable),
	}
	AST.Evaluate(*st)
}

func TestDefineConstant(t *testing.T) {
	code := `
define N = 1      // Constant N
const x = 1        // Immutable variable x
Println(N)      // This is actually: print(1)
Println(x)
`
	out := captureOutput(func() {
		runSource(code)
	})

	expected := "1\n1"
	if out != expected {
		t.Fatalf("Esperava saída:\n%s\nObteve:\n%s", expected, out)
	}
}

func TestParserErrorOnConstantAssignment(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava Parser Error em N = 1")
		}
		errMsg := fmt.Sprint(r)
		if !strings.Contains(errMsg, "[PARSER]") {
			t.Fatalf("Esperava erro contendo [PARSER], mas obteve: %s", errMsg)
		}
	}()

	code := `
define N = 1
N = 1
`
	runSource(code)
}

func TestSemanticErrorOnImmutableVariable(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava Semantic Error em x = 2")
		}
		errMsg := fmt.Sprint(r)
		if !strings.Contains(errMsg, "cannot change the value of x") {
			t.Fatalf("Esperava erro contendo 'cannot change the value of x', mas obteve: %s", errMsg)
		}
	}()

	code := `
const x = 1
x = 2
`
	runSource(code)
}

func TestConstWithoutInitialization(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava Parser Error para const sem inicialização")
		}
		errMsg := fmt.Sprint(r)
		if !strings.Contains(errMsg, "[PARSER]") {
			t.Fatalf("Esperava erro [PARSER], mas obteve: %s", errMsg)
		}
	}()

	code := `
const x
`
	runSource(code)
}

func TestMultipleDefinesAndConstExpression(t *testing.T) {
	code := `
define A = 10
define B = 20
const c = A + B
Println(c)
`
	out := captureOutput(func() {
		runSource(code)
	})

	expected := "30"
	if out != expected {
		t.Fatalf("Esperava %s, mas obteve: %s", expected, out)
	}
}
