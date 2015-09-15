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
		if doesVariableMatch(itemVar, p.GetVariable()) {
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
			if p.GetPolicy() == oproto.RetentionPolicyItem_DROP {
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
	valueStartAge := now - value.GetTimestamp()
	if value.EndTimestamp == nil {
		value.EndTimestamp = value.Timestamp
	}
	valueEndAge := now - value.GetEndTimestamp()
	for _, item := range policies {
		// Look for policies that match the variable age
		if item.MinAge != nil {
			if valueStartAge < item.GetMinAge() || valueEndAge < item.GetMinAge() {
				continue
			}
		}
		if item.MaxAge != nil {
			if valueStartAge > item.GetMaxAge() || valueEndAge > item.GetMaxAge() {
				continue
			}
		}
		return item
	}
	return nil
}
