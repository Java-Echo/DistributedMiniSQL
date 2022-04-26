package parser

import (
	"miniSQL/src/Interpreter/types"
	"strings"
	"testing"
)

var sql_strings = []string{
	"create table cxz(" +
		"afsdfsad int unique," +
		"what char(30) not null," +
		"primary key (what)" +
		");",
	"select a,b,c,d,e,f,g from cxz where a=123 and b=456 or c=234;",
}

func TestParser(t *testing.T) {
	ch := make(chan types.DStatements)
	go func() {
		// 通过通道通知main的goroutine
		for _, sql := range sql_strings {
			io := strings.NewReader(sql) // 组装io
			Parse(io, ch)
		}
	}()

	// for stmt := range ch {
	// 	log.Println(stmt.GetOperationType())
	// 	log.Println("有了")
	// }
	// data := <-ch
	// log.Print("能不能输出")
	// log.Print(data.GetOperationType())
	// cdb := data.(types.CreateTableStatement)
	// log.Print(cdb.TableName)
	t.Errorf("不要慌，我只是想要一个log而已\n")
}
