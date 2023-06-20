package security

import (
	"fmt"
	"strings"
)

type AuthToken struct {
	Exp               int64          `json:"exp"`
	Iat               int64          `json:"iat"`
	AuthTime          int64          `json:"auth_time"`
	Jti               string         `json:"jti"`
	Iss               string         `json:"iss"`
	Aud               string         `json:"aud"`
	Sub               string         `json:"sub"`
	Typ               string         `json:"typ"`
	Azp               string         `json:"azp"`
	Authorization     *Authorization `json:"authorization"`
	Scope             string         `json:"scope"`
	EmailVerified     bool           `json:"email_verified"`
	Name              string         `json:"name"`
	Groups            []string       `json:"groups"`
	PreferredUsername string         `json:"preferred_username"`
	GivenName         string         `json:"given_name"`
	FamilyName        string         `json:"family_name"`
	Email             string         `json:"email"`
}

type Authorization struct {
	Permissions    []Permission `json:"permissions"`
	PermissionsMap map[string][]string
}

type Permission struct {
	Scopes      []string `json:"scopes"`
	Rsid        string   `json:"rsid"`
	Rsname      string   `json:"rsname"`
	DisplayName string
}

func (permission *Permission) getDisplayName() (displayName string) {
	if permission.DisplayName != "" {
		return permission.DisplayName
	}
	displayNames := strings.Split(permission.Rsname, ".")
	fmt.Println("DisplayNames slice::::", displayNames)
	permission.DisplayName = displayNames[len(displayNames)-1]
	return permission.DisplayName
}

func (authorization *Authorization) getPermissionsMap() (permissionsMap map[string][]string) {
	if authorization.PermissionsMap == nil {
		authorization.PermissionsMap = make(map[string][]string)
	}
	for _, p := range authorization.Permissions {
		authorization.PermissionsMap[p.getDisplayName()] = p.Scopes
	}
	return authorization.PermissionsMap
}
