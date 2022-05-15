package lexer

import (
	"io"
	"log"
)

type Token int

// 工具类：接受由Scanner传入的经过基本处理的token，通过lookahead来赋予该token更多信息
type Tokenizer interface {
	FromStrLit(lit string, TokenType Token, lastToken int) int
}

// 工具类：初步处理输入，得到基础的token
type Scanner interface {
	Scan() (tok Token, lit string)
}

// 模块类：将输入解析为token的总体工具类
type LexerImpl struct {
	scanner   Scanner   // 处理输入的工具类
	tokenizer Tokenizer // 进一步赋予token信息的工具类
	Result    interface{}
}

// 存储类：存储token解析的结果，最终的Lex()主要利用这个对象来返回结果
type LexerResult struct {
	Token   int
	Literal string
}

// 新建一个“将输入解析为token”的总体工具类的实例化对象
func NewLexerImpl(r io.Reader) *LexerImpl {
	return &LexerImpl{
		scanner:   NewScanner(r),
		tokenizer: NewTokenizer(),
	}
}

func (li *LexerImpl) Lex(lastToken int) (*LexerResult, error) {
	result := &LexerResult{}

	tok, lit := li.scanner.Scan() // 这里的scanner已经完成了输入流的基础token化

	switch tok {
	case T_EOF:
		// Stop lex
	case T_IDENT, T_INTEGER, T_FLOAT, T_STRING, T_LEFT_PARENTHESIS, T_RIGHT_PARENTHESIS, T_COMMA, T_SEMICOLON, T_EQUAL, T_ANGLE_LEFT, T_ANGLE_RIGHT, T_ANGLE_LEFT_EQUAL, T_ANGLE_RIGHT_EQUAL, T_NOT_EQUAL, T_ASTERISK, T_POINT:
		result.Literal = lit
	default:
		log.Printf("UnexpectedToken: tok is %d, lit is %s\n", tok, lit)
		// return nil, UnexpectedTokenErr		// ToDo:这里的错误机制需要完善
		return nil, nil
	}

	result.Token = li.tokenizer.FromStrLit(lit, tok, lastToken)

	return result, nil
}
