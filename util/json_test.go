package util

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncode(t *testing.T) {
	type payload struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name       string
		statusCode int
		value      payload
		want       string
	}{
		{
			name:       "encode valid payload",
			statusCode: http.StatusOK,
			value:      payload{Message: "Hello, world!"},
			want:       `{"message":"Hello, world!"}` + "\n", // json.Encoder adds a newline
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			err := Encode(w, r, tt.statusCode, tt.value)
			if err != nil {
				t.Errorf("Encode() error = %v", err)
			}

			if w.Code != tt.statusCode {
				t.Errorf("Encode() statusCode = %v, want %v", w.Code, tt.statusCode)
			}

			if got := w.Body.String(); got != tt.want {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type payload struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name    string
		body    string
		want    payload
		wantErr bool
	}{
		{
			name:    "decode valid payload",
			body:    `{"message":"Hello, world!"}`,
			want:    payload{Message: "Hello, world!"},
			wantErr: false,
		},
		{
			name:    "decode invalid payload",
			body:    `{"message":}`,
			want:    payload{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))

			got, err := Decode[payload](r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
