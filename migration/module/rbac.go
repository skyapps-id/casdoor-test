package rbac

import (
	"fmt"
	"log"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

// CasdoorConfig holds Casdoor configuration
type CasdoorConfig struct {
	Endpoint         string
	ClientID         string
	ClientSecret     string
	Certificate      string
	OrganizationName string
	ApplicationName  string
}

// CasdoorMigration handles Casdoor RBAC setup
type CasdoorMigration struct {
	client *casdoorsdk.Client
	config CasdoorConfig
}

// NewCasdoorMigration creates a new migration instance
func NewCasdoorMigration(config CasdoorConfig) (*CasdoorMigration, error) {
	// Initialize Casdoor SDK client
	casdoorsdk.InitConfig(
		config.Endpoint,
		config.ClientID,
		config.ClientSecret,
		config.Certificate,
		config.OrganizationName,
		config.ApplicationName,
	)

	return &CasdoorMigration{
		config: config,
	}, nil
}

// MigrateRoles creates default roles for RBAC
func (m *CasdoorMigration) MigrateRoles() error {
	log.Println("Starting role migration...")

	roles := []casdoorsdk.Role{
		{
			Owner:       m.config.OrganizationName,
			Name:        "admin",
			CreatedTime: time.Now().Format("2006-01-02 15:04:05.000"),
			DisplayName: "Administrator",
			Description: "Full system access with all permissions",
			Users:       []string{},
			Roles:       []string{},
			Domains:     []string{},
			IsEnabled:   true,
		},
		{
			Owner:       m.config.OrganizationName,
			Name:        "manager",
			CreatedTime: time.Now().Format("2006-01-02 15:04:05.000"),
			DisplayName: "Manager",
			Description: "Manage users and content",
			Users:       []string{},
			Roles:       []string{},
			Domains:     []string{},
			IsEnabled:   true,
		},
		{
			Owner:       m.config.OrganizationName,
			Name:        "user",
			CreatedTime: time.Now().Format("2006-01-02 15:04:05.000"),
			DisplayName: "Regular User",
			Description: "Basic user access",
			Users:       []string{},
			Roles:       []string{},
			Domains:     []string{},
			IsEnabled:   true,
		},
	}

	for _, role := range roles {
		// Check if role exists
		existingRole, err := casdoorsdk.GetRole(role.Name)
		if err != nil {
			log.Printf("Error checking role %s: %v", role.Name, err)
		}

		if existingRole == nil || existingRole.Name == "" {
			// Create new role
			affected, err := casdoorsdk.AddRole(&role)
			if err != nil {
				return fmt.Errorf("failed to create role %s: %v", role.Name, err)
			}
			log.Printf("Created role: %s (affected: %v)", role.Name, affected)
		} else {
			log.Printf("Role already exists: %s", role.Name)
		}
	}

	log.Println("Role migration completed successfully")
	return nil
}

// MigratePermissions creates default permissions for RBAC
func (m *CasdoorMigration) MigratePermissions() error {
	log.Println("Starting permission migration...")

	permissions := []casdoorsdk.Permission{
		{
			Owner:        m.config.OrganizationName,
			Name:         "user-read",
			CreatedTime:  time.Now().Format("2006-01-02 15:04:05.000"),
			DisplayName:  "Read Users",
			Description:  "Permission to read user data",
			Users:        []string{},
			Roles:        []string{m.config.OrganizationName + "/admin", m.config.OrganizationName + "/manager", m.config.OrganizationName + "/user"},
			Domains:      []string{},
			Model:        "",
			ResourceType: "Custom",
			Resources:    []string{"users"},
			Actions:      []string{"read"},
			Effect:       "Allow",
			IsEnabled:    true,
			ApproveTime:  time.Now().Format("2006-01-02 15:04:05.000"),
		},
		{
			Owner:        m.config.OrganizationName,
			Name:         "user-write",
			CreatedTime:  time.Now().Format("2006-01-02 15:04:05.000"),
			DisplayName:  "Write Users",
			Description:  "Permission to create and update users",
			Users:        []string{},
			Roles:        []string{m.config.OrganizationName + "/admin", m.config.OrganizationName + "/manager"},
			Domains:      []string{},
			Model:        "",
			ResourceType: "Custom",
			Resources:    []string{"users"},
			Actions:      []string{"write"},
			Effect:       "Allow",
			IsEnabled:    true,
			ApproveTime:  time.Now().Format("2006-01-02 15:04:05.000"),
		},
		{
			Owner:        m.config.OrganizationName,
			Name:         "user-delete",
			CreatedTime:  time.Now().Format("2006-01-02 15:04:05.000"),
			DisplayName:  "Delete Users",
			Description:  "Permission to delete users",
			Users:        []string{},
			Roles:        []string{m.config.OrganizationName + "/admin"},
			Domains:      []string{},
			Model:        "",
			ResourceType: "Custom",
			Resources:    []string{"users"},
			Actions:      []string{"delete"},
			Effect:       "Allow",
			IsEnabled:    true,
			ApproveTime:  time.Now().Format("2006-01-02 15:04:05.000"),
		},
		{
			Owner:        m.config.OrganizationName,
			Name:         "role-manage",
			CreatedTime:  time.Now().Format("2006-01-02 15:04:05.000"),
			DisplayName:  "Manage Roles",
			Description:  "Permission to manage roles and permissions",
			Users:        []string{},
			Roles:        []string{m.config.OrganizationName + "/admin"},
			Domains:      []string{},
			Model:        "",
			ResourceType: "Custom",
			Resources:    []string{"roles", "permissions"},
			Actions:      []string{"read", "write", "delete"},
			Effect:       "Allow",
			IsEnabled:    true,
			ApproveTime:  time.Now().Format("2006-01-02 15:04:05.000"),
		},
	}

	for _, perm := range permissions {
		// Check if permission exists
		existingPerm, err := casdoorsdk.GetPermission(perm.Name)
		if err != nil {
			log.Printf("Error checking permission %s: %v", perm.Name, err)
		}

		if existingPerm == nil || existingPerm.Name == "" {
			// Create new permission
			affected, err := casdoorsdk.AddPermission(&perm)
			if err != nil {
				return fmt.Errorf("failed to create permission %s: %v", perm.Name, err)
			}
			log.Printf("Created permission: %s (affected: %v)", perm.Name, affected)
		} else {
			log.Printf("Permission already exists: %s", perm.Name)
		}
	}

	log.Println("Permission migration completed successfully")
	return nil
}

// MigrateModel creates Casbin model for RBAC
func (m *CasdoorMigration) MigrateModel() error {
	log.Println("Starting model migration...")

	// Casbin RBAC model definition
	modelText := `[request_definition]
r = subOwner, subName, method, urlPath, objOwner, objName

[policy_definition]
p = subOwner, subName, method, urlPath, objOwner, objName

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (r.subOwner == p.subOwner || p.subOwner == "*") && \
    (r.subName == p.subName || p.subName == "*" || r.subName != "anonymous" && p.subName == "!anonymous") && \
    (r.method == p.method || p.method == "*") && \
    (r.urlPath == p.urlPath || p.urlPath == "*") && \
    (r.objOwner == p.objOwner || p.objOwner == "*") && \
    (r.objName == p.objName || p.objName == "*") || \
    (r.subOwner == r.objOwner && r.subName == r.objName)
`

	model := casdoorsdk.Model{
		Owner:       m.config.OrganizationName,
		Name:        "rbac-model",
		CreatedTime: time.Now().Format("2006-01-02 15:04:05.000"),
		DisplayName: "RBAC Model",
		ModelText:   modelText,
	}

	// Check if model exists
	existingModel, err := casdoorsdk.GetModel(model.Name)
	if err != nil {
		log.Printf("Error checking model: %v", err)
	}

	if existingModel == nil || existingModel.Name == "" {
		// Create new model
		affected, err := casdoorsdk.AddModel(&model)
		if err != nil {
			return fmt.Errorf("failed to create model: %v", err)
		}
		log.Printf("Created model: %s (affected: %v)", model.Name, affected)
	} else {
		log.Printf("Model already exists: %s", model.Name)
		model.Key = existingModel.Key // must set ID to update!
		affected, err := casdoorsdk.UpdateModel(&model)
		if err != nil {
			return fmt.Errorf("failed to update model: %v", err)
		}
		log.Printf("Updated existing RBAC model: %s (affected: %v)", model.Name, affected)
	}

	log.Println("Model migration completed successfully")
	return nil
}

// MigrateAdapter creates adapter for enforcer
func (m *CasdoorMigration) MigrateAdapter() error {
	log.Println("Starting adapter migration...")

	adapter := casdoorsdk.Adapter{
		Owner:       m.config.OrganizationName,
		Name:        "rbac-adapter",
		CreatedTime: time.Now().Format("2006-01-02 15:04:05.000"),
		Table:       "casbin_rule",
		UseSameDb:   true,
	}

	// Check if adapter exists
	existingAdapter, err := casdoorsdk.GetAdapter(adapter.Name)
	if err != nil {
		log.Printf("Error checking adapter: %v", err)
	}

	if existingAdapter == nil || existingAdapter.Name == "" {
		// Create new adapter
		affected, err := casdoorsdk.AddAdapter(&adapter)
		if err != nil {
			return fmt.Errorf("failed to create adapter: %v", err)
		}
		log.Printf("Created adapter: %s (affected: %v)", adapter.Name, affected)
	} else {
		log.Printf("Adapter already exists: %s", adapter.Name)
	}

	log.Println("Adapter migration completed successfully")
	return nil
}

// MigrateEnforcer creates enforcer for RBAC
func (m *CasdoorMigration) MigrateEnforcer() error {
	log.Println("Starting enforcer migration...")

	enforcer := casdoorsdk.Enforcer{
		Owner:       m.config.OrganizationName,
		Name:        "rbac-enforcer",
		CreatedTime: time.Now().Format("2006-01-02 15:04:05.000"),
		DisplayName: "RBAC Enforcer",
		Description: "Main enforcer for RBAC system",
		Model:       m.config.OrganizationName + "/rbac-model",
		Adapter:     m.config.OrganizationName + "/rbac-adapter",
		IsEnabled:   true,
	}

	// Check if enforcer exists
	existingEnforcer, err := casdoorsdk.GetEnforcer(enforcer.Name)
	if err != nil {
		log.Printf("Error checking enforcer: %v", err)
	}

	if existingEnforcer == nil || existingEnforcer.Name == "" {
		// Create new enforcer
		affected, err := casdoorsdk.AddEnforcer(&enforcer)
		if err != nil {
			return fmt.Errorf("failed to create enforcer: %v", err)
		}
		log.Printf("Created enforcer: %s (affected: %v)", enforcer.Name, affected)
	} else {
		log.Printf("Enforcer already exists: %s", enforcer.Name)
	}

	log.Println("Enforcer migration completed successfully")
	return nil
}

// MigratePolicies creates Casbin policies mapped to URLs
func (m *CasdoorMigration) MigratePolicies() error {
	log.Println("Starting policy migration...")

	type PolicyDef struct {
		Role     string
		Resource string
		Action   string
	}

	policies := []PolicyDef{
		// USERS permissions
		{"admin", "/api/users", "GET"},
		{"admin", "/api/users", "POST"},
		{"admin", "/api/users/*", "GET"},
		{"admin", "/api/users/*", "PUT"},
		{"admin", "/api/users/*", "DELETE"},

		{"manager", "/api/users", "GET"},
		{"manager", "/api/users/*", "PUT"},

		{"user", "/api/users", "GET"}, // user boleh list profiles

		// PRODUCTS permissions
		{"admin", "/api/products", "GET"},
		{"admin", "/api/products", "POST"},
		{"admin", "/api/products/*", "GET"},
		{"admin", "/api/products/*", "PUT"},
		{"admin", "/api/products/*", "DELETE"},

		{"manager", "/api/products", "GET"},
		{"manager", "/api/products", "POST"},
		{"manager", "/api/products/*", "GET"},
		{"manager", "/api/products/*", "PUT"},
		{"manager", "/api/products/*", "DELETE"},

		{"user", "/api/products", "GET"},
		{"user", "/api/products/*", "GET"},
	}

	// Fetch enforcer
	enforcer, err := casdoorsdk.GetEnforcer("rbac-enforcer")
	if err != nil {
		return fmt.Errorf("failed to get enforcer: %v", err)
	}

	owner := m.config.OrganizationName

	for _, policy := range policies {
		rule := &casdoorsdk.CasbinRule{
			Ptype: "p",
			V0:    "web-apps",      // subOwner
			V1:    policy.Role,     // subName
			V2:    policy.Action,   // method
			V3:    policy.Resource, // urlPath
			V4:    owner,           // objOwner
			V5:    "*",             // objName
		}

		affected, err := casdoorsdk.AddPolicy(enforcer, rule)
		if err != nil {
			log.Printf("Failed: %s %s %s → %v",
				policy.Role, policy.Resource, policy.Action, err)
			continue
		}

		if affected {
			log.Printf("Added policy: %s %s %s",
				policy.Role, policy.Resource, policy.Action)
		}
	}

	log.Println("Policy migration completed successfully")
	return nil
}

// Run executes all migrations
func (m *CasdoorMigration) Run() error {
	log.Println("=== Starting Casdoor RBAC Migration ===")

	// Step 1: Create Casbin Model
	if err := m.MigrateModel(); err != nil {
		return fmt.Errorf("model migration failed: %v", err)
	}

	// Step 2: Create Adapter
	if err := m.MigrateAdapter(); err != nil {
		return fmt.Errorf("adapter migration failed: %v", err)
	}

	// Step 3: Create Roles
	if err := m.MigrateRoles(); err != nil {
		return fmt.Errorf("role migration failed: %v", err)
	}

	// Step 4: Create Permissions
	// if err := m.MigratePermissions(); err != nil {
	// 	return fmt.Errorf("permission migration failed: %v", err)
	// }

	// Step 5: Create Enforcer
	if err := m.MigrateEnforcer(); err != nil {
		return fmt.Errorf("enforcer migration failed: %v", err)
	}

	// Step 6: Create Policies (URL & Method mappings)
	if err := m.MigratePolicies(); err != nil {
		return fmt.Errorf("policy migration failed: %v", err)
	}

	log.Println("=== Casdoor RBAC Migration Completed Successfully ===")
	return nil
}

// Rollback removes all created resources (use with caution)
func (m *CasdoorMigration) Rollback() error {
	log.Println("=== Starting Casdoor RBAC Rollback ===")

	// Delete policies first
	log.Println("Deleting policies...")
	enforcerObj, err := casdoorsdk.GetEnforcer("rbac-enforcer")
	if err != nil {
		log.Printf("Failed to fetch enforcer: %v", err)
		return err
	}
	roles := []string{"admin", "manager", "user"}
	for _, role := range roles {
		roleID := m.config.OrganizationName + "/" + role

		rule := &casdoorsdk.CasbinRule{
			Ptype: "p",
			V0:    roleID, // subject (role)
			// Kosongkan V1 dan V2 → remove semua policies untuk subject ini
		}
		affected, err := casdoorsdk.RemovePolicy(enforcerObj, rule)
		if err != nil {
			log.Printf("Error removing policies for role %s: %v", role, err)
		} else if affected {
			log.Printf("Removed policies for role: %s", role)
		} else {
			log.Printf("No policies found for role: %s", role)
		}
	}

	// Delete permissions
	permissions := []string{"user-read", "user-write", "user-delete", "role-manage"}
	for _, name := range permissions {
		affected, err := casdoorsdk.DeletePermission(&casdoorsdk.Permission{
			Owner: m.config.OrganizationName,
			Name:  name,
		})
		if err != nil {
			log.Printf("Error deleting permission %s: %v", name, err)
		} else {
			log.Printf("Deleted permission: %s (affected: %v)", name, affected)
		}
	}

	// Delete roles
	for _, name := range roles {
		affected, err := casdoorsdk.DeleteRole(&casdoorsdk.Role{
			Owner: m.config.OrganizationName,
			Name:  name,
		})
		if err != nil {
			log.Printf("Error deleting role %s: %v", name, err)
		} else {
			log.Printf("Deleted role: %s (affected: %v)", name, affected)
		}
	}

	// Delete enforcer
	affected, err := casdoorsdk.DeleteEnforcer(&casdoorsdk.Enforcer{
		Owner: m.config.OrganizationName,
		Name:  "rbac-enforcer",
	})
	if err != nil {
		log.Printf("Error deleting enforcer: %v", err)
	} else {
		log.Printf("Deleted enforcer (affected: %v)", affected)
	}

	// Delete adapter
	affected, err = casdoorsdk.DeleteAdapter(&casdoorsdk.Adapter{
		Owner: m.config.OrganizationName,
		Name:  "rbac-adapter",
	})
	if err != nil {
		log.Printf("Error deleting adapter: %v", err)
	} else {
		log.Printf("Deleted adapter (affected: %v)", affected)
	}

	// Delete model
	affected, err = casdoorsdk.DeleteModel(&casdoorsdk.Model{
		Owner: m.config.OrganizationName,
		Name:  "rbac-model",
	})
	if err != nil {
		log.Printf("Error deleting model: %v", err)
	} else {
		log.Printf("Deleted model (affected: %v)", affected)
	}

	log.Println("=== Casdoor RBAC Rollback Completed ===")
	return nil
}
