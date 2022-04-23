package lexer

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
	scanner   *Scanner  // 处理输入的工具类
	tokenizer Tokenizer // 进一步赋予token信息的工具类
	Result    interface{}
}

// 存储类：存储token解析的结果，最终的Lex()主要利用这个对象来返回结果
type LexerResult struct {
	Token   int
	Literal string
}
