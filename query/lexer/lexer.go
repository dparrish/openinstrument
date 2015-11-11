
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
	NumStates = 84
	NumSymbols = 104
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
18: 'p'
19: 'e'
20: 'r'
21: 'c'
22: 'e'
23: 'n'
24: 't'
25: 'i'
26: 'l'
27: 'e'
28: 's'
29: 'u'
30: 'm'
31: 'm'
32: 'e'
33: 'a'
34: 'n'
35: 'm'
36: 'e'
37: 'd'
38: 'i'
39: 'a'
40: 'n'
41: 'm'
42: 'i'
43: 'n'
44: 'm'
45: 'a'
46: 'x'
47: 's'
48: 't'
49: 'd'
50: 'd'
51: 'e'
52: 'v'
53: 'r'
54: 'a'
55: 't'
56: 'e'
57: 'r'
58: 'a'
59: 't'
60: 'e'
61: '_'
62: 's'
63: 'i'
64: 'g'
65: 'n'
66: 'e'
67: 'd'
68: 'd'
69: 'e'
70: 'l'
71: 't'
72: 'a'
73: 'l'
74: 'a'
75: 't'
76: 'e'
77: 's'
78: 't'
79: '['
80: ':'
81: ']'
82: '{'
83: '}'
84: '{'
85: '}'
86: '='
87: '/'
88: '/'
89: '\n'
90: '/'
91: '*'
92: '*'
93: '*'
94: '/'
95: ' '
96: '\t'
97: '\n'
98: '\r'
99: '1'-'9'
100: 'a'-'z'
101: 'A'-'Z'
102: '0'-'9'
103: .

*/
