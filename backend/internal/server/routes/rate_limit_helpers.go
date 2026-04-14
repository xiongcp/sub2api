package routes

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	ratelimitmiddleware "github.com/Wei-Shaw/sub2api/internal/middleware"
	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
)

func clientIPRateLimitKey(c *gin.Context) (string, bool) {
	if c == nil {
		return "", false
	}
	ip := strings.TrimSpace(c.ClientIP())
	if ip == "" {
		return "", false
	}
	return "ip:" + ip, true
}

func authSubjectRateLimitKey(c *gin.Context) (string, bool) {
	if c == nil {
		return "", false
	}
	if subject, ok := servermiddleware.GetAuthSubjectFromContext(c); ok && subject.UserID > 0 {
		return "user:" + strconv.FormatInt(subject.UserID, 10), true
	}
	return clientIPRateLimitKey(c)
}

func clientIPAndJSONFieldRateLimitKey(field string) ratelimitmiddleware.RateLimitKeyFunc {
	field = strings.TrimSpace(field)
	return func(c *gin.Context) (string, bool) {
		base, ok := clientIPRateLimitKey(c)
		if !ok {
			return "", false
		}
		text, ok := readJSONFieldRateLimitValue(c, field)
		if !ok {
			return base, true
		}
		return base + ":" + text, true
	}
}

func clientIPAndHashedJSONFieldRateLimitKey(field string) ratelimitmiddleware.RateLimitKeyFunc {
	field = strings.TrimSpace(field)
	return func(c *gin.Context) (string, bool) {
		base, ok := clientIPRateLimitKey(c)
		if !ok {
			return "", false
		}
		text, ok := readJSONFieldRateLimitValue(c, field)
		if !ok {
			return base, true
		}
		sum := sha256.Sum256([]byte(text))
		return base + ":sha256:" + hex.EncodeToString(sum[:]), true
	}
}

func readJSONFieldRateLimitValue(c *gin.Context, field string) (string, bool) {
	if field == "" || c == nil || c.Request == nil || c.Request.Body == nil {
		return "", false
	}

	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", false
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(raw))
	if len(raw) == 0 {
		return "", false
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return "", false
	}
	value, ok := payload[field]
	if !ok {
		return "", false
	}
	text := strings.TrimSpace(strings.TrimSpace(toRateLimitString(value)))
	if text == "" {
		return "", false
	}
	return text, true
}

func toRateLimitString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.Itoa(v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}
