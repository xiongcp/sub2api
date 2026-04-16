package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	requestBodyReadFailedCode              = "request_body_read_failed"
	requestBodyReadFailedMessage           = "Request body upload was interrupted. Please retry."
	requestBodyReadFailedRetryAfterSeconds = 1
)

type apiErrorFormat int

const (
	apiErrorFormatOpenAI apiErrorFormat = iota
	apiErrorFormatAnthropic
	apiErrorFormatClaudeCompat
	apiErrorFormatGoogle
)

func handleRetryableRequestBodyReadError(c *gin.Context, reqLog *zap.Logger, err error, format apiErrorFormat) {
	logRequestBodyReadFailure(c, reqLog, err)
	if c == nil {
		return
	}

	c.Header("Retry-After", strconv.Itoa(requestBodyReadFailedRetryAfterSeconds))

	switch format {
	case apiErrorFormatGoogle:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": requestBodyReadFailedMessage,
				"status":  "INVALID_ARGUMENT",
				"details": []gin.H{
					{
						"reason":              strings.ToUpper(requestBodyReadFailedCode),
						"retry_after_seconds": requestBodyReadFailedRetryAfterSeconds,
					},
				},
			},
		})
	case apiErrorFormatAnthropic:
		c.JSON(http.StatusBadRequest, gin.H{
			"type": "error",
			"error": gin.H{
				"type":    "invalid_request_error",
				"code":    requestBodyReadFailedCode,
				"message": requestBodyReadFailedMessage,
			},
		})
	case apiErrorFormatClaudeCompat:
		c.JSON(http.StatusBadRequest, gin.H{
			"type": "error",
			"error": gin.H{
				"type":    "invalid_request_error",
				"code":    requestBodyReadFailedCode,
				"message": requestBodyReadFailedMessage,
			},
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"type":    "invalid_request_error",
				"code":    requestBodyReadFailedCode,
				"message": requestBodyReadFailedMessage,
			},
		})
	}
}

func logRequestBodyReadFailure(c *gin.Context, reqLog *zap.Logger, err error) {
	if reqLog == nil {
		reqLog = requestLogger(c, "handler.request_body_read")
	}

	fields := make([]zap.Field, 0, 8)
	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	if c != nil && c.Request != nil {
		fields = append(fields,
			zap.String("method", c.Request.Method),
			zap.String("path", requestURLPath(c)),
			zap.Int64("content_length", c.Request.ContentLength),
			zap.Strings("transfer_encoding", append([]string(nil), c.Request.TransferEncoding...)),
			zap.String("content_type", strings.TrimSpace(c.GetHeader("Content-Type"))),
			zap.String("user_agent", strings.TrimSpace(c.Request.UserAgent())),
		)
	}
	if clientIP := strings.TrimSpace(ip.GetClientIP(c)); clientIP != "" {
		fields = append(fields, zap.String("client_ip", clientIP))
	}

	reqLog.Warn("request_body_read_failed", fields...)
}

func requestURLPath(c *gin.Context) string {
	if c == nil || c.Request == nil || c.Request.URL == nil {
		return ""
	}
	return c.Request.URL.Path
}
