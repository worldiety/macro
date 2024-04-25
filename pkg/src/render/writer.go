// Copyright 2020 Torben Schinke
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package render

import (
	"bytes"
	"fmt"
	"strings"
)

// Writer introduces a Printf contract.
type Writer interface {
	Printf(str string, args ...interface{})
}

// A BufferedWriter is just a strings builder with a Printf.
type BufferedWriter bytes.Buffer

func (b *BufferedWriter) Printf(format string, args ...interface{}) {
	if len(args) == 0 {
		(*bytes.Buffer)(b).WriteString(format)
		return
	}

	(*bytes.Buffer)(b).WriteString(fmt.Sprintf(format, args...))
}

func (b *BufferedWriter) Print(a ...interface{}) {
	(*bytes.Buffer)(b).WriteString(fmt.Sprint(a...))
}

// String returns the builders text.
func (b *BufferedWriter) String() string {
	return (*bytes.Buffer)(b).String()
}

// Bytes returns the backing slice.
func (b *BufferedWriter) Bytes() []byte {
	return (*bytes.Buffer)(b).Bytes()
}

// WithLineNumbers puts a 1 based line number to the left and returns the text.
func WithLineNumbers(text string) string {
	sb := &strings.Builder{}
	for i, line := range strings.Split(text, "\n") {
		sb.WriteString(fmt.Sprintf("%4d: %s\n", i+1, line))
	}
	return sb.String()
}
