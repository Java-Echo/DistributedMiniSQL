package lexer

import (
	"testing"
)

func TestFromStrLit(t *testing.T) {
	cases := []struct {
		lit       string
		TokenType Token
		lastToken int
		expect    int
	}{
		{
			"10",
			T_INTEGER,
			0,
			decimal_value,
		},
		{
			"0xff",
			T_INTEGER,
			0,
			hex_value,
		},
		// Not hex value
		{
			"ff",
			T_INTEGER,
			0,
			0,
		},
	}

	tk := keywordTokenizer{}
	for _, c := range cases {
		actual := tk.FromStrLit(c.lit, c.TokenType, c.lastToken)
		if actual != c.expect {
			t.Errorf("Expected: %v, but actual: %v\n", c.expect, actual)
		}
	}
}
