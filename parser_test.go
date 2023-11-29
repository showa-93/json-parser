package jsonparser

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    any
		wantErr bool
	}{
		{
			name:    "文字列",
			input:   "\"test\"",
			want:    Value{VString, "test"},
			wantErr: false,
		},
		{
			name:    "数値:123",
			input:   "123",
			want:    Value{VNumber, float64(123)},
			wantErr: false,
		},
		{
			name:    "数値:123.456",
			input:   "123.456",
			want:    Value{VNumber, float64(123.456)},
			wantErr: false,
		},
		{
			name:    "数値:123.456e+1",
			input:   "123.456e+1",
			want:    Value{VNumber, float64(123.456e+1)},
			wantErr: false,
		},
		{
			name:    "真偽値:true",
			input:   "true",
			want:    Value{VBoolean, true},
			wantErr: false,
		},
		{
			name:    "null",
			input:   "null",
			want:    Value{VNull, nil},
			wantErr: false,
		},
		{
			name:    "invalid token",
			input:   "}",
			wantErr: true,
		},
		{
			name:    "empty object",
			input:   "{}",
			want:    Value{VObject, map[string]Value{}},
			wantErr: false,
		},
		{
			name:  "object",
			input: `{"test1": "value", "test2":	123.456, "test3": {"test4": null}}`,
			want: Value{VObject, map[string]Value{
				"test1": {VString, "value"},
				"test2": {VNumber, float64(123.456)},
				"test3": {VObject, map[string]Value{
					"test4": {VNull, nil},
				}},
			}},
			wantErr: false,
		},
		{
			name:    "empty array",
			input:   `[]`,
			want:    Value{VArray, []Value{}},
			wantErr: false,
		},
		{
			name:  "array",
			input: `["1", 123, [null, 123.456]]`,
			want: Value{VArray,
				[]Value{
					{VString, "1"},
					{VNumber, float64(123)},
					{VArray,
						[]Value{
							{VNull, nil},
							{VNumber, float64(123.456)},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)
			got, err := p.Parse()
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("invalid parser want=%v, got=%v", tt.want, got)
			}
		})
	}
}
