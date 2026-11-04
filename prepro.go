package main

import (
	"regexp"
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

	codeSemComentarios := stringLimpa.String()

	// Processa diretivas "define <name> = <val>"
	reDefine := regexp.MustCompile(`^\s*define\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*(.+)$`)
	lines := strings.Split(codeSemComentarios, "\n")
	var resultLines []string

	for lineIdx := 0; lineIdx < len(lines); lineIdx++ {
		line := lines[lineIdx]
		match := reDefine.FindStringSubmatch(line)
		if len(match) == 3 {
			defName := strings.TrimSpace(match[1])
			defVal := strings.TrimSpace(match[2])
			if defName != "" && defVal != "" {
				reWord := regexp.MustCompile(`\b` + regexp.QuoteMeta(defName) + `\b`)
				for nextIdx := lineIdx + 1; nextIdx < len(lines); nextIdx++ {
					lines[nextIdx] = reWord.ReplaceAllString(lines[nextIdx], defVal)
				}
				// Linha de define é consumida no pré-processamento
				continue
			}
		}
		resultLines = append(resultLines, line)
	}

	return strings.Join(resultLines, "\n")
}