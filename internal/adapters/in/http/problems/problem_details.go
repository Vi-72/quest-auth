package problems

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ProblemDetails RFC 7807
type ProblemDetails struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func (p *ProblemDetails) Error() string {
	return fmt.Sprintf("%d: %s - %s", p.Status, p.Title, p.Detail)
}

func (p *ProblemDetails) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		return
	}
}

// Factory methods for common HTTP errors

// NewBadRequest creates a 400 Bad Request problem
func NewBadRequest(detail string) *ProblemDetails {
	return &ProblemDetails{
		Type:   "bad-request",
		Title:  "Bad Request",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}

// NewNotFound creates a 404 Not Found problem
func NewNotFound(detail string) *ProblemDetails {
	return &ProblemDetails{
		Type:   "not-found",
		Title:  "Not Found",
		Status: http.StatusNotFound,
		Detail: detail,
	}
}
