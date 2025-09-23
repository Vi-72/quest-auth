package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRefreshTokensPreservesClaims(t *testing.T) {
	service := NewService("secret", time.Minute, time.Hour)

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	email := "user@example.com"
	name := "John Doe"
	phone := "+1234567890"
	createdAt := time.Unix(1700000000, 0).UTC()

	pair, err := service.GenerateTokenPair(userID, email, name, phone, createdAt)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	refreshedPair, err := service.RefreshTokens(pair.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshTokens() error = %v", err)
	}

	claims, err := service.ValidateAccessToken(refreshedPair.AccessToken)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if claims.Name != name {
		t.Fatalf("expected name %q, got %q", name, claims.Name)
	}

	if claims.Phone != phone {
		t.Fatalf("expected phone %q, got %q", phone, claims.Phone)
	}

	if !claims.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected createdAt %v, got %v", createdAt, claims.CreatedAt)
	}
}
