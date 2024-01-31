package rummykyc

type Kyc struct {
	UserID       int          `json:"userID"`
	ProductID    string       `json:"productID"`
	AddressProof AddressProof `json:"addressProof"`
	PanProof     PanProof     `json:"panProof"`
	Profile      Profile      `json:"profile"`
}
type AddressProof struct {
	DocumentID string `json:"documentID"`
	DocType    string `json:"docType"`
	Status     string `json:"status"`
	CreatedAt  int    `json:"createdAt"`
	UpdatedAt  int    `json:"updatedAt"`
	ModifiedBy string `json:"modifiedBy"`
}
type PanProof struct {
	DocumentID string `json:"documentID"`
	DocType    string `json:"docType"`
	Status     string `json:"status"`
	CreatedAt  int    `json:"createdAt"`
	UpdatedAt  int    `json:"updatedAt"`
	ModifiedBy string `json:"modifiedBy"`
}
type Profile struct {
	FirstName   string `json:"firstName"`
	MiddleName  string `json:"middleName"`
	LastName    string `json:"lastName"`
	DateOfBirth int64  `json:"dateOfBirth"`
	Pin         string `json:"pin"`
	Address     string `json:"address"`
	Address2    string `json:"address2"`
	City        string `json:"city"`
	State       string `json:"state"`
	Gender      string `json:"gender"`
}
