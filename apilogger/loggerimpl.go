package apilogger

import (
	"context"
	"encoding/json"

	"bitbucket.org/junglee_games/getsetgo/eventqueue"
	"bitbucket.org/junglee_games/getsetgo/logger"
)

const (
	RUMMY_GAME_AUDIT         = "rummy_game_audit"
	KYC_API_USAGE_EVENT_TYPE = "kyc_api_usage"
)

type apiUsageLoggerImpl struct {
	eventPublisher eventqueue.Publisher
}

type Kye struct {
	UserID int64 `json:"userId"`
}

func (a *apiUsageLoggerImpl) Log(ctx context.Context, data *ApiData) error {

	if data == nil {
		logger.Error(ctx, "API usage data is nil")
		return nil
	}

	// Create a key for the event
	key := &Kye{
		UserID: data.UserID,
	}

	// Convert the ApiData to JSON
	message, err := json.Marshal(data)
	if err != nil {
		logger.Error(ctx, "Failed to marshal API usage data: %v", err)
		return err
	}

	a.eventPublisher.PublishAsync(ctx, key, message)
	return nil
}
