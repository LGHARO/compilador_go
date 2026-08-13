//colocamos o nome do package no topo para indicar que esse é um arquivo executavel

package main

//podemos importar as bibliotecas dessa forma
import (
	"fmt"
	"os"
	"strconv"
	"strings"
)


func main(){
	
	
	resultado := 0;

	// os.Args[1] pega a operacao
	lista_operacoes := strings.Split(os.Args[1], " ")

	// criamos o slice que vai guardar as tres operacoes 
	var operacoes []string
	
	for i := 0; i < len(lista_operacoes) || len(operacoes) == 3; {
		// array completo
		
		if len(operacoes) == 3{


			if operacoes[1] == "-"{

				num1, _ := strconv.Atoi(operacoes[0]);
				num2, _ := strconv.Atoi(operacoes[2]);

				//calculamos o resultado
				resultado = num1 - num2
				
				
				// redimencionamos o slice para que volte a ter apenas um elemento
				operacoes= operacoes[:0]

				operacoes = append(operacoes, strconv.Itoa(resultado))

			}else if operacoes[1] == "+"{

				num1, _ := strconv.Atoi(operacoes[0]);
				num2, _ := strconv.Atoi(operacoes[2]);
				//calculamos o resultado
				resultado = num1 + num2
				 
				
				// redimencionamos o slice para que volte a ter apenas um elemento
				operacoes= operacoes[:0]

				operacoes = append(operacoes, strconv.Itoa(resultado))
			}

			
			
		// array faltando o ultimo numero
		}else if len(operacoes) == 2{
			if lista_operacoes[i] == "+" || lista_operacoes[i] == "-"{
				panic("erro: dois operadores em seguida")
			}else{
				operacoes = append(operacoes, lista_operacoes[i])
			}
			i++
		// array faltando o operador 
		}else if len(operacoes) == 1{

			if lista_operacoes[i] == "+"{
				operacoes = append(operacoes, lista_operacoes[i])
			}else if lista_operacoes[i] == "-"{
				operacoes = append(operacoes, lista_operacoes[i])
			}else{
				panic("erro: dois numeros em sequencia")
			}
			i++
		// adicao primeiro numero 
		}else if len(operacoes) == 0 {
			if lista_operacoes[i] == "+" || lista_operacoes[i] == "-"{
					panic("erro: comecou com operador")
				}else{
					operacoes = append(operacoes, lista_operacoes[i])
				}
				i++
		}

	}
	fmt.Println(resultado)

}