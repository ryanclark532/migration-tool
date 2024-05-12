package paser

import (
	"fmt"
	"ryanclark532/migration-tool/internal/lexer"
)

type Parser struct {
	Tokens       []lexer.Token
	curr         lexer.Token
	position     int
	readPosition int
}

func isInitiator(token lexer.Token) bool {
	return token.Type_ == lexer.Create || token.Type_ == lexer.Alter
}

func (p *Parser) readNext() {
	if p.readPosition >= len(p.Tokens) {
		p.curr = lexer.Token{Type_: lexer.Eof}
	} else {
		p.curr = p.Tokens[p.readPosition]
	}

	p.position = p.readPosition
	p.readPosition++
}

func (p *Parser) GetNextQuery() Query {
	if p.curr.Type_ == lexer.Eof {
		return Query{Action: lexer.Token{Type_: lexer.Eof}}
	}

	if isInitiator(p.curr) {
		action := p.curr

		p.readNext()
		p.readNext()

		tableName := p.curr

		for {
			p.readNext()
			if p.curr.Type_ == lexer.Eof || isInitiator(p.curr) {
				break
			}

			if !lexer.IsKeyword(p.curr) {
				continue
			}
			fmt.Println(p.curr.Literal)
		}

		return Query{Action: action, TableName: tableName}
	}
	p.readNext()
	return Query{Action: lexer.Token{Type_: lexer.Illegal}}
}

type Query struct {
	Action    lexer.Token
	TableName lexer.Token
	Columns   []Column
}

type Column struct {
	Name   lexer.Token
	Action lexer.Token
}

func CreateParser(tokenizer *lexer.Tokenizer) Parser {
	var tokens []lexer.Token
	for {
		token := tokenizer.GetNextToken()
		if token.Type_ == lexer.Eof {
			break
		}

		if token.Type_ == lexer.Illegal {
			continue
		}

		if token.Literal == "dbo" {
			continue
		}
		tokens = append(tokens, token)
	}
	parser := Parser{
		position:     0,
		readPosition: 0,
		Tokens:       tokens,
	}
	parser.readNext()
	return parser
}
