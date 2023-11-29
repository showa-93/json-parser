package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	jsonparser "github.com/showa-93/json-parser"
)

const INDENT = 4

func main() {
	if len(os.Args) != 2 {
		panic("引数にファイルパスをおくれ。")
	}

	filePath := os.Args[1]
	log.Println("read file path: ", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	l, err := jsonparser.NewLexer(f)
	if err != nil {
		panic(err)
	}
	p := jsonparser.NewParser(l)
	value, err := p.Parse()
	if err != nil {
		panic(err)
	}

	output(value, 0, false)
	fmt.Println()
}

func output(value jsonparser.Value, indent int, special bool) {
	switch value.Type {
	case jsonparser.VNumber:
		fmt.Print(fmt.Sprintf("%f", value.Data))
	case jsonparser.VBoolean:
		fmt.Print(fmt.Sprintf("%t", value.Data))
	case jsonparser.VString:
		fmt.Printf("\"%s\"", green(value.Data.(string)))
	case jsonparser.VArray:
		if special {
			fmt.Println("[")
		} else {
			fmt.Printf("%s[\n", strings.Repeat(" ", indent))
		}

		values := value.Data.([]jsonparser.Value)
		for i, v := range values {
			fmt.Printf("%s", strings.Repeat(" ", indent+INDENT))
			switch v.Type {
			case jsonparser.VArray, jsonparser.VObject:
				output(v, indent+INDENT, true)
			default:
				output(v, indent+INDENT, false)
			}

			if i != len(values)-1 {
				fmt.Print(",")
			}
			fmt.Println()
		}
		fmt.Printf("%s]", strings.Repeat(" ", indent))
	case jsonparser.VObject:
		if special {
			fmt.Println("{")
		} else {
			fmt.Printf("%s{\n", strings.Repeat(" ", indent))
		}

		var i int
		values := value.Data.(map[string]jsonparser.Value)
		for k, v := range values {
			i++
			fmt.Printf("%s\"%s\": ", strings.Repeat(" ", indent+INDENT), yellow(k))
			switch v.Type {
			case jsonparser.VArray, jsonparser.VObject:
				output(v, indent+INDENT, true)
			default:
				output(v, indent+INDENT, false)
			}

			if i != len(values) {
				fmt.Print(",")
			}
			fmt.Println()
		}
		fmt.Printf("%s}", strings.Repeat(" ", indent))
	case jsonparser.VNull:
		fmt.Print(red("null"))
	}
}

func red(s string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[m", s)
}
func green(s string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[m", s)
}
func yellow(s string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[m", s)
}
