
package parser

import "github.com/dparrish/openinstrument/query/ast"

type (
	//TODO: change type and variable names to be consistent with other tables
	ProdTab      [numProductions]ProdTabEntry
	ProdTabEntry struct {
		String     string
		Id         string
		NTType     int
		Index int
		NumSymbols int
		ReduceFunc func([]Attrib) (Attrib, error)
	}
	Attrib interface {
	}
)

var productionsTable = ProdTab {
	ProdTabEntry{
		String: `S' : Query	<<  >>`,
		Id: "S'",
		NTType: 0,
		Index: 0,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Query : Variable	<< ast.NewQuery(X[0]) >>`,
		Id: "Query",
		NTType: 1,
		Index: 1,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewQuery(X[0])
		},
	},
	ProdTabEntry{
		String: `Query : Mutation	<< ast.NewQuery(X[0]) >>`,
		Id: "Query",
		NTType: 1,
		Index: 2,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewQuery(X[0])
		},
	},
	ProdTabEntry{
		String: `Query : Aggregation	<< ast.NewQuery(X[0]) >>`,
		Id: "Query",
		NTType: 1,
		Index: 3,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewQuery(X[0])
		},
	},
	ProdTabEntry{
		String: `Aggregation : AggregationType "(" VariableList ")"	<< ast.NewAggregation(X[0], nil, X[2]) >>`,
		Id: "Aggregation",
		NTType: 2,
		Index: 4,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewAggregation(X[0], nil, X[2])
		},
	},
	ProdTabEntry{
		String: `Aggregation : AggregationType "(" Mutation ")"	<< ast.NewAggregation(X[0], nil, X[2]) >>`,
		Id: "Aggregation",
		NTType: 2,
		Index: 5,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewAggregation(X[0], nil, X[2])
		},
	},
	ProdTabEntry{
		String: `Aggregation : AggregationType "by" "(" LabelList ")" "(" VariableList ")"	<< ast.NewAggregation(X[0], X[3], X[6]) >>`,
		Id: "Aggregation",
		NTType: 2,
		Index: 6,
		NumSymbols: 8,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewAggregation(X[0], X[3], X[6])
		},
	},
	ProdTabEntry{
		String: `Aggregation : AggregationType "by" "(" LabelList ")" "(" Mutation ")"	<< ast.NewAggregation(X[0], X[3], X[6]) >>`,
		Id: "Aggregation",
		NTType: 2,
		Index: 7,
		NumSymbols: 8,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewAggregation(X[0], X[3], X[6])
		},
	},
	ProdTabEntry{
		String: `Mutation : MutationType "(" Duration "," VariableList ")"	<< ast.NewMutation(X[0], X[2], X[4]) >>`,
		Id: "Mutation",
		NTType: 3,
		Index: 8,
		NumSymbols: 6,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewMutation(X[0], X[2], X[4])
		},
	},
	ProdTabEntry{
		String: `LabelList : varlabel	<< ast.NewLabelList(X[0]) >>`,
		Id: "LabelList",
		NTType: 4,
		Index: 9,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLabelList(X[0])
		},
	},
	ProdTabEntry{
		String: `LabelList : LabelList "," varlabel	<< ast.AppendStringToList(X[0], X[2]) >>`,
		Id: "LabelList",
		NTType: 4,
		Index: 10,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.AppendStringToList(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Percentile : "percentile" "(" Int64 ")"	<< ast.NewPercentile(X[2]) >>`,
		Id: "Percentile",
		NTType: 5,
		Index: 11,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewPercentile(X[2])
		},
	},
	ProdTabEntry{
		String: `AggregationType : "sum"	<<  >>`,
		Id: "AggregationType",
		NTType: 6,
		Index: 12,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `AggregationType : "mean"	<<  >>`,
		Id: "AggregationType",
		NTType: 6,
		Index: 13,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `AggregationType : "median"	<<  >>`,
		Id: "AggregationType",
		NTType: 6,
		Index: 14,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `AggregationType : "min"	<<  >>`,
		Id: "AggregationType",
		NTType: 6,
		Index: 15,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `AggregationType : "max"	<<  >>`,
		Id: "AggregationType",
		NTType: 6,
		Index: 16,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `AggregationType : "stddev"	<<  >>`,
		Id: "AggregationType",
		NTType: 6,
		Index: 17,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `AggregationType : Percentile	<<  >>`,
		Id: "AggregationType",
		NTType: 6,
		Index: 18,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `MutationType : "rate"	<<  >>`,
		Id: "MutationType",
		NTType: 7,
		Index: 19,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `MutationType : "rate_signed"	<<  >>`,
		Id: "MutationType",
		NTType: 7,
		Index: 20,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `MutationType : "delta"	<<  >>`,
		Id: "MutationType",
		NTType: 7,
		Index: 21,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `MutationType : "latest"	<<  >>`,
		Id: "MutationType",
		NTType: 7,
		Index: 22,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Variable : varname	<< ast.NewVariable(X[0], nil) >>`,
		Id: "Variable",
		NTType: 8,
		Index: 23,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewVariable(X[0], nil)
		},
	},
	ProdTabEntry{
		String: `Variable : varname "{}"	<< ast.NewVariable(X[0], nil) >>`,
		Id: "Variable",
		NTType: 8,
		Index: 24,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewVariable(X[0], nil)
		},
	},
	ProdTabEntry{
		String: `Variable : varname "{" KeyValues "}"	<< ast.NewVariable(X[0], X[2]) >>`,
		Id: "Variable",
		NTType: 8,
		Index: 25,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewVariable(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `KeyValues : KeyValue	<< ast.NewKeyValueList(X[0]) >>`,
		Id: "KeyValues",
		NTType: 9,
		Index: 26,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewKeyValueList(X[0])
		},
	},
	ProdTabEntry{
		String: `KeyValues : KeyValues "," KeyValue	<< ast.AppendKeyValueList(X[0], X[2]) >>`,
		Id: "KeyValues",
		NTType: 9,
		Index: 27,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.AppendKeyValueList(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `KeyValue : varlabel "=" varlabel	<< ast.NewKeyValue(X[0], X[2]) >>`,
		Id: "KeyValue",
		NTType: 10,
		Index: 28,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewKeyValue(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `VariableList : Variable	<< ast.NewVariableList(X[0]) >>`,
		Id: "VariableList",
		NTType: 11,
		Index: 29,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewVariableList(X[0])
		},
	},
	ProdTabEntry{
		String: `VariableList : VariableList "," Variable	<< ast.AddVariableToVariableList(X[0], X[2]) >>`,
		Id: "VariableList",
		NTType: 11,
		Index: 30,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.AddVariableToVariableList(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Duration : varlabel	<< ast.NewDuration(X[0]) >>`,
		Id: "Duration",
		NTType: 12,
		Index: 31,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewDuration(X[0])
		},
	},
	ProdTabEntry{
		String: `Int64 : int64	<< ast.NewInt64(X[0]) >>`,
		Id: "Int64",
		NTType: 13,
		Index: 32,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewInt64(X[0])
		},
	},
	
}
