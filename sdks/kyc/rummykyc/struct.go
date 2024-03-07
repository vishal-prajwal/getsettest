package rummykyc

import "time"

type Kyc struct {
	UserID         int          `json:"userID,omitempty"`
	Medium         string       `json:"medium,omitempty"`
	ProductID      string       `json:"productID,omitempty"`
	CreatedAt      time.Time    `json:"createdAt,omitempty"`
	UpdatedAt      time.Time    `json:"updatedAt,omitempty"`
	IDProof        IDProof      `json:"idProof,omitempty"`
	AddressProof   AddressProof `json:"addressProof,omitempty"`
	PanProof       PanProof     `json:"panProof,omitempty"`
	PlatformName   string       `json:"PlatformName,omitempty"`
	DocumentType   string       `json:"DocumentType,omitempty"`
	IsSoftApproved bool         `json:"isSoftApproved,omitempty"`
}
type StatusReason struct {
	OcrIntegrity      string `json:"ocrIntegrity,omitempty"`
	ProfileStatus     string `json:"profileStatus,omitempty"`
	IsRestricted      bool   `json:"isRestricted,omitempty"`
	RestrictionReason string `json:"restrictionReason,omitempty"`
}
type IDProof struct {
	RecordID           string       `json:"recordID,omitempty"`
	DocumentID         string       `json:"documentID,omitempty"`
	DocType            string       `json:"docType,omitempty"`
	Status             string       `json:"status,omitempty"`
	StatusMessage      string       `json:"statusMessage,omitempty"`
	StatusReason       StatusReason `json:"statusReason,omitempty"`
	MismatchedFields   any          `json:"mismatchedFields,omitempty"`
	DbValidation       string       `json:"dbValidation,omitempty"`
	DbValidationReason string       `json:"dbValidationReason,omitempty"`
	CreatedAt          time.Time    `json:"createdAt,omitempty"`
	UpdatedAt          time.Time    `json:"updatedAt,omitempty"`
	ModifiedBy         string       `json:"modifiedBy,omitempty"`
	PlatformName       string       `json:"PlatformName,omitempty"`
	DuplicateKycUsers  any          `json:"duplicateKycUsers,omitempty"`
	ManualDeclineMsg   string       `json:"manualDeclineMsg,omitempty"`
}
type AddressProof struct {
	RecordID           string       `json:"recordID,omitempty"`
	DocumentID         string       `json:"documentID,omitempty"`
	DocType            string       `json:"docType,omitempty"`
	Status             string       `json:"status,omitempty"`
	StatusMessage      string       `json:"statusMessage,omitempty"`
	StatusReason       StatusReason `json:"statusReason,omitempty"`
	MismatchedFields   any          `json:"mismatchedFields,omitempty"`
	DbValidation       string       `json:"dbValidation,omitempty"`
	DbValidationReason string       `json:"dbValidationReason,omitempty"`
	CreatedAt          time.Time    `json:"createdAt,omitempty"`
	UpdatedAt          time.Time    `json:"updatedAt,omitempty"`
	ModifiedBy         string       `json:"modifiedBy,omitempty"`
	PlatformName       string       `json:"PlatformName,omitempty"`
	DuplicateKycUsers  any          `json:"duplicateKycUsers,omitempty"`
	ManualDeclineMsg   string       `json:"manualDeclineMsg,omitempty"`
}
type PanProof struct {
	RecordID           string       `json:"recordID,omitempty"`
	DocumentID         string       `json:"documentID,omitempty"`
	DocType            string       `json:"docType,omitempty"`
	Status             string       `json:"status,omitempty"`
	StatusMessage      string       `json:"statusMessage,omitempty"`
	StatusReason       StatusReason `json:"statusReason,omitempty"`
	MismatchedFields   any          `json:"mismatchedFields,omitempty"`
	DbValidation       string       `json:"dbValidation,omitempty"`
	DbValidationReason string       `json:"dbValidationReason,omitempty"`
	CreatedAt          time.Time    `json:"createdAt,omitempty"`
	UpdatedAt          time.Time    `json:"updatedAt,omitempty"`
	ModifiedBy         string       `json:"modifiedBy,omitempty"`
	PlatformName       string       `json:"PlatformName,omitempty"`
	DuplicateKycUsers  any          `json:"duplicateKycUsers,omitempty"`
	ManualDeclineMsg   string       `json:"manualDeclineMsg,omitempty"`
}
