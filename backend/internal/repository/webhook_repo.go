package repository

import (
	"caiyun/internal/models"
	"time"

	"gorm.io/gorm"
)

// WebhookRepository Webhook端点数据访问层
type WebhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) WithContext(ctx interface{}) *WebhookRepository {
	return r
}

func (r *WebhookRepository) Create(endpoint *models.WebhookEndpoint) error {
	return r.db.Create(endpoint).Error
}

func (r *WebhookRepository) GetByID(id uint) (*models.WebhookEndpoint, error) {
	var endpoint models.WebhookEndpoint
	err := r.db.First(&endpoint, id).Error
	return &endpoint, err
}

func (r *WebhookRepository) GetByUserID(userID uint, page, pageSize int) ([]*models.WebhookEndpoint, int64, error) {
	var endpoints []*models.WebhookEndpoint
	var total int64

	query := r.db.Model(&models.WebhookEndpoint{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("updated_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&endpoints).Error
	return endpoints, total, err
}

func (r *WebhookRepository) GetActiveByUser(userID uint) ([]*models.WebhookEndpoint, error) {
	var endpoints []*models.WebhookEndpoint
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).
		Find(&endpoints).Error
	return endpoints, err
}

func (r *WebhookRepository) Update(endpoint *models.WebhookEndpoint, updates map[string]interface{}) error {
	return r.db.Model(endpoint).Updates(updates).Error
}

func (r *WebhookRepository) Delete(id uint) error {
	return r.db.Delete(&models.WebhookEndpoint{}, id).Error
}

func (r *WebhookRepository) UpdateLastTriggered(id uint) error {
	now := time.Now()
	return r.db.Model(&models.WebhookEndpoint{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_triggered_at": now,
			"updated_at":        now,
		}).Error
}

func (r *WebhookRepository) IncrementFailCount(id uint) error {
	return r.db.Model(&models.WebhookEndpoint{}).Where("id = ?", id).
		Update("fail_count", gorm.Expr("fail_count + 1")).Error
}

func (r *WebhookRepository) ResetFailCount(id uint) error {
	return r.db.Model(&models.WebhookEndpoint{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"fail_count": 0,
			"updated_at": time.Now(),
		}).Error
}
