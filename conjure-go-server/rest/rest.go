// Copyright (c) 2018 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rest

import (
	"net/http"

	"github.com/palantir/conjure-go-runtime/conjure-go-contract/codecs"
)

// WriteJSONResponse marshals the provided object to JSON using a JSON encoder with SetEscapeHTML(false) and writes the
// resulting JSON as a JSON response to the provided http.ResponseWriter with the provided status code. If marshaling
// the provided object as JSON results in an error, writes a 500 response with the text content of the error.
func WriteJSONResponse(w http.ResponseWriter, obj interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := codecs.JSON.Encode(w, obj); err != nil {
		// if JSON encode failed, send error response. If JSON encode succeeded but write failed, then this
		// should be a no-op since the socket failed anyway.
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
