package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SaveReport 将 Agent 输出的 Markdown 文本保存到本地文件
func SaveReport(content string) (string, error) {
	// 尝试从内容中提取公司和岗位作为文件名
	title := "面经汇总报告"
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "### 📋 面经汇总：") {
			title = strings.TrimSpace(strings.TrimPrefix(line, "### 📋 面经汇总："))
			// 清理文件名中的非法字符
			title = strings.ReplaceAll(title, "/", "_")
			title = strings.ReplaceAll(title, "\\", "_")
			title = strings.ReplaceAll(title, " ", "_")
			title = strings.ReplaceAll(title, "-", "_")
			title = strings.ReplaceAll(title, "__", "_")
			break
		}
	}

	// 创建 reports 目录
	outDir := "reports"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", fmt.Errorf("创建报告目录失败: %w", err)
	}

	// 生成带时间戳的文件名
	filename := fmt.Sprintf("%s_%s.md", title, time.Now().Format("20060102_150405"))
	filePath := filepath.Join(outDir, filename)

	// 写入文件
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return filePath, nil
}
