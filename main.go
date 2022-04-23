package main

import (
	. "miniSQL/src/Interpreter/lexer"
	"strings"
)

// // 解析的结果
// type LexerResult struct {
// 	Token   int
// 	Literal string
// }

func main() {
	s := strings.NewReader("creat table tabel_name where key <= 1")
	l := NewScanner(s)
	p := NewTokenizer()
	result := &LexerResult{}
	lastToken := 0
	for {
		tok, str := l.Scan()
		// fmt.Print(tok)
		// fmt.Print("  ")
		// fmt.Println(str)
		// if tok == 1 {
		// 	break
		// }
		switch tok {
		case T_EOF:
			// Stop lex
		case T_IDENT, T_INTEGER, T_FLOAT, T_STRING, T_LEFT_PARENTHESIS, T_RIGHT_PARENTHESIS, T_COMMA, T_SEMICOLON, T_EQUAL, T_ANGLE_LEFT, T_ANGLE_RIGHT, T_ANGLE_LEFT_EQUAL, T_ANGLE_RIGHT_EQUAL, T_NOT_EQUAL, T_ASTERISK, T_POINT:
			result.Literal = str
			// default:
			// 	log.Printf("UnexpectedToken: tok is %d, lit is %s\n", tok, lit)
			// 	return nil, UnexpectedTokenErr
		}

		result.Token = p.FromStrLit(str, tok, lastToken)
	}
}
