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
	TokenInitiator
	TokenKeyword
	TokenOperator
	TokenDelimiter
	TokenString
	TokenNumber
	TokenComment
	TokenWhitespace
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value string
}

type Query struct {
	Resource     string
	ResourceName string
	Action       string
	Columns      []QueryColumn
}
type QueryColumn struct {
	ColumnName string
	Action     string
}

var initators = map[string]struct{}{
	"SELECT": {},
	"DELETE": {},
	"UPDATE": {},
	"INSERT": {},
	"CREATE": {},
	"ALTER":  {},
	"DROP":   {},
}

var keywords = map[string]struct{}{
	"FROM":       {},
	"WHERE":      {},
	"GROUP BY":   {},
	"ORDER BY":   {},
	"INTO":       {},
	"VALUES":     {},
	"SET":        {},
	"TABLE":      {},
	"ADD":        {},
	"PRIMARY":    {},
	"KEY":        {},
	"FOREIGN":    {},
	"REFERENCES": {},
	"DATABASE":   {},
	"INDEX":      {},
	"PROCEDURE":  {},
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

	return append(tokens, Token{Type: TokenEOF})
}

func getTokenType(token string) TokenType {
	if _, ok := initators[token]; ok {
		return TokenInitiator
	}
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

	queries, err := GroupQuerys(tokens)
	if err != nil {
		panic(err)
	}

	for _, query := range queries {
		fmt.Print(query.Action+" ", query.Resource+" ", query.ResourceName+" ")
		for _, resource := range query.Columns {
			fmt.Print(resource.Action+" ", resource.ColumnName+" ")
		}
	}
}

func GroupQuerys(tokens []Token) ([]Query, error) {
	var queries []Query
	for i := 0; i < len(tokens); i++ {
		elem := tokens[i]
		if elem.Type != TokenInitiator {
			continue
		}

		tableName, err := getElementIfExists(tokens, i+1)
		if err != nil {
			return nil, err
		}

		resource, err := getElementIfExists(tokens, i+2)
		if err != nil {
			return nil, err
		}

		i += 2
		var columns []QueryColumn
		for {
			elem, err := getElementIfExists(tokens, i)
			if err != nil {
				break
			}
			if elem.Type != TokenKeyword {
				i++
				continue
			}

			name, err := getElementIfExists(tokens, i+1)
			if err != nil {
				break
			}
			columns = append(columns, QueryColumn{Action: elem.Value, ColumnName: name.Value})
			i += 2
		}
		queries = append(queries, Query{Action: elem.Value, Resource: tableName.Value, ResourceName: resource.Value, Columns: columns})
	}
	return queries, nil
}

func getElementIfExists[T any](arr []T, i int) (*T, error) {
	if i < 0 || i >= len(arr) {
		return nil, fmt.Errorf("error accessing index: %v", i)
	}
	return &arr[i], nil
}
