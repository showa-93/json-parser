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

func (p *Parser) Parse() (Value, error) {
	switch p.curToken.Type {
	case TLBrace:
		return p.parseObject()
	case TLBracket:
		return p.parseArray()
	case TString:
		return Value{VString, p.curToken.Literal}, nil
	case TNumber:
		v, _ := strconv.ParseFloat(p.curToken.Literal, 64)
		return Value{VNumber, v}, nil
	case TBoolean:
		v, _ := strconv.ParseBool(p.curToken.Literal)
		return Value{VBoolean, v}, nil
	case TNull:
		return Value{VNull, nil}, nil
	default:
		return Value{}, NewParserError("token must start { or [ or string or number or boolean or null")
	}
}

func (p *Parser) parseObject() (Value, error) {
	if err := p.nextToken(); err != nil {
		return Value{}, err
	}
	if p.curTokenIs(TRBrace) {
		return Value{VObject, make(map[string]Value)}, nil
	}

	v := make(map[string]Value)
	for {
		// トークンが "key:" の並びになっていることを確認
		if !p.curTokenIs(TString) {
			return Value{}, NewParserError("key expected string token")
		}
		key := p.curToken.Literal
		if err := p.nextToken(); err != nil {
			return Value{}, err
		}
		if !p.curTokenIs(TColon) {
			return Value{}, NewParserError("expected colon token after key")
		}
		if err := p.nextToken(); err != nil {
			return Value{}, err
		}

		// valueをパースする
		value, err := p.Parse()
		if err != nil {
			return Value{}, err
		}

		v[key] = value

		// , か }であることを確認する
		if err := p.nextToken(); err != nil {
			return Value{}, err
		}
		if p.curTokenIs(TRBrace) {
			return Value{VObject, v}, nil
		}
		if !p.curTokenIs(TComma) {
			return Value{}, NewParserError("expected comma or } in object")
		}
		if err := p.nextToken(); err != nil {
			return Value{}, err
		}
	}
}

func (p *Parser) parseArray() (Value, error) {
	if err := p.nextToken(); err != nil {
		return Value{}, err
	}
	if p.curTokenIs(TRBracket) {
		return Value{VArray, make([]Value, 0)}, nil
	}

	values := make([]Value, 0)
	for {
		value, err := p.Parse()
		if err != nil {
			return Value{}, err
		}

		values = append(values, value)
		if err := p.nextToken(); err != nil {
			return Value{}, err
		}

		if p.curTokenIs(TRBracket) {
			return Value{VArray, values}, nil
		}
		if !p.curTokenIs(TComma) {
			return Value{}, NewParserError("expected comma or ] in array")
		}

		if err := p.nextToken(); err != nil {
			return Value{}, err
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
