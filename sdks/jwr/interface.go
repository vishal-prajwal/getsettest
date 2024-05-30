package jwr

import (
	"context"

	"github.com/sony/gobreaker/v2"
)

const (
	GetUserProfilePath      = "/user"
	UpdateUserProfilePath   = "/user"
	UpdateUserProfilePathV2 = "/internal-api/user_profile/"
)

type JWR interface {
	GetUserProfile(ctx context.Context, userID int, apiTimeOut int) (*UserProfile, error)
	FullUpdateProfile(ctx context.Context, userID int, userProfile UserProfile, apiTimeOut int) error
	FullUpdateProfileV2(ctx context.Context, userID int, userProfile UserProfile, apiTimeOut int, retries int, cb *gobreaker.CircuitBreaker[[]byte]) error
}

func New(config JWRSDKConfig) (JWR, error) {

	return &JWRImpl{
		BaseURL:           config.BaseURL,
		Token:             config.Token,
		DefaultAPITimeout: config.APITimeout,
	}, nil
}
