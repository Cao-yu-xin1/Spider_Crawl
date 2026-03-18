// Package utils HTTP 工具函数
// ⚠️ 此文件包含故意设置的 Goroutine 泄漏问题，用于 bugfix 练习
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient HTTP 客户端
// 问题: 未设置 Timeout，导致连接不释放
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient 创建 HTTP 客户端
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Get 发送 GET 请求
func (c *HTTPClient) Get(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get request failed: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Post 发送 POST 请求
func (c *HTTPClient) Post(url string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("post request failed: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// GetWithTimeout 带超时的 GET 请求
// 问题: 虽然 context 有超时，但 client 本身无超时
func (c *HTTPClient) GetWithTimeout(url string, timeout time.Duration) ([]byte, error) {
	// 问题: 即使这里使用了带超时的请求，但底层 client 仍然没有默认超时
	// 如果服务器响应慢，仍然会泄漏
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get request failed: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// BatchGet 批量 GET 请求
// 问题: 并发请求时，如果服务器慢，会创建大量泄漏的 Goroutine
func (c *HTTPClient) BatchGet(urls []string) (map[string][]byte, error) {
	results := make(map[string][]byte)
	errChan := make(chan error, len(urls))

	for _, url := range urls {
		go func(u string) {
			data, err := c.Get(u)
			if err != nil {
				errChan <- err
				return
			}
			results[u] = data
			errChan <- nil
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		if err := <-errChan; err != nil {
			return nil, err
		}
	}

	return results, nil
}

// LongPolling 长轮询
// 问题: 长轮询场景下，如果没有超时，会一直等待
func (c *HTTPClient) LongPolling(url string) ([]byte, error) {
	// 问题: 无限等待服务器响应
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("long polling failed: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// StreamRequest 流式请求
// 问题: 流式请求如果不设置超时，可能永久阻塞
func (c *HTTPClient) StreamRequest(url string, handler func([]byte) error) error {
	resp, err := c.client.Get(url)
	if err != nil {
		return fmt.Errorf("stream request failed: %w", err)
	}
	defer resp.Body.Close()

	buffer := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read stream failed: %w", err)
		}

		if err := handler(buffer[:n]); err != nil {
			return err
		}
	}

	return nil
}

// RetryGet 带重试的 GET 请求
// 问题: 重试时如果服务器一直不响应，会累积泄漏
func (c *HTTPClient) RetryGet(url string, maxRetries int) ([]byte, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		data, err := c.Get(url)
		if err == nil {
			return data, nil
		}
		lastErr = err

		// 指数退避
		time.Sleep(time.Second * time.Duration(1<<i))
	}

	return nil, fmt.Errorf("retry failed after %d attempts: %w", maxRetries, lastErr)
}
