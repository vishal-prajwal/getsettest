package domain

type UserByPanRequest struct {
	PanNumber  []string
	XProductID string
}

type UserByPanResponse struct {
	Error       string        `json:"error,omitempty"`
	UserPanInfo []UserPanInfo `json:"data,omitempty"`
}

type UserPanInfo struct {
	UserID int    `json:"userID"`
	PanNo  string `json:"panNo"`
}
