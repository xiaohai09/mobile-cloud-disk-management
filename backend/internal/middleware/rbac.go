package middleware

import (
	"net/http"
	"strings"

	"caiyun/internal/models"
	"caiyun/internal/repository"

	"github.com/gin-gonic/gin"
)

const (
	// PermissionsKey context key for user permissions
	PermissionsKey = "permissions"
	// PermissionCodesKey context key for permission codes slice
	PermissionCodesKey = "permission_codes"
)

// PermissionMiddleware 加载用户权限并注入上下文。
// 对已认证请求：加载角色权限 + 用户特有权限；
// 对未认证请求：注入空权限集合，由后续业务自行决定是否放行。
func PermissionMiddleware(permRepo *repository.PermissionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserID(c)
		if userID == 0 {
			c.Set(PermissionsKey, []models.Permission{})
			c.Set(PermissionCodesKey, []string{})
			c.Next()
			return
		}

		codes, err := permRepo.FindByUserID(userID)
		if err != nil {
			// 权限加载失败时，按最小权限原则拒绝
			abortWithError(c, http.StatusInternalServerError, "权限加载失败")
			return
		}

		// 预加载 Permission 对象（可选，用于审计日志）
		perms := make([]models.Permission, 0, len(codes))
		// 简化处理：仅传递 codes，按需从 repo 查详情
		_ = perms

		c.Set(PermissionsKey, []models.Permission{})
		c.Set(PermissionCodesKey, codes)
		c.Next()
	}
}

// GetPermissionCodes 从上下文中获取用户权限标识列表（只读）。
func GetPermissionCodes(c *gin.Context) []string {
	v, _ := c.Get(PermissionCodesKey)
	if v == nil {
		return nil
	}
	codes, ok := v.([]string)
	if !ok {
		return nil
	}
	return codes
}

// HasPermission 检查当前用户是否具备指定权限标识。
func HasPermission(c *gin.Context, code string) bool {
	codes := GetPermissionCodes(c)
	if len(codes) == 0 {
		return false
	}
	for _, perm := range codes {
		if perm == code {
			return true
		}
		// 支持通配符，如 "user:*"
		if strings.HasSuffix(perm, ":*") {
			prefix := strings.TrimSuffix(perm, "*")
			if strings.HasPrefix(code, prefix) {
				return true
			}
		}
	}
	return false
}

// RequirePermission 要求当前用户具备指定权限标识，否则返回 403。
// 用于替代或增强现有的 AdminMiddleware，提供细粒度权限控制。
func RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasPermission(c, code) {
			abortWithError(c, http.StatusForbidden, "权限不足")
			return
		}
		c.Next()
	}
}

// RequireAnyPermission 要求当前用户具备至少一个指定权限标识。
func RequireAnyPermission(codes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, code := range codes {
			if HasPermission(c, code) {
				c.Next()
				return
			}
		}
		abortWithError(c, http.StatusForbidden, "权限不足")
	}
}

// getUserID 从已认证上下文中安全地读取用户 ID。
func getUserID(c *gin.Context) uint {
	v, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	switch id := v.(type) {
	case uint:
		return id
	case int:
		return uint(id)
	case float64:
		return uint(id)
	default:
		return 0
	}
}
