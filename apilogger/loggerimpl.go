package apilogger

import (
	"context"
	"encoding/json"

	"bitbucket.org/junglee_games/getsetgo/eventqueue"
	"bitbucket.org/junglee_games/getsetgo/logger"
)

const (
	RUMMY_GAME_AUDIT_SOURCE  = "rummy_game_audit"
	KYC_API_USAGE_EVENT_TYPE = "kyc_api_usage"
	RUMMY_PRODUCT_ID         = 3 // Assuming a fixed product ID for Rummy
	CONTEXT_USER_ID          = "X-User-Id"
	CONTEXT_REQUEST_ID       = "X-Request-Id"
)

type apiUsageLoggerImpl struct {
	eventPublisher eventqueue.Publisher
	Config         Config
}

func (a *apiUsageLoggerImpl) Log(ctx context.Context, data *ApiData) error {

	if data == nil {
		logger.Error(ctx, "API usage data is nil")
		return nil
	}

	// Convert the ApiData to JSON
	message, err := json.Marshal(data)
	if err != nil {
		logger.Error(ctx, "Failed to marshal API usage data: %v", err)
		return err
	}

	a.eventPublisher.PublishAsync(ctx, data.UserID, message)
	return nil
}

func (a *apiUsageLoggerImpl) GetConfig() Config {
	return a.Config
}
