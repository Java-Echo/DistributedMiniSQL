package main

import (
	"fmt"
	"strings"

	"miniSQL/src/Interpreter/lexer"
)

func main() {
	s := strings.NewReader("creat table tabel_name where key <= 1")
	l := lexer.NewScanner(s)
	for {
		tok, str := l.Scan()
		fmt.Print(tok)
		fmt.Print("  ")
		fmt.Println(str)
		if tok == 1 {
			break
		}
	}
}
