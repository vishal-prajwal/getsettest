package howzatkyc

type KycResponse struct {
	Code    string `json:"code"`
	Status  string `json:"status"`
	Data    Data   `json:"data"`
	Message string `json:"message,omitempty"`
}

type Data struct {
	UserID               int     `json:"userId,omitempty"`
	KycType              string  `json:"kycType,omitempty"`
	Attempts             int     `json:"attempts,omitempty"`
	DocumentType         string  `json:"documentType,omitempty"`
	DocumentNumber       string  `json:"documentNumber,omitempty"`
	ReferenceID          string  `json:"referenceId,omitempty"`
	Status               string  `json:"status,omitempty"`
	Remark               string  `json:"remark,omitempty"`
	CreatedAt            float64 `json:"createdAt,omitempty"`
	UpdatedAt            float64 `json:"updatedAt,omitempty"`
	LastUpdatedChannelID int     `json:"lastUpdatedChannelId,omitempty"`
	ChannelID            int     `json:"channelId,omitempty"`
	FirstName            string  `json:"firstName,omitempty"`
	LastName             string  `json:"lastName,omitempty"`
	Dob                  []int   `json:"dob,omitempty"`
	Pincode              int     `json:"pincode,omitempty"`
	Source               string  `json:"source,omitempty"`
	KycLiteBlocked       bool    `json:"kycLiteBlocked,omitempty"`
}

type PanByUserResponse struct {
	Code   string `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
	Data   struct {
		UserID               int     `json:"userId,omitempty"`
		KycType              string  `json:"kycType,omitempty"`
		Attempts             int     `json:"attempts,omitempty"`
		DocumentType         string  `json:"documentType,omitempty"`
		DocumentNumber       string  `json:"documentNumber,omitempty"`
		ReferenceID          string  `json:"referenceId,omitempty"`
		Status               string  `json:"status,omitempty"`
		Remark               string  `json:"remark,omitempty"`
		CreatedAt            float64 `json:"createdAt,omitempty"`
		UpdatedAt            float64 `json:"updatedAt,omitempty"`
		LastUpdatedChannelID int     `json:"lastUpdatedChannelId,omitempty"`
		ChannelID            int     `json:"channelId,omitempty"`
		Source               string  `json:"source,omitempty"`
		KycLiteBlocked       bool    `json:"kycLiteBlocked,omitempty"`
	} `json:"data,omitempty"`
}