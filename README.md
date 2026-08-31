# compilador-ling-par

[![Compilation Status](https://compiler-tester.insper-comp.com.br/svg/emilybrtt/compilador-ling-par)](https://compiler-tester.insper-comp.com.br/svg/emilybrtt/compilador-ling-par)

This repository is monitored by Compiler Tester for automatic compilation status.

Comando pra rodar: 
```bash
go run main.go [argumento]
```

## EBNF

```ebnf
EXPRESSION = TERM, { ("+" | "-"), TERM } ;
TERM = FACTOR, { ("*" | "/"), FACTOR } ;
FACTOR = NUMBER | ("+" | "-"), FACTOR | "(", EXPRESSION, ")" ;
NUMBER = DIGIT, { DIGIT } ;
DIGIT = "0" | "1" | ... | "9" ;
```

![Diagrama Sintático](https://compiler-tester.insper-comp.com.br/ds?version=v1.1)
