
package lexer

import (
	
	// "fmt"
	// "github.com/dparrish/openinstrument/query/util"
	
	"io/ioutil"
	"unicode/utf8"
	"github.com/dparrish/openinstrument/query/token"
)

const(
	NoState = -1
	NumStates = 82
	NumSymbols = 100
) 

type Lexer struct {
	src             []byte
	pos             int
	line            int
	column          int
}

func NewLexer(src []byte) *Lexer {
	lexer := &Lexer{
		src:    src,
		pos:    0,
		line:   1,
		column: 1,
	}
	return lexer
}

func NewLexerFile(fpath string) (*Lexer, error) {
	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return NewLexer(src), nil
}

func (this *Lexer) Scan() (tok *token.Token) {
	
	// fmt.Printf("Lexer.Scan() pos=%d\n", this.pos)
	
	tok = new(token.Token)
	if this.pos >= len(this.src) {
		tok.Type = token.EOF
		tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = this.pos, this.line, this.column
		return
	}
	start, end := this.pos, 0
	tok.Type = token.INVALID
	state, rune1, size := 0, rune(-1), 0
	for state != -1 {
	
		// fmt.Printf("\tpos=%d, line=%d, col=%d, state=%d\n", this.pos, this.line, this.column, state)
	
		if this.pos >= len(this.src) {
			rune1 = -1
		} else {
			rune1, size = utf8.DecodeRune(this.src[this.pos:])
			this.pos += size
		}
		switch rune1 {
		case '\n':
			this.line++
			this.column = 1
		case '\r':
			this.column = 1
		case '\t':
			this.column += 4
		default:
			this.column++
		}

	
		// Production start
		if rune1 != -1 {
			state = TransTab[state](rune1)
		} else {
			state = -1
		}
		// Production end

		// Debug start
		// nextState := -1
		// if rune1 != -1 {
		// 	nextState = TransTab[state](rune1)
		// }
		// fmt.Printf("\tS%d, : tok=%s, rune == %s(%x), next state == %d\n", state, token.TokMap.Id(tok.Type), util.RuneToString(rune1), rune1, nextState)
		// fmt.Printf("\t\tpos=%d, size=%d, start=%d, end=%d\n", this.pos, size, start, end)
		// if nextState != -1 {
		// 	fmt.Printf("\t\taction:%s\n", ActTab[nextState].String())
		// }
		// state = nextState
		// Debug end
	

		if state != -1 {
			switch {
			case ActTab[state].Accept != -1:
				tok.Type = ActTab[state].Accept
				// fmt.Printf("\t Accept(%s), %s(%d)\n", string(act), token.TokMap.Id(tok), tok)
				end = this.pos
			case ActTab[state].Ignore != "":
				// fmt.Printf("\t Ignore(%s)\n", string(act))
				start = this.pos
				state = 0
				if start >= len(this.src) {
					tok.Type = token.EOF
				}

			}
		} else {
			if tok.Type == token.INVALID {
				end = this.pos
			}
		}
	}
	if end > start {
		this.pos = end
		tok.Lit = this.src[start:end]
	} else {
		tok.Lit = []byte{}
	}
	tok.Pos.Offset = start
	tok.Pos.Column = this.column
	tok.Pos.Line = this.line
	return
}

func (this *Lexer) Reset() {
	this.pos = 0
}

/*
Lexer symbols:
0: '-'
1: '/'
2: '/'
3: '_'
4: '-'
5: '/'
6: '/'
7: '_'
8: '-'
9: '('
10: ')'
11: 'b'
12: 'y'
13: ','
14: 'p'
15: 'e'
16: 'r'
17: 'c'
18: 'e'
19: 'n'
20: 't'
21: 'i'
22: 'l'
23: 'e'
24: 's'
25: 'u'
26: 'm'
27: 'm'
28: 'e'
29: 'a'
30: 'n'
31: 'm'
32: 'e'
33: 'd'
34: 'i'
35: 'a'
36: 'n'
37: 'm'
38: 'i'
39: 'n'
40: 'm'
41: 'a'
42: 'x'
43: 's'
44: 't'
45: 'd'
46: 'd'
47: 'e'
48: 'v'
49: 'r'
50: 'a'
51: 't'
52: 'e'
53: 'r'
54: 'a'
55: 't'
56: 'e'
57: '_'
58: 's'
59: 'i'
60: 'g'
61: 'n'
62: 'e'
63: 'd'
64: 'd'
65: 'e'
66: 'l'
67: 't'
68: 'a'
69: 'l'
70: 'a'
71: 't'
72: 'e'
73: 's'
74: 't'
75: '['
76: ':'
77: ']'
78: '{'
79: '}'
80: '{'
81: '}'
82: '='
83: '/'
84: '/'
85: '\n'
86: '/'
87: '*'
88: '*'
89: '*'
90: '/'
91: ' '
92: '\t'
93: '\n'
94: '\r'
95: '1'-'9'
96: 'a'-'z'
97: 'A'-'Z'
98: '0'-'9'
99: .

*/
