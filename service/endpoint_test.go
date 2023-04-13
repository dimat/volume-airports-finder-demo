package service

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockEndpoint struct {
	CallFunc func(interface{}) (interface{}, error)
}

func (m *MockEndpoint) Call(req interface{}) (interface{}, error) {
	return m.CallFunc(req)
}

func TestHandlerFunc(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        io.Reader
		mockCallFunc   func(interface{}) (interface{}, error)
		expectedStatus int
		expectedResp   string
		method         string
	}{
		{
			name:           "Malformed JSON",
			method:         http.MethodPost,
			reqBody:        bytes.NewBufferString(`{"key": "value",}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid HTTP method",
			method: http.MethodGet,
			reqBody: bytes.NewBufferString(`{
				"key": "value"
			}`),
			mockCallFunc: func(interface{}) (interface{}, error) {
				return map[string]string{"result": "success"}, nil
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "Endpoint returns error",
			method: http.MethodPost,
			reqBody: bytes.NewBufferString(`{
				"key": "value"
			}`),
			mockCallFunc: func(interface{}) (interface{}, error) {
				return nil, errors.New("test error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Valid request and response",
			method: http.MethodPost,
			reqBody: bytes.NewBufferString(`{
				"key": "value"
			}`),
			mockCallFunc: func(interface{}) (interface{}, error) {
				return map[string]string{"result": "success"}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResp:   "{\"result\":\"success\"}\n",
		},
		{
			name:   "Response not encodable",
			method: http.MethodPost,
			reqBody: bytes.NewBufferString(`{
				"key": "value"
			}`),
			mockCallFunc: func(interface{}) (interface{}, error) {
				return func() {}, nil // A function is not JSON encodable
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockEndpoint := &MockEndpoint{CallFunc: tc.mockCallFunc}
			handler := HandlerFunc[interface{}, interface{}](mockEndpoint)

			req := httptest.NewRequest(tc.method, "/", tc.reqBody)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rec.Code)
			}

			if tc.expectedResp != "" {
				resp, err := io.ReadAll(rec.Body)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if string(resp) != tc.expectedResp {
					t.Errorf("Expected response %v, got %v", tc.expectedResp, string(resp))
				}
			}
		})
	}
}
