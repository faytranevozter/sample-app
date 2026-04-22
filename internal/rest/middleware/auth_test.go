package middleware

import (
	"app/domain/model/auth"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestMiddleware(secret string) *appMiddleware {
	return &appMiddleware{
		secret: secret,
	}
}

func generateTestToken(secret string, claims auth.JWTClaimUser) string {
	if claims.RegisteredClaims.ID == "" {
		claims.RegisteredClaims.ID = "test-id"
	}
	if claims.RegisteredClaims.Issuer == "" {
		claims.RegisteredClaims.Issuer = "member"
	}
	if claims.RegisteredClaims.IssuedAt == nil {
		claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Now())
	}
	if claims.RegisteredClaims.NotBefore == nil {
		claims.RegisteredClaims.NotBefore = jwt.NewNumericDate(time.Now())
	}
	if claims.RegisteredClaims.ExpiresAt == nil {
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestAuth_NoHeader(t *testing.T) {
	mdl := newTestMiddleware("secret")

	r := gin.New()
	r.GET("/test", mdl.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	msg, _ := resp["message"].(string)
	if msg != "Unauthorized: Header authorization is required" {
		t.Errorf("message = %q", msg)
	}
}

func TestAuth_InvalidBearerFormat(t *testing.T) {
	mdl := newTestMiddleware("secret")

	r := gin.New()
	r.GET("/test", mdl.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Token something")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_MalformedToken(t *testing.T) {
	mdl := newTestMiddleware("secret")

	r := gin.New()
	r.GET("/test", mdl.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-jwt-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_WrongSignature(t *testing.T) {
	mdl := newTestMiddleware("correct-secret")

	// Generate token with different secret
	tokenString := generateTestToken("wrong-secret", auth.JWTClaimUser{
		UserID: "user-123",
	})

	r := gin.New()
	r.GET("/test", mdl.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_ExpiredToken(t *testing.T) {
	secret := "test-secret"
	mdl := newTestMiddleware(secret)

	claims := auth.JWTClaimUser{
		UserID: "user-123",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "test-id",
			Issuer:    "member",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))

	r := gin.New()
	r.GET("/test", mdl.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	msg, _ := resp["message"].(string)
	if msg != "Unauthorized: Token expired" {
		t.Errorf("message = %q, want %q", msg, "Unauthorized: Token expired")
	}
}

func TestAuth_ValidToken(t *testing.T) {
	secret := "test-secret"
	mdl := newTestMiddleware(secret)

	tokenString := generateTestToken(secret, auth.JWTClaimUser{
		UserID: "user-123",
	})

	var capturedClaims auth.JWTClaimUser
	r := gin.New()
	r.GET("/test", mdl.Auth(), func(c *gin.Context) {
		capturedClaims = c.MustGet("token_data").(auth.JWTClaimUser)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	if capturedClaims.UserID != "user-123" {
		t.Errorf("UserID = %q, want %q", capturedClaims.UserID, "user-123")
	}
}

func TestAuth_EmptyBearerToken(t *testing.T) {
	mdl := newTestMiddleware("secret")

	r := gin.New()
	r.GET("/test", mdl.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
