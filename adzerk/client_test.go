package adzerk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"encoding/json"

	"golang.org/x/net/context"
)

type testCase struct {
	IP               string
	response         string
	resolvedResponse map[string]interface{}
}

func TestDo(t *testing.T) {
	tcs := []testCase{
		{
			IP:               "127.0.0.1",
			response:         `{"response":"response"}`,
			resolvedResponse: map[string]interface{}{"response": "response"},
		},
	}
	for _, tc := range tcs {
		runTestDo(t, tc)
	}
}

func runTestDo(t *testing.T, tc testCase) {
	// Adzerk mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Expected nil err, got %s", err.Error())
		}

		var actualReq Request
		err = json.Unmarshal(reqBody, &actualReq)
		if err != nil {
			t.Errorf("Expected nil json unmarshal err, instead got %s", err)
		}

		if actualReq.IP != tc.IP {
			t.Errorf("Expected forwarded IP: %s, got %s", tc.IP, actualReq.IP)
		}

		fmt.Fprintln(w, tc.response)
	}))

	c := NewClient(nil)
	c.URL = ts.URL
	req, err := c.NewRequest(RequestData{
		IP: tc.IP,
	})
	if err != nil {
		t.Errorf("Expected nil err from NewRequest, got %s", err)
	}
	var resp map[string]interface{}
	response, err := c.Do(context.Background(), req, &resp)
	if err != nil {
		t.Errorf("Expected nil err from Do, got %s", err)
	}
	if !reflect.DeepEqual(resp, tc.resolvedResponse) {
		t.Errorf("Expected response: %v, got %v", tc.resolvedResponse, &resp)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}
