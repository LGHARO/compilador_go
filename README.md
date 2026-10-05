# compilador_go

[![Compilation Status](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)

![Diagrama Sintático](https://compiler-tester.insper-comp.com.br/ds?version=v2.1)

This repository is monitored by Compiler Tester for automatic compilation status.

```ebnf

PROGRAM = { STATEMENT } ;
STATEMENT = ( ε | ASSIGNMENT | PRINT | IF | WHILE | BLOCK ), EOL ;
BLOCK = "{", { STATEMENT }, "}" ;
IF = "if", BOOL_EXPRESSION, BLOCK, [ "else", ( IF | BLOCK ) ] ;
WHILE = "for", BOOL_EXPRESSION, BLOCK ;
ASSIGNMENT = IDENTIFIER, "=", BOOL_EXPRESSION ;
PRINT = "Println", "(", BOOL_EXPRESSION, ")" ;
BOOL_EXPRESSION = BOOL_TERM, { "||", BOOL_TERM } ;
BOOL_TERM = REL_EXPRESSION, { "&&", REL_EXPRESSION } ;
REL_EXPRESSION = EXPRESSION, [ ("==" | "<" | ">"), EXPRESSION ] ;
EXPRESSION = TERM, { ("+" | "-"), TERM } ;
TERM = FACTOR, { ("*" | "/"), FACTOR } ;
FACTOR = (("+" | "-" | "!"), FACTOR) | NUMBER | READ | "(", BOOL_EXPRESSION, ")" | IDENTIFIER ;
READ = "Scanln", "(", ")" ;
IDENTIFIER = LETTER, { LETTER | DIGIT | "_" } ;
NUMBER = DIGIT, { DIGIT } ;
LETTER = ( a | ... | z | A | ... | Z ) ;
DIGIT = ( 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 ) ;
EOL = "\n" ;

```