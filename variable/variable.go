package variable

import (
	"encoding/csv"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
)

// Variable represents a variable name with all labels.
type Variable struct {
	Variable     string
	Labels       map[string]string
	MinTimestamp int64
	MaxTimestamp int64
}

func valueNeedsQuotes(str string) bool {
	return strings.ContainsAny(str, ",[]{}")
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// String returns the complete variable as a string, including all labels correctly quoted.
func (v *Variable) String() string {
	str := v.Variable
	if len(v.Labels) > 0 {
		str += "{"
		first := true
		labels := sortedKeys(v.Labels)
		for _, key := range labels {
			value := v.Labels[key]
			if !first {
				str += ","
			}
			first = false
			str += key
			str += "="
			needsQuotes := valueNeedsQuotes(value)
			if needsQuotes {
				str += "\""
			}
			str += value
			if needsQuotes {
				str += "\""
			}
		}
		str += "}"
	}
	if v.MinTimestamp != 0 || v.MaxTimestamp != 0 {
		str += fmt.Sprintf("[%d", v.MinTimestamp)
		if v.MaxTimestamp != 0 {
			str += fmt.Sprintf(":%d", v.MaxTimestamp)
		}
		str += "]"
	}
	return str
}

// AsProto returns the Variable encoded as a new StreamVariable protobuf.
func (v *Variable) AsProto() (p *oproto.StreamVariable) {
	p = new(oproto.StreamVariable)
	v.ToProto(p)
	return
}

// AsProto encodes the Variable into an existing protobuf.
func (v *Variable) ToProto(p *oproto.StreamVariable) {
	p.Reset()
	p.Name = v.Variable
	p.Label = v.Labels
	p.MinTimestamp = v.MinTimestamp
	p.MaxTimestamp = v.MaxTimestamp
}

// ParseFromString extracts the variable name and all labels from a string.
func (v *Variable) ParseFromString(textvar string) error {
	re, err := regexp.Compile("^(.+?)({(.*)})?(?:\\[(-?\\d+):?(-?\\d+)?\\])?$")
	if err != nil {
		log.Panic("Invalid regexp")
	}
	matches := re.FindStringSubmatch(textvar)
	v.Variable = matches[1]
	if matches[3] != "" {
		reader := csv.NewReader(strings.NewReader(matches[3]))
		record, err := reader.Read()
		if err != nil {
			return fmt.Errorf("Error parsing labels from %s: %s", textvar, err)
		}
		v.Labels = make(map[string]string, 0)
		for _, item := range record {
			substrings := strings.SplitN(item, "=", 2)
			if len(substrings) != 2 {
				return fmt.Errorf("Error splitting labels from %s", textvar)
			}
			v.Labels[strings.TrimSpace(substrings[0])] = strings.TrimSpace(substrings[1])
		}
	}
	if matches[4] != "" {
		ts, err := strconv.Atoi(matches[4])
		if err != nil {
			return fmt.Errorf("Invalid Min Timestamp")
		}
		v.MinTimestamp = int64(ts)
	}
	if matches[5] != "" {
		ts, err := strconv.Atoi(matches[5])
		if err != nil {
			return fmt.Errorf("Invalid Max Timestamp")
		}
		v.MaxTimestamp = int64(ts)
	}
	if v.MaxTimestamp != 0 && v.MaxTimestamp < v.MinTimestamp {
		return fmt.Errorf("Max timstamp is before min timestamp")
	}
	return nil
}

// ParseFromProto extracts the details from a protobuf.
func (v *Variable) ParseFromProto(p *oproto.StreamVariable) error {
	v.Variable = p.Name
	v.MinTimestamp = p.MinTimestamp
	v.MaxTimestamp = p.MaxTimestamp
	// Copy labels
	v.Labels = p.Label
	return nil
}

// Match checks that the supplied variable matches this one, using the matching rules defined above.
func (v *Variable) Match(match *Variable) bool {
	if strings.HasSuffix(match.Variable, "*") {
		if match.Variable != "*" {
			m := match.Variable[0 : len(match.Variable)-2]
			if !strings.HasPrefix(v.Variable, m) {
				// Invalid prefix match
				return false
			}
		}
	} else if match.Variable != v.Variable {
		// Variable name does not match
		return false
	}
	for key, value := range match.Labels {
		if value == "" {
			// Key must not exist
			_, ok := v.Labels[key]
			if ok {
				return false
			}
		} else if value == "*" {
			// Key must exist
			val, ok := v.Labels[key]
			if !ok {
				return false
			}
			if val == "" {
				return false
			}
		} else {
			val, ok := v.Labels[key]
			if !ok {
				return false
			}
			if val != value {
				return false
			}
		}
	}
	return true
}

func (v *Variable) TimestampInsideRange(timestamp uint64) bool {
	minTimestamp := v.MinTimestamp
	if minTimestamp < 0 {
		minTimestamp = int64(openinstrument.NowMs()) + minTimestamp
	}
	maxTimestamp := v.MaxTimestamp
	if maxTimestamp < 0 {
		maxTimestamp = int64(openinstrument.NowMs()) + maxTimestamp
	}
	if maxTimestamp != 0 && int64(timestamp) > maxTimestamp {
		// Too new
		return false
	}
	if int64(timestamp) < minTimestamp {
		// Too old
		return false
	}
	return true
}

// NewFromString creates a new Variable from the supplied string.
func NewFromString(textvar string) *Variable {
	v := &Variable{}
	if err := v.ParseFromString(textvar); err != nil {
		log.Println(err)
		return nil
	}
	return v
}

// NewFromProto creates a new Variable from the supplied protobuf.
func NewFromProto(p *oproto.StreamVariable) *Variable {
	v := &Variable{}
	if p == nil {
		return v
	}
	if err := v.ParseFromProto(p); err != nil {
		log.Println(err)
		return nil
	}
	return v
}

func ProtoToString(p *oproto.StreamVariable) string {
	return NewFromProto(p).String()
}
