package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prannvs/campus-leave-system/internal/auth"
	"github.com/prannvs/campus-leave-system/internal/core"
	"github.com/prannvs/campus-leave-system/internal/models"
)

func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			core.ErrorResponse(c, http.StatusUnauthorized,
				models.ErrUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			core.ErrorResponse(c, http.StatusUnauthorized,
				models.ErrUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			core.ErrorResponse(c, http.StatusUnauthorized,
				models.ErrUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			core.ErrorResponse(c, http.StatusUnauthorized,
				models.ErrUnauthorized, "Role not found in context")
			c.Abort()
			return
		}

		userRole := role.(models.Role)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		core.ErrorResponse(c, http.StatusForbidden,
			models.ErrInvalidRole, "Insufficient permissions")
		c.Abort()
	}
}

func GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, models.ErrUnauthorized
	}
	return userID.(uint), nil
}

func GetUserRole(c *gin.Context) (models.Role, error) {
	role, exists := c.Get("role")
	if !exists {
		return "", models.ErrUnauthorized
	}
	return role.(models.Role), nil
}

func IsAdmin(c *gin.Context) bool {
	role, _ := GetUserRole(c)
	return role == models.RoleAdmin
}

func IsFaculty(c *gin.Context) bool {
	role, _ := GetUserRole(c)
	return role == models.RoleFaculty
}

func IsWarden(c *gin.Context) bool {
	role, _ := GetUserRole(c)
	return role == models.RoleWarden
}
