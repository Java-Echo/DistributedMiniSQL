package main

import (
	"fmt"
	"strings"

	"miniSQL/src/parser"
)

func main() {
	s := strings.NewReader("creat table tabel_name where key <= 1")
	l := parser.NewScanner(s)
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
