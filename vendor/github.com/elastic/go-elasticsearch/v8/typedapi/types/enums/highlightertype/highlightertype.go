// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated from the elasticsearch-specification DO NOT EDIT.
// https://github.com/elastic/elasticsearch-specification/tree/2f823ff6fcaa7f3f0f9b990dc90512d8901e5d64

// Package highlightertype
package highlightertype

import "strings"

// https://github.com/elastic/elasticsearch-specification/blob/2f823ff6fcaa7f3f0f9b990dc90512d8901e5d64/specification/_global/search/_types/highlighting.ts#L175-L190
type HighlighterType struct {
	Name string
}

var (
	Plain = HighlighterType{"plain"}

	Fastvector = HighlighterType{"fvh"}

	Unified = HighlighterType{"unified"}
)

func (h HighlighterType) MarshalText() (text []byte, err error) {
	return []byte(h.String()), nil
}

func (h *HighlighterType) UnmarshalText(text []byte) error {
	switch strings.ReplaceAll(strings.ToLower(string(text)), "\"", "") {

	case "plain":
		*h = Plain
	case "fvh":
		*h = Fastvector
	case "unified":
		*h = Unified
	default:
		*h = HighlighterType{string(text)}
	}

	return nil
}

func (h HighlighterType) String() string {
	return h.Name
}
