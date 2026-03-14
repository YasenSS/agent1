package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"myagent/agent"
	"myagent/config"
	"myagent/output"
)

func main() {
	cfg := config.Load()
	if cfg.OpenAIKey == "" || cfg.OpenAIKey == "your_openai_api_key_here" {
		log.Fatal("请设置有效的 OPENAI_API_KEY 环境变量")
	}

	executor, err := agent.NewExecutor(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== 面经 Scout Agent ===")
	fmt.Println("我可以帮你收集和汇总互联网大厂的后端面经。")
	fmt.Println("示例指令：")
	fmt.Println("  - 帮我搜索字节跳动 Go 后端面经")
	fmt.Println("  - 收集阿里巴巴 Java 后端近半年的面经并总结高频考点")
	fmt.Println("输入 quit 退出程序。")
	fmt.Println(strings.Repeat("=", 40))

	for {
		fmt.Print("\n你: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if strings.EqualFold(input, "quit") || strings.EqualFold(input, "exit") {
			fmt.Println("再见！祝面试顺利！")
			break
		}

		fmt.Println("\n[Agent 思考中...]")
		result, err := agent.Run(ctx, executor, input)
		if err != nil {
			fmt.Printf("执行出错: %v\n", err)
			continue
		}

		fmt.Printf("\n面经Scout:\n%s\n", result)

		// 尝试自动保存为 Markdown 报告文件
		if strings.Contains(result, "面经汇总") || strings.Contains(result, "### 📋") {
			savedPath, err := output.SaveReport(result)
			if err != nil {
				fmt.Printf("⚠️ 报告保存失败: %v\n", err)
			} else if savedPath != "" {
				fmt.Printf("\n✅ 已自动为您生成精美的 Markdown 报告，保存在: %s\n", savedPath)
			}
		}
	}
}
