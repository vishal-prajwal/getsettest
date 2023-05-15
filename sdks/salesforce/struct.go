package salesforce

type AccessTokenRequest struct {
	GrantType    string `json:"grant_type" validate:"required"`
	Assertion    string `json:"assertion" validate:"required"`
	ClientId     string `json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

type CreateTaskRequest struct {
	Status        string `json:"Status" validate:"required"`
	Subject       string `json:"Subject" validate:"required"`
	Priority      string `json:"Priority" validate:"required"`
	UserID        string `json:"UserId__c" validate:"required"`
	SocialNetwork string `json:"social_network__c" validate:"required"`
	AccessToken   string `json:",omitempty"`
	BaseURL       string `json:",omitempty"`
	Name          string `json:"name"`
}

type SaleForceCreateTaskHTTPRequest struct {
	UserID        string `json:"User_ID__c" validate:"required"`
	SocialNetwork string `json:"Social_Network__c" validate:"required"`
	AccessToken   string `json:",omitempty"`
	BaseURL       string `json:",omitempty"`
}

type CreateTaskResponse struct {
	ID      string   `json:"id"`
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}

type SalesforceSDKConfig struct {
	BaseURL    string
	Token      string
	APITimeout int
}
