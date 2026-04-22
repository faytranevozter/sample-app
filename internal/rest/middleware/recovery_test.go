package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecovery_NoPanic(t *testing.T) {
	mdl := newTestMiddleware("secret")

	r := gin.New()
	r.Use(mdl.Recovery())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRecovery_WithPanic(t *testing.T) {
	mdl := newTestMiddleware("secret")

	r := gin.New()
	r.Use(mdl.Recovery())
	r.GET("/test", func(c *gin.Context) {
		panic("something went wrong")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestRecovery_WithPanicError(t *testing.T) {
	mdl := newTestMiddleware("secret")

	r := gin.New()
	r.Use(mdl.Recovery())
	r.GET("/test", func(c *gin.Context) {
		panic(42)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
