package infrastructure

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"caiyun/internal/domain/entity"
)

// ExportService handles data export operations
type ExportService struct {
	exportDir string
}

// NewExportService creates a new export service
func NewExportService(exportDir string) *ExportService {
	if exportDir == "" {
		exportDir = "/tmp/caiyun/exports"
	}
	return &ExportService{exportDir: exportDir}
}

// ExportData exports data in the specified format
func (s *ExportService) ExportData(job *domain.ExportJob) (string, error) {
	if err := os.MkdirAll(s.exportDir, 0755); err != nil {
		return "", fmt.Errorf("create export dir failed: %w", err)
	}

	filename := fmt.Sprintf("export_%s_%d_%d.%s",
		job.Type,
		job.UserID,
		time.Now().Unix(),
		job.Format,
	)
	filepath := filepath.Join(s.exportDir, filename)

	var data []byte
	var err error

	switch job.Format {
	case domain.ExportFormatCSV:
		data, err = s.generateCSV(job)
	case domain.ExportFormatJSON:
		data, err = s.generateJSON(job)
	case domain.ExportFormatXLSX:
		// XLSX requires external library, fallback to CSV
		job.Format = domain.ExportFormatCSV
		filepath = strings.TrimSuffix(filepath, ".xlsx") + ".csv"
		data, err = s.generateCSV(job)
	default:
		return "", fmt.Errorf("unsupported export format: %s", job.Format)
	}

	if err != nil {
		return "", fmt.Errorf("generate export failed: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return "", fmt.Errorf("write export file failed: %w", err)
	}

	return filepath, nil
}

// generateCSV generates CSV export
func (s *ExportService) generateCSV(job *domain.ExportJob) ([]byte, error) {
	// This is a simplified version - in production, fetch real data from repository
	headers := []string{"id", "created_at", "status", "message"}

	var records [][]string
	records = append(records, headers)

	// Sample data - replace with actual repository calls
	for i := 1; i <= 10; i++ {
		records = append(records, []string{
			fmt.Sprintf("%d", i),
			time.Now().Format(time.RFC3339),
			"success",
			fmt.Sprintf("Sample record %d", i),
		})
	}

	var buf strings.Builder
	writer := csv.NewWriter(&buf)
	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}
	writer.Flush()
	return []byte(buf.String()), nil
}

// generateJSON generates JSON export
func (s *ExportService) generateJSON(job *domain.ExportJob) ([]byte, error) {
	// Simplified - replace with actual data fetching
	data := map[string]interface{}{
		"type":      job.Type,
		"user_id":   job.UserID,
		"exported_at": time.Now().Format(time.RFC3339),
		"records":   []string{"sample1", "sample2"},
	}
	return json.MarshalIndent(data, "", "  ")
}

// CleanupExpired removes expired export files
func (s *ExportService) CleanupExpired(maxAge time.Duration) error {
	if err := os.MkdirAll(s.exportDir, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(s.exportDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-maxAge)
	var lastErr error
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			lastErr = err
			continue
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(filepath.Join(s.exportDir, entry.Name())); err != nil {
				lastErr = err
			}
		}
	}
	return lastErr
}
