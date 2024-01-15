package howzat

type Howzat interface {
	FetchUserByPan(userByPanRequest UserByPanRequest) (UserByPanResponse, error)
}
