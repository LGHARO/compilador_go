# compilador_go

[![Compilation Status](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)

![Diagrama Sintático](https://compiler-tester.insper-comp.com.br/ds?version=v2.3)

This repository is monitored by Compiler Tester for automatic compilation status.

```ebnf

PROGRAM = { FUNCTION | STATEMENT } ;
FUNCTION = "func", IDENTIFIER, "(", [ PARAMETER, { ",", PARAMETER } ], ")", [ TYPE ], BLOCK ;
PARAMETER = IDENTIFIER, TYPE ;
STATEMENT = ( "var", IDENTIFIER, TYPE, [ "=", BOOLEXPRESSION ] | "Println", "(", BOOLEXPRESSION, ")" | "for", BOOLEXPRESSION, BLOCK | "if", BOOLEXPRESSION, BLOCK, [ "else", ( IF | BLOCK ) ] | "return", BOOLEXPRESSION | FUNCCALL | BLOCK | IDENTIFIER, "=", BOOLEXPRESSION | Ε ), "\n" ;
FUNCCALL = IDENTIFIER, "(", [ BOOLEXPRESSION, { ",", BOOLEXPRESSION } ], ")" ;
BLOCK = "{", "\n", { STATEMENT }, "}" ;
IF = "if", BOOLEXPRESSION, BLOCK, [ "else", ( IF | BLOCK ) ] ;
WHILE = "for", BOOLEXPRESSION, BLOCK ;
ASSIGNMENT = IDENTIFIER, "=", BOOLEXPRESSION ;
PRINT = "Println", "(", BOOLEXPRESSION, ")" ;

BOOLEXPRESSION = BOOLTERM, { "||", BOOLTERM } ;
BOOLTERM = RELEXPRESSION, { "&&", RELEXPRESSION } ;
RELEXPRESSION = EXPRESSION, { ("<" | "==" | ">"), EXPRESSION } ;

EXPRESSION = TERM, { ("+" | "-"), TERM } ;
TERM = FACTOR, { ("*" | "/"), FACTOR } ;
FACTOR = NUMBER | BOOLEAN | STRING | FUNCCALL | IDENTIFIER | ( "!" | "+" | "-" ), FACTOR | "(", BOOLEXPRESSION, ")" | "Scanln", "(", ")" ;

IDENTIFIER = LETTER, { LETTER | DIGIT | "_" } ;
NUMBER = DIGIT, { DIGIT } ;
STRING = '"', {"..."}, '"' ;
BOOLEAN = "true" | "false" ;
TYPE = "string" | "int" | "bool" ;
LETTER = "a" | "b" | ... | "z" | "A" | "B" | ... | "Z" ;
DIGIT = "0" | "1" | ... | "9" ;
EOL = "\n" ;

```

## Tipos suportados

| Tipo     | Exemplos                   | Operações                          |
|----------|----------------------------|------------------------------------|
| `int`    | `42`, `Scanln()`           | `+ - * / == < >`                   |
| `string` | `"hello"`, `ScanlnStr()`   | `+` (concatenação), `==`           |
| `bool`   | `true`, `false`            | `&&` `\|\|` `!` `==`               |

## Declaração de variáveis

```
var <nome> <tipo>
var <nome> <tipo> = <expressão>
```

---

## Versão 2.3 - Funções e Escopo de Variáveis

Nesta versão, a linguagem passa a suportar declaração e chamada de funções com parâmetros tipados, tipo de retorno opcional, retorno com verificação de tipo e escopo léxico de variáveis encadeado.

### Principais Funcionalidades:

1. **Tokens Adicionados (`lexer.go`)**:
   - `func` -> token `FUNC`
   - `return` -> token `RETURN`
   - `,` -> token `COMMA`

2. **Nós da AST (`parser.go`)**:
   - `FuncDec`: armazena declaração de função (`n+2` filhos: `IdentifierOp` com o nome, `n` nós `VarDec` representando os parâmetros, e `Block` com o corpo). O tipo de retorno é armazenado em `value`.
   - `FuncCall`: chamada de função com `value` igual ao nome e `children` contendo as expressões passadas como argumentos.
   - `Return`: nó de retorno com um único filho (expressão retornada).

3. **Escopo Encadeado (`symboltable.go`)**:
   - `SymbolTable` agora possui o ponteiro `parent *SymbolTable`, formando uma lista ligada reversa de escopos.
   - `getter` e `setter` buscam variáveis recursivamente nos escopos pais.
   - `createVariable` declara variáveis no escopo local, permitindo sombreamento (*shadowing*) de variáveis mais externas sem sobrescrevê-las.
   - `createFunction` registra funções na tabela de símbolos raiz (global).
   - `Variable` armazena as flags `isFunc` e `isReturn`.

4. **Controle de Fluxo e Execução (`parser.go`)**:
   - `Block`, `If` e `While` criam novos escopos locais e propagam o valor de retorno caso uma instrução `return` seja executada.
   - Chamadas de função (`FuncCall`) avaliam argumentos no escopo do chamador, criam um novo escopo com acesso aos símbolos globais, validam tipos e quantidade de argumentos, executam o corpo e verificam se o valor retornado corresponde à assinatura.
   - A função `main` é obrigatória: o método `Parser.run()` insere automaticamente um nó `FuncCall("main")` na raiz da AST.

### Exemplo de Sintaxe:

```go
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
```