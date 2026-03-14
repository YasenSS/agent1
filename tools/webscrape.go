package tools

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

type WebScrapeTool struct {
	MaxContentLength int
}

func (t WebScrapeTool) Name() string {
	return "web_scraper"
}

func (t WebScrapeTool) Description() string {
	return `读取指定 URL 的网页正文内容。
输入必须是一个合法的 http 或 https URL。
输出是网页的纯文本正文（已去除 HTML 标签）。
注意：如果页面无法访问，会返回错误信息，请尝试其他 URL。`
}

func (t WebScrapeTool) Call(ctx context.Context, input string) (string, error) {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "http") {
		return "", fmt.Errorf("输入必须是合法的 URL，收到: %s", input)
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", input, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求页面失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("页面返回状态码 %d，无法读取内容。", resp.StatusCode), nil
	}

	// 处理网页编码，防止中文乱码
	utf8Reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		// 如果转换失败，回退到原始 Body
		utf8Reader = resp.Body
	}

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		return "", fmt.Errorf("解析 HTML 失败: %w", err)
	}

	// 移除无用元素
	doc.Find("script, style, nav, footer, header, .sidebar, .comment").Remove()

	var sb strings.Builder
	doc.Find("h1, h2, h3, h4, p, li, pre, code").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			sb.WriteString(text)
			sb.WriteString("\n")
		}
	})

	content := sb.String()

	// 文本截断：防止超出 LLM Token 限制
	maxLen := t.MaxContentLength
	if maxLen <= 0 {
		maxLen = 4000
	}
	if utf8.RuneCountInString(content) > maxLen {
		runes := []rune(content)
		content = string(runes[:maxLen]) + "\n...[内容已截断]"
	}

	if strings.TrimSpace(content) == "" {
		return "页面内容为空或无法提取有效正文。", nil
	}
	return content, nil
}
