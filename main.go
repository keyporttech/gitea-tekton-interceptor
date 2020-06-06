/*
 Copyright 2020 Keyporttech Inc.
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

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	// Environment variable containing gitea secret token
	envSecret = "GITEA_SECRET_TOKEN"
	// sha1Prefix is the prefix used by gitea before the HMAC hexdigest.
	sha1Prefix = "sha1"
	// sha256Prefix and sha512Prefix are provided for future compatibility.
	sha256Prefix = "sha256"
	sha512Prefix = "sha512"
	// signatureHeader is the Gitea header key used to pass the HMAC hexdigest.
	signatureHeader = "X-Gitea-Signature"
	// eventTypeHeader is the Gitea header key used to pass the event type.
	eventTypeHeader = "X-Gitea-Event"
	// deliveryIDHeader is the Gitea header key used to pass the unique ID for the webhook event.
	deliveryIDHeader = "X-Gitea-Delivery"
)

// DeliveryID returns the unique delivery ID of webhook request r.
//
func DeliveryID(r *http.Request) string {
	return r.Header.Get(deliveryIDHeader)
}

// genMAC generates the HMAC signature for a message provided the secret key
// and hashFunc.
func genMAC(message, key []byte, hashFunc func() hash.Hash) []byte {
	mac := hmac.New(hashFunc, key)
	mac.Write(message)
	return mac.Sum(nil)
}

// checkMAC reports whether messageMAC is a valid HMAC tag for message.
func checkMAC(message, messageMAC, key []byte, hashFunc func() hash.Hash) bool {
	expectedMAC := genMAC(message, key, hashFunc)
	fmt.Printf("expecting %x: %x", messageMAC, expectedMAC)
	return hmac.Equal(messageMAC, expectedMAC)
}

// messageMAC returns the hex-decoded HMAC tag from the signature and its
// corresponding hash function.
func messageMAC(signature string) ([]byte, func() hash.Hash, error) {
	if signature == "" {
		return nil, nil, errors.New("missing signature")
	}

	var hashFunc func() hash.Hash

	hashFunc = sha256.New

	buf, err := hex.DecodeString(signature)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding signature %q: %v", signature, err)
	}
	return buf, hashFunc, nil
}

// ValidateSignature validates the signature for the given payload.
// signature is the gitea hash signature delivered in the X-Gitea-Signature header.
// payload is the JSON payload sent by gitea Webhooks.
// secretToken is the gitea Webhook secret token.
//
func ValidateSignature(signature string, payload, secretToken []byte) error {
	messageMAC, hashFunc, err := messageMAC(signature)
	if err != nil {
		return err
	}
	fmt.Printf("signature=%o", messageMAC )
	if !checkMAC(payload, messageMAC, secretToken, hashFunc) {
		return errors.New("payload signature check failed")
	}
	return nil
}

func ValidatePayload(r *http.Request, secretToken []byte) (payload []byte, err error) {
	var body []byte // Raw body that gitea uses to calculate the signature.

	switch ct := r.Header.Get("Content-Type"); ct {
	case "application/json":
		var err error
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, err
		}

		// If the content type is application/json,
		// the JSON payload is just the original body.
		payload = body

	case "application/x-www-form-urlencoded":
		// payloadFormParam is the name of the form parameter that the JSON payload
		// will be in if a webhook has its content type set to application/x-www-form-urlencoded.
		const payloadFormParam = "payload"

		var err error
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, err
		}

		// If the content type is application/x-www-form-urlencoded,
		// the JSON payload will be under the "payload" form param.
		form, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		payload = []byte(form.Get(payloadFormParam))

	default:
		return nil, fmt.Errorf("Webhook request has unsupported Content-Type %q", ct)
	}

	// Only validate the signature if a secret token exists. This is intended for
	// local development only and all webhooks should ideally set up a secret token.
	if len(secretToken) > 0 {
		sig := r.Header.Get(signatureHeader)
		if err := ValidateSignature(sig, body, secretToken); err != nil {
			return nil, err
		}
	}

	return payload, nil
}

// main function
func main() {
	secretToken := os.Getenv(envSecret)
	if secretToken == "" {
		log.Fatalf("No secret token given")
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//TODO: We should probably send over the EL eventID as a X-Tekton-Event-Id header as well
		payload, err := ValidatePayload(request, []byte(secretToken))
		id := DeliveryID(request)
		if err != nil {
			http.Error(writer, fmt.Sprint(err), http.StatusBadRequest)
		}
		n, err := writer.Write(payload)
		if err != nil {
			log.Printf("Failed to write response for gitea event ID: %s. Bytes writted: %d. Error: %q", id, n, err)
		}
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil))
}
