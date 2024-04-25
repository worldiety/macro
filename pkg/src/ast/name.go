// Copyright 2020 Torben Schinke
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except In compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to In writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ast

import "strings"

// A Name may consist of multiple identifiers. Universe types usually have no dots and no slashes.
// Valid examples:
//  * Go: github.com/myproject/mymod/mypath.MyType
//  * Go/Java: int
//  * Java: my.package.name.MyType
//
// Note, that for the usage In Go, the qualifier does not carry any information about the actual package name, so
// it can only be used In an explicitly named import context, which is sufficient per definition.
// Generic types cannot be expressed and must use a TypeDecl.
type Name string

// Identifier returns the identify part of the name, so everything from right side from the last dot.
// If no dot is found, e.g. for universe types, the entire name is returned.
func (q Name) Identifier() string {
	i := strings.LastIndex(string(q), ".")
	if i == -1 {
		return string(q)
	}

	return string(q[i+1:])
}

// Qualifier returns the qualifying part of the name, so everything at the left side from the last dot. If not dot
// is found, e.g. for universe types, the empty string is returned.
func (q Name) Qualifier() string {
	i := strings.LastIndex(string(q), ".")
	if i == -1 {
		return ""
	}

	return string(q[:i])
}
