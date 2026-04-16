package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestHandleRetryableRequestBodyReadError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		format           apiErrorFormat
		wantTopLevelType string
		wantErrorCode    string
		wantErrorType    string
	}{
		{
			name:          "openai",
			format:        apiErrorFormatOpenAI,
			wantErrorCode: requestBodyReadFailedCode,
			wantErrorType: "invalid_request_error",
		},
		{
			name:             "anthropic",
			format:           apiErrorFormatAnthropic,
			wantTopLevelType: "error",
			wantErrorCode:    requestBodyReadFailedCode,
			wantErrorType:    "invalid_request_error",
		},
		{
			name:             "claude_compat",
			format:           apiErrorFormatClaudeCompat,
			wantTopLevelType: "error",
			wantErrorCode:    requestBodyReadFailedCode,
			wantErrorType:    "invalid_request_error",
		},
		{
			name:          "google",
			format:        apiErrorFormatGoogle,
			wantErrorCode: "400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			req := httptest.NewRequest(http.MethodPost, "/responses", nil)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handleRetryableRequestBodyReadError(c, nil, errors.New("unexpected EOF"), tt.format)

			require.Equal(t, http.StatusBadRequest, rec.Code)
			require.Equal(t, "1", rec.Header().Get("Retry-After"))
			require.Equal(t, requestBodyReadFailedMessage, gjson.Get(rec.Body.String(), "error.message").String())
			if tt.wantTopLevelType != "" {
				require.Equal(t, tt.wantTopLevelType, gjson.Get(rec.Body.String(), "type").String())
			}
			if tt.wantErrorType != "" {
				require.Equal(t, tt.wantErrorType, gjson.Get(rec.Body.String(), "error.type").String())
			}
			require.Equal(t, tt.wantErrorCode, gjson.Get(rec.Body.String(), "error.code").String())
		})
	}
}
