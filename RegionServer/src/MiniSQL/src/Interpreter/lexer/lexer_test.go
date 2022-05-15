package lexer

import (
	"log"
	"strings"
	"testing"
)

var sql_strings = []string{
	// "create table 1_a",
	" ",
	"create table cxz(" +
		"afsdfsad int unique," +
		"what char(30) not null," +
		"primary key (what)" +
		");",
	"select a,b,c,d,e,f,g from cxz where a=123 and b=456 or c=234;",
}

//DoTo:这里有个bug就是我们无法正确地停止输入
func TestLexerLex(t *testing.T) {
	for _, str := range sql_strings {
		LastToken := 0
		io := strings.NewReader(str) // 组装io
		impl := NewLexerImpl(io)     // 组装待测试的LexerImpl
		for LastToken != int(T_EOF) {
			r, _ := impl.Lex(LastToken)

			tokVal := r.Token
			literal := r.Literal
			LastToken = tokVal

			// log.Print(tokVal)
			// log.Print("  ")
			log.Print(literal)
			// log.Print("  ")
			// log.Print(LastToken)
			if tokVal == 0 {
				// 检测是否到达输入末尾
				break
			}
		}

	}
	t.Errorf("不要慌，我只是想要一个log而已\n")
}
