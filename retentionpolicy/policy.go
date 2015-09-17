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

func (policy *RetentionPolicy) Apply(itemVar *variable.Variable, input <-chan *oproto.Value) <-chan *oproto.Value {
	output := make(chan *oproto.Value, 1000)
	matchingPolicies := make([]*oproto.RetentionPolicyItem, 0)
	for _, p := range policy.policy.Policy {
		if doesVariableMatch(itemVar, p.Variable) {
			matchingPolicies = append(matchingPolicies, p)
		}
	}
	if len(matchingPolicies) == 0 {
		log.Printf("No matching policies for variable %s", itemVar)
	}

	go func() {
		for value := range input {
			if len(matchingPolicies) == 0 {
				// No matching policies, drop everything
				continue
			}

			// Find first matching policy
			p := findFirstMatchingPolicy(value, matchingPolicies)
			if p == nil {
				log.Printf("No matching policies for variable, dropping value")
				continue
			}
			log.Printf("Found matching policy %s", p)
			if p.Policy == oproto.RetentionPolicyItem_DROP {
				// Matching policy is DROP, so don't output this value
				continue
			}

			output <- value
		}
		close(output)
	}()
	return output
}

func doesVariableMatch(itemVar *variable.Variable, policyVars []*oproto.StreamVariable) bool {
	for _, v := range policyVars {
		policyVar := variable.NewFromProto(v)
		if itemVar.Match(policyVar) {
			return true
		}
		log.Printf("Stream variable %s doesn't match policy variable %s", itemVar.String(), policyVar.String())
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
		// Look for policies that match the variable age
		if item.MinAge != 0 {
			if valueStartAge < item.MinAge || valueEndAge < item.MinAge {
				continue
			}
		}
		if item.MaxAge != 0 {
			if valueStartAge > item.MaxAge || valueEndAge > item.MaxAge {
				continue
			}
		}
		return item
	}
	return nil
}
