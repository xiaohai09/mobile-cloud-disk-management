package infrastructure

import (
	"fmt"
	"time"

	"caiyun/internal/domain/entity"
)

// PlatformAdapter defines the interface for cloud platform adapters
type PlatformAdapter interface {
	// GetPlatform returns the platform type
	GetPlatform() domain.PlatformType

	// ValidateAuth validates the auth data for this platform
	ValidateAuth(authData map[string]interface{}) error

	// SyncAccount syncs account information from the platform
	SyncAccount(account *domain.PlatformAccount) error

	// GetDailyTasks returns the daily task types for this platform
	GetDailyTasks() []string

	// ExecuteTask executes a specific task type for the account
	ExecuteTask(account *domain.PlatformAccount, taskType string) (map[string]interface{}, error)
}

// PlatformAdapterRegistry manages platform adapters
type PlatformAdapterRegistry struct {
	adapters map[domain.PlatformType]PlatformAdapter
}

// NewPlatformAdapterRegistry creates a new platform adapter registry
func NewPlatformAdapterRegistry() *PlatformAdapterRegistry {
	return &PlatformAdapterRegistry{
		adapters: make(map[domain.PlatformType]PlatformAdapter),
	}
}

// Register registers a platform adapter
func (r *PlatformAdapterRegistry) Register(adapter PlatformAdapter) {
	r.adapters[adapter.GetPlatform()] = adapter
}

// Get returns the adapter for the specified platform
func (r *PlatformAdapterRegistry) Get(platform domain.PlatformType) (PlatformAdapter, error) {
	adapter, ok := r.adapters[platform]
	if !ok {
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
	return adapter, nil
}

// List returns all registered platform types
func (r *PlatformAdapterRegistry) List() []domain.PlatformType {
	platforms := make([]domain.PlatformType, 0, len(r.adapters))
	for p := range r.adapters {
		platforms = append(platforms, p)
	}
	return platforms
}

// Cloud189Adapter implements PlatformAdapter for China Telecom Cloud 189
type Cloud189Adapter struct{}

func (a *Cloud189Adapter) GetPlatform() domain.PlatformType {
	return domain.PlatformTypeCloud189
}

func (a *Cloud189Adapter) ValidateAuth(authData map[string]interface{}) error {
	// Validate required fields for 189 cloud
	if _, ok := authData["account"]; !ok {
		return fmt.Errorf("missing account field")
	}
	if _, ok := authData["password"]; !ok {
		return fmt.Errorf("missing password field")
	}
	return nil
}

func (a *Cloud189Adapter) SyncAccount(account *domain.PlatformAccount) error {
	// Sync account info from 189 cloud API
	account.LastSyncAt = &time.Time{}
	*account.LastSyncAt = time.Now()
	return nil
}

func (a *Cloud189Adapter) GetDailyTasks() []string {
	return []string{"signin", "shake", "redpacket", "cloudbattle"}
}

func (a *Cloud189Adapter) ExecuteTask(account *domain.PlatformAccount, taskType string) (map[string]interface{}, error) {
	// Execute task on 189 cloud platform
	return map[string]interface{}{
		"task_type": taskType,
		"status":    "success",
		"reward":    "10",
	}, nil
}

// Cloud115Adapter implements PlatformAdapter for 115 Cloud
type Cloud115Adapter struct{}

func (a *Cloud115Adapter) GetPlatform() domain.PlatformType {
	return domain.PlatformTypeCloud115
}

func (a *Cloud115Adapter) ValidateAuth(authData map[string]interface{}) error {
	if _, ok := authData["cookie"]; !ok {
		return fmt.Errorf("missing cookie field")
	}
	return nil
}

func (a *Cloud115Adapter) SyncAccount(account *domain.PlatformAccount) error {
	account.LastSyncAt = &time.Time{}
	*account.LastSyncAt = time.Now()
	return nil
}

func (a *Cloud115Adapter) GetDailyTasks() []string {
	return []string{"signin", "shake", "daily_task"}
}

func (a *Cloud115Adapter) ExecuteTask(account *domain.PlatformAccount, taskType string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"task_type": taskType,
		"status":    "success",
		"reward":    "5",
	}, nil
}

// AliyunAdapter implements PlatformAdapter for Alibaba Cloud
type AliyunAdapter struct{}

func (a *AliyunAdapter) GetPlatform() domain.PlatformType {
	return domain.PlatformTypeAliyun
}

func (a *AliyunAdapter) ValidateAuth(authData map[string]interface{}) error {
	if _, ok := authData["access_key_id"]; !ok {
		return fmt.Errorf("missing access_key_id field")
	}
	if _, ok := authData["access_key_secret"]; !ok {
		return fmt.Errorf("missing access_key_secret field")
	}
	return nil
}

func (a *AliyunAdapter) SyncAccount(account *domain.PlatformAccount) error {
	account.LastSyncAt = &time.Time{}
	*account.LastSyncAt = time.Now()
	return nil
}

func (a *AliyunAdapter) GetDailyTasks() []string {
	return []string{"signin", "resource_task"}
}

func (a *AliyunAdapter) ExecuteTask(account *domain.PlatformAccount, taskType string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"task_type": taskType,
		"status":    "success",
		"reward":    "15",
	}, nil
}

// QuarkAdapter implements PlatformAdapter for Quark Drive
type QuarkAdapter struct{}

func (a *QuarkAdapter) GetPlatform() domain.PlatformType {
	return domain.PlatformTypeQuark
}

func (a *QuarkAdapter) ValidateAuth(authData map[string]interface{}) error {
	if _, ok := authData["cookie"]; !ok {
		return fmt.Errorf("missing cookie field")
	}
	return nil
}

func (a *QuarkAdapter) SyncAccount(account *domain.PlatformAccount) error {
	account.LastSyncAt = &time.Time{}
	*account.LastSyncAt = time.Now()
	return nil
}

func (a *QuarkAdapter) GetDailyTasks() []string {
	return []string{"signin", "draw", "task"}
}

func (a *QuarkAdapter) ExecuteTask(account *domain.PlatformAccount, taskType string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"task_type": taskType,
		"status":    "success",
		"reward":    "3",
	}, nil
}

// BaiduAdapter implements PlatformAdapter for Baidu Netdisk
type BaiduAdapter struct{}

func (a *BaiduAdapter) GetPlatform() domain.PlatformType {
	return domain.PlatformTypeBaidu
}

func (a *BaiduAdapter) ValidateAuth(authData map[string]interface{}) error {
	if _, ok := authData["bdstoken"]; !ok {
		return fmt.Errorf("missing bdstoken field")
	}
	return nil
}

func (a *BaiduAdapter) SyncAccount(account *domain.PlatformAccount) error {
	account.LastSyncAt = &time.Time{}
	*account.LastSyncAt = time.Now()
	return nil
}

func (a *BaiduAdapter) GetDailyTasks() []string {
	return []string{"signin", "task"}
}

func (a *BaiduAdapter) ExecuteTask(account *domain.PlatformAccount, taskType string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"task_type": taskType,
		"status":    "success",
		"reward":    "2",
	}, nil
}

// RegisterDefaultAdapters registers all default platform adapters
func RegisterDefaultAdapters(registry *PlatformAdapterRegistry) {
	registry.Register(&Cloud189Adapter{})
	registry.Register(&Cloud115Adapter{})
	registry.Register(&AliyunAdapter{})
	registry.Register(&QuarkAdapter{})
	registry.Register(&BaiduAdapter{})
}
