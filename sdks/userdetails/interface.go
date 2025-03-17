package userdetails

type UserDetailsSDK interface {
	GetUserDetailsFromMobileNumber(mobile string, platform Platform) (*UserDetailFromNumber, error)
	GetUserDetailsFromUserId(userId int64, platform Platform) (*UserDetailFromUserId, error)
}

type UserDetailsConfig interface {
	GetUserDetailsBaseURL() string
}
