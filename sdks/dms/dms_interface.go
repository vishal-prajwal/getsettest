package dms

import dmsService "bitbucket.org/junglee_games/japi-dms-service/models/dto"

type DMS interface {
	Initiate(req IntiateRequest) (IntiateResponse, error)
	SaveKycInfo(req dmsService.KycInfo) (*dmsService.KycInfoResponse, error)
}
