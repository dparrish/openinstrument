package query

import (
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/query/lexer"
	"github.com/dparrish/openinstrument/query/parser"
)

func Parse(q string) (*oproto.Query, error) {
	p := parser.NewParser()
	l := lexer.NewLexer([]byte(q))
	v, err := p.Parse(l)
	if err == nil {
		return v.(*oproto.Query), err
	}
	return nil, err
}
