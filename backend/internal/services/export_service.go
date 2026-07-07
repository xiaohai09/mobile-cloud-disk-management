package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"caiyun/internal/models"
	"caiyun/internal/repository"
)

// ExportService 数据导出服务
type ExportService struct {
	exportRepo   *repository.ExportHistoryRepository
	accountRepo  *repository.AccountRepository
	taskRepo     *repository.TaskLogRepository
	exchangeRepo *repository.ExchangeRecordRepository
}

func NewExportService(
	exportRepo *repository.ExportHistoryRepository,
	accountRepo *repository.AccountRepository,
	taskRepo *repository.TaskLogRepository,
	exchangeRepo *repository.ExchangeRecordRepository,
) *ExportService {
	return &ExportService{
		exportRepo:   exportRepo,
		accountRepo:  accountRepo,
		taskRepo:     taskRepo,
		exchangeRepo: exchangeRepo,
	}
}

// ExportRequest 导出请求
type ExportRequest struct {
	ExportType string                 `json:"export_type" binding:"required"` // accounts/tasks/exchange/records
	Format     string                 `json:"format" binding:"required,oneof=csv json"` // csv/json
	Filters    map[string]interface{} `json:"filters"`
	UserID     uint
}

// ExportResult 导出结果
type ExportResult struct {
	ID        uint   `json:"id"`
	FilePath  string `json:"file_path"`
	FileSize  int64  `json:"file_size"`
	Status    string `json:"status"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

// ExportData 执行数据导出
func (s *ExportService) ExportData(req *ExportRequest) (*ExportResult, error) {
	history := &models.ExportHistory{
		UserID:     req.UserID,
		ExportType: req.ExportType,
		Format:     req.Format,
		Status:     "pending",
	}
	if req.Filters != nil {
		filtersJSON, _ := json.Marshal(req.Filters)
		history.Filters = string(filtersJSON)
	}

	if err := s.exportRepo.Create(history); err != nil {
		return nil, err
	}

	go func(historyID uint) {
		_ = s.ExecuteExportJob(historyID, req)
	}(history.ID)

	return &ExportResult{
		ID:     history.ID,
		Status: "pending",
	}, nil
}

// ExecuteExportJob executes an export job for an existing history record
func (s *ExportService) ExecuteExportJob(historyID uint, req *ExportRequest) error {
	_ = s.exportRepo.UpdateStatus(historyID, "processing", "", 0, "")

	var filePath string
	var fileSize int64
	var err error

	switch req.ExportType {
	case "accounts":
		filePath, fileSize, err = s.exportAccounts(req)
	case "tasks":
		filePath, fileSize, err = s.exportTasks(req)
	case "exchange", "records":
		filePath, fileSize, err = s.exportExchangeRecords(req)
	default:
		err = fmt.Errorf("不支持的导出类型: %s", req.ExportType)
	}

	status := "success"
	errorMsg := ""
	if err != nil {
		status = "failed"
		errorMsg = err.Error()
	}

	_ = s.exportRepo.UpdateStatus(historyID, status, filePath, fileSize, errorMsg)
	return err
}

func (s *ExportService) exportAccounts(req *ExportRequest) (string, int64, error) {
	accounts, err := s.accountRepo.FindByUserID(req.UserID)
	if err != nil {
		return "", 0, err
	}

	if req.Format == "csv" {
		return s.exportToCSV(accounts, "accounts")
	}
	return s.exportToJSON(accounts, "accounts")
}

func (s *ExportService) exportTasks(req *ExportRequest) (string, int64, error) {
	_, logs, err := s.taskRepo.FindByUserID(req.UserID, 1, 1000)
	if err != nil {
		return "", 0, err
	}

	if req.Format == "csv" {
		return s.exportToCSV(logs, "tasks")
	}
	return s.exportToJSON(logs, "tasks")
}

func (s *ExportService) exportExchangeRecords(req *ExportRequest) (string, int64, error) {
	_, records, err := s.exchangeRepo.FindByUserID(req.UserID, 1, 1000)
	if err != nil {
		return "", 0, err
	}

	if req.Format == "csv" {
		return s.exportToCSV(records, "exchange")
	}
	return s.exportToJSON(records, "exchange")
}

func (s *ExportService) exportToCSV(data interface{}, prefix string) (string, int64, error) {
	_ = prefix
	_ = data
	return "", 0, errors.New("CSV导出暂未实现")
}

func (s *ExportService) exportToJSON(data interface{}, prefix string) (string, int64, error) {
	_ = prefix
	_ = data
	return "", 0, errors.New("JSON导出暂未实现")
}

// GetExportHistory 获取导出历史
func (s *ExportService) GetExportHistory(userID uint, page, pageSize int) ([]*models.ExportHistory, int64, error) {
	return s.exportRepo.GetByUserID(userID, page, pageSize)
}
