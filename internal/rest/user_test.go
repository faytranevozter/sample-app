package rest

import (
	"app/domain/model/auth"
	request_model "app/domain/model/request"
	"app/domain/model/response"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// mockUserService implements UserService
type mockUserService struct {
	loginFn            func(ctx context.Context, payload request_model.LoginRequest) (int, response.Base)
	registerFn         func(ctx context.Context, payload request_model.RegisterRequest) (int, response.Base)
	getMeFn            func(ctx context.Context, claim auth.JWTClaimUser) (int, response.Base)
	sampleUserListFn   func(ctx context.Context, claim auth.JWTClaimUser, query url.Values) (int, response.Base)
	sampleUserDetailFn func(ctx context.Context, claim auth.JWTClaimUser, id string) (int, response.Base)
	sampleUserExportFn func(ctx context.Context, claim auth.JWTClaimUser, query url.Values) (int, response.Base)
}

func (m *mockUserService) Login(ctx context.Context, payload request_model.LoginRequest) (int, response.Base) {
	if m.loginFn != nil {
		return m.loginFn(ctx, payload)
	}
	return http.StatusOK, response.Success(nil)
}

func (m *mockUserService) Register(ctx context.Context, payload request_model.RegisterRequest) (int, response.Base) {
	if m.registerFn != nil {
		return m.registerFn(ctx, payload)
	}
	return http.StatusOK, response.Success(nil)
}

func (m *mockUserService) GetMe(ctx context.Context, claim auth.JWTClaimUser) (int, response.Base) {
	if m.getMeFn != nil {
		return m.getMeFn(ctx, claim)
	}
	return http.StatusOK, response.Success(nil)
}

func (m *mockUserService) SampleUserList(ctx context.Context, claim auth.JWTClaimUser, query url.Values) (int, response.Base) {
	if m.sampleUserListFn != nil {
		return m.sampleUserListFn(ctx, claim, query)
	}
	return http.StatusOK, response.Success(nil)
}

func (m *mockUserService) SampleUserDetail(ctx context.Context, claim auth.JWTClaimUser, id string) (int, response.Base) {
	if m.sampleUserDetailFn != nil {
		return m.sampleUserDetailFn(ctx, claim, id)
	}
	return http.StatusOK, response.Success(nil)
}

func (m *mockUserService) SampleUserExport(ctx context.Context, claim auth.JWTClaimUser, query url.Values) (int, response.Base) {
	if m.sampleUserExportFn != nil {
		return m.sampleUserExportFn(ctx, claim, query)
	}
	return http.StatusOK, response.Success(nil)
}

// mockMiddleware implements middleware.Middleware
type mockMiddleware struct{}

func (m *mockMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("token_data", auth.JWTClaimUser{UserID: "test-user-id"})
		c.Next()
	}
}

func (m *mockMiddleware) Cors() gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func (m *mockMiddleware) Logger(writer io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func (m *mockMiddleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func (m *mockMiddleware) Cache(expiry ...time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func setupRouter(svc UserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api")
	mdl := &mockMiddleware{}
	NewUserHandler(group, svc, mdl)
	return r
}

func TestLogin_Handler_Success(t *testing.T) {
	svc := &mockUserService{
		loginFn: func(ctx context.Context, payload request_model.LoginRequest) (int, response.Base) {
			if payload.Email != "test@example.com" {
				t.Errorf("Email = %q, want %q", payload.Email, "test@example.com")
			}
			return http.StatusOK, response.Success(map[string]string{"token": "abc123"})
		},
	}

	r := setupRouter(svc)
	body := `{"email":"test@example.com","password":"pass123"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp response.Base
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ErrorCode != 0 {
		t.Errorf("ErrorCode = %d, want 0", resp.ErrorCode)
	}
}

func TestLogin_Handler_InvalidJSON(t *testing.T) {
	svc := &mockUserService{}
	r := setupRouter(svc)

	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRegister_Handler_Success(t *testing.T) {
	svc := &mockUserService{
		registerFn: func(ctx context.Context, payload request_model.RegisterRequest) (int, response.Base) {
			if payload.Name != "Test User" {
				t.Errorf("Name = %q, want %q", payload.Name, "Test User")
			}
			return http.StatusOK, response.Success(nil)
		},
	}

	r := setupRouter(svc)
	body := `{"name":"Test User","email":"test@example.com","password":"pass123"}`
	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRegister_Handler_InvalidJSON(t *testing.T) {
	svc := &mockUserService{}
	r := setupRouter(svc)

	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader("bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetMe_Handler_Success(t *testing.T) {
	svc := &mockUserService{
		getMeFn: func(ctx context.Context, claim auth.JWTClaimUser) (int, response.Base) {
			if claim.UserID != "test-user-id" {
				t.Errorf("UserID = %q, want %q", claim.UserID, "test-user-id")
			}
			return http.StatusOK, response.Success(map[string]string{"name": "Test"})
		},
	}

	r := setupRouter(svc)
	req := httptest.NewRequest("GET", "/api/auth/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestUserList_Handler_Success(t *testing.T) {
	svc := &mockUserService{
		sampleUserListFn: func(ctx context.Context, claim auth.JWTClaimUser, query url.Values) (int, response.Base) {
			return http.StatusOK, response.Success(response.List{
				List:  []any{},
				Page:  1,
				Limit: 10,
				Total: 0,
			})
		},
	}

	r := setupRouter(svc)
	req := httptest.NewRequest("GET", "/api/sample/user/list?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestUserDetail_Handler_Success(t *testing.T) {
	svc := &mockUserService{
		sampleUserDetailFn: func(ctx context.Context, claim auth.JWTClaimUser, id string) (int, response.Base) {
			if id != "123e4567-e89b-12d3-a456-426614174000" {
				t.Errorf("id = %q", id)
			}
			return http.StatusOK, response.Success(nil)
		},
	}

	r := setupRouter(svc)
	req := httptest.NewRequest("GET", "/api/sample/user/detail/123e4567-e89b-12d3-a456-426614174000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestUserExport_Handler_Success(t *testing.T) {
	svc := &mockUserService{
		sampleUserExportFn: func(ctx context.Context, claim auth.JWTClaimUser, query url.Values) (int, response.Base) {
			return http.StatusOK, response.Success(map[string]string{"base64": "data"})
		},
	}

	r := setupRouter(svc)
	req := httptest.NewRequest("GET", "/api/sample/user/export", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}
