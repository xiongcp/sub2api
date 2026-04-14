package handler

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ip"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	publicKeyUsageDefaultPageSize = 20
	publicKeyUsageMaxPageSize     = 50
)

type publicKeyUsageQueryRequest struct {
	APIKey    string `json:"api_key" binding:"required"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	Timezone  string `json:"timezone"`
}

type publicUsageLogItem struct {
	ID                  int64   `json:"id"`
	RequestID           string  `json:"request_id"`
	Model               string  `json:"model"`
	RequestType         string  `json:"request_type"`
	ServiceTier         *string `json:"service_tier,omitempty"`
	ReasoningEffort     *string `json:"reasoning_effort,omitempty"`
	InboundEndpoint     *string `json:"inbound_endpoint,omitempty"`
	InputTokens         int     `json:"input_tokens"`
	OutputTokens        int     `json:"output_tokens"`
	CacheCreationTokens int     `json:"cache_creation_tokens"`
	CacheReadTokens     int     `json:"cache_read_tokens"`
	TotalTokens         int     `json:"total_tokens"`
	ActualCost          float64 `json:"actual_cost"`
	DurationMs          *int    `json:"duration_ms,omitempty"`
	CreatedAt           string  `json:"created_at"`
}

// PublicKeyUsageQuery handles public API key usage lookup for the /key-usage page.
// POST /api/v1/public/key-usage/query
func (h *GatewayHandler) PublicKeyUsageQuery(c *gin.Context) {
	var req publicKeyUsageQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	apiKey, subject, ok := h.authenticatePublicUsageKey(c, req.APIKey)
	if !ok {
		return
	}

	startTime, endTime, err := parsePublicUsageDateRange(req.StartDate, req.EndDate, req.Timezone)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	page, pageSize := normalizePublicUsagePagination(req.Page, req.PageSize)

	payload, err := h.buildUsageResponsePayload(c.Request.Context(), c, apiKey, subject, startTime, endTime)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	usageLogs, err := h.buildPublicUsageLogsPage(c.Request.Context(), apiKey.ID, subject.UserID, startTime, endTime, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	payload["usage_logs"] = usageLogs
	response.Success(c, payload)
}

func (h *GatewayHandler) authenticatePublicUsageKey(c *gin.Context, rawKey string) (*service.APIKey, middleware2.AuthSubject, bool) {
	key := strings.TrimSpace(rawKey)
	if key == "" {
		response.BadRequest(c, "api_key is required")
		return nil, middleware2.AuthSubject{}, false
	}

	apiKey, err := h.apiKeyService.GetByKey(c.Request.Context(), key)
	if err != nil {
		if errors.Is(err, service.ErrAPIKeyNotFound) {
			response.Unauthorized(c, "Invalid API key")
			return nil, middleware2.AuthSubject{}, false
		}
		response.InternalError(c, "Failed to validate API key")
		return nil, middleware2.AuthSubject{}, false
	}

	if !apiKey.IsActive() &&
		apiKey.Status != service.StatusAPIKeyExpired &&
		apiKey.Status != service.StatusAPIKeyQuotaExhausted {
		response.Unauthorized(c, "API key is disabled")
		return nil, middleware2.AuthSubject{}, false
	}

	if len(apiKey.IPWhitelist) > 0 || len(apiKey.IPBlacklist) > 0 {
		clientIP := ip.GetTrustedClientIP(c)
		allowed, _ := ip.CheckIPRestrictionWithCompiledRules(clientIP, apiKey.CompiledIPWhitelist, apiKey.CompiledIPBlacklist)
		if !allowed {
			response.Forbidden(c, "Access denied")
			return nil, middleware2.AuthSubject{}, false
		}
	}

	if apiKey.User == nil {
		response.Unauthorized(c, "User associated with API key not found")
		return nil, middleware2.AuthSubject{}, false
	}
	if !apiKey.User.IsActive() {
		response.Unauthorized(c, "User account is not active")
		return nil, middleware2.AuthSubject{}, false
	}

	return apiKey, middleware2.AuthSubject{
		UserID:      apiKey.User.ID,
		Concurrency: apiKey.User.Concurrency,
	}, true
}

func parsePublicUsageDateRange(startDate, endDate, userTZ string) (time.Time, time.Time, error) {
	now := timezone.NowInUserLocation(userTZ)
	startTime := now.AddDate(0, 0, -30)
	endTime := now

	if strings.TrimSpace(startDate) != "" {
		t, err := timezone.ParseInUserLocation("2006-01-02", startDate, userTZ)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("Invalid start_date format, use YYYY-MM-DD")
		}
		startTime = t
	}

	if strings.TrimSpace(endDate) != "" {
		t, err := timezone.ParseInUserLocation("2006-01-02", endDate, userTZ)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("Invalid end_date format, use YYYY-MM-DD")
		}
		endTime = t.AddDate(0, 0, 1)
	}

	if !endTime.After(startTime) {
		return time.Time{}, time.Time{}, errors.New("end_date must be later than or equal to start_date")
	}

	return startTime, endTime, nil
}

func normalizePublicUsagePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = publicKeyUsageDefaultPageSize
	}
	if pageSize > publicKeyUsageMaxPageSize {
		pageSize = publicKeyUsageMaxPageSize
	}
	return page, pageSize
}

func (h *GatewayHandler) buildUsageResponsePayload(ctx context.Context, c *gin.Context, apiKey *service.APIKey, subject middleware2.AuthSubject, startTime, endTime time.Time) (gin.H, error) {
	usageData := h.buildUsageData(ctx, apiKey.ID)

	var modelStats any
	if h.usageService != nil {
		stats, err := h.usageService.GetAPIKeyModelStats(ctx, apiKey.ID, startTime, endTime)
		if err != nil {
			return nil, err
		}
		if len(stats) > 0 {
			modelStats = stats
		}
	}

	if apiKey.Quota > 0 || apiKey.HasRateLimits() {
		return h.buildUsageQuotaLimitedResponse(ctx, apiKey, usageData, modelStats), nil
	}

	return h.buildUsageUnrestrictedResponse(ctx, c, apiKey, subject, usageData, modelStats)
}

func (h *GatewayHandler) buildUsageQuotaLimitedResponse(ctx context.Context, apiKey *service.APIKey, usageData gin.H, modelStats any) gin.H {
	resp := gin.H{
		"mode":    "quota_limited",
		"isValid": apiKey.Status == service.StatusAPIKeyActive || apiKey.Status == service.StatusAPIKeyQuotaExhausted || apiKey.Status == service.StatusAPIKeyExpired,
		"status":  apiKey.Status,
	}

	if apiKey.Quota > 0 {
		remaining := apiKey.GetQuotaRemaining()
		resp["quota"] = gin.H{
			"limit":     apiKey.Quota,
			"used":      apiKey.QuotaUsed,
			"remaining": remaining,
			"unit":      "USD",
		}
		resp["remaining"] = remaining
		resp["unit"] = "USD"
	}

	if apiKey.HasRateLimits() && h.apiKeyService != nil {
		rateLimitData, err := h.apiKeyService.GetRateLimitData(ctx, apiKey.ID)
		if err == nil && rateLimitData != nil {
			rateLimits := make([]gin.H, 0, 3)
			if apiKey.RateLimit5h > 0 {
				used := rateLimitData.EffectiveUsage5h()
				entry := gin.H{
					"window":       "5h",
					"limit":        apiKey.RateLimit5h,
					"used":         used,
					"remaining":    max(0, apiKey.RateLimit5h-used),
					"window_start": rateLimitData.Window5hStart,
				}
				if rateLimitData.Window5hStart != nil && !service.IsWindowExpired(rateLimitData.Window5hStart, service.RateLimitWindow5h) {
					entry["reset_at"] = rateLimitData.Window5hStart.Add(service.RateLimitWindow5h)
				}
				rateLimits = append(rateLimits, entry)
			}
			if apiKey.RateLimit1d > 0 {
				used := rateLimitData.EffectiveUsage1d()
				entry := gin.H{
					"window":       "1d",
					"limit":        apiKey.RateLimit1d,
					"used":         used,
					"remaining":    max(0, apiKey.RateLimit1d-used),
					"window_start": rateLimitData.Window1dStart,
				}
				if rateLimitData.Window1dStart != nil && !service.IsWindowExpired(rateLimitData.Window1dStart, service.RateLimitWindow1d) {
					entry["reset_at"] = rateLimitData.Window1dStart.Add(service.RateLimitWindow1d)
				}
				rateLimits = append(rateLimits, entry)
			}
			if apiKey.RateLimit7d > 0 {
				used := rateLimitData.EffectiveUsage7d()
				entry := gin.H{
					"window":       "7d",
					"limit":        apiKey.RateLimit7d,
					"used":         used,
					"remaining":    max(0, apiKey.RateLimit7d-used),
					"window_start": rateLimitData.Window7dStart,
				}
				if rateLimitData.Window7dStart != nil && !service.IsWindowExpired(rateLimitData.Window7dStart, service.RateLimitWindow7d) {
					entry["reset_at"] = rateLimitData.Window7dStart.Add(service.RateLimitWindow7d)
				}
				rateLimits = append(rateLimits, entry)
			}
			if len(rateLimits) > 0 {
				resp["rate_limits"] = rateLimits
			}
		}
	}

	if apiKey.ExpiresAt != nil {
		resp["expires_at"] = apiKey.ExpiresAt
		resp["days_until_expiry"] = apiKey.GetDaysUntilExpiry()
	}

	if usageData != nil {
		resp["usage"] = usageData
	}
	if modelStats != nil {
		resp["model_stats"] = modelStats
	}

	return resp
}

func (h *GatewayHandler) buildUsageUnrestrictedResponse(ctx context.Context, c *gin.Context, apiKey *service.APIKey, subject middleware2.AuthSubject, usageData gin.H, modelStats any) (gin.H, error) {
	if apiKey.Group != nil && apiKey.Group.IsSubscriptionType() {
		resp := gin.H{
			"mode":     "unrestricted",
			"isValid":  true,
			"planName": apiKey.Group.Name,
			"unit":     "USD",
		}

		subscription, ok := middleware2.GetSubscriptionFromContext(c)
		if !ok && h.subscriptionService != nil {
			sub, err := h.subscriptionService.GetActiveSubscription(ctx, subject.UserID, apiKey.Group.ID)
			switch {
			case err == nil:
				subscription = sub
				ok = true
			case errors.Is(err, service.ErrSubscriptionNotFound):
			default:
				return nil, err
			}
		}

		if ok && subscription != nil {
			remaining := h.calculateSubscriptionRemaining(apiKey.Group, subscription)
			resp["remaining"] = remaining
			resp["subscription"] = gin.H{
				"daily_usage_usd":   subscription.DailyUsageUSD,
				"weekly_usage_usd":  subscription.WeeklyUsageUSD,
				"monthly_usage_usd": subscription.MonthlyUsageUSD,
				"daily_limit_usd":   apiKey.Group.DailyLimitUSD,
				"weekly_limit_usd":  apiKey.Group.WeeklyLimitUSD,
				"monthly_limit_usd": apiKey.Group.MonthlyLimitUSD,
				"expires_at":        subscription.ExpiresAt,
			}
		}

		if usageData != nil {
			resp["usage"] = usageData
		}
		if modelStats != nil {
			resp["model_stats"] = modelStats
		}
		return resp, nil
	}

	latestUser, err := h.userService.GetByID(ctx, subject.UserID)
	if err != nil {
		return nil, err
	}

	resp := gin.H{
		"mode":      "unrestricted",
		"isValid":   true,
		"planName":  "钱包余额",
		"remaining": latestUser.Balance,
		"unit":      "USD",
		"balance":   latestUser.Balance,
	}
	if usageData != nil {
		resp["usage"] = usageData
	}
	if modelStats != nil {
		resp["model_stats"] = modelStats
	}
	return resp, nil
}

func (h *GatewayHandler) buildPublicUsageLogsPage(ctx context.Context, apiKeyID, userID int64, startTime, endTime time.Time, page, pageSize int) (gin.H, error) {
	if h.usageService == nil {
		return gin.H{
			"items":     []publicUsageLogItem{},
			"total":     0,
			"page":      page,
			"page_size": pageSize,
			"pages":     1,
		}, nil
	}

	logs, pageResult, err := h.usageService.ListWithFilters(ctx, pagination.PaginationParams{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    "created_at",
		SortOrder: pagination.SortOrderDesc,
	}, usagestats.UsageLogFilters{
		UserID:    userID,
		APIKeyID:  apiKeyID,
		StartTime: &startTime,
		EndTime:   &endTime,
	})
	if err != nil {
		return nil, err
	}

	items := make([]publicUsageLogItem, 0, len(logs))
	for i := range logs {
		log := logs[i]
		requestType := log.EffectiveRequestType().String()
		items = append(items, publicUsageLogItem{
			ID:                  log.ID,
			RequestID:           log.RequestID,
			Model:               log.Model,
			RequestType:         requestType,
			ServiceTier:         log.ServiceTier,
			ReasoningEffort:     log.ReasoningEffort,
			InboundEndpoint:     log.InboundEndpoint,
			InputTokens:         log.InputTokens,
			OutputTokens:        log.OutputTokens,
			CacheCreationTokens: log.CacheCreationTokens,
			CacheReadTokens:     log.CacheReadTokens,
			TotalTokens:         log.TotalTokens(),
			ActualCost:          log.ActualCost,
			DurationMs:          log.DurationMs,
			CreatedAt:           log.CreatedAt.Format(time.RFC3339),
		})
	}

	pages := 1
	total := int64(0)
	if pageResult != nil {
		pages = pageResult.Pages
		total = pageResult.Total
		page = pageResult.Page
		pageSize = pageResult.PageSize
	}

	return gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"pages":     pages,
	}, nil
}
