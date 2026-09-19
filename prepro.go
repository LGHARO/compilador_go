package main

import (
	"strings"
)

type Prepro struct{}

func (prepro *Prepro) filter(code string) string {
	var stringLimpa strings.Builder
	i := 0

	// Quando acharmos "//" percorremos e retiramos até a quebra de linha
	for i < len(code) {
		
		//unico jeito de comparar é oçhar os dois chars
		if i+1 < len(code) && code[i] == '/' && code[i+1] == '/' {
			
			//percorre até achar o final ou encontrar a quebra de linha
			for i < len(code) && code[i] != '\n' {
				i++
			}
			continue // Pula para a próxima iteração
		}

		// Adiciona o byte atual (usando WriteByte ou WriteRune)
		stringLimpa.WriteByte(code[i])
		i++
	}

	return stringLimpa.String()
}