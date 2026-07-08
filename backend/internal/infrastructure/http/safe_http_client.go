package http

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SafeHTTPClient 防止 SSRF 的 HTTP 客户端
// 规则：
// 1. 禁止访问私有/保留 IPv4 段：10.0.0.0/8、172.16.0.0/12、192.168.0.0/16、169.254.0.0/16
// 2. 禁止访问环回：127.0.0.0/8、::1
// 3. 禁止访问 IPv6 链路本地/唯一本地：fe80::/10、fc00::/7
// 4. 禁止访问 0.0.0.0/8
// 5. 仅允许端口 80、443、8080（可按需扩展）
type SafeHTTPClient struct {
	client     *http.Client
	allowPorts map[int]struct{}
}

// NewSafeHTTPClient 创建默认安全客户端（允许 80/443/8080，超时 10s）
func NewSafeHTTPClient() *SafeHTTPClient {
	return &SafeHTTPClient{
		client: &http.Client{Timeout: 10 * time.Second},
		allowPorts: map[int]struct{}{
			80:   {},
			443:  {},
			8080: {},
		},
	}
}

// NewSafeHTTPClientWithPorts 创建自定义端口限制的安全客户端
func NewSafeHTTPClientWithPorts(timeout time.Duration, ports ...int) *SafeHTTPClient {
	m := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		m[p] = struct{}{}
	}
	return &SafeHTTPClient{
		client:     &http.Client{Timeout: timeout},
		allowPorts: m,
	}
}

// Do 执行 HTTP 请求，请求前验证目标地址不命中私有/环回/IPv6 本地段
func (c *SafeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil {
		return nil, fmt.Errorf("无效请求")
	}
	if err := c.validateURL(req.URL); err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

// validateURL 解析主机并校验 IP 段与端口
func (c *SafeHTTPClient) validateURL(u *url.URL) error {
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("空目标主机")
	}

	// 先校验端口
	if err := c.validatePort(u.Port()); err != nil {
		return err
	}

	// DNS 解析，校验每一个解析到的 A/AAAA 记录
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		// 如果解析失败，尝试把 host 当 IP 直接校验
		ip := net.ParseIP(host)
		if ip == nil {
			return fmt.Errorf("无法解析目标主机: %s", host)
		}
		ips = []net.IP{ip}
	}

	for _, ip := range ips {
		if err := c.validateIP(ip); err != nil {
			return err
		}
	}
	return nil
}

// validatePort 仅允许白名单端口
func (c *SafeHTTPClient) validatePort(portStr string) error {
	if portStr == "" {
		// 无显式端口时使用 scheme 默认值
		return nil
	}
	var p int
	_, err := fmt.Sscanf(portStr, "%d", &p)
	if err != nil {
		return fmt.Errorf("无效端口: %s", portStr)
	}
	if _, ok := c.allowPorts[p]; !ok {
		return fmt.Errorf("端口 %d 不在允许范围内", p)
	}
	return nil
}

// validateIP 校验单条 IP 是否命中私有/环回/保留段
func (c *SafeHTTPClient) validateIP(ip net.IP) error {
	ip = ip.To4()
	if ip == nil {
		// IPv6：仅允许非本地、非保留地址（严格模式默认拒绝，可后续扩展白名单）
		return fmt.Errorf("禁止 IPv6 地址")
	}

	if ip.IsLoopback() {
		return fmt.Errorf("禁止环回地址 127.0.0.0/8")
	}
	if ip.IsPrivate() {
		return fmt.Errorf("禁止私有地址 10.0.0.0/8、172.16.0.0/12、192.168.0.0/16")
	}
	if ip.IsLinkLocalUnicast() {
		return fmt.Errorf("禁止链路本地地址 169.254.0.0/16")
	}
	if ip.IsLinkLocalMulticast() {
		return fmt.Errorf("禁止链路本地多播地址")
	}
	if ip.IsMulticast() {
		return fmt.Errorf("禁止多播地址")
	}

	// 0.0.0.0/8
	if strings.HasPrefix(ip.String(), "0.") {
		return fmt.Errorf("禁止 0.0.0.0/8")
	}
	return nil
}
