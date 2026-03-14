package tools

import (
	"context"
	"testing"
	"unicode/utf8"
)

func TestWebScrapeTool_Call(t *testing.T) {
	tool := WebScrapeTool{MaxContentLength: 500}
	ctx := context.Background()

	// 测试正常页面 (以一个常见的技术文章页面为例)
	url := "https://go.dev/doc/"
	result, err := tool.Call(ctx, url)
	if err != nil {
		t.Fatalf("抓取失败: %v", err)
	}
	if result == "" {
		t.Fatal("抓取内容为空")
	}
	t.Logf("成功抓取页面，内容长度: %d", utf8.RuneCountInString(result))

	// 测试内容截断
	if utf8.RuneCountInString(result) > 520 { // 500 + 截断提示的长度
		t.Fatalf("内容未被正确截断，当前长度: %d", utf8.RuneCountInString(result))
	}

	// 测试非法 URL
	_, err = tool.Call(ctx, "not-a-url")
	if err == nil {
		t.Fatal("对于非法 URL 应该返回错误，但返回了 nil")
	}
}
