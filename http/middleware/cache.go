package middleware

import (
	"bytes"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/web"
	"github.com/gin-gonic/gin"
)

const (
	DefaultCacheTTL         = 5 * time.Minute
	DefaultMaxCacheBodySize = 2 * 1024 * 1024 // 2 MB max body size for Redis storage
	DefaultCacheKeyPrefix   = "cache:"
)

func formatCacheKey(key string) string {
	if strings.HasPrefix(key, DefaultCacheKeyPrefix) {
		return key
	}
	return DefaultCacheKeyPrefix + key
}

type CacheParam struct {
	Key           string
	TTL           time.Duration
	ExcludedPaths []string
	ShouldSkip    func(ctx *web.Context) bool
	Methods       []string
	StatusCodes   []int
	MaxBodySize   int64
}

type CachedResponse struct {
	Status int                 `json:"status"`
	Header map[string][]string `json:"header"`
	Body   []byte              `json:"body"`
}

type Cache struct {
	cacheService contracts.CacheService
}

func NewCache(cacheService contracts.CacheService) *Cache {
	return &Cache{
		cacheService: cacheService,
	}
}

var defaultExcludedPrefixes = []string{
	"/auth",
	"/login",
	"/logout",
	"/csrf-cookie",
	"/sso",
	"/oauth",
	"/callback",
	"/swagger",
	"/doc",
}

var sensitiveHeadersToStrip = map[string]bool{
	"set-cookie":          true,
	"cookie":              true,
	"authorization":       true,
	"proxy-authenticate":  true,
	"proxy-authorization": true,
	"x-cache":             true,
	"vary":                true,
	"connection":          true,
	"keep-alive":          true,
	"transfer-encoding":   true,
	"upgrade":             true,
}

func isDefaultExcludedPath(path string) bool {
	lower := strings.ToLower(path)
	for _, prefix := range defaultExcludedPrefixes {
		if strings.HasPrefix(lower, prefix) || strings.Contains(lower, "/swagger") || strings.Contains(lower, "/doc") {
			return true
		}
	}
	return false
}

func isSensitiveOrCORSHeader(name string) bool {
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "access-control-") {
		return true
	}
	return sensitiveHeadersToStrip[lower]
}

func cleanHeadersForCache(h http.Header) map[string][]string {
	cleaned := make(map[string][]string)
	for k, values := range h {
		if isSensitiveOrCORSHeader(k) {
			continue
		}
		cleaned[k] = values
	}
	return cleaned
}

func isCacheableResponse(header http.Header) bool {
	cc := strings.ToLower(header.Get("Cache-Control"))
	if strings.Contains(cc, "no-store") || strings.Contains(cc, "no-cache") || strings.Contains(cc, "private") {
		return false
	}
	if strings.ToLower(header.Get("Pragma")) == "no-cache" {
		return false
	}
	return true
}

type cachedResponseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *cachedResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *cachedResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (w *cachedResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *cachedResponseWriter) Status() int {
	if w.status != 0 {
		return w.status
	}
	return w.ResponseWriter.Status()
}

func (c *Cache) Handle(param interface{}) web.HandlerFunc {
	var p CacheParam
	switch v := param.(type) {
	case string:
		p.Key = v
	case CacheParam:
		p = v
	case *CacheParam:
		if v != nil {
			p = *v
		}
	case time.Duration:
		p.TTL = v
	}

	if p.TTL <= 0 {
		p.TTL = DefaultCacheTTL
	}
	if p.MaxBodySize <= 0 {
		p.MaxBodySize = DefaultMaxCacheBodySize
	}
	if len(p.Methods) == 0 {
		p.Methods = []string{http.MethodGet, http.MethodHead}
	}
	if len(p.StatusCodes) == 0 {
		p.StatusCodes = []int{http.StatusOK}
	}

	return func(ctx *web.Context) {
		// If key is empty, cache service is nil, or Redis is not configured, bypass caching
		if p.Key == "" || c.cacheService == nil || c.cacheService.Client() == nil {
			ctx.Next()
			return
		}

		// Only cache configured HTTP methods (default GET, HEAD)
		if !slices.Contains(p.Methods, ctx.Request.Method) {
			ctx.Next()
			return
		}

		path := ctx.Request.URL.Path

		// Exclude auth and swagger endpoints by default
		if isDefaultExcludedPath(path) {
			ctx.Next()
			return
		}

		// Check custom excluded paths
		for _, excluded := range p.ExcludedPaths {
			if strings.HasPrefix(path, excluded) {
				ctx.Next()
				return
			}
		}

		// Custom skip predicate
		if p.ShouldSkip != nil && p.ShouldSkip(ctx) {
			ctx.Next()
			return
		}

		key := formatCacheKey(p.Key)

		// Check if response exists in cache
		var cached CachedResponse
		err := c.cacheService.Get(ctx.Request.Context(), key, &cached)
		if err == nil && len(cached.Body) > 0 {
			for k, values := range cached.Header {
				if isSensitiveOrCORSHeader(k) {
					continue
				}
				ctx.Writer.Header()[k] = values
			}
			ctx.Writer.Header().Set("X-Cache", "HIT")
			ctx.Writer.WriteHeader(cached.Status)
			_, _ = ctx.Writer.Write(cached.Body)
			ctx.Abort()
			return
		}

		// Cache MISS: intercept response writer
		writer := &cachedResponseWriter{
			ResponseWriter: ctx.Writer,
			body:           &bytes.Buffer{},
		}
		ctx.Writer = writer
		ctx.Writer.Header().Set("X-Cache", "MISS")

		ctx.Next()

		// Save response to cache if status code matches, Cache-Control allows, and size is within limit
		status := writer.Status()
		if slices.Contains(p.StatusCodes, status) && isCacheableResponse(writer.Header()) {
			bodyBytes := writer.body.Bytes()
			if int64(len(bodyBytes)) <= p.MaxBodySize {
				resp := CachedResponse{
					Status: status,
					Header: cleanHeadersForCache(writer.Header()),
					Body:   bodyBytes,
				}
				_ = c.cacheService.Set(ctx.Request.Context(), key, resp, p.TTL)
			}
		}
	}
}
