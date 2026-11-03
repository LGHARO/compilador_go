# compilador_go

[![Compilation Status](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)](https://compiler-tester.insper-comp.com.br/svg/LGHARO/compilador_go)

![Diagrama Sintático](https://compiler-tester.insper-comp.com.br/ds?version=v2.3)

This repository is monitored by Compiler Tester for automatic compilation status.

```ebnf

PROGRAM =  {(FUNCDEC|VARDEC|Ε),"\n"} ;
FUNCDEC = "func",IDENTIFIER,"(",(IDENTIFIER,TYPE,{",",IDENTIFIER,TYPE}|Ε),")",(TYPE|Ε),BLOCK ;
VARDEC = "var",IDENTIFIER,TYPE,("=",BOOLEXPRESSION|Ε) ;
BLOCK = "{","\n",{STATEMENT,},"}" ;
STATEMENT = ("return",BOOLEXPRESSION|("Println","(",BOOLEXPRESSION,")")|("for",BOOLEXPRESSION,STATEMENT)|("if",BOOLEXPRESSION,BLOCK,(Ε|"else",BLOCK))|(IDENTIFIER,("=",BOOLEXPRESSION|"(",(BOOLEXPRESSION,{",",BOOLEXPRESSION}|Ε),")"))|BLOCK|VARDEC|Ε),"\n" ;
BOOLEXPRESSION = BOOLTERM, {"||", BOOLTERM} ;
BOOLTERM = RELEXPRESSION, {"&&", RELEXPRESSION} ;
RELEXPRESSION = EXPRESSION, { ("==" | ">" | "<"), EXPRESSION } ;
EXPRESSION = TERM, { ("+" | "-"), TERM } ;
TERM = FACTOR, { ("*" | "/"), FACTOR } ;
FACTOR = NUMBER|STRING|BOOLEAN|IDENTIFIER,("(",(BOOLEXPRESSION,{",",BOOLEXPRESSION}|Ε),")"|Ε)|("!"|"+"|"-"),FACTOR|"(",BOOLEXPRESSION,")"|"Scanln","(",")" ;
TYPE = "string" | "int" | "bool" ;
BOOLEAN = "true" | "false" ;
STRING = '"', {"..."}, '"' ;
IDENTIFIER = LETTER, { LETTER | DIGIT | "_" } ;
NUMBER = DIGIT, { DIGIT } ;
LETTER = "a"|...|"z"|"A"|...|"Z" ;
DIGIT = "0"|...|"9" ;

```
