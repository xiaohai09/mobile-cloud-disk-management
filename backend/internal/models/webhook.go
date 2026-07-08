package models

import "time"

// WebhookEndpoint Webhook端点
type WebhookEndpoint struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	URL            string    `gorm:"size:500;not null" json:"url"`
	Events         string    `gorm:"type:json;not null" json:"events"` // JSON array
	Secret         string    `gorm:"size:100" json:"-"`
	Headers        string    `gorm:"type:json" json:"headers"` // JSON object
	IsActive       bool      `gorm:"not null;default:1" json:"is_active"`
	LastTriggeredAt *time.Time `json:"last_triggered_at"`
	FailCount      int       `gorm:"not null;default:0" json:"fail_count"`
	Status         string    `gorm:"size:20;not null;default:'active'" json:"status"`
	NextRetryAt    *time.Time `json:"next_retry_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 指定表名
func (WebhookEndpoint) TableName() string {
	return "webhook_endpoints"
}

// WebhookDelivery Webhook投递日志
type WebhookDelivery struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	EndpointID  uint      `gorm:"not null;index" json:"endpoint_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	EventType   string    `gorm:"size:50;not null" json:"event_type"`
	Payload     string    `gorm:"type:json" json:"payload"`
	StatusCode  *int      `json:"status_code"`
	ResponseBody string   `gorm:"type:text" json:"response_body"`
	ErrorMsg    string    `gorm:"size:500" json:"error_msg"`
	DurationMs  *int      `json:"duration_ms"`
	Status      string    `gorm:"size:20;not null;default:'pending'" json:"status"`
	RetryCount  int       `gorm:"not null;default:0" json:"retry_count"`
	NextRetryAt *time.Time `json:"next_retry_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 指定表名
func (WebhookDelivery) TableName() string {
	return "webhook_deliveries"
}
