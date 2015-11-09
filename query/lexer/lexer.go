
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
	NumStates = 77
	NumSymbols = 96
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
0: '/'
1: '/'
2: '_'
3: '-'
4: '/'
5: '/'
6: '_'
7: '-'
8: '('
9: ')'
10: 'b'
11: 'y'
12: ','
13: 'p'
14: 'e'
15: 'r'
16: 'c'
17: 'e'
18: 'n'
19: 't'
20: 'i'
21: 'l'
22: 'e'
23: 's'
24: 'u'
25: 'm'
26: 'm'
27: 'e'
28: 'a'
29: 'n'
30: 'm'
31: 'e'
32: 'd'
33: 'i'
34: 'a'
35: 'n'
36: 'm'
37: 'i'
38: 'n'
39: 'm'
40: 'a'
41: 'x'
42: 's'
43: 't'
44: 'd'
45: 'd'
46: 'e'
47: 'v'
48: 'r'
49: 'a'
50: 't'
51: 'e'
52: 'r'
53: 'a'
54: 't'
55: 'e'
56: '_'
57: 's'
58: 'i'
59: 'g'
60: 'n'
61: 'e'
62: 'd'
63: 'd'
64: 'e'
65: 'l'
66: 't'
67: 'a'
68: 'l'
69: 'a'
70: 't'
71: 'e'
72: 's'
73: 't'
74: '{'
75: '}'
76: '{'
77: '}'
78: '='
79: '/'
80: '/'
81: '\n'
82: '/'
83: '*'
84: '*'
85: '*'
86: '/'
87: ' '
88: '\t'
89: '\n'
90: '\r'
91: '1'-'9'
92: 'a'-'z'
93: 'A'-'Z'
94: '0'-'9'
95: .

*/
