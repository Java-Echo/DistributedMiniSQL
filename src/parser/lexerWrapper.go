package parser

import (
	"miniSQL/src/lexer"
)

// lex的包装类
type lexerWrapper struct {
	impl *lexer.LexerImpl
	// channelSend chan<- types.DStatements
	lastLiteral string // 向前看一位
	err         error
}
