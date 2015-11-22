
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
	NumStates = 99
	NumSymbols = 120
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
6: '_'
7: '-'
8: '*'
9: '/'
10: '_'
11: '-'
12: '*'
13: '('
14: ')'
15: 'b'
16: 'y'
17: ','
18: 's'
19: 'u'
20: 'm'
21: 'm'
22: 'e'
23: 'a'
24: 'n'
25: 'm'
26: 'e'
27: 'd'
28: 'i'
29: 'a'
30: 'n'
31: 'm'
32: 'i'
33: 'n'
34: 'm'
35: 'a'
36: 'x'
37: 's'
38: 't'
39: 'd'
40: 'd'
41: 'e'
42: 'v'
43: 'r'
44: 'a'
45: 't'
46: 'e'
47: 'r'
48: 'a'
49: 't'
50: 'e'
51: '_'
52: 's'
53: 'i'
54: 'g'
55: 'n'
56: 'e'
57: 'd'
58: 'l'
59: 'a'
60: 't'
61: 'e'
62: 's'
63: 't'
64: 'r'
65: 'e'
66: 'g'
67: 'u'
68: 'l'
69: 'a'
70: 'r'
71: 'i'
72: 'z'
73: 'e'
74: 'i'
75: 'n'
76: 't'
77: 'e'
78: 'r'
79: 'p'
80: 'o'
81: 'l'
82: 'a'
83: 't'
84: 'e'
85: 'p'
86: 'e'
87: 'r'
88: 'c'
89: 'e'
90: 'n'
91: 't'
92: 'i'
93: 'l'
94: 'e'
95: '['
96: ':'
97: ']'
98: '{'
99: '}'
100: '{'
101: '}'
102: '='
103: '/'
104: '/'
105: '\n'
106: '/'
107: '*'
108: '*'
109: '*'
110: '/'
111: ' '
112: '\t'
113: '\n'
114: '\r'
115: '1'-'9'
116: 'a'-'z'
117: 'A'-'Z'
118: '0'-'9'
119: .

*/
