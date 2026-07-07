package repository

import (
	"caiyun/internal/models"
	"gorm.io/gorm"
)

// PermissionRepository 权限仓库
type PermissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓库
func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// FindByCode 根据权限标识查找权限
func (r *PermissionRepository) FindByCode(code string) (*models.Permission, error) {
	var perm models.Permission
	if err := r.db.Where("code = ?", code).First(&perm).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

// FindByRole 查找角色的所有权限标识
func (r *PermissionRepository) FindByRole(role string) ([]string, error) {
	var codes []string
	if err := r.db.Table("permissions").
		Select("permissions.code").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role = ?", role).
		Pluck("permissions.code", &codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}

// FindByUserID 查找用户的所有权限标识（角色权限 + 用户特有权限）
func (r *PermissionRepository) FindByUserID(userID uint) ([]string, error) {
	var codes []string
	if err := r.db.Table("permissions").
		Select("DISTINCT permissions.code").
		Joins("LEFT JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("LEFT JOIN user_permissions ON user_permissions.permission_id = permissions.id").
		Where("role_permissions.role = (SELECT role FROM users WHERE id = ?) OR user_permissions.user_id = ?", userID, userID).
		Pluck("permissions.code", &codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}
