# compilador_go

[![Compilation Status](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)

![Diagrama Sintático](https://compiler-tester.insper-comp.com.br/ds?version=v1.1)

This repository is monitored by Compiler Tester for automatic compilation status.

```ebnf
EXPRESSION = TERM, { ("+" | "-"), TERM } ;
TERM = FACTOR, { ("*" | "/"), FACTOR } ;
FACTOR = NUMBER | ("+" | "-"), FACTOR | "(", EXPRESSION, ")" ;
NUMBER = DIGIT, { DIGIT } ;
DIGIT = "0" | "1" | ... | "9" ;
```
