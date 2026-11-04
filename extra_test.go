package main

import (
	"os"
	"strings"
	"testing"
)

func compileSource(code string) string {
	filter := &Prepro{}
	dadosFiltrados := filter.filter(code)

	l := &Lexer{
		source:   dadosFiltrados,
		position: 0,
	}

	p := &Parser{
		lexer: *l,
	}

	resetNodeId()
	CodeReset()

	AST := p.run()

	st := &SymbolTable{
		table:        make(map[string]Variable),
		currentShift: 0,
	}

	AST.Generate(st)

	tmpFile := "test_out.asm"
	defer os.Remove(tmpFile)
	CodeDump(tmpFile)

	asmBytes, err := os.ReadFile(tmpFile)
	if err != nil {
		panic(err)
	}
	return string(asmBytes)
}

func TestFunctionAssemblyStructure(t *testing.T) {
	code := `
func soma(x int, y int) int {
  var a int
  a = x + y 
  return a
}

func main() {
  var b int
  b = soma(3, 4)
  Println(b)
}
`
	asm := compileSource(code)

	// Verifica se _start chama main e termina com int 0x80
	if !strings.Contains(asm, "_start:") {
		t.Fatal("Esperava rótulo _start:")
	}
	if !strings.Contains(asm, "call main") {
		t.Fatal("Esperava call main dentro de _start")
	}
	if !strings.Contains(asm, "int 0x80") {
		t.Fatal("Esperava interrupção int 0x80")
	}

	// Verifica prólogo e epílogo de soma
	if !strings.Contains(asm, "soma:\npush ebp\nmov ebp, esp") {
		t.Fatal("Esperava rótulo soma com prólogo 'push ebp; mov ebp, esp'")
	}
	if !strings.Contains(asm, "mov esp, ebp\npop ebp\nret") {
		t.Fatal("Esperava epílogo 'mov esp, ebp; pop ebp; ret'")
	}

	// Verifica acesso aos parâmetros x ([ebp+8]) e y ([ebp+12])
	if !strings.Contains(asm, "[ebp+8]") {
		t.Fatal("Esperava acesso ao primeiro parâmetro x em [ebp+8]")
	}
	if !strings.Contains(asm, "[ebp+12]") {
		t.Fatal("Esperava acesso ao segundo parâmetro y em [ebp+12]")
	}

	// Verifica acesso à variável local a em [ebp-4]
	if !strings.Contains(asm, "[ebp-4]") {
		t.Fatal("Esperava alocação e acesso à variável local a em [ebp-4]")
	}

	// Verifica chamada a soma com limpeza da pilha (add esp, 8)
	if !strings.Contains(asm, "call soma\nadd esp, 8") {
		t.Fatal("Esperava 'call soma' seguido de 'add esp, 8'")
	}
}

func TestFunctionMultipleParameters(t *testing.T) {
	code := `
func calc(a int, b int, c int, d int) int {
  return a + b + c + d
}

func main() {
  Println(calc(1, 2, 3, 4))
}
`
	asm := compileSource(code)

	if !strings.Contains(asm, "[ebp+8]") || !strings.Contains(asm, "[ebp+12]") ||
		!strings.Contains(asm, "[ebp+16]") || !strings.Contains(asm, "[ebp+20]") {
		t.Fatal("Esperava offsets [ebp+8], [ebp+12], [ebp+16], [ebp+20] para 4 parâmetros")
	}

	if !strings.Contains(asm, "call calc\nadd esp, 16") {
		t.Fatal("Esperava 'call calc' seguido de 'add esp, 16'")
	}
}

func TestFunctionRecursion(t *testing.T) {
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
	asm := compileSource(code)

	if !strings.Contains(asm, "call fat\nadd esp, 4") {
		t.Fatal("Esperava chamada recursiva 'call fat' com 'add esp, 4'")
	}
}

func TestGlobalVariableWithFunctions(t *testing.T) {
	code := `
var g int = 42

func test() int {
  return g
}

func main() {
  Println(test())
}
`
	asm := compileSource(code)

	// Verifica se _var_g foi declarado em .data
	if !strings.Contains(asm, "_var_g: dd 0") {
		t.Fatal("Esperava _var_g: dd 0 na seção .data")
	}

	// Verifica se g é acessado via [_var_g]
	if !strings.Contains(asm, "[_var_g]") {
		t.Fatal("Esperava acesso a g via [_var_g]")
	}
}

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
	compileSource(code)
}

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
	compileSource(code)
}

func TestErrorMissingMain(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por ausência da função main")
		}
	}()

	code := `
func soma(x int, y int) int {
  return x + y
}
`
	compileSource(code)
}

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
	compileSource(code)
}

func TestTrailingCommaInParams(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por vírgula trailing na declaração de parâmetros")
		}
	}()

	code := `
func test(x int,) int {
  return 2
}
func main() {
  Println(test(5))
}
`
	compileSource(code)
}

func TestTrailingCommaInCall(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Esperava panic por vírgula trailing na chamada de função")
		}
	}()

	code := `
func test(x int, y int) int {
  return x + y
}
func main() {
  Println(test(2, ))
}
`
	compileSource(code)
}

func evaluateSource(code string) string {
	filter := &Prepro{}
	dadosFiltrados := filter.filter(code)

	l := &Lexer{
		source:   dadosFiltrados,
		position: 0,
	}

	p := &Parser{
		lexer: *l,
	}

	resetNodeId()

	AST := p.run()

	st := &SymbolTable{
		table:        make(map[string]Variable),
		currentShift: 0,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	AST.Evaluate(*st)

	w.Close()
	os.Stdout = oldStdout

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	return strings.TrimSpace(string(buf[:n]))
}

func TestExamplePromptEvaluate(t *testing.T) {
	code := `
var b int = 5

func soma(x int, y int) int {
  var a int
  a = x + y 
  Println(a)
  return a
}

func main()  {
  var a int
  {
    var b int
    a = 3
    b = soma(a, 4)
    Println(b)
  }
  Println(a)
  Println(b)
}
`
	out := evaluateSource(code)
	expected := "7\n7\n3\n5"
	if out != expected {
		t.Fatalf("Esperava:\n%s\nObteve:\n%s", expected, out)
	}
}

func TestRecursionEvaluate(t *testing.T) {
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
	out := evaluateSource(code)
	expected := "120"
	if out != expected {
		t.Fatalf("Esperava:\n%s\nObteve:\n%s", expected, out)
	}
}

