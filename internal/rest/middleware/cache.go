package middleware

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func md5String(url string) string {
	h := md5.New()
	io.WriteString(h, url)
	return hex.EncodeToString(h.Sum(nil))
}

type CachedData struct {
	Status   int
	Body     []byte
	Header   http.Header
	CachedAt time.Time
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (c *appMiddleware) _GetCache(key string) (*CachedData, error) {
	if data, err := c.cache.store.Get(context.Background(), key).Result(); err == nil {
		var cch *CachedData
		dec := gob.NewDecoder(bytes.NewBuffer([]byte(data)))
		dec.Decode(&cch)
		return cch, nil
	} else {
		return nil, err
	}
}

func (c *appMiddleware) _SetCache(key string, cch *CachedData, expired time.Duration) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	if err := enc.Encode(*cch); err != nil {
		panic(err)
	}

	return c.cache.store.Set(
		context.Background(),
		key,
		b.String(),
		expired,
	).Err()
}

func (m *appMiddleware) Cache(expiry ...time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		if !m.cache.enabled {
			c.Next()
			return
		}

		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		tohash := c.Request.URL.RequestURI()
		for _, k := range m.cache.headerKeys {
			if v, ok := c.Request.Header[k]; ok {
				tohash += k + strings.Join(v, "")
			}
		}

		cacheKey := m.cache.cachePrefix + md5String(tohash)

		if data, err := m._GetCache(cacheKey); err == nil {
			start := time.Now()
			c.Writer.WriteHeader(data.Status)
			for k, val := range data.Header {
				for _, v := range val {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.Header().Add("X-Gin-Cache", fmt.Sprintf("%f ms", time.Since(start).Seconds()*1000))
			c.Writer.Header().Add("X-Cached-At", data.CachedAt.String())
			c.Writer.Write(data.Body)
			c.Abort()
			return
		}

		// using separate writer to capture response
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		c.Next()

		response := w.body.Bytes()
		responseStatus := c.Writer.Status()

		if responseStatus == http.StatusOK {
			cacheTTL := m.cache.storeTTL
			if len(expiry) > 0 {
				cacheTTL = expiry[0]
			}
			if err := m._SetCache(cacheKey, &CachedData{
				Status:   responseStatus,
				Body:     response,
				Header:   http.Header(w.Header()),
				CachedAt: time.Now(),
			}, cacheTTL); err != nil {
				log.Printf("failed to set cache %v", err)
			}
		}
	}
}
