package jsonparser

import (
	"strings"
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Token
		wantErr bool
	}{
		{
			name:  "object",
			input: `{"test1": "value", "test2":	123.456, "test3": null}`,
			want: []Token{
				{TLBrace, "{"},
				{TString, "test1"},
				{TColon, ":"},
				{TString, "value"},
				{TComma, ","},
				{TString, "test2"},
				{TColon, ":"},
				{TNumber, "123.456"},
				{TComma, ","},
				{TString, "test3"},
				{TColon, ":"},
				{TNull, "null"},
				{TRBrace, "}"},
			},
			wantErr: false,
		},
		{
			name:  "array",
			input: `[true, false]`,
			want: []Token{
				{TLBracket, "["},
				{TBoolean, "true"},
				{TComma, ","},
				{TBoolean, "false"},
				{TRBracket, "]"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(strings.NewReader(tt.input))
			gots := make([]Token, 0, len(tt.want))
			for {
				got, err := l.NextToken()
				if err != nil {
					if tt.wantErr {
						break
					}
					t.Fatal(err)
				}
				if got.Type == TEOF {
					break
				}
				gots = append(gots, got)
			}
			if len(gots) != len(tt.want) {
				t.Errorf("invalid token length want=%v, got=%v", len(tt.want), len(gots))
			}
			for i, got := range gots {
				if len(tt.want)-1 < i {
					break
				}
				if got != tt.want[i] {
					t.Errorf("invalid token want=%v, got=%v", tt.want[i], got)
				}
			}
		})
	}
}

func TestReadString(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    Token
		wantErr bool
	}{
		{"empty", "\"\"", Token{TString, ""}, false},
		{"string", "\"test\"", Token{TString, "test"}, false},
		{"escape \\", "\"test\\\\\"", Token{TString, "test\\\\"}, false},
		{"escape \\/", "\"test\\/\"", Token{TString, "test\\/"}, false},
		{"escape \\r", "\"test\\r\"", Token{TString, "test\\r"}, false},
		{"escape \\n", "\"test\\n\"", Token{TString, "test\\n"}, false},
		{"escape \\t", "\"test\\t\"", Token{TString, "test\\t"}, false},
		{"escape \\b", "\"test\\b\"", Token{TString, "test\\b"}, false},
		{"escape \\f", "\"test\\f\"", Token{TString, "test\\f"}, false},
		{"escape unicode", "\"\\u6628\\u65E5open\\u3057\\u305F\\uD842\\uDFB7\\u91CE\\u5C4B\"", Token{TString, "\\u6628\\u65E5open\\u3057\\u305F\\uD842\\uDFB7\\u91CE\\u5C4B"}, false},
		{"not quaote", "\"test", Token{}, true},
		{"not unicode", "\"\\u662\"", Token{}, true},
		{"not escape char", "\"\\x\"", Token{}, true},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(strings.NewReader(tt.input))
			l.readChar()
			got, err := l.readString()
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("invalid token want=%v, got=%v", tt.want, got)
			}
		})
	}
}

func TestReadNull(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    Token
		wantErr bool
	}{
		{"null", "null", Token{TNull, "null"}, false},
		{"nulx", "nulx", Token{}, true},
		{"empty", "n", Token{}, true},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(strings.NewReader(tt.input))
			l.readChar()
			got, err := l.readNull()
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("invalid token want=%v, got=%v", tt.want, got)
			}
		})
	}
}

func TestReadBoolean(t *testing.T) {
	testCases := []struct {
		name    string
		target  bool
		input   string
		want    Token
		wantErr bool
	}{
		{"true", true, "true", Token{TBoolean, "true"}, false},
		{"false", false, "false", Token{TBoolean, "false"}, false},
		{"tr", true, "tr", Token{}, true},
		{"fa", false, "fa", Token{}, true},
		{"empty", true, "", Token{}, true},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(strings.NewReader(tt.input))
			l.readChar()
			got, err := l.readBoolean(tt.target)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("invalid token want=%v, got=%v", tt.want, got)
			}
		})
	}
}

func TestReadNumber(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    Token
		wantErr bool
	}{
		{"0", "0", Token{TNumber, "0"}, false},
		{"123", "123", Token{TNumber, "123"}, false},
		{"0123", "0123", Token{TNumber, "0"}, false},
		{"+0.456", "+0.456", Token{TNumber, "+0.456"}, false},
		{"-123.456", "-123.456", Token{TNumber, "-123.456"}, false},
		{"2e+0", "2e+0", Token{TNumber, "2e+0"}, false},
		{"2e-1", "2e-01", Token{TNumber, "2e-01"}, false},
		{"1.2e+012", "1.2e+012", Token{TNumber, "1.2e+012"}, false},
		{"a", "a", Token{}, true},
		{"1.a", "1.a", Token{}, true},
		{"1e1", "1e1", Token{}, true},
		{"1e+", "1e+", Token{}, true},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(strings.NewReader(tt.input))
			ch, err := l.readChar()
			if err != nil {
				t.Fatal(err)
			}
			got, err := l.readNumber(ch)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("invalid token want=%v, got=%v", tt.want, got)
			}
		})
	}
}
