package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func runCode(code string) (string, error) {
	filter := &Prepro{}
	dados_filtrados := filter.filter(code)

	l := &Lexer{
		source:   dados_filtrados,
		position: 0,
	}

	p := &Parser{
		lexer: *l,
	}

	AST := p.run()

	st := &SymbolTable{
		table: make(map[string]Variable),
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	AST.Evaluate(*st)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return strings.TrimSpace(buf.String()), nil
}

func TestExtraV22UserExample(t *testing.T) {
	code := `var y float = 1.6
var x int = (int) y
y = x + y
Println(x)
Println(y)
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "2\n3.6"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV22Arithmetic(t *testing.T) {
	code := `var a int = 2
var b float = 3.5
Println(a + b)
Println(b - a)
Println(a * b)
Println(7.0 / 2)
Println(7 / 2)
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "5.5\n1.5\n7\n3.5\n3"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV22Comparisons(t *testing.T) {
	code := `Println(2.0 == 2)
Println(2.5 > 2)
Println(2.5 < 2)
Println(2 < 2.5)
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "true\ntrue\nfalse\ntrue"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV22Casting(t *testing.T) {
	code := `Println((int) 1.6)
Println((int) 1.4)
Println((int) -1.6)
Println((int) -1.4)
Println((string) 3.6)
Println((int) "10")
Println((float) "2.5")
Println((bool) 1.6)
Println((bool) 0.0)
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "2\n1\n-2\n-1\n3.6\n10\n2.5\ntrue\nfalse"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV22Precedence(t *testing.T) {
	code := `var y float = 1.6
Println((int) y * 10)
Println(10 * (int) y)
Println(-(int) y)
Println((int) -y)
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "20\n20\n-2\n-2"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV22NestedCast(t *testing.T) {
	code := `Println((int) (float) (int) 3.8)
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "4"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}
