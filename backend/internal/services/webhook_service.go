package services

import (
	"bytes"
	"caiyun/internal/models"
	"caiyun/internal/repository"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

// WebhookService Webhook通知服务
type WebhookService struct {
	webhookRepo *repository.WebhookRepository
	deliveryRepo *repository.WebhookDeliveryRepository
	db          *gorm.DB
}

func NewWebhookService(
	webhookRepo *repository.WebhookRepository,
	deliveryRepo *repository.WebhookDeliveryRepository,
	db *gorm.DB,
) *WebhookService {
	return &WebhookService{
		webhookRepo:  webhookRepo,
		deliveryRepo: deliveryRepo,
		db:           db,
	}
}

// CreateEndpointRequest 创建Webhook端点请求
type CreateEndpointRequest struct {
	UserID   uint                   `json:"user_id"`
	Name     string                 `json:"name" binding:"required"`
	URL      string                 `json:"url" binding:"required,url"`
	Events   []string               `json:"events"`
	Secret   string                 `json:"secret"`
	Headers  map[string]string      `json:"headers"`
	IsActive *bool                  `json:"is_active"`
}

// UpdateEndpointRequest 更新Webhook端点请求
type UpdateEndpointRequest struct {
	Name     *string                `json:"name"`
	URL      *string                `json:"url"`
	Events   []string               `json:"events"`
	Secret   *string                `json:"secret"`
	Headers  map[string]string      `json:"headers"`
	IsActive *bool                  `json:"is_active"`
}

// ListEndpoints 获取Webhook端点列表
func (s *WebhookService) ListEndpoints(userID uint, page, pageSize int) ([]*models.WebhookEndpoint, int64, error) {
	return s.webhookRepo.GetByUserID(userID, page, pageSize)
}

// CreateEndpoint 创建Webhook端点
func (s *WebhookService) CreateEndpoint(req *CreateEndpointRequest) (*models.WebhookEndpoint, error) {
	endpoint := &models.WebhookEndpoint{
		UserID:   req.UserID,
		Name:     req.Name,
		URL:      req.URL,
		Events:   s.marshalEvents(req.Events),
		IsActive: true,
	}
	if req.Secret != "" {
		endpoint.Secret = req.Secret
	}
	if req.Headers != nil {
		headersJSON, _ := json.Marshal(req.Headers)
		endpoint.Headers = string(headersJSON)
	}
	if req.IsActive != nil {
		endpoint.IsActive = *req.IsActive
	}

	if err := s.webhookRepo.Create(endpoint); err != nil {
		return nil, err
	}
	return endpoint, nil
}

// UpdateEndpoint 更新Webhook端点
func (s *WebhookService) UpdateEndpoint(id uint, userID uint, req *UpdateEndpointRequest) (*models.WebhookEndpoint, error) {
	endpoint, err := s.webhookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Webhook端点不存在")
		}
		return nil, err
	}
	if endpoint.UserID != userID {
		return nil, errors.New("无权操作此Webhook端点")
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.URL != nil {
		updates["url"] = *req.URL
	}
	if req.Events != nil {
		updates["events"] = s.marshalEvents(req.Events)
	}
	if req.Secret != nil {
		updates["secret"] = *req.Secret
	}
	if req.Headers != nil {
		headersJSON, _ := json.Marshal(req.Headers)
		updates["headers"] = string(headersJSON)
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := s.webhookRepo.Update(endpoint, updates); err != nil {
		return nil, err
	}
	return s.webhookRepo.GetByID(id)
}

// DeleteEndpoint 删除Webhook端点
func (s *WebhookService) DeleteEndpoint(id uint, userID uint) error {
	endpoint, err := s.webhookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Webhook端点不存在")
		}
		return err
	}
	if endpoint.UserID != userID {
		return errors.New("无权操作此Webhook端点")
	}
	return s.webhookRepo.Delete(id)
}

// TestEndpoint 测试Webhook端点
func (s *WebhookService) TestEndpoint(id uint, userID uint) error {
	endpoint, err := s.webhookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Webhook端点不存在")
		}
		return err
	}
	if endpoint.UserID != userID {
		return errors.New("无权操作此Webhook端点")
	}
	if !endpoint.IsActive {
		return errors.New("Webhook端点已停用")
	}

	payload := map[string]interface{}{
		"event":     "test",
		"timestamp": time.Now().Unix(),
		"message":   "这是一条测试消息",
	}
	return s.deliverWebhook(endpoint, "test", payload)
}

// TriggerWebhook 触发Webhook通知
func (s *WebhookService) TriggerWebhook(userID uint, eventType string, payload map[string]interface{}) error {
	endpoints, err := s.webhookRepo.GetActiveByUser(userID)
	if err != nil {
		return err
	}

	if len(endpoints) == 0 {
		return nil
	}

	for _, endpoint := range endpoints {
		if !s.isEventSubscribed(endpoint, eventType) {
			continue
		}
		_ = s.deliverWebhook(endpoint, eventType, payload)
	}
	return nil
}

func (s *WebhookService) deliverWebhook(endpoint *models.WebhookEndpoint, eventType string, payload map[string]interface{}) error {
	startTime := time.Now()
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint.URL, bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "caiyun-webhook/1.0")
	req.Header.Set("X-Caiyun-Event", eventType)

	if endpoint.Secret != "" {
		signature := s.computeSignature(payloadBytes, endpoint.Secret)
		req.Header.Set("X-Caiyun-Signature", signature)
	}

	if endpoint.Headers != "" {
		var customHeaders map[string]string
		if err := json.Unmarshal([]byte(endpoint.Headers), &customHeaders); err == nil {
			for k, v := range customHeaders {
				req.Header.Set(k, v)
			}
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	duration := int(time.Since(startTime).Milliseconds())

	delivery := &models.WebhookDelivery{
		EndpointID: endpoint.ID,
		UserID:     endpoint.UserID,
		EventType:  eventType,
		Payload:    string(payloadBytes),
		DurationMs: &duration,
	}
	if resp != nil {
		statusCode := resp.StatusCode
		delivery.StatusCode = &statusCode
		_ = resp.Body.Close()
	}

	if err != nil {
		delivery.ErrorMsg = err.Error()
		_ = s.increaseFailCount(endpoint.ID)
	} else if delivery.StatusCode != nil && *delivery.StatusCode >= 200 && *delivery.StatusCode < 300 {
		_ = s.resetFailCount(endpoint.ID)
		_ = s.webhookRepo.UpdateLastTriggered(endpoint.ID)
	} else {
		delivery.ErrorMsg = fmt.Sprintf("HTTP %d", *delivery.StatusCode)
		_ = s.increaseFailCount(endpoint.ID)
	}

	_ = s.deliveryRepo.Create(delivery)
	return err
}

func (s *WebhookService) computeSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return "sha256=" + hex.EncodeToString(h.Sum(nil))
}

func (s *WebhookService) isEventSubscribed(endpoint *models.WebhookEndpoint, eventType string) bool {
	var events []string
	if err := json.Unmarshal([]byte(endpoint.Events), &events); err != nil {
		return false
	}
	for _, event := range events {
		if strings.EqualFold(event, eventType) {
			return true
		}
	}
	return false
}

func (s *WebhookService) marshalEvents(events []string) string {
	if events == nil {
		return "[]"
	}
	bytes, _ := json.Marshal(events)
	return string(bytes)
}

func (s *WebhookService) increaseFailCount(endpointID uint) error {
	return s.webhookRepo.IncrementFailCount(endpointID)
}

func (s *WebhookService) resetFailCount(endpointID uint) error {
	return s.webhookRepo.ResetFailCount(endpointID)
}
