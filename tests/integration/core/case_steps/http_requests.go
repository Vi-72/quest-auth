package casesteps

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
)

// HTTPRequest represents a test HTTP request
type HTTPRequest struct {
	Method      string
	URL         string
	Body        interface{}
	Headers     map[string]string
	ContentType string
}

// HTTPResponse represents a test HTTP response
type HTTPResponse struct {
	StatusCode int
	Body       string
	Headers    http.Header
}

// ExecuteHTTPRequest performs an HTTP request against provided handler
func ExecuteHTTPRequest(ctx context.Context, handler http.Handler, req HTTPRequest) (*HTTPResponse, error) {
	var body io.Reader

	if req.Body != nil {
		switch v := req.Body.(type) {
		case string:
			body = bytes.NewReader([]byte(v))
		case json.RawMessage:
			body = bytes.NewReader([]byte(v))
		default:
			b, err := json.Marshal(req.Body)
			if err != nil {
				return nil, err
			}
			body = bytes.NewReader(b)
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, body)
	if err != nil {
		return nil, err
	}

	if req.ContentType != "" {
		httpReq.Header.Set("Content-Type", req.ContentType)
	} else if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httpReq)

	return &HTTPResponse{
		StatusCode: rr.Code,
		Body:       rr.Body.String(),
		Headers:    rr.Header(),
	}, nil
}

// RegisterHTTPRequest builds request for user registration
func RegisterHTTPRequest(body interface{}) HTTPRequest {
	return HTTPRequest{
		Method:      http.MethodPost,
		URL:         "/api/v1/auth/register",
		Body:        body,
		ContentType: "application/json",
	}
}

// LoginHTTPRequest builds request for user login
func LoginHTTPRequest(body interface{}) HTTPRequest {
	return HTTPRequest{
		Method:      http.MethodPost,
		URL:         "/api/v1/auth/login",
		Body:        body,
		ContentType: "application/json",
	}
}
