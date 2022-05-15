package interpreter

import (
	"miniSQL/src/Interpreter/parser"
	"miniSQL/src/Interpreter/types"
	"strings"
)

func CreateInterpreter(input chan string) chan types.DStatements {
	output := make(chan types.DStatements)
	go func() {
		for {
			sql := <-input
			io := strings.NewReader(sql)
			parser.Parse(io, output)
		}
	}()
	return output
}
