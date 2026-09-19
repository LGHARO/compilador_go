# compilador_go

[![Compilation Status](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)

![Diagrama Sintático](https://compiler-tester.insper-comp.com.br/ds?version=v2.0)

This repository is monitored by Compiler Tester for automatic compilation status.

```ebnf

PROGRAM = { STATEMENT } ;
STATEMENT = ("Println","(",EXPRESSION,")"|IDENTIFIER,"=",EXPRESSION|Ε),"\n";
ASSIGNMENT = IDENTIFIER, "=", EXPRESSION ;
PRINT = "Println", "(", EXPRESSION, ")" ;
EXPRESSION = TERM, { ("+" | "-"), TERM } ;
TERM = FACTOR, { ("*" | "/"), FACTOR } ;
FACTOR = NUMBER|IDENTIFIER|("+"|"-"),FACTOR|"(",EXPRESSION,")" ;
IDENTIFIER = LETTER, { LETTER | DIGIT | "_" } ;
NUMBER = DIGIT, { DIGIT } ;
LETTER = "a"|"b"|...|"z"|"A"|"B"|...|"Z" ;
DIGIT = "0"|"1"|...|"9" ;
EOL = "\n" ;

```