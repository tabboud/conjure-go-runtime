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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/witchcraft-go-error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go-runtime/conjure-go-contract/errors"
)

func TestHandlerFunc_ServeHTTP(t *testing.T) {
	for _, test := range []struct {
		Name    string
		Handler func(w http.ResponseWriter, r *http.Request) error
		Verify  func(t *testing.T, resp *http.Response, err error)
	}{
		{
			Name: "Basic 200",
			Handler: func(w http.ResponseWriter, r *http.Request) error {
				WriteJSONResponse(w, "hello world", 200)
				return nil
			},
			Verify: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				body := readBody(t, resp)

				assert.Equal(t, "\"hello world\"\n", string(body))
				assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
				assert.Equal(t, 200, resp.StatusCode)
			},
		},
		{
			Name: "Basic 404",
			Handler: func(w http.ResponseWriter, r *http.Request) error {
				return werror.Wrap(errors.NewNotFound(errors.SafeParam("resource", "foo")), "wrap1", werror.UnsafeParam("unsafe", "secret"))
			},
			Verify: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				conjureErr, err := errors.ErrorFromResponse(resp)
				require.NoError(t, err)
				assert.Equal(t, errors.NotFound, conjureErr.Code())
				assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
				assert.Equal(t, 404, resp.StatusCode)
				assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"resource": objmatcher.NewEqualsMatcher("foo"),
				}).Matches(conjureErr.Parameters()))
			},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			srv := httptest.NewServer(HandlerFunc(test.Handler))
			defer srv.Close()
			resp, err := http.Get(srv.URL)
			test.Verify(t, resp, err)
		})
	}
}

func readBody(t *testing.T, resp *http.Response) []byte {
	require.NotNil(t, resp)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	return body
}
