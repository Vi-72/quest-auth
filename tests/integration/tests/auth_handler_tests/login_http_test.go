package auth_handler_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (s *Suite) TestLoginUser_Success() {
	// First register
	reg := map[string]any{
		"email":    "login_int@example.com",
		"phone":    "+1234567880",
		"name":     "Login Int",
		"password": "securepassword123",
	}
	rb, _ := json.Marshal(reg)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	s.TestDIContainer.HTTPRouter.ServeHTTP(rr, req)
	s.Require().Equal(http.StatusCreated, rr.Code, rr.Body.String())

	// Now login
	body := map[string]any{
		"email":    reg["email"],
		"password": "securepassword123",
	}
	b, _ := json.Marshal(body)
	lreq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	lreq.Header.Set("Content-Type", "application/json")
	lrr := httptest.NewRecorder()
	s.TestDIContainer.HTTPRouter.ServeHTTP(lrr, lreq)

	s.Require().Equal(http.StatusOK, lrr.Code, lrr.Body.String())
	var resp struct {
		User struct {
			ID    string  `json:"id"`
			Email string  `json:"email"`
			Name  string  `json:"name"`
			Phone *string `json:"phone"`
		} `json:"user"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}
	s.Require().NoError(json.Unmarshal(lrr.Body.Bytes(), &resp))
	s.Assert().NotEmpty(resp.User.ID)
	s.Assert().Equal(reg["email"], resp.User.Email)
	s.Assert().Equal(reg["name"], resp.User.Name)
	s.Require().NotNil(resp.User.Phone)
	s.Assert().Equal(reg["phone"], *resp.User.Phone)
	s.Assert().NotEmpty(resp.AccessToken)
	s.Assert().NotEmpty(resp.RefreshToken)
}

func (s *Suite) TestLoginUser_Unauthorized() {
	// Non-existing user
	body := map[string]any{
		"email":    "nouser@example.com",
		"password": "wrongpass",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	s.TestDIContainer.HTTPRouter.ServeHTTP(rr, req)
	s.Assert().Equal(http.StatusUnauthorized, rr.Code)
}
