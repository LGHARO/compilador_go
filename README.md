# compilador_go

[![Compilation Status](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)

![Diagrama Sintático](https://compiler-tester.insper-comp.com.br/ds?version=v2.1)

This repository is monitored by Compiler Tester for automatic compilation status.

```ebnf

PROGRAM = { STATEMENT } ;
STATEMENT = ( ("Println","(",BOOLEXPRESSION,")") | ("for",BOOLEXPRESSION,BLOCK) | ("if",BOOLEXPRESSION,BLOCK,("else",BLOCK) | ε) | (IDENTIFIER,"=",BOOLEXPRESSION) | BLOCK | ε ), "\n" ;
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
FACTOR = NUMBER | IDENTIFIER | ("!" | "+" | "-"), FACTOR | "(", BOOLEXPRESSION, ")" | "Scanln", "(", ")" ;

IDENTIFIER = LETTER, { LETTER | DIGIT | "_" } ;
NUMBER = DIGIT, { DIGIT } ;
LETTER = "a" | "b" | ... | "z" | "A" | "B" | ... | "Z" ;
DIGIT = "0" | "1" | ... | "9" ;
EOL = "\n" ;

```
