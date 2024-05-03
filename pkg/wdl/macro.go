package wdl

import (
	"encoding/json"
	"fmt"
	"regexp"
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
		m.jsonParams = "{" + fnAndParams[1] + "}"
		buf := []byte(m.jsonParams)
		if err := json.Unmarshal(buf, &m.debugParsedParams); err != nil {
			return nil, fmt.Errorf("error parsing macro invocation params %s: %v", m.name, err)
		}
	} else {
		m.jsonParams = "{}"
	}

	return m, nil
}

func (m *MacroInvocation) UnmarshalParams(dst any) error {
	return json.Unmarshal([]byte(m.jsonParams), dst)
}
