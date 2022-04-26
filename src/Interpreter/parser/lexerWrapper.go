package parser

import (
	"log"
	"miniSQL/src/Interpreter/lexer"
	"miniSQL/src/Interpreter/types"
)

// lex的包装类
type lexerWrapper struct {
	impl        *lexer.LexerImpl
	channelSend chan<- types.DStatements
	lastLiteral string // 向前看一位
	err         error
}

func newLexerWrapper(li *lexer.LexerImpl, channel chan<- types.DStatements) *lexerWrapper {
	return &lexerWrapper{
		impl:        li,
		channelSend: channel,
	}
}

// 真正实现了parser需要的Lex接口
func (l *lexerWrapper) Lex(lval *yySymType) int {
	r, err := l.impl.Lex(lval.LastToken)
	if err != nil {
		log.Fatal(err)
	}
	l.lastLiteral = r.Literal

	tokVal := r.Token
	lval.str = r.Literal
	lval.LastToken = tokVal
	return tokVal
}

// 真正实现了parser需要的Error接口
// ToDo:完善错误提示机制
func (l *lexerWrapper) Error(errStr string) {
	// l.err = wrapParseError(l.lastLiteral, errStr)
}
