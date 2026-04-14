package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestClientIPAndHashedJSONFieldRateLimitKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/public/key-usage/query", strings.NewReader(`{"api_key":"sk-test-123"}`))
	c.Request.RemoteAddr = "198.51.100.2:12345"

	keyFunc := clientIPAndHashedJSONFieldRateLimitKey("api_key")
	key, ok := keyFunc(c)
	require.True(t, ok)

	sum := sha256.Sum256([]byte("sk-test-123"))
	expected := "ip:198.51.100.2:sha256:" + hex.EncodeToString(sum[:])
	require.Equal(t, expected, key)

	raw, err := io.ReadAll(c.Request.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{"api_key":"sk-test-123"}`, string(raw))
}

func TestClientIPAndHashedJSONFieldRateLimitKeyFallsBackWhenFieldMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/public/key-usage/query", strings.NewReader(`{"unexpected":"value"}`))
	c.Request.RemoteAddr = "203.0.113.9:4567"

	keyFunc := clientIPAndHashedJSONFieldRateLimitKey("api_key")
	key, ok := keyFunc(c)
	require.True(t, ok)
	require.Equal(t, "ip:203.0.113.9", key)
}
