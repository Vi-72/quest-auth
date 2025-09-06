package auth_handler_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type registerRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (s *Suite) TestRegisterUser_Success() {
	ctx := context.Background()
	_ = ctx

	body := registerRequest{
		Email:    "newuser_int@example.com",
		Phone:    "+1234567899",
		Name:     "Integration User",
		Password: "securepassword123",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	s.TestDIContainer.HTTPRouter.ServeHTTP(rr, req)

	s.Require().Equal(http.StatusCreated, rr.Code, rr.Body.String())

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
	s.Require().NoError(json.Unmarshal(rr.Body.Bytes(), &resp))

	s.Assert().NotEmpty(resp.User.ID)
	s.Assert().Equal(body.Email, resp.User.Email)
	s.Assert().Equal(body.Name, resp.User.Name)
	s.Require().NotNil(resp.User.Phone)
	s.Assert().Equal(body.Phone, *resp.User.Phone)
	s.Assert().NotEmpty(resp.AccessToken)
	s.Assert().NotEmpty(resp.RefreshToken)
	s.Assert().Equal("Bearer", resp.TokenType)
	s.Assert().Greater(resp.ExpiresIn, 0)
}

func (s *Suite) TestRegisterUser_BadRequest() {
	// Missing required fields
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	s.TestDIContainer.HTTPRouter.ServeHTTP(rr, req)
	s.Assert().Equal(http.StatusBadRequest, rr.Code)
}
