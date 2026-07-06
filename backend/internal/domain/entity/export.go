package domain

import "time"

// ExportFormat represents supported export formats
type ExportFormat string

const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
	ExportFormatXLSX ExportFormat = "xlsx"
)

// ExportType represents what data to export
type ExportType string

const (
	ExportTypeTaskLogs      ExportType = "task_logs"
	ExportTypeCloudStats    ExportType = "cloud_stats"
	ExportTypeExchangeRecords ExportType = "exchange_records"
	ExportTypeAccounts      ExportType = "accounts"
	ExportTypeAll           ExportType = "all"
)

// ExportJob represents an export job request
type ExportJob struct {
	ID          uint
	UserID      uint
	Type        ExportType
	Format      ExportFormat
	Filters     map[string]interface{}
	Status      string
	FilePath    string
	FileSize    int64
	ExpiresAt   time.Time
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// WebhookEventType represents webhook event types
type WebhookEventType string

const (
	WebhookEventTaskSuccess  WebhookEventType = "task.success"
	WebhookEventTaskFailure  WebhookEventType = "task.failure"
	WebhookEventExchangeHit  WebhookEventType = "exchange.hit"
	WebhookEventSystemAlert  WebhookEventType = "system.alert"
)

// WebhookEndpoint represents a webhook endpoint
type WebhookEndpoint struct {
	ID          uint
	UserID      uint
	Name        string
	URL         string
	Events      []WebhookEventType
	Secret      string
	IsActive    bool
	Headers     map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PlatformType represents supported cloud platforms
type PlatformType string

const (
	PlatformTypeCloud189   PlatformType = "cloud_189"
	PlatformTypeCloud115   PlatformType = "cloud_115"
	PlatformTypeAliyun     PlatformType = "aliyun"
	PlatformTypeQuark      PlatformType = "quark"
	PlatformTypeBaidu      PlatformType = "baidu"
)

// PlatformAccount represents a cloud platform account
type PlatformAccount struct {
	ID          uint
	UserID      uint
	Platform    PlatformType
	Username    string
	DisplayName string
	AuthData    map[string]interface{} // encrypted auth tokens
	IsActive    bool
	LastSyncAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
