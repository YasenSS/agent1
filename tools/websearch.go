package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type WebSearchTool struct {
	ApiKey string
}

func (t WebSearchTool) Name() string {
	return "web_search"
}

func (t WebSearchTool) Description() string {
	return `在互联网上搜索与后端面经相关的文章。
输入应该是搜索关键词，例如 "字节跳动 Go 后端面经 2024"。
输出是搜索结果列表，包含标题和 URL。
使用建议：组合 "公司名 + 岗位 + 面经" 作为搜索词效果最佳。`
}

func (t WebSearchTool) Call(ctx context.Context, input string) (string, error) {
	if t.ApiKey == "" {
		return "", fmt.Errorf("TAVILY_API_KEY 未配置")
	}

	apiURL := "https://api.tavily.com/search"
	
	requestBody, err := json.Marshal(map[string]interface{}{
		"api_key": t.ApiKey,
		"query":   input,
		"max_results": 5,
	})
	if err != nil {
		return "", fmt.Errorf("构建请求体失败: %w", err)
	}

	resp, err := http.Post(apiURL, "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		return "", fmt.Errorf("搜索请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Results []struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析搜索结果失败: %w, raw body: %s", err, string(body))
	}

	if len(result.Results) == 0 {
		return fmt.Sprintf("未找到相关结果，Tavily API 响应: %s", string(body)), nil
	}

	var sb strings.Builder
	for i, r := range result.Results {
		sb.WriteString(fmt.Sprintf("%d. %s\n   URL: %s\n", i+1, r.Title, r.URL))
	}

	if sb.Len() == 0 {
		return "未找到相关结果，请尝试更换关键词。", nil
	}
	return sb.String(), nil
}
