package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewRequest(t *testing.T, method, path string, payload interface{}) *http.Request {
	var body io.Reader

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("could not marshal request body: %v", err)
		}

		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, path, body)
	assert.NoError(t, err)

	return req
}
