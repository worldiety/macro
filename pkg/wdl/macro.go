package wdl

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var MacroRegex = regexp.MustCompile(`#\[.+\]`)

// A MacroName is something like go.TaggedUnion
type MacroName string

// MacroInvocation is parsed from texts in the following formats (see also [MacroRegex]]):
//   - #[<identifier>.<identifier> <json object body>?]
//   - for example #[go.TaggedUnion]
//   - for example #[go.TaggedUnion "tag":"$_type"]
type MacroInvocation struct {
	name              MacroName
	debugParsedParams map[string]any
	jsonParams        string
	pos               Pos
}

func (m *MacroInvocation) Name() MacroName {
	return m.name
}

func (m *MacroInvocation) DebugParsedParams() map[string]any {
	return m.debugParsedParams
}

func (m *MacroInvocation) JsonParams() string {
	return m.jsonParams
}

// Value returns some duck typing string from the json which looks like a value.
func (m *MacroInvocation) Value() string {
	if m.debugParsedParams != nil {
		if v, ok := m.debugParsedParams["value"]; ok {
			return fmt.Sprintf("%v", v)
		}

		if v, ok := m.debugParsedParams["values"]; ok {
			if slice, ok := v.([]string); ok && len(slice) > 0 {
				return slice[0]
			}
		}
	}

	for _, a := range m.debugParsedParams {
		return fmt.Sprintf("%v", a)
	}

	return ""
}

func (m *MacroInvocation) Pos() Pos {
	return m.pos
}

func ParseMacroInvocation(text string, pos Pos) (*MacroInvocation, error) {
	text = strings.TrimSpace(text)
	if !(strings.HasPrefix(text, "#[") && strings.HasSuffix(text, "]")) {
		return nil, fmt.Errorf("not a macro invocation: %s", text)
	}

	m := &MacroInvocation{pos: pos}
	text = text[2:]
	text = text[:len(text)-1]
	text = strings.TrimSpace(text)
	fnAndParams := strings.SplitN(text, " ", 2)
	m.name = MacroName(fnAndParams[0])

	if len(fnAndParams) == 2 {
		raw := fnAndParams[1]
		parsers := []func(string) (string, error){
			jsonObj,
			jsonSlice,
			jsonStr,
		}

		var lastErr error
		for _, parser := range parsers {
			jStr, err := parser(raw)
			if err == nil {
				if err := json.Unmarshal([]byte(jStr), &m.debugParsedParams); err != nil {
					return nil, fmt.Errorf("unexpected double decoding error %s: %v", m.name, err)
				}
				m.jsonParams = jStr
				lastErr = nil
				break
			} else {
				lastErr = err
			}

		}

		if lastErr != nil {
			return nil, fmt.Errorf("cannot interpret macro params into any json format %s: %v", m.name, lastErr)
		}

	} else {
		m.jsonParams = "{}"
	}

	return m, nil
}

func jsonObj(text string) (string, error) {
	jStr := "{" + text + "}"
	var tmp any
	err := json.Unmarshal([]byte(jStr), &tmp)
	return jStr, err
}

func jsonSlice(text string) (string, error) {
	jStr := "{\"values\":" + text + "}"
	var tmp any
	err := json.Unmarshal([]byte(jStr), &tmp)
	return jStr, err
}

func jsonStr(text string) (string, error) {
	jStr := "{\"value\":" + strconv.Quote(text) + "}"
	var tmp any
	err := json.Unmarshal([]byte(jStr), &tmp)
	return jStr, err
}

func (m *MacroInvocation) UnmarshalParams(dst any) error {
	return json.Unmarshal([]byte(m.jsonParams), dst)
}
