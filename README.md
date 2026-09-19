# compilador_go

[![Compilation Status](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)

![Diagrama Sintático](https://compiler-tester.insper-comp.com.br/ds?version=v2.0)

This repository is monitored by Compiler Tester for automatic compilation status.

```ebnf

PROGRAM = { STATEMENT } ;
STATEMENT = ( ε | ASSIGNMENT | PRINT), EOL ;
ASSIGNMENT = IDENTIFIER, "=", EXPRESSION ;
PRINT = "Println", "(", EXPRESSION, ")" ;
EXPRESSION = TERM, { ("+" | "-"), TERM } ;
TERM = FACTOR, { ("*" | "/"), FACTOR } ;
FACTOR = (("+" | "-"), FACTOR) | NUMBER | "(", EXPRESSION, ")" | IDENTIFIER ;
IDENTIFIER = LETTER, { LETTER | DIGIT | "_" } ;
NUMBER = DIGIT, { DIGIT } ;
LETTER = ( a | ... | z | A | ... | Z ) ;
DIGIT = ( 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 ) ;
EOL = "\n" ;

```