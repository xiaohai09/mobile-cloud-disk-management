package models

import "time"

// ExportHistory 导出历史记录
type ExportHistory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	ExportType string   `gorm:"size:50;not null" json:"export_type"` // accounts/tasks/exchange/records
	Format    string    `gorm:"size:10;not null" json:"format"`       // csv/json
	Filters   string    `gorm:"type:json" json:"filters"`
	FilePath  string    `gorm:"size:255" json:"file_path"`
	FileSize  int64     `json:"file_size"`
	Status    string    `gorm:"size:20;not null;default:'pending'" json:"status"` // pending/success/failed
	ErrorMsg  string    `gorm:"size:500" json:"error_msg"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (ExportHistory) TableName() string {
	return "export_history"
}
