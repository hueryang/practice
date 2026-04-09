package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://open.bigmodel.cn/api/paas/v4"
	chatPath       = "/chat/completions"
)

// EnvAPIKey is the primary environment variable for the BigModel API key.
const EnvAPIKey = "BIGMODEL_API_KEY"

// EnvAPIKeyAlt is an alternate env name some users configure for Zhipu/BigModel.
const EnvAPIKeyAlt = "ZHIPU_API_KEY"

// EnvBaseURL overrides the API base URL (trailing slash stripped). Used for tests.
const EnvBaseURL = "BIGMODEL_BASE_URL"

// Client calls the BigModel chat completions API.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClientFromEnv builds a client using BIGMODEL_API_KEY or ZHIPU_API_KEY.
func NewClientFromEnv() (*Client, error) {
	key := strings.TrimSpace(os.Getenv(EnvAPIKey))
	if key == "" {
		key = strings.TrimSpace(os.Getenv(EnvAPIKeyAlt))
	}
	if key == "" {
		return nil, fmt.Errorf("未设置 API Key：请设置环境变量 %s 或 %s", EnvAPIKey, EnvAPIKeyAlt)
	}
	return NewClient(key), nil
}

// NewClient returns a client with the given API key and default endpoint and timeout.
func NewClient(apiKey string) *Client {
	base := strings.TrimSuffix(strings.TrimSpace(os.Getenv(EnvBaseURL)), "/")
	if base == "" {
		base = defaultBaseURL
	}
	return &Client{
		BaseURL: base,
		APIKey:  strings.TrimSpace(apiKey),
		HTTPClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *Client) endpoint() string {
	return strings.TrimSuffix(c.BaseURL, "/") + chatPath
}

func (c *Client) httpDo() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

// ChatCompletion calls POST /chat/completions. See https://docs.bigmodel.cn/
func (c *Client) ChatCompletion(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("llm.Client: API Key 为空")
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化 chat 请求: %w", err)
	}
	httpReq, err := http.NewRequest(http.MethodPost, c.endpoint(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("构造 HTTP 请求: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpDo().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求 BigModel API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BigModel API HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var out ChatCompletionResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("解析 BigModel 响应 JSON: %w", err)
	}
	if len(out.Choices) == 0 {
		return nil, fmt.Errorf("BigModel 响应中 choices 为空")
	}
	return &out, nil
}
