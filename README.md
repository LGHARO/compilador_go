# compilador_go

[![Compilation Status](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)

![Diagrama Sintático](https://compiler-tester.insper-comp.com.br/ds?version=v3.0)

This repository is monitored by Compiler Tester for automatic compilation status.

```ebnf

PROGRAM = { STATEMENT } ;
STATEMENT = ( "var", IDENTIFIER, TYPE, [ "=", BOOLEXPRESSION ]
            | "Println", "(", BOOLEXPRESSION, ")"
            | "for", BOOLEXPRESSION, BLOCK
            | "if", BOOLEXPRESSION, BLOCK, [ "else", ( IF | BLOCK ) ]
            | BLOCK
            | IDENTIFIER, "=", BOOLEXPRESSION
            | Ε ), EOL ;
BLOCK = "{", EOL, { STATEMENT }, "}" ;
IF = "if", BOOLEXPRESSION, BLOCK, [ "else", ( IF | BLOCK ) ] ;
WHILE = "for", BOOLEXPRESSION, BLOCK ;
ASSIGNMENT = IDENTIFIER, "=", BOOLEXPRESSION ;
PRINT = "Println", "(", BOOLEXPRESSION, ")" ;

BOOLEXPRESSION = BOOLTERM, { "||", BOOLTERM } ;
BOOLTERM = RELEXPRESSION, { "&&", RELEXPRESSION } ;
RELEXPRESSION = EXPRESSION, { ("<" | "==" | ">"), EXPRESSION } ;

EXPRESSION = TERM, { ("+" | "-"), TERM } ;
TERM = FACTOR, { ("*" | "/"), FACTOR } ;
FACTOR = NUMBER | BOOLEAN | STRING | IDENTIFIER | ( "!" | "+" | "-" ), FACTOR | "(", BOOLEXPRESSION, ")" | "Scanln", "(", ")" | "ScanlnStr", "(", ")" ;

IDENTIFIER = LETTER, { LETTER | DIGIT | "_" } ;
NUMBER = DIGIT, { DIGIT } ;
STRING = '"', { ANY_CHAR_EXCEPT_QUOTE_OR_NEWLINE }, '"' ;
BOOLEAN = "true" | "false" ;
TYPE = "string" | "int" | "bool" ;
LETTER = "a" | "b" | ... | "z" | "A" | "B" | ... | "Z" ;
DIGIT = "0" | "1" | ... | "9" ;
EOL = "\n" ;

```

