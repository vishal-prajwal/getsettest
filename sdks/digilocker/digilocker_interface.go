package digilocker

import (
	"context"
)

type Digilocker interface {
	GetRedirectURL(ctx context.Context) string
	StartKYC(ctx context.Context, transactionId, referenceId, redirectURL string) (*KYCStartDetails, error)
	CheckAccountstatus(ctx context.Context, mobile, aadhaar string) (*AccountStatusDetails, error)
	GetAddharDetails(ctx context.Context, transactionId, referenceId string) (*AadhaarDetails, error)
	Healthcheck(ctx context.Context) (*HealthcheckResult, error)
	GetPanDetails(ctx context.Context, refId string, panNumber string, fullName string) (*PanDetails, error)
	GetPanDigilockerDoc(ctx context.Context, refId string) (*PanDetails, error)
}
