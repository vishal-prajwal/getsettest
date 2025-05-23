package jwr

import (
	"strings"
	"time"

	"bitbucket.org/junglee_games/getsetgo/circuit_breaker/gobreaker"
)

type UserProfile struct {
	FirstName     string `json:"firstName"`
	MiddleName    string `json:"middleName"`
	LastName      string `json:"lastName"`
	DOB           int64  `json:"dateOfBirth"`
	StringDOB     string `json:"stringDateOfBirth"`
	Pin           string `json:"pin"`
	Address       string `json:"address"`
	Address2      string `json:"address2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Gender        string `json:"gender"`
	Updatedby     string `json:"updatedBy"`
	IsRDCSyncUser bool   `json:"rdcsyncUser"`
}

type UpdateUserProfileRequest struct {
	FirstName     *string `json:"firstName,omitempty"`
	MiddleName    *string `json:"middleName,omitempty"`
	LastName      *string `json:"lastName,omitempty"`
	DOB           *int64  `json:"dateOfBirth,omitempty"`
	Pin           *string `json:"pin,omitempty"`
	Address       *string `json:"address,omitempty"`
	Address2      *string `json:"address2,omitempty"`
	City          *string `json:"city,omitempty"`
	State         *string `json:"state,omitempty"`
	Gender        *string `json:"gender,omitempty"`
	Updatedby     *string `json:"updatedBy,omitempty"`
	IsRDCSyncUser *bool   `json:"rdcsyncUser,omitempty"`
}

type JWRSDKConfig struct {
	BaseURL      string
	InternalURL  string
	Token        string
	APITimeout   int
	GobreakerCfg *gobreaker.GobreakerCfg
}

func (u UserProfile) GetFullName() string {
	if strings.ToLower(u.FirstName) == "fnu" {
		u.FirstName = ""
	}
	if strings.ToLower(u.LastName) == "lnu" {
		u.LastName = ""
	}
	if u.MiddleName == "" {
		return strings.TrimSpace(u.FirstName + " " + u.LastName)
	}
	return strings.TrimSpace(u.FirstName + " " + u.MiddleName + " " + u.LastName)
}

func (u UserProfile) GetDOB() (string, error) {
	t := time.Unix(0, int64(u.DOB)*int64(time.Millisecond))

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return "", err
	}

	t = t.In(loc)
	s := t.Format("02/01/2006")
	return s, nil
}

func (u UserProfile) GetHyphenDOB() (string, error) {
	t := time.Unix(0, int64(u.DOB)*int64(time.Millisecond))

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return "", err
	}

	t = t.In(loc)
	s := t.Format("02-01-2006")
	return s, nil
}
