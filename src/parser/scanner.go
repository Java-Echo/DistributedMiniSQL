package parser

import (
	"bufio"
	"bytes"
	"io"
)

type Token int

// 词法分析中的初步结果(部分内容可以经过tokenize来形成语义更加明确的token)
const (
	// 特殊标记
	ILLEGAL Token = iota
	EOF
	WS // 空白字符
	// 常规类型数据
	IDENT   // ID，此时我们并不区分关键词，而是归类到同一类
	INTEGER // 整数
	FLOAT   // 浮点数
	STRING  // 字符串
	// 其他标记
	ASTERISK          // *
	COMMA             // ,
	LEFT_PARENTHESIS  // (
	RIGHT_PARENTHESIS // )
	SEMICOLON         // ;
	EQUAL             // =
	ANGLE_LEFT        // <
	ANGLE_LEFT_EQUAL  //<=
	ANGLE_RIGHT_EQUAL //>=
	ANGLE_RIGHT       // >
	NOT_EQUAL         // <> or !=
	POINT             //  .
)

type State int // 状态机的状态

const (
	STATE_INIT State = iota
	STATE_INTEGER
	STATE_POINT
	STATE_FRACTION
	STATE_IDENT
	STATE_ANGLE_LEFT
	STATE_ANGLE_RIGHT
	STATE_END
)

type CharType int // 单个字符的数据类型

const (
	NUM CharType = iota
	CHAR
	SPECIAL_SYMBOL
	ILLEGAL_SYMBOL
	SPACE
	UNDERLINE
)

// eof represents a marker rune for the end of the reader.
var eof = rune(0)

type InputScanner struct {
	r          *bufio.Reader
	apostropne bool // apostropne is true means
}

func NewScanner(r io.Reader) *InputScanner {
	return &InputScanner{r: bufio.NewReader(r), apostropne: false}
}

// scanner不断从输入流中读取数据，尝试拼接出一个个初步解析的token
func (s *InputScanner) Scan() (tok Token, lit string) {
	ch := s.read()
	var buf bytes.Buffer
	state := STATE_INIT
	for state != STATE_END {
		if checkCharType(ch) == ILLEGAL_SYMBOL {
			return ILLEGAL, string(ch)
		}
		// buf.WriteRune(ch)
		switch state {
		case STATE_INIT:
			switch checkCharType(ch) {
			case NUM:
				buf.WriteRune(ch)
				state = STATE_INTEGER
			case CHAR:
				buf.WriteRune(ch)
				state = STATE_IDENT
			case SPECIAL_SYMBOL:
				switch ch {
				case eof:
					return EOF, ""
				case '.':
					return POINT, string(ch)
				case '*':
					return ASTERISK, string(ch)
				case ',':
					return COMMA, string(ch)
				case '(':
					return LEFT_PARENTHESIS, string(ch)
				case ')':
					return RIGHT_PARENTHESIS, string(ch)
				case ';':
					return SEMICOLON, string(ch)
				case '=':
					return EQUAL, string(ch)
				case '<':
					buf.WriteRune(ch)
					state = STATE_ANGLE_LEFT
				case '>':
					buf.WriteRune(ch)
					state = STATE_ANGLE_RIGHT
				}
			case SPACE:
			case UNDERLINE:
				return ILLEGAL, string(ch)
			}
		case STATE_INTEGER:
			switch checkCharType(ch) {
			case NUM:
				buf.WriteRune(ch)
			case CHAR, SPACE, UNDERLINE:
				s.unread()
				return INTEGER, buf.String()
			case SPECIAL_SYMBOL:
				if ch == '.' {
					buf.WriteRune(ch)
					state = STATE_POINT
				} else {
					s.unread()
					return INTEGER, buf.String()
				}
			}
		case STATE_POINT:
			switch checkCharType(ch) {
			case NUM:
				buf.WriteRune(ch)
				state = STATE_FRACTION
			case CHAR, SPECIAL_SYMBOL, SPACE, UNDERLINE:
				return ILLEGAL, string(ch)
			}
		case STATE_FRACTION:
			switch checkCharType(ch) {
			case NUM:
				buf.WriteRune(ch)
			case CHAR, SPECIAL_SYMBOL, SPACE, UNDERLINE:
				s.unread()
				return FLOAT, buf.String()
			}
		case STATE_IDENT:
			switch checkCharType(ch) {
			case NUM, CHAR, UNDERLINE:
				buf.WriteRune(ch)
			case SPECIAL_SYMBOL, SPACE:
				s.unread()
				return IDENT, buf.String()
			}
		case STATE_ANGLE_LEFT:
			switch checkCharType(ch) {
			case NUM, CHAR, SPACE:
				s.unread()
				return ANGLE_LEFT, buf.String()
			case SPECIAL_SYMBOL:
				// ch = s.read()
				if ch == '=' {
					return ANGLE_LEFT_EQUAL, "<="
				} else if ch == '>' {
					return NOT_EQUAL, "<>"
				} else {
					s.unread()
					return ANGLE_LEFT, buf.String()
				}
			}
		case STATE_ANGLE_RIGHT:
			switch checkCharType(ch) {
			case NUM, CHAR, SPACE:
				s.unread()
				return ANGLE_RIGHT, buf.String()
			case SPECIAL_SYMBOL:
				// ch = s.read()
				if ch == '=' {
					return ANGLE_RIGHT_EQUAL, ">="
				} else {
					s.unread()
					return ANGLE_RIGHT, buf.String()
				}
			}
		}
		ch = s.read()
	}

	return ILLEGAL, string(ch)
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *InputScanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *InputScanner) unread() { _ = s.r.UnreadRune() }

func checkCharType(ch rune) CharType {
	if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' {
		return CHAR
	} else if ch >= '0' && ch <= '9' {
		// fmt.Println("检测到数字")
		return NUM
	} else if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
		return SPACE
	} else if ch == '.' || ch == '*' || ch == ',' || ch == '(' || ch == ')' || ch == ';' || ch == '=' || ch == '<' || ch == '>' || ch == eof {
		return SPECIAL_SYMBOL
	} else if ch == '_' {
		return UNDERLINE
	} else {
		return ILLEGAL_SYMBOL
	}
}
