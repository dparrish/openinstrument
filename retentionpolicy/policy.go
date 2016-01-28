package retentionpolicy

import (
	"log"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
)

type RetentionPolicy struct {
	policy *oproto.RetentionPolicy
}

func New(inputPolicy *oproto.RetentionPolicy) *RetentionPolicy {
	policy := &RetentionPolicy{
		policy: inputPolicy,
	}
	return policy
}

func (policy *RetentionPolicy) Apply(input *oproto.ValueStream) *oproto.ValueStream {
	itemVar := variable.NewFromProto(input.Variable)
	matchingPolicies := make([]*oproto.RetentionPolicyItem, 0)
	for _, policy := range policy.policy.Policy {
		if doesVariableMatch(itemVar, policy.Variable) {
			matchingPolicies = append(matchingPolicies, policy)
		}
	}
	if len(matchingPolicies) == 0 {
		log.Printf("No matching policies for variable %s", itemVar)
	}

	output := &oproto.ValueStream{
		Variable: input.Variable,
		Value:    []*oproto.Value{},
	}

	for _, value := range input.Value {
		// Find first matching policy
		if policy := findFirstMatchingPolicy(value, matchingPolicies); policy != nil {
			if policy.Policy == oproto.RetentionPolicyItem_DROP {
				// Matching policy is DROP, so don't output this value
				continue
			}

			if len(policy.Mutation) == 0 {
				output.Value = append(output.Value, value)
				continue
			}

			/*
				// TODO(dparrish): Fix this!
				for _, mutation := range policy.Mutation {
					//log.Printf("Applying mutation %s", mutation)
					var outStream &oproto.ValueStream
					switch mutation.Type {
					case oproto.StreamMutation_MEAN:
						outStream = mutations.Mean(stream)
					case oproto.StreamMutation_INTERPOLATE:
						outStream = mutations.Interpolate(uint64(mutation.Param), stream)
					case oproto.StreamMutation_MIN:
						outStream = mutations.Min(uint64(mutation.Param), stream)
					case oproto.StreamMutation_MAX:
						outStream = mutations.Max(uint64(mutation.Param), stream)
					case oproto.StreamMutation_FIRST:
						outStream = mutations.First(uint64(mutation.Param), stream)
					case oproto.StreamMutation_LAST:
						outStream = mutations.Last(uint64(mutation.Param), stream)
					case oproto.StreamMutation_RATE:
						outStream = mutations.Rate(stream)
					case oproto.StreamMutation_ROOT:
						outStream = mutations.Root(mutation.Param, stream)
					case oproto.StreamMutation_POWER:
						outStream = mutations.Power(mutation.Param, stream)
					case oproto.StreamMutation_ADD:
						outStream = mutations.Add(mutation.Param, stream)
					case oproto.StreamMutation_MULTIPLY:
						outStream = mutations.Multiply(mutation.Param, stream)
					case oproto.StreamMutation_RATE_SIGNED:
						outStream = mutations.SignedRate(stream)
					case oproto.StreamMutation_MOVING_AVERAGE:
						outStream = mutations.MovingAverage(uint64(mutation.Param), stream)
					}
					if outStream == nil {
						log.Printf("No stream returned from mutation")
						continue
					}
					outStream.Variable = stream.Variable
				}
			*/
		}
	}
	return output
}

func doesVariableMatch(itemVar *variable.Variable, policyVars []*oproto.StreamVariable) bool {
	if len(policyVars) == 0 {
		//log.Printf("Stream variable %s matches default policy", itemVar.String())
		return true
	}
	for _, v := range policyVars {
		policyVar := variable.NewFromProto(v)
		if itemVar.Match(policyVar) {
			//log.Printf("Stream variable %s matches policy variable %s", itemVar.String(), policyVar.String())
			return true
		}
		//log.Printf("Stream variable %s doesn't match policy variable %s", itemVar.String(), policyVar.String())
	}
	return false
}

func findFirstMatchingPolicy(value *oproto.Value, policies []*oproto.RetentionPolicyItem) *oproto.RetentionPolicyItem {
	now := openinstrument.NowMs()
	valueStartAge := now - value.Timestamp
	if value.EndTimestamp == 0 && value.Timestamp != 0 {
		value.EndTimestamp = value.Timestamp
	}
	valueEndAge := now - value.EndTimestamp
	for _, item := range policies {
		if len(item.Variable) == 0 {
			// No variables supplied, this matches everything
			return item
		}
		for _, i := range item.Variable {
			// Look for policies that match the variable age
			v := variable.NewFromProto(i)
			if v.TimestampInsideRange(valueStartAge) || v.TimestampInsideRange(valueEndAge) {
				return item
			}
		}
	}
	return nil
}
