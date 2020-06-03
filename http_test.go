package main

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

func TestValidateRequired(t *testing.T) {
	err := validateRequired("a", "val")
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	err = validateRequired("a", "")
	if err == nil {
		t.Fatal("Expected error but found nil")
	}
	expected := "missing required param a"
	if expected != err.Error() {
		t.Errorf("expected %q, found %q", expected, err.Error())
	}
}

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
	registerRoutes(r, ic, validAuthToken)
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
