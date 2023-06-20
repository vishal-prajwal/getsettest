package security

import (
	"strings"
)

func HasHierarchialPermissionForResource(authToken *AuthToken, resource string) (hasAccess bool) {

	if authToken == nil {
		return false
	}
	if authToken.Authorization == nil {
		return false
	}
	resources := strings.Split(resource, "\\.")
	if len(resources) > 1 {
		for _, v := range resources {
			if HasParentPermission(authToken, v) {
				return true
			}
		}
	}
	return HasParentPermission(authToken, resource)
}

func HasHierarchialPermissionForScope(authToken *AuthToken, resource string, scope string) (hasAccess bool) {
	if authToken == nil {
		return false
	}
	if authToken.Authorization == nil {
		return false
	}

	resources := strings.Split(resource, "\\.")
	if len(resources) > 1 {
		for _, v := range resources {
			if HasParentPermission(authToken, v) {
				return true
			}
		}
	}
	return HasParentPermission(authToken, resource) || HasSpecificPermission(authToken, resource, scope)
}

func HasParentPermission(authToken *AuthToken, resource string) (hasAccess bool) {
	if authToken == nil {
		return false
	}
	if authToken.Authorization == nil {
		return false
	}
	scopes, present := authToken.Authorization.getPermissionsMap()[resource]
	// resource permission should not have any scopes
	return present && len(scopes) == 0
}

func HasSpecificPermission(authToken *AuthToken, resource string, scope string) (hasAccess bool) {
	if authToken == nil {
		return false
	}
	if authToken.Authorization == nil {
		return false
	}
	scopes := authToken.Authorization.getPermissionsMap()[resource]
	// resource permission should not have any scopes
	for _, v := range scopes {
		if v == scope {
			return true
		}
	}
	return false
}

func HasGroupAccess(authToken *AuthToken, group string) (result bool) {
	if authToken == nil {
		return false
	}
	for _, v := range authToken.Groups {
		if v == group {
			return true
		}
	}
	return false
}
