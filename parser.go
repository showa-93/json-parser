package jsonparser

import (
	"fmt"
	"strconv"
)

type ParserError struct {
	msg string
}

func NewParserError(msg string) error {
	return &ParserError{
		msg: msg,
	}
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("Parser Error: %s", e.msg)
}

type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Parse() (any, error) {
	switch p.curToken.Type {
	case TLBrace:
		return p.parseObject()
	case TLBracket:
		return p.parseArray()
	case TString:
		return p.curToken.Literal, nil
	case TNumber:
		v, _ := strconv.ParseFloat(p.curToken.Literal, 64)
		return v, nil
	case TBoolean:
		v, _ := strconv.ParseBool(p.curToken.Literal)
		return v, nil
	case TNull:
		return nil, nil
	default:
		return nil, NewParserError("token must start { or [ or string or number or boolean or null")
	}
}

func (p *Parser) parseObject() (any, error) {
	if err := p.nextToken(); err != nil {
		return nil, err
	}
	if p.curTokenIs(TRBrace) {
		return make(map[string]any), nil
	}

	v := make(map[string]any)
	for {
		// トークンが "key:" の並びになっていることを確認
		if !p.curTokenIs(TString) {
			return nil, NewParserError("key expected string token")
		}
		key := p.curToken.Literal
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		if !p.curTokenIs(TColon) {
			return nil, NewParserError("expected colon token after key")
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}

		// valueをパースする
		value, err := p.Parse()
		if err != nil {
			return nil, err
		}

		v[key] = value

		// , か }であることを確認する
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		if p.curTokenIs(TRBrace) {
			return v, nil
		}
		if !p.curTokenIs(TComma) {
			return nil, NewParserError("expected comma or } in object")
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
	}
}

func (p *Parser) parseArray() (any, error) {
	if err := p.nextToken(); err != nil {
		return nil, err
	}
	if p.curTokenIs(TRBracket) {
		return make([]any, 0), nil
	}

	values := make([]any, 0)
	for {
		value, err := p.Parse()
		if err != nil {
			return nil, err
		}

		values = append(values, value)
		if err := p.nextToken(); err != nil {
			return nil, err
		}

		if p.curTokenIs(TRBracket) {
			return values, nil
		}
		if !p.curTokenIs(TComma) {
			return nil, NewParserError("expected comma or ] in array")
		}

		if err := p.nextToken(); err != nil {
			return nil, err
		}
	}
}

func (p *Parser) nextToken() (err error) {
	p.curToken = p.peekToken
	p.peekToken, err = p.lexer.NextToken()
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}
