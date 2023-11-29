package jsonparser

import (
	"fmt"
	"io"
	"strings"
)

type LexerError struct {
	msg string
}

func NewLexerError(msg string) error {
	return &LexerError{
		msg: msg,
	}
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("Lexer Error: %s", e.msg)
}

type Lexer struct {
	r    io.Reader
	cur  byte
	peek byte
	eof  bool
}

func NewLexer(reader io.Reader) *Lexer {
	l := &Lexer{
		r: reader,
	}
	l.readChar()

	return l
}

func (l *Lexer) NextToken() (Token, error) {
	ch, err := l.readChar()
	if err != nil {
		return Token{}, err
	}

	// 今回はwhitespaceはskip
	for isWhitespace(ch) {
		ch, err = l.readChar()
		if err != nil {
			return Token{}, err
		}
	}

	switch ch {
	case 0:
		return NewToken(TEOF, ""), nil
	// Object
	case '{':
		return NewToken(TLBrace, "{"), nil
	case '}':
		return NewToken(TRBrace, "}"), nil
	// Array
	case '[':
		return NewToken(TLBracket, "["), nil
	case ']':
		return NewToken(TRBracket, "]"), nil
	// delimiter
	case ',':
		return NewToken(TComma, ","), nil
	case ':':
		return NewToken(TColon, ":"), nil
	// Type
	case '"':
		return l.readString()
	case 'n':
		return l.readNull()
	case 't':
		return l.readBoolean(true)
	case 'f':
		return l.readBoolean(false)
	default:
		return l.readNumber(ch)
	}
}

func (l *Lexer) readString() (Token, error) {
	// https://www.rfc-editor.org/rfc/rfc8259#section-7
	var sb strings.Builder
	for {
		ch, err := l.readChar()
		if err != nil {
			return Token{}, err
		}

		switch ch {
		case '"':
			return NewToken(TString, sb.String()), nil
		case '\\':
			sb.WriteByte(ch)
			// エスケープ文字を処理
			ch, err = l.readChar()
			if err != nil {
				return Token{}, err
			}
			sb.WriteByte(ch)

			if isEscaped(ch) {
				// エスケープ対象のため読み取りを継続
			} else if ch == 'u' {
				// サロゲートペアの処理はしていない
				// http://hirakun.blog57.fc2.com/blog-entry-215.html
				for i := 0; i < 4; i++ {
					ch, err = l.readChar()
					if err != nil {
						return Token{}, err
					}
					if !isHexDigit(ch) {
						return Token{}, NewLexerError(fmt.Sprintf("unexpected characters %c", ch))
					}
					sb.WriteByte(ch)
				}
			} else {
				return Token{}, NewLexerError(fmt.Sprintf("unexpected characters %c", ch))
			}
		case 0:
			return Token{}, NewLexerError("missing \"")
		default:
			sb.WriteByte(ch)
		}
	}
}

func (l *Lexer) readNull() (Token, error) {
	ch1, err := l.readChar()
	if err != nil {
		return Token{}, err
	}
	ch2, err := l.readChar()
	if err != nil {
		return Token{}, err
	}
	ch3, err := l.readChar()
	if err != nil {
		return Token{}, err
	}

	if ch1 == 'u' && ch2 == 'l' && ch3 == 'l' {
		return NewToken(TNull, "null"), nil
	}

	return Token{}, NewLexerError(fmt.Sprintf("unexpected characters %c", 'n'))
}

