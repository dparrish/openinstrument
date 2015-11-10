// These functions are used by the parser to generate the object tree returned from query.Parse

package ast

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/query/token"
)

// Variable parsing
type VariableLabel struct {
	key   string
	value string
}

type VariableLabelList struct {
	l []*VariableLabel
}

func NewVariable(name, labels interface{}) (*oproto.StreamVariable, error) {
	v := &oproto.StreamVariable{
		Name:  string(name.(*token.Token).Lit),
		Label: map[string]string{},
	}

	if labels != nil {
		for _, label := range labels.(*VariableLabelList).l {
			v.Label[label.key] = label.value
		}
	}

	return v, nil
}

func NewKeyValue(key, value interface{}) (*VariableLabel, error) {
	return &VariableLabel{string(key.(*token.Token).Lit), string(value.(*token.Token).Lit)}, nil
}

func NewKeyValueList(vl interface{}) (*VariableLabelList, error) {
	return &VariableLabelList{
		[]*VariableLabel{
			vl.(*VariableLabel),
		},
	}, nil
}

func AppendKeyValueList(l, vl interface{}) (*VariableLabelList, error) {
	l.(*VariableLabelList).l = append(l.(*VariableLabelList).l, vl.(*VariableLabel))
	return l.(*VariableLabelList), nil
}

// Variable List
type VariableList []*oproto.StreamVariable

func NewVariableList(v interface{}) (VariableList, error) {
	return VariableList{v.(*oproto.StreamVariable)}, nil
}

func AddVariableToVariableList(list, v interface{}) (VariableList, error) {
	return append(list.(VariableList), v.(*oproto.StreamVariable)), nil
}

// Query parsing
func NewQuery(input interface{}) (*oproto.Query, error) {
	query := &oproto.Query{}
	switch v := input.(type) {
	case *oproto.StreamMutation:
		query.Mutation = append(query.Mutation, v)
	case *oproto.StreamAggregation:
		query.Aggregation = append(query.Aggregation, v)
	case *oproto.StreamVariable:
		query.Variable = append(query.Variable, v)
	}
	return query, nil
}

func AddVariableListToQuery(a, v interface{}) (*oproto.Query, error) {
	for _, variable := range v.(VariableList) {
		a.(*oproto.Query).Variable = append(a.(*oproto.Query).Variable, variable)
	}
	return a.(*oproto.Query), nil
}

func NewLabelList(s interface{}) ([]string, error) {
	str := string(s.(*token.Token).Lit)
	return []string{str}, nil
}

func AppendStringToList(l, s interface{}) ([]string, error) {
	str := string(s.(*token.Token).Lit)
	return append(l.([]string), str), nil
}

func NewMutation(t, d, v interface{}) (*oproto.StreamMutation, error) {
	m := &oproto.StreamMutation{}
	duration := d.(time.Duration)
	m.SampleFrequency = uint32(duration.Nanoseconds() / 1000000)

	f, ok := oproto.StreamMutation_SampleType_value[strings.ToUpper(string(t.(*token.Token).Lit))]
	if !ok {
		return nil, fmt.Errorf("Invalid mutation type %s", string(t.(*token.Token).Lit))
	}
	m.SampleType = oproto.StreamMutation_SampleType(f)

	for _, variable := range v.(VariableList) {
		m.Variable = append(m.Variable, variable)
	}
	return m, nil
}

func NewDuration(digit interface{}) (time.Duration, error) {
	duration, err := time.ParseDuration(string(digit.(*token.Token).Lit))
	if err != nil {
		return 0, err
	}
	return duration, nil
}

func NewInt64(digits interface{}) (int64, error) {
	i, err := strconv.Atoi(string(digits.(*token.Token).Lit))
	if err != nil {
		return 0, err
	}
	return int64(i), nil
}

type Percentile struct {
	percentile int32
}

func NewPercentile(i interface{}) (Percentile, error) {
	return Percentile{int32(i.(int64))}, nil
}

func NewAggregation(aggType, by, vars interface{}) (*oproto.StreamAggregation, error) {
	agg := &oproto.StreamAggregation{}
	switch v := aggType.(type) {
	case *token.Token:
		f, ok := oproto.StreamAggregation_AggregateType_value[strings.ToUpper(string(v.Lit))]
		if !ok {
			return nil, fmt.Errorf("Invalid aggregation type %s", v)
		}
		agg.Type = oproto.StreamAggregation_AggregateType(f)
	case Percentile:
		agg.Type = oproto.StreamAggregation_PERCENTILE
		agg.Percentile = uint32(v.percentile)
	case string:
		f, ok := oproto.StreamAggregation_AggregateType_value[strings.ToUpper(v)]
		if !ok {
			return nil, fmt.Errorf("Invalid aggregation type %s", v)
		}
		agg.Type = oproto.StreamAggregation_AggregateType(f)
	}

	if by != nil {
		for _, x := range by.([]string) {
			agg.Label = append(agg.Label, x)
		}
	}

	switch v := vars.(type) {
	case *oproto.StreamMutation:
		agg.Mutation = append(agg.Mutation, v)
	case VariableList:
		for _, variable := range v {
			agg.Variable = append(agg.Variable, variable)
		}
	}

	return agg, nil
}

func AddRangeToVariable(v, from, to interface{}) (*oproto.StreamVariable, error) {
	variable := v.(*oproto.StreamVariable)
	if from != nil {
		value, err := strconv.Atoi(string(from.(*token.Token).Lit))
		if err != nil {
			return nil, err
		}
		variable.MinTimestamp = int64(value)
	}
	if to != nil {
		value, err := strconv.Atoi(string(to.(*token.Token).Lit))
		if err != nil {
			return nil, err
		}
		variable.MaxTimestamp = int64(value)
	}
	return variable, nil
}
