# compilador_go

[![Compilation Status](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)

![Diagrama Sintático](https://compiler-tester.insper-comp.com.br/ds?version=v2.2)

This repository is monitored by Compiler Tester for automatic compilation status.

```ebnf

PROGRAM = { STATEMENT } ;
STATEMENT = ( VARDEC | ("Println","(",BOOLEXPRESSION,")") | ("for",BOOLEXPRESSION,BLOCK) | ("if",BOOLEXPRESSION,BLOCK,("else",( IF | BLOCK )) | ε) | (IDENTIFIER,"=",BOOLEXPRESSION) | BLOCK | ε ), "\n" ;
BLOCK = "{", "\n", { STATEMENT }, "}" ;
VARDEC = "var", IDENTIFIER, TYPE, ["=", BOOLEXPRESSION] ;
IF = "if", BOOLEXPRESSION, BLOCK, [ "else", ( IF | BLOCK ) ] ;
WHILE = "for", BOOLEXPRESSION, BLOCK ;
ASSIGNMENT = IDENTIFIER, "=", BOOLEXPRESSION ;
PRINT = "Println", "(", BOOLEXPRESSION, ")" ;

BOOLEXPRESSION = BOOLTERM, { "||", BOOLTERM } ;
BOOLTERM = RELEXPRESSION, { "&&", RELEXPRESSION } ;
RELEXPRESSION = EXPRESSION, { ("<" | "==" | ">"), EXPRESSION } ;

EXPRESSION = TERM, { ("+" | "-"), TERM } ;
TERM = FACTOR, { ("*" | "/"), FACTOR } ;
FACTOR = NUMBER | STRING | BOOLEAN | IDENTIFIER | ("!" | "+" | "-"), FACTOR | "(", BOOLEXPRESSION, ")" | "Scanln", "(", ")" | "ScanlnStr", "(", ")" ;

IDENTIFIER = LETTER, { LETTER | DIGIT | "_" } ;
NUMBER = DIGIT, { DIGIT } ;
STRING = '"', { ... }, '"' ;
BOOLEAN = "true" | "false" ;
TYPE = "int" | "string" | "bool" ;
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

## Mudanças para suporte a string em if/while

Para testar `if`/`for` com entrada de string, foi adicionada a função `ScanlnStr()` ao lexer e ao parser:

- **Lexer**: novo token `READ_STR` para a palavra-chave `ScanlnStr`
- **Parser**: novo nó `ReadStringOp` que lê uma string com `fmt.Scan`
- **SymbolTable**: `Variable` agora armazena `interface{}` e `varType string`
- A condição do `if`/`for` aceita `bool` ou `int`; para usar string no controle de fluxo,
  compare a string com `==` (retorna `bool`) e use o resultado como condição.