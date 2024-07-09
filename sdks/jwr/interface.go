package jwr

import (
	"context"

	gb "bitbucket.org/junglee_games/getsetgo/circuit_breaker/gobreaker"
	"github.com/sony/gobreaker/v2"
)

const (
	GetUserProfilePath      = "/user"
	UpdateUserProfilePath   = "/user"
	UpdateUserProfilePathV2 = "/user_profile/"
)

type JWR interface {
	GetUserProfile(ctx context.Context, userID int, apiTimeOut int) (*UserProfile, error)
	FullUpdateProfile(ctx context.Context, userID int, userProfile UserProfile, apiTimeOut int) error
	FullUpdateProfileV2(ctx context.Context, userID int, userProfile UserProfile, apiTimeOut int, retries int) error
}

func New(config JWRSDKConfig) (JWR, error) {

	var cb *gobreaker.CircuitBreaker[[]byte]
	var err error
	if config.GobreakerCfg != nil &&
		config.GobreakerCfg.Enabled {
		cb, err = gb.GetCircutBreaker(config.GobreakerCfg)
		if err != nil {
			return nil, err
		}
	}

	return &JWRImpl{
		BaseURL:           config.BaseURL,
		Token:             config.Token,
		InternalURL:       config.InternalURL,
		DefaultAPITimeout: config.APITimeout,
		cb:                cb,
	}, nil
}
