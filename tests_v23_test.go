package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput executa uma função e captura sua saída padrão
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

// runCode executa uma string de código no compilador
func runCode(code string) {
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

// TestExamplePrompt verifica o exemplo oficial fornecido no enunciado
func TestExamplePrompt(t *testing.T) {
	code := `
var b int = 5 // Variável Global

func soma(x int, y int) int {
  var a int // Variável Local
  a = x + y 
  Println(a) // Imprime 7
  return a
}

func main()  {
  var a int
  {
    var b int
    a = 3
    b = soma(a, 4)
    Println(b) // Imprime 7
  }
  Println(a) // Imprime 3
  Println(b) // Imprime 5
}
`
	out := captureOutput(func() {
		runCode(code)
	})

	expected := "7\n7\n3\n5"
	if out != expected {
		t.Fatalf("Esperava saída:\n%s\nMas obteve:\n%s", expected, out)
	}
}

// TestRecursion testa função recursiva (fatorial)
func TestRecursion(t *testing.T) {
	code := `
func fat(n int) int {
  if n == 0 || n == 1 {
    return 1
  }
  return n * fat(n - 1)
}

func main() {
  Println(fat(5))
}
`
	out := captureOutput(func() {
		runCode(code)
	})

	expected := "120"
	if out != expected {
		t.Fatalf("Esperava saída:\n%s\nMas obteve:\n%s", expected, out)
	}
}

// TestMutualRecursion testa funções mutualmente recursivas
func TestMutualRecursion(t *testing.T) {
	code := `
func isEven(n int) bool {
  if n == 0 {
    return true
  }
  return isOdd(n - 1)
}

func isOdd(n int) bool {
  if n == 0 {
    return false
  }
  return isEven(n - 1)
}

func main() {
  Println(isEven(4))
  Println(isOdd(4))
}
`
	out := captureOutput(func() {
		runCode(code)
	})

	expected := "true\nfalse"
	if out != expected {
		t.Fatalf("Esperava saída:\n%s\nMas obteve:\n%s", expected, out)
	}
}

// TestScopeVariable testa isolamento e shadowing de variáveis em blocos
func TestScopeVariable(t *testing.T) {
	code := `
var x int = 100

func main() {
  var x int = 10
  {
    var x int = 1
    Println(x)
  }
  Println(x)
}
`
	out := captureOutput(func() {
		runCode(code)
	})

	expected := "1\n10"
	if out != expected {
		t.Fatalf("Esperava saída:\n%s\nMas obteve:\n%s", expected, out)
	}
}

// TestErrorIncorrectArgCount testa erro ao passar quantidade incorreta de argumentos
func TestErrorIncorrectArgCount(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por número incorreto de argumentos")
		}
	}()

	code := `
func soma(x int, y int) int {
  return x + y
}

func main() {
  soma(1)
}
`
	runCode(code)
}

// TestErrorDifferentArgType testa erro ao passar tipos de argumentos diferentes
func TestErrorDifferentArgType(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por incompatibilidade de tipo no argumento")
		}
	}()

	code := `
func soma(x int, y int) int {
  return x + y
}

func main() {
  soma(1, "dois")
}
`
	runCode(code)
}

// TestErrorFunctionDoesNotExist testa erro ao chamar função inexistente
func TestErrorFunctionDoesNotExist(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic ao chamar função inexistente")
		}
	}()

	code := `
func main() {
  foo(1)
}
`
	runCode(code)
}

// TestErrorVariableOutOfScope testa erro ao acessar variável fora do escopo
func TestErrorVariableOutOfScope(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por variável fora de escopo")
		}
	}()

	code := `
func main() {
  {
    var bloco int = 42
  }
  Println(bloco)
}
`
	runCode(code)
}

// TestErrorVoidFunctionReturn testa erro ao tentar retornar valor em função void
func TestErrorVoidFunctionReturn(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por return em função void")
		}
	}()

	code := `
func semRetorno() {
  return 10
}

func main() {
  semRetorno()
}
`
	runCode(code)
}

// TestErrorNoReturnInTypedFunction testa erro quando função tipada não retorna valor
func TestErrorNoReturnInTypedFunction(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por ausência de return em função tipada")
		}
	}()

	code := `
func comRetorno() int {
  var x int = 1
}

func main() {
  comRetorno()
}
`
	runCode(code)
}

// TestErrorFunctionInsideFunction testa erro ao tentar declarar função dentro de função
func TestErrorFunctionInsideFunction(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava erro ao declarar função dentro de função")
		}
	}()

	code := `
func main() {
  func interna() {
  }
}
`
	runCode(code)
}

// TestErrorMissingMain testa erro quando a função main não é declarada
func TestErrorMissingMain(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic pela ausência da função main")
		}
	}()

	code := `
func soma(x int, y int) int {
  return x + y
}
`
	runCode(code)
}

// TestMisalignedCloseBraceIf (Test 54): if 1 == 1 { } sem quebra de linha deve dar erro
func TestMisalignedCloseBraceIf(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por '}' na mesma linha que '{' em if")
		}
	}()

	code := `
func main() {
  if 1 == 1 { }
}
`
	runCode(code)
}

// TestMisalignedCloseBraceElse (Test 56): else { } sem quebra de linha deve dar erro
func TestMisalignedCloseBraceElse(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por '}' na mesma linha que '{' em else")
		}
	}()

	code := `
func main() {
  if 1 == 1 {
  } else { }
}
`
	runCode(code)
}

// TestMisalignedCloseBraceWhile (Test 59): for 1 == 1 { } sem quebra de linha deve dar erro
func TestMisalignedCloseBraceWhile(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por '}' na mesma linha que '{' em for")
		}
	}()

	code := `
func main() {
  for 1 == 1 { }
}
`
	runCode(code)
}

// TestValidEmptyBlockWithNewline verifica que bloco vazio com quebra de linha é aceito
func TestValidEmptyBlockWithNewline(t *testing.T) {
	code := `
func main() {
  if 1 == 1 {
  }
}
`
	runCode(code)
}

