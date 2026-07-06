package repository

import (
	"caiyun/internal/models"

	"gorm.io/gorm"
)

// WebhookDeliveryRepository Webhook投递日志数据访问层
type WebhookDeliveryRepository struct {
	db *gorm.DB
}

func NewWebhookDeliveryRepository(db *gorm.DB) *WebhookDeliveryRepository {
	return &WebhookDeliveryRepository{db: db}
}

func (r *WebhookDeliveryRepository) WithContext(ctx interface{}) *WebhookDeliveryRepository {
	return r
}

func (r *WebhookDeliveryRepository) Create(delivery *models.WebhookDelivery) error {
	return r.db.Create(delivery).Error
}

func (r *WebhookDeliveryRepository) GetByEndpointID(endpointID uint, page, pageSize int) ([]*models.WebhookDelivery, int64, error) {
	var deliveries []*models.WebhookDelivery
	var total int64

	query := r.db.Model(&models.WebhookDelivery{}).Where("endpoint_id = ?", endpointID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&deliveries).Error
	return deliveries, total, err
}

func (r *WebhookDeliveryRepository) GetRecentByUser(userID uint, limit int) ([]*models.WebhookDelivery, error) {
	var deliveries []*models.WebhookDelivery
	if limit <= 0 {
		limit = 20
	}
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&deliveries).Error
	return deliveries, err
}
