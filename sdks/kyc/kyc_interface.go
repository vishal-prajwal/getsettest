package kyc

type KYC interface {
	FetchUserByPan(userByPanRequest UserByPanRequest) (UserByPanResponse, error)
}