func (l *Lexer) readBoolean(b bool) (Token, error) {
	if b {
		ch1, err := l.readChar()
		if err != nil {
			return Token{}, err
		}
		ch2, err := l.readChar()
		if err != nil {
			return Token{}, err
		}
		ch3, err := l.readChar()
		if err != nil {
			return Token{}, err
		}
		if ch1 == 'r' && ch2 == 'u' && ch3 == 'e' {
			return NewToken(TBoolean, "true"), nil
		}

		return Token{}, NewLexerError(fmt.Sprintf("unexpected characters %c", 't'))
	} else {
		ch1, err := l.readChar()
		if err != nil {
			return Token{}, err
		}
		ch2, err := l.readChar()
		if err != nil {
			return Token{}, err
		}
		ch3, err := l.readChar()
		if err != nil {
			return Token{}, err
		}
		ch4, err := l.readChar()
		if err != nil {
			return Token{}, err
		}
		if ch1 == 'a' && ch2 == 'l' && ch3 == 's' && ch4 == 'e' {
			return NewToken(TBoolean, "false"), nil
		}

		return Token{}, NewLexerError(fmt.Sprintf("unexpected characters %c", 'f'))
	}
}

func (l *Lexer) readNumber(ch byte) (Token, error) {
	// https://www.rfc-editor.org/rfc/rfc8259#section-6
	var (
		sb  strings.Builder
		err error
	)
	sb.WriteByte(ch)

	if isSign(ch) {
		ch, err = l.readChar()
		if err != nil {
			return Token{}, err
		}
		sb.WriteByte(ch)
	}

	if !isDigit(ch) {
		return Token{}, NewLexerError(fmt.Sprintf("unexpected characters %c", ch))
	}

	if ch != '0' {
		for isDigit(l.peekChar()) {
			ch, err = l.readChar()
			if err != nil {
				return Token{}, err
			}
			sb.WriteByte(ch)
		}
	}

	if isDot(l.peekChar()) {
		ch, err = l.readChar()
		if err != nil {
			return Token{}, err
		}
		sb.WriteByte(ch)

		ch, err = l.readChar()
		if err != nil {
			return Token{}, err
		}
		sb.WriteByte(ch)
		if !isDigit(ch) {
			return Token{}, NewLexerError(fmt.Sprintf("unexpected characters %c", ch))
		}
		for isDigit(l.peekChar()) {
			ch, err = l.readChar()
			if err != nil {
				return Token{}, err
			}
			sb.WriteByte(ch)
		}
	}

	if l.peekChar() == 'e' {
		ch, err = l.readChar()
		if err != nil {
			return Token{}, err
		}
		sb.WriteByte(ch)
		ch, err = l.readChar()
		if err != nil {
			return Token{}, err
		}
		sb.WriteByte(ch)
		if !isSign(ch) {
			return Token{}, NewLexerError(fmt.Sprintf("unexpected characters %c", ch))
		}
		ch, err = l.readChar()
		if err != nil {
			return Token{}, err
		}
		sb.WriteByte(ch)
		if !isDigit(ch) {
			return Token{}, NewLexerError(fmt.Sprintf("unexpected characters %c", ch))
		}
		for isDigit(l.peekChar()) {
			ch, err = l.readChar()
			if err != nil {
				return Token{}, err
			}
			sb.WriteByte(ch)
		}
	}

	return NewToken(TNumber, sb.String()), nil
}

func (l *Lexer) readChar() (byte, error) {
	if l.eof {
		return byte(0), nil
	}

	var ch [1]byte
	n, err := l.r.Read(ch[:])
	if err != nil {
		if err != io.EOF {
			return byte(0), err
		}
		l.eof = true
	}
	if n == 0 {
		l.eof = true
	}

	l.cur = l.peek
	l.peek = ch[0]

	return l.cur, nil
}

func (l *Lexer) peekChar() byte {
	if l.eof {
		return byte(0)
	}

	return l.peek
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isEscaped(ch byte) bool {
	return ch == '"' || ch == '\\' || ch == '/' || ch == 'b' || ch == 'f' || ch == 'r' || ch == 'n' || ch == 't'
}

func isHexDigit(ch byte) bool {
	return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isSign(ch byte) bool {
	return ch == '+' || ch == '-'
}

func isDot(ch byte) bool {
	return ch == '.'
}
