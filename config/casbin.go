package config

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

var Enforcer *casbin.Enforcer

func InitCasbin() error {
	// Model definition (RBAC with domains/tenants)
	modelText := `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
`

	m, err := model.NewModelFromString(modelText)
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	// File adapter untuk menyimpan policies
	adapter := fileadapter.NewAdapter("rbac_policy.csv")

	Enforcer, err = casbin.NewEnforcer(m, adapter)
	if err != nil {
		return fmt.Errorf("failed to create enforcer: %w", err)
	}

	// Load policy dari file
	if err := Enforcer.LoadPolicy(); err != nil {
		log.Printf("Warning: failed to load policy: %v", err)
	}

	return nil
}

// Sync RBAC dari Casdoor ke Casbin
func SyncRBACFromCasdoor() error {
	// Hapus semua policy yang ada
	Enforcer.ClearPolicy()

	// 1. Get semua roles dari Casdoor
	roles, err := CasdoorClient.GetRoles()
	if err != nil {
		return fmt.Errorf("failed to get roles: %w", err)
	}

	// 2. Get semua users dari Casdoor
	users, err := CasdoorClient.GetUsers()
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	// 3. Definisikan permissions untuk setiap role
	// Format: role, domain, resource, action
	rolePermissions := map[string][][]string{
		"admin": {
			{"admin", "skyapps", "users", "read"},
			{"admin", "skyapps", "users", "write"},
			{"admin", "skyapps", "users", "delete"},
			{"admin", "skyapps", "roles", "read"},
			{"admin", "skyapps", "roles", "write"},
			{"admin", "skyapps", "roles", "delete"},
			{"admin", "skyapps", "rbac", "write"},
		},
		"manager": {
			{"manager", "skyapps", "users", "read"},
			{"manager", "skyapps", "users", "write"},
			{"manager", "skyapps", "roles", "read"},
		},
		"user": {
			{"user", "skyapps", "users", "read"},
		},
	}

	// 4. Tambahkan permissions untuk setiap role
	for _, role := range roles {
		if role.Owner != "skyapps" {
			continue
		}

		permissions, exists := rolePermissions[role.Name]
		if !exists {
			continue
		}

		for _, perm := range permissions {
			if _, err := Enforcer.AddPolicy(perm); err != nil {
				log.Printf("Failed to add policy for role %s: %v", role.Name, err)
			}
		}
	}

	// 5. Assign roles ke users
	for _, user := range users {
		if user.Owner != "skyapps" {
			continue
		}

		for _, role := range user.Roles {
			// Format: user, role, domain
			if _, err := Enforcer.AddGroupingPolicy(user.Name, role.Name, "skyapps"); err != nil {
				log.Printf("Failed to add role %s to user %s: %v", role.Name, user.Name, err)
			}
		}
	}

	// 6. Save policy ke file
	if err := Enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save policy: %w", err)
	}

	log.Println("RBAC synced successfully from Casdoor")
	return nil
}

// Check permission
func CheckPermission(username, resource, action string) (bool, error) {
	return Enforcer.Enforce(username, "skyapps", resource, action)
}
