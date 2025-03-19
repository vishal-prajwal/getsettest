package kyc

import (
	"bitbucket.org/junglee_games/getsetgo/sdks/kyc/domain"
)

type KycService interface {
	FetchUserByPan(userByPanRequest domain.UserByPanRequest) (*domain.UserByPanResponse, error)
	FetchPanByUserID(userID int, productID string) (*domain.PanByUserResponse, error)
}
