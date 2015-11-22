package query

import (
	"log"

	"github.com/dparrish/openinstrument"
	"github.com/dparrish/openinstrument/datastore"
	"github.com/dparrish/openinstrument/mutations"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/query/lexer"
	"github.com/dparrish/openinstrument/query/parser"
	"github.com/dparrish/openinstrument/variable"
)

type Query struct {
	query *oproto.Query
}

func Parse(q string) (*Query, error) {
	p := parser.NewParser()
	l := lexer.NewLexer([]byte(q))
	v, err := p.Parse(l)
	if err != nil {
		return nil, err
	}
	return &Query{v.(*oproto.Query)}, nil
}

func NewFromProto(q *oproto.Query) *Query {
	//log.Println(openinstrument.ProtoText(q))
	return &Query{q}
}

func (q *Query) AsProto() *oproto.Query {
	return q.query
}

func (q *Query) Run(store datastore.ReadableStore) (chan *oproto.ValueStream, error) {
	return RunQuery(q.query, store)
}

func RunQuery(query *oproto.Query, store datastore.ReadableStore) (chan *oproto.ValueStream, error) {
	log.Printf("Running query %v", query)
	output := make(chan *oproto.ValueStream, 100)
	go func() {
		defer close(output)
		for _, v := range query.Variable {
			log.Printf("Returning variable %s", variable.ProtoToString(v))
			for stream := range store.Reader(variable.NewFromProto(v)) {
				outputStream := &oproto.ValueStream{Variable: stream.Variable}
				stv := variable.NewFromProto(stream.Variable)
				for _, v := range stream.Value {
					if stv.TimestampInsideRange(v.Timestamp) {
						outputStream.Value = append(outputStream.Value, v)
					}
				}
				output <- outputStream
			}
		}
		for _, child := range query.Aggregation {
			log.Printf("Running child aggregation")
			input := []*oproto.ValueStream{}
			for _, q := range child.Query {
				o, err := RunQuery(q, store)
				if err != nil {
					return
				}
				for stream := range o {
					input = append(input, stream)
				}
			}
			log.Printf("Child aggregation returned output")
			for _, stream := range input {
				log.Println(openinstrument.ProtoText(stream))
				output <- stream
			}
		}
		for _, child := range query.Mutation {
			log.Printf("Running child mutation")
			input, err := RunQuery(child.Query, store)
			if err != nil {
				log.Printf("Error in child mutation: %s", err)
				return
			}
			for stream := range input {
				var outStream *oproto.ValueStream
				switch child.Type {
				case oproto.StreamMutation_MEAN:
					outStream = mutations.Mean(stream)
				case oproto.StreamMutation_INTERPOLATE:
					outStream = mutations.Interpolate(uint64(child.Param), stream)
				case oproto.StreamMutation_MIN:
					outStream = mutations.Min(uint64(child.Param), stream)
				case oproto.StreamMutation_MAX:
					outStream = mutations.Max(uint64(child.Param), stream)
				case oproto.StreamMutation_FIRST:
					outStream = mutations.First(uint64(child.Param), stream)
				case oproto.StreamMutation_LAST:
					outStream = mutations.Last(uint64(child.Param), stream)
				case oproto.StreamMutation_RATE:
					outStream = mutations.Rate(stream)
				case oproto.StreamMutation_ROOT:
					outStream = mutations.Root(child.Param, stream)
				case oproto.StreamMutation_POWER:
					outStream = mutations.Power(child.Param, stream)
				case oproto.StreamMutation_ADD:
					outStream = mutations.Add(child.Param, stream)
				case oproto.StreamMutation_MULTIPLY:
					outStream = mutations.Multiply(child.Param, stream)
				case oproto.StreamMutation_RATE_SIGNED:
					outStream = mutations.SignedRate(stream)
				case oproto.StreamMutation_MOVING_AVERAGE:
					outStream = mutations.MovingAverage(uint64(child.Param), stream)
				}
				if outStream == nil {
					log.Printf("No stream returned from mutation")
					continue
				}
				outStream.Variable = stream.Variable
				output <- outStream
			}
		}
	}()
	return output, nil
}
