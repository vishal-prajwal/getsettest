package jwr

import "context"

const (
	GetUserProfilePath    = "/user"
	UpdateUserProfilePath = "/user"
)

type JWR interface {
	GetUserProfile(ctx context.Context, userID int, apiTimeOut int) (*UserProfile, error)
	FullUpdateProfile(ctx context.Context, userID int, userProfile UserProfile, apiTimeOut int) error
}

func New(config JWRSDKConfig) (JWR, error) {

	return &JWRImpl{
		BaseURL:           config.BaseURL,
		Token:             config.Token,
		DefaultAPITimeout: config.APITimeout,
	}, nil
}
