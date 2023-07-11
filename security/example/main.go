package main

import (
	"fmt"

	"bitbucket.org/junglee_games/getsetgo/security"
)

func main() {
	permissions := make([]security.Permission, 0)
	scopes := make([]string, 1)
	scopes[0] = "test"
	permissions = append(permissions, security.Permission{
		Scopes: scopes,
		Rsname: "All.Hierarchy1.Hierarchy2",
		Rsid:   "Test",
	})
	token := security.AuthToken{
		Authorization: &security.Authorization{
			Permissions: permissions,
		},
	}
	result := security.HasHierarchialPermissionForResource(&token, "All.Hierarchy1.Hierarchy2")
	fmt.Printf("Result HasHierarchialPermissionForResource: %v", result)
	result = security.HasHierarchialPermissionForScope(&token, "All.Hierarchy1.Hierarchy2", "test")
	fmt.Printf("Result HasHierarchialPermissionForScope: %v", result)

}
