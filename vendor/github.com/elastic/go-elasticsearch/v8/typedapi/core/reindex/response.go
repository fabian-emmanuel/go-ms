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

package reindex

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// Response holds the response body struct for the package reindex
//
// https://github.com/elastic/elasticsearch-specification/blob/2f823ff6fcaa7f3f0f9b990dc90512d8901e5d64/specification/_global/reindex/ReindexResponse.ts#L26-L45
type Response struct {
	Batches              *int64                           `json:"batches,omitempty"`
	Created              *int64                           `json:"created,omitempty"`
	Deleted              *int64                           `json:"deleted,omitempty"`
	Failures             []types.BulkIndexByScrollFailure `json:"failures,omitempty"`
	Noops                *int64                           `json:"noops,omitempty"`
	RequestsPerSecond    *float32                         `json:"requests_per_second,omitempty"`
	Retries              *types.Retries                   `json:"retries,omitempty"`
	SliceId              *int                             `json:"slice_id,omitempty"`
	Task                 types.TaskId                     `json:"task,omitempty"`
	ThrottledMillis      *int64                           `json:"throttled_millis,omitempty"`
	ThrottledUntilMillis *int64                           `json:"throttled_until_millis,omitempty"`
	TimedOut             *bool                            `json:"timed_out,omitempty"`
	Took                 *int64                           `json:"took,omitempty"`
	Total                *int64                           `json:"total,omitempty"`
	Updated              *int64                           `json:"updated,omitempty"`
	VersionConflicts     *int64                           `json:"version_conflicts,omitempty"`
}

// NewResponse returns a Response
func NewResponse() *Response {
	r := &Response{}
	return r
}
