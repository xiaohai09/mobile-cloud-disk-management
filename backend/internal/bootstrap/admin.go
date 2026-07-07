package bootstrap

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"caiyun/internal/models"
)

const (
	// DefaultAdminRole 默认管理员角色
	DefaultAdminRole = "admin"
	// DefaultAdminUsername 默认管理员用户名
	DefaultAdminUsername = "admin"
)

// SeedDefaultAdmin 初始化默认管理员账号
// 仅在管理员不存在时创建，不会覆盖现有用户
func SeedDefaultAdmin(repos Repositories) error {
	username := GetEnv("DEFAULT_ADMIN_USERNAME", DefaultAdminUsername)
	password := GetEnv("DEFAULT_ADMIN_PASSWORD", "")
	if password == "" {
		return fmt.Errorf("DEFAULT_ADMIN_PASSWORD 未设置，必须显式提供默认管理员密码")
	}

	existing, err := repos.User.FindByUsername(username)
	if err != nil {
		return fmt.Errorf("查询默认管理员失败: %w", err)
	}
	if existing != nil {
		log.Printf("[admin] 默认管理员 %s 已存在，跳过初始化", username)
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成密码哈希失败: %w", err)
	}

	now := time.Now()
	user := &models.User{
		Username:    username,
		Password:    string(hashedPassword),
		Email:       fmt.Sprintf("%s@example.com", username),
		Role:        DefaultAdminRole,
		TokenVersion: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := repos.User.Create(user); err != nil {
		return fmt.Errorf("创建默认管理员失败: %w", err)
	}

	log.Printf("[admin] 默认管理员初始化完成: username=%s", username)
	return nil
}
