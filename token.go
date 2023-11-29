package jsonparser

type TokenType int

const (
	TEOF TokenType = iota

	// Type
	TString
	TNumber
	TBoolean
	TNull

	// Object
	TLBrace
	TRBrace
	// Array
	TLBracket
	TRBracket

	// delimiter
	TComma
	TColon
)

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
	}
}
