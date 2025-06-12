package apilogger

import (
	"context"

	"bitbucket.org/junglee_games/getsetgo/eventqueue"
	"bitbucket.org/junglee_games/getsetgo/logger"
)

const (
	RUMMY_GAME_AUDIT_SOURCE  = "rummy_game_audit"
	KYC_API_USAGE_EVENT_TYPE = "kyc_api_usage"
	CONTEXT_USER_ID          = "X-User-Id"
	CONTEXT_REQUEST_ID       = "X-Request-Id"
)

const (
	UNKNOWN_USER_ID  = iota
	HZ_PRODUCT_ID
	OTHER_PRODUCT_ID
	RUMMY_PRODUCT_ID
)

type apiUsageLoggerImpl struct {
	eventPublisher eventqueue.Publisher
	Config         *Config
}

func (a *apiUsageLoggerImpl) Log(ctx context.Context, data *ApiData) error {

	if data == nil {
		logger.Error(ctx, "API usage data is nil")
		return nil
	}

	a.eventPublisher.PublishAsync(ctx, data.UserID, data)
	return nil
}

func (a *apiUsageLoggerImpl) GetConfig() *Config {
	return a.Config
}
