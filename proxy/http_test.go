/*
Copyright 2020 Splunk Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

const (
	validAuthToken string = "authtoken"
)

func TestPingForbidden(t *testing.T) {
	var tests = []struct {
		useHeader       bool
		authHeaderValue string
		expectedStatus  int
	}{
		{false, "", http.StatusForbidden},
		{true, "", http.StatusForbidden},
		{true, "invalid", http.StatusForbidden},
		{true, validAuthToken, http.StatusOK},
	}
	gin.DefaultWriter = ioutil.Discard
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	mockCtrl := gomock.NewController(t)
	ic := NewMockistioer(mockCtrl)
	RegisterRoutes(r, ic, validAuthToken)
	ts := httptest.NewServer(r)
	defer ts.Close()

	client := &http.Client{}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/ping", ts.URL), nil)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if test.useHeader {
				req.Header.Add("Authorization", fmt.Sprintf("token %s", test.authHeaderValue))
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("unexpected 	error %v", err)
			}
			if test.expectedStatus != resp.StatusCode {
				t.Errorf("expected http status to be %d but found %d", test.expectedStatus, resp.StatusCode)
			}
		})
	}
}
