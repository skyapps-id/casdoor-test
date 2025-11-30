package handlers

import (
	"fmt"
	"net/http"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/labstack/echo/v4"
	"github.com/skyapps-id/casdoor-test/config"
)

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

func GetLoginURL(c echo.Context) error {
	url := config.CasdoorClient.GetSigninUrl("http://localhost:9000/callback")
	return c.JSON(http.StatusOK, map[string]string{
		"url": url,
	})
}

func HandleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	token, err := config.CasdoorClient.GetOAuthToken(code, state)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to get token",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token.AccessToken,
		"expires_in": token.Expiry,
	})
}

func GetCurrentUser(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User ID missing or invalid",
		})
	}

	userEmail, ok := c.Get("user_email").(string)
	if !ok || userEmail == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User email missing or invalid",
		})
	}

	// Get full user info dari Casdoor
	user, err := config.CasdoorClient.GetUser(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user info",
		})
	}

	fmt.Println(user, "_+_+_+_+")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"email":   userEmail,
		// "display_name": user.DisplayName,
		// "roles":        user.Roles,
	})
}

func ListUsers(c echo.Context) error {
	users, err := config.CasdoorClient.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get users",
		})
	}

	// Filter hanya user dari skyapps
	skyappsUsers := []interface{}{}
	for _, user := range users {
		if user.Owner == "skyapps" {
			skyappsUsers = append(skyappsUsers, map[string]interface{}{
				"username":     user.Name,
				"email":        user.Email,
				"display_name": user.DisplayName,
				"roles":        user.Roles,
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": skyappsUsers,
		"total": len(skyappsUsers),
	})
}

func AddUser(c echo.Context) error {
	var req struct {
		Username    string `json:"username" validate:"required"`
		DisplayName string `json:"display_name" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		Password    string `json:"password" validate:"required,min=8"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	user := &casdoorsdk.User{
		Owner:       "skyapps",
		Name:        req.Username,
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Password:    req.Password,
	}

	affected, err := config.CasdoorClient.AddUser(user)
	if err != nil || !affected {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"user": map[string]string{
			"username": req.Username,
			"email":    req.Email,
		},
	})
}

func UpdateUser(c echo.Context) error {
	username := c.Param("username")

	var req struct {
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	user := &casdoorsdk.User{
		Owner:       "skyapps",
		Name:        username,
		DisplayName: req.DisplayName,
		Email:       req.Email,
	}

	affected, err := config.CasdoorClient.UpdateUser(user)
	if err != nil || !affected {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update user",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User updated successfully",
	})
}

func DeleteUser(c echo.Context) error {
	username := c.Param("username")

	user := &casdoorsdk.User{
		Owner: "skyapps",
		Name:  username,
	}

	affected, err := config.CasdoorClient.DeleteUser(user)
	if err != nil || !affected {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}
