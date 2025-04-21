package userdetails

type UserDetailsResponseFromNumber struct {
	Code     string                `json:"code"`
	Status   string                `json:"status"`
	Message  string                `json:"message"`
	UserData *UserDetailFromNumber `json:"data"`
}

type UserDetailFromNumber struct {
	UserId       int    `json:"userId"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile"`
	UserName     string `json:"userName"`
	Status       bool   `json:"status"`
}

type UserDetailsResponseFromUserId struct {
	Code     string                `json:"code"`
	Status   string                `json:"status"`
	Message  string                `json:"message"`
	UserData *UserDetailFromUserId `json:"data"`
}

type UserDetailFromUserId struct {
	MobileNumber string `json:"mobileNumber"`
}
