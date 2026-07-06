package repository

import (
	"caiyun/internal/models"
	"time"

	"gorm.io/gorm"
)

// ExportHistoryRepository 导出历史数据访问层
type ExportHistoryRepository struct {
	db *gorm.DB
}

func NewExportHistoryRepository(db *gorm.DB) *ExportHistoryRepository {
	return &ExportHistoryRepository{db: db}
}

func (r *ExportHistoryRepository) WithContext(ctx interface{}) *ExportHistoryRepository {
	return r
}

func (r *ExportHistoryRepository) Create(history *models.ExportHistory) error {
	return r.db.Create(history).Error
}

func (r *ExportHistoryRepository) GetByID(id uint) (*models.ExportHistory, error) {
	var history models.ExportHistory
	err := r.db.First(&history, id).Error
	return &history, err
}

func (r *ExportHistoryRepository) GetByUserID(userID uint, page, pageSize int) ([]*models.ExportHistory, int64, error) {
	var histories []*models.ExportHistory
	var total int64

	query := r.db.Model(&models.ExportHistory{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&histories).Error
	return histories, total, err
}

func (r *ExportHistoryRepository) UpdateStatus(id uint, status string, filePath string, fileSize int64, errorMsg string) error {
	updates := map[string]interface{}{
		"status":      status,
		"updated_at":  time.Now(),
	}
	if filePath != "" {
		updates["file_path"] = filePath
	}
	if fileSize > 0 {
		updates["file_size"] = fileSize
	}
	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}
	return r.db.Model(&models.ExportHistory{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ExportHistoryRepository) GetRecentByUser(userID uint, limit int) ([]*models.ExportHistory, error) {
	var histories []*models.ExportHistory
	if limit <= 0 {
		limit = 10
	}
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}
