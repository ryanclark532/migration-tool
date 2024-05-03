package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TokenType int

const (
	TokenIdentifier TokenType = iota
	TokenKeyword
	TokenOperator
	TokenDelimiter
	TokenString
	TokenNumber
	TokenComment
	TokenWhitespace
)

type Token struct {
	Type  TokenType
	Value string
}

var initators = map[string]struct{}{}

var keywords = map[string]struct{}{
	"SELECT":     {},
	"FROM":       {},
	"WHERE":      {},
	"GROUP BY":   {},
	"ORDER BY":   {},
	"INSERT":     {},
	"INTO":       {},
	"VALUES":     {},
	"UPDATE":     {},
	"SET":        {},
	"DELETE":     {},
	"CREATE":     {},
	"TABLE":      {},
	"DROP":       {},
	"ALTER":      {},
	"ADD":        {},
	"PRIMARY":    {},
	"KEY":        {},
	"FOREIGN":    {},
	"REFERENCES": {},
}

var operators = map[string]struct{}{
	"+":   {},
	"-":   {},
	"*":   {},
	"/":   {},
	"=":   {},
	"!=":  {},
	">":   {},
	"<":   {},
	">=":  {},
	"<=":  {},
	"AND": {},
	"OR":  {},
	"NOT": {},
}

var delimiters = map[string]struct{}{
	"(": {},
	")": {},
	";": {},
	",": {},
}

var keywordInverse = map[string]string{
	"ADD":            "DROP COLUMN",
	"ADD CONSTRAINT": "DROP CONSTRAINT",
}

func isIdentifierChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '_' ||
		ch == '.'
}

func isOperatorChar(ch byte) bool {
	return strings.ContainsRune("=><!+-*/", rune(ch))
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func tokenize(input string) []Token {
	var tokens []Token
	var token strings.Builder

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if isWhitespace(ch) {
			if token.Len() > 0 {
				tokens = append(tokens, Token{Type: getTokenType(token.String()), Value: token.String()})
				token.Reset()
			}
			continue
		}

		if ch == '-' && input[i+1] == '-' {
			sub := input[i:]
			x := strings.IndexByte(sub, '\n') + 1
			i = i + x
			continue
		}

		if ch == '/' && i+1 < len(input) && input[i+1] == '*' {
			sub := input[i:]
			x := strings.Index(sub, "*/") + 1
			i = i + x
			continue
		}

		if isIdentifierChar(ch) {
			token.WriteByte(ch)
			continue
		}

		if isOperatorChar(ch) {
			if token.Len() > 0 {
				tokens = append(tokens, Token{Type: getTokenType(token.String()), Value: token.String()})
				token.Reset()
			}
			token.WriteByte(ch)
			if i+1 < len(input) && isOperatorChar(input[i+1]) {
				i++
				token.WriteByte(input[i])
			}
			tokens = append(tokens, Token{Type: TokenOperator, Value: token.String()})
			token.Reset()
			continue
		}

		if _, ok := delimiters[string(ch)]; ok {
			if token.Len() > 0 {
				tokens = append(tokens, Token{Type: getTokenType(token.String()), Value: token.String()})
				token.Reset()
			}
			tokens = append(tokens, Token{Type: TokenDelimiter, Value: string(ch)})
			continue
		}

	}

	if token.Len() > 0 {
		tokens = append(tokens, Token{Type: getTokenType(token.String()), Value: token.String()})
	}

	return tokens
}

func getTokenType(token string) TokenType {
	if _, ok := keywords[token]; ok {
		return TokenKeyword
	}
	if _, ok := operators[token]; ok {
		return TokenOperator
	}
	if _, ok := delimiters[token]; ok {
		return TokenDelimiter
	}
	return TokenIdentifier
}

func main() {
	file, err := os.Open("example-update.sql")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var input strings.Builder
	for scanner.Scan() {
		input.WriteString(scanner.Text())
		input.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		os.Exit(1)
	}

	tokens := tokenize(input.String())

	generateDown(tokens)
}

func generateDown(tokens []Token) {
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type != TokenKeyword {
			continue
		}

		switch tokens[i].Value {
		case "CREATE":
			switch tokens[i+1].Value {
			case "TABLE":
				tableName := tokens[i+2].Value
				fmt.Printf("Delete Table %s \n", tableName)
			}
		case "ALTER":
			switch tokens[i+1].Value {
			case "TABLE":
				tableName := tokens[i+2].Value
				fmt.Printf("Alter Table %s \n", tableName)
				var opp string
				if tokens[i+4].Type == TokenKeyword {
					opp = fmt.Sprintf("%s %s", tokens[i+3].Value, tokens[i+4].Value)
				} else {
					opp = tokens[i+3].Value
				}
				fmt.Printf("%s: %s\n", opp, keywordInverse[opp])
			}
		}

	}
}
