/*
 Copyright 2020 Keyporttech Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the Licens}ze at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// tests validate payload
func TestValidatePayload(t *testing.T) {
	const defaultBody = `{"hey":true}` // All tests below use the default request body and signature.
	const defaultSignature = "126f2c800419c60137ce748d7672e77b65cf16d6"
	secretKey := []byte("0123456789abcdef")
	tests := []struct {
		signature   string
		eventID     string
		event       string
		wantEventID string
		wantEvent   string
		wantPayload string
	}{
		// The following tests generate expected errors:
		{},                         // Missing signature
		{signature: "yo"},     // Signature not hex string
		{signature: "012345"}, // Invalid signature
		// The following tests expect err=nil:
		{
			signature:   defaultSignature,
			eventID:     "caesar-salad",
			event:       "ping",
			wantEventID: "caesar-salad",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		},
		{
			signature:   defaultSignature,
			event:       "ping",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		},
		{
			signature:   "b1f8020f5b4cd42042f807dd939015c4a418bc1ff7f604dd55b0a19b5d953d9b",
			event:       "ping",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		}
	}

	for _, test := range tests {
		buf := bytes.NewBufferString(defaultBody)
		req, err := http.NewRequest("GET", "http://localhost/event", buf)
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		if test.signature != "" {
			req.Header.Set(signatureHeader, test.signature)
		}
		req.Header.Set("Content-Type", "application/json")

		got, err := ValidatePayload(req, secretKey)
		if err != nil {
			if test.wantPayload != "" {
				t.Errorf("ValidatePayload(%#v): err = %v, want nil", test, err)
			}
			continue
		}
		if string(got) != test.wantPayload {
			t.Errorf("ValidatePayload = %q, want %q", got, test.wantPayload)
		}
	}
}
