package jwt_helper

import (
	"app/domain/model/auth"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGetJwtCredential(t *testing.T) {
	// Set env vars for testing
	os.Setenv("JWT_MEMBER_SECRET_KEY", "test-secret-key")
	os.Setenv("JWT_MEMBER_TTL", "60")
	defer func() {
		os.Unsetenv("JWT_MEMBER_SECRET_KEY")
		os.Unsetenv("JWT_MEMBER_TTL")
	}()

	cred := GetJwtCredential()

	if cred.Member.Secret != "test-secret-key" {
		t.Errorf("Secret = %q, want %q", cred.Member.Secret, "test-secret-key")
	}
	if cred.Member.TLL != 60*time.Minute {
		t.Errorf("TLL = %v, want %v", cred.Member.TLL, 60*time.Minute)
	}
	if cred.Member.Algo != jwt.SigningMethodHS256 {
		t.Errorf("Algo = %v, want HS256", cred.Member.Algo)
	}
}

func TestGetJwtCredential_DefaultTTL(t *testing.T) {
	os.Setenv("JWT_MEMBER_SECRET_KEY", "secret")
	os.Unsetenv("JWT_MEMBER_TTL")
	defer os.Unsetenv("JWT_MEMBER_SECRET_KEY")

	cred := GetJwtCredential()

	// When TTL env is not set, it defaults to 0
	if cred.Member.TLL != 0 {
		t.Errorf("TLL = %v, want 0", cred.Member.TLL)
	}
}

func TestGenerateJWTToken_Success(t *testing.T) {
	jwtCred := JWT{
		Secret: "test-secret-key-for-testing",
		TLL:    30 * time.Minute,
		Algo:   jwt.SigningMethodHS256,
	}

	claims := auth.JWTClaimUser{
		UserID: "test-user-id-123",
	}

	tokenString, err := GenerateJWTToken(jwtCred, claims)
	if err != nil {
		t.Fatalf("GenerateJWTToken() error = %v", err)
	}
	if tokenString == "" {
		t.Error("GenerateJWTToken() returned empty token")
	}

	// Parse back the token to verify claims
	parsed, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaimUser{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtCred.Secret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse generated token: %v", err)
	}

	parsedClaims, ok := parsed.Claims.(*auth.JWTClaimUser)
	if !ok {
		t.Fatal("Failed to cast claims")
	}
	if parsedClaims.UserID != "test-user-id-123" {
		t.Errorf("UserID = %q, want %q", parsedClaims.UserID, "test-user-id-123")
	}
	if parsedClaims.Issuer != "member" {
		t.Errorf("Issuer = %q, want %q", parsedClaims.Issuer, "member")
	}
	if parsedClaims.ExpiresAt == nil {
		t.Error("ExpiresAt should not be nil")
	}
	if parsedClaims.IssuedAt == nil {
		t.Error("IssuedAt should not be nil")
	}
	if parsedClaims.NotBefore == nil {
		t.Error("NotBefore should not be nil")
	}
	if parsedClaims.ID == "" {
		t.Error("ID should not be empty")
	}
}

func TestGenerateJWTToken_ZeroTTL(t *testing.T) {
	jwtCred := JWT{
		Secret: "test-secret-key",
		TLL:    0,
		Algo:   jwt.SigningMethodHS256,
	}

	claims := auth.JWTClaimUser{
		UserID: "user-123",
	}

	tokenString, err := GenerateJWTToken(jwtCred, claims)
	if err != nil {
		t.Fatalf("GenerateJWTToken() error = %v", err)
	}
	if tokenString == "" {
		t.Error("GenerateJWTToken() returned empty token")
	}

	// Parse back and verify no expiry
	parsed, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaimUser{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtCred.Secret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse generated token: %v", err)
	}

	parsedClaims := parsed.Claims.(*auth.JWTClaimUser)
	if parsedClaims.ExpiresAt != nil {
		t.Error("ExpiresAt should be nil when TTL is 0s")
	}
}

func TestGenerateJWTToken_PresetClaims(t *testing.T) {
	jwtCred := JWT{
		Secret: "test-secret-key",
		TLL:    30 * time.Minute,
		Algo:   jwt.SigningMethodHS256,
	}

	now := time.Now()
	claims := auth.JWTClaimUser{
		UserID: "user-456",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "custom-id",
			Issuer:    "custom-issuer",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
		},
	}

	tokenString, err := GenerateJWTToken(jwtCred, claims)
	if err != nil {
		t.Fatalf("GenerateJWTToken() error = %v", err)
	}

	parsed, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaimUser{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtCred.Secret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	parsedClaims := parsed.Claims.(*auth.JWTClaimUser)
	if parsedClaims.ID != "custom-id" {
		t.Errorf("ID = %q, want %q", parsedClaims.ID, "custom-id")
	}
	if parsedClaims.Issuer != "custom-issuer" {
		t.Errorf("Issuer = %q, want %q", parsedClaims.Issuer, "custom-issuer")
	}
}

func TestGenerateJWTToken_UnsupportedClaims(t *testing.T) {
	jwtCred := JWT{
		Secret: "test-secret-key",
		TLL:    30 * time.Minute,
		Algo:   jwt.SigningMethodHS256,
	}

	// Pass a non-JWTClaimUser type
	claims := jwt.RegisteredClaims{
		Issuer: "test",
	}

	_, err := GenerateJWTToken(jwtCred, claims)
	if err == nil {
		t.Error("GenerateJWTToken() should return error for unsupported claims type")
	}
	if err.Error() != "claim data not supported" {
		t.Errorf("error = %q, want %q", err.Error(), "claim data not supported")
	}
}
