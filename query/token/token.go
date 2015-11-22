
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
		"sum",
		"mean",
		"median",
		"min",
		"max",
		"stddev",
		"rate",
		"rate_signed",
		"latest",
		"regularize",
		"interpolate",
		"percentile",
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
		"sum": 7,
		"mean": 8,
		"median": 9,
		"min": 10,
		"max": 11,
		"stddev": 12,
		"rate": 13,
		"rate_signed": 14,
		"latest": 15,
		"regularize": 16,
		"interpolate": 17,
		"percentile": 18,
		"[": 19,
		"int64": 20,
		":": 21,
		"]": 22,
		"varname": 23,
		"{}": 24,
		"{": 25,
		"}": 26,
		"=": 27,
	},
}

