
package token

import(
	"fmt"
)

type Token struct {
	Type
	Lit []byte
	Pos
}

type Type int

const(
	INVALID Type = iota
	EOF
)

type Pos struct {
	Offset int
	Line int
	Column int
}

func (this Pos) String() string {
	return fmt.Sprintf("Pos(offset=%d, line=%d, column=%d)", this.Offset, this.Line, this.Column)
}

type TokenMap struct {
	typeMap  []string
	idMap map[string]Type
}

func (this TokenMap) Id(tok Type) string {
	if int(tok) < len(this.typeMap) {
		return this.typeMap[tok]
	}
	return "unknown"
}

func (this TokenMap) Type(tok string) Type {
	if typ, exist := this.idMap[tok]; exist {
		return typ
	}
	return INVALID
}

func (this TokenMap) TokenString(tok *Token) string {
	//TODO: refactor to print pos & token string properly
	return fmt.Sprintf("%s(%d,%s)", this.Id(tok.Type), tok.Type, tok.Lit)
}

func (this TokenMap) StringType(typ Type) string {
	return fmt.Sprintf("%s(%d)", this.Id(typ), typ)
}

var TokMap = TokenMap{
	typeMap: []string{
		"INVALID",
		"$",
		"(",
		")",
		"by",
		",",
		"varlabel",
		"percentile",
		"sum",
		"mean",
		"median",
		"min",
		"max",
		"stddev",
		"rate",
		"rate_signed",
		"delta",
		"latest",
		"[",
		"int64",
		":",
		"]",
		"varname",
		"{}",
		"{",
		"}",
		"=",
	},

	idMap: map[string]Type {
		"INVALID": 0,
		"$": 1,
		"(": 2,
		")": 3,
		"by": 4,
		",": 5,
		"varlabel": 6,
		"percentile": 7,
		"sum": 8,
		"mean": 9,
		"median": 10,
		"min": 11,
		"max": 12,
		"stddev": 13,
		"rate": 14,
		"rate_signed": 15,
		"delta": 16,
		"latest": 17,
		"[": 18,
		"int64": 19,
		":": 20,
		"]": 21,
		"varname": 22,
		"{}": 23,
		"{": 24,
		"}": 25,
		"=": 26,
	},
}

