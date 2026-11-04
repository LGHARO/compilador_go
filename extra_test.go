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

func TestExtraV21UserExample(t *testing.T) {
	code := `x = 3
n = if x > 2 { 5 } else { 3 } // if x > 2 then 5 else 3
for i = 0; i < n; i = i + 1 {
  Println(i)
}
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "0\n1\n2\n3\n4"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV21ProfessorTest1(t *testing.T) {
	code := `// For implementation
n = 5
y = 1
i = 0
for (i = 0; i < n; i = i + 1) {
  y = y * (i + 1)
}
Println(if i == n { y } else { n })
Println(if i == n+1 { y } else { n })
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "120\n5"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV21ElseBranch(t *testing.T) {
	code := `x = 1
n = if x > 2 { 5 } else { 3 }
for i = 0; i < n; i = i + 1 {
  Println(i)
}
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "0\n1\n2"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV21Precedence(t *testing.T) {
	code := `x = 3
n = 10 * if x > 2 { 5 } else { 3 } + 2
Println(n)
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "52"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV21Unary(t *testing.T) {
	code := `x = 3
n = -if x > 2 { 5 } else { 3 }
Println(n)
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "-5"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV21NestedIfExpr(t *testing.T) {
	code := `x = 3
y = 1
n = if x > 2 { if y > 0 { 42 } else { 99 } } else { 0 }
Println(n)
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "42"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV21WhileCompatibility(t *testing.T) {
	code := `i = 0
for i < 3 {
  Println(i)
  i = i + 1
}
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "0\n1\n2"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV21EmptyInit(t *testing.T) {
	code := `i = 0
for ; i < 3; i = i + 1 {
  Println(i)
}
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "0\n1\n2"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV21MisalignedBraceFor(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic on misaligned brace in for")
		}
	}()

	code := `for 1 == 1 { }
`
	runCode(code)
}

func TestExtraV21MoreCases(t *testing.T) {
	code := `Println(if 1 < 2 { 100 } else { 200 })
for i = 0; i < (if 1 == 1 { 3 } else { 5 }); i = i + 1 {
  Println(i)
}
for i = 5; i > 0; i = i - 2 {
  Println(i)
}
`
	out, err := runCode(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "100\n0\n1\n2\n5\n3\n1"
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestExtraV21MisalignedBraceIf(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic on misaligned brace in if")
		}
	}()

	code := `if 1 == 1 { }
`
	runCode(code)
}
