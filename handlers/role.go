package handlers

import (
	"net/http"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/labstack/echo/v4"
	"github.com/skyapps-id/casdoor-test/config"
)

func ListRoles(c echo.Context) error {
	roles, err := config.CasdoorClient.GetRoles()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get roles",
		})
	}

	skyappsRoles := []interface{}{}
	for _, role := range roles {
		if role.Owner == "skyapps" {
			skyappsRoles = append(skyappsRoles, map[string]interface{}{
				"name":         role.Name,
				"display_name": role.DisplayName,
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"roles": skyappsRoles,
	})
}

func AddRole(c echo.Context) error {
	var req struct {
		Name        string `json:"name" validate:"required"`
		DisplayName string `json:"display_name" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	role := &casdoorsdk.Role{
		Owner:       "skyapps",
		Name:        req.Name,
		DisplayName: req.DisplayName,
	}

	affected, err := config.CasdoorClient.AddRole(role)
	if err != nil || !affected {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create role",
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Role created successfully",
	})
}

func UpdateRole(c echo.Context) error {
	roleName := c.Param("role")

	var req struct {
		DisplayName string `json:"display_name"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	role := &casdoorsdk.Role{
		Owner:       "skyapps",
		Name:        roleName,
		DisplayName: req.DisplayName,
	}

	affected, err := config.CasdoorClient.UpdateRole(role)
	if err != nil || !affected {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update role",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Role updated successfully",
	})
}

func DeleteRole(c echo.Context) error {
	roleName := c.Param("role")

	role := &casdoorsdk.Role{
		Owner: "skyapps",
		Name:  roleName,
	}

	affected, err := config.CasdoorClient.DeleteRole(role)
	if err != nil || !affected {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete role",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Role deleted successfully",
	})
}

func AssignRole(c echo.Context) error {
	username := c.Param("username")

	var req struct {
		Role string `json:"role" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	// Get user
	user, err := config.CasdoorClient.GetUser(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	// Add role to user
	user.Roles = append(user.Roles, &casdoorsdk.Role{
		Owner: "skyapps",
		Name:  req.Role,
	})

	affected, err := config.CasdoorClient.UpdateUser(user)
	if err != nil || !affected {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to assign role",
		})
	}

	// Sync RBAC
	config.SyncRBACFromCasdoor()

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Role assigned successfully",
	})
}

func RemoveRole(c echo.Context) error {
	username := c.Param("username")
	roleName := c.Param("role")

	user, err := config.CasdoorClient.GetUser(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	// Remove role
	newRoles := []*casdoorsdk.Role{}
	for _, role := range user.Roles {
		if role.Name != roleName {
			newRoles = append(newRoles, role)
		}
	}
	user.Roles = newRoles

	affected, err := config.CasdoorClient.UpdateUser(user)
	if err != nil || !affected {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to remove role",
		})
	}

	// Sync RBAC
	config.SyncRBACFromCasdoor()

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Role removed successfully",
	})
}

func SyncRBAC(c echo.Context) error {
	if err := config.SyncRBACFromCasdoor(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to sync RBAC",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "RBAC synced successfully",
	})
}
