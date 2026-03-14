package agent

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	lc_tools "github.com/tmc/langchaingo/tools"

	"myagent/config"
	"myagent/tools"
)

const systemPrompt = `你是一个资深的后端技术专家和面试顾问，名叫"面经Scout"。
你的核心任务是帮助用户收集并分析互联网大厂（如阿里、腾讯、字节跳动、美团、百度等）的后端开发面试经验。

## 工作流程
1. 根据用户指定的公司和岗位方向，使用 web_search 工具搜索相关面经文章
2. 从搜索结果中挑选最相关的 2-3 篇文章
3. 使用 web_scraper 工具逐一读取文章内容
4. 综合所有文章内容，生成结构化的面经汇总报告

## 输出格式要求
请严格按照以下 Markdown 格式输出汇总报告：

### 📋 面经汇总：[公司名] - [岗位方向]

#### 一、面试概况
- 来源文章数量与链接
- 面试时间范围
- 面试轮次概况（几轮技术面 + HR面）

#### 二、高频考点 TOP 5
按出现频率从高到低列出，每个考点附带具体面试题示例。

#### 三、详细面试题整理
按知识领域分类（如：编程语言基础、数据结构与算法、数据库、缓存、消息队列、系统设计、项目经验等）

#### 四、面试建议
基于以上面经内容给出的备考建议。

## 注意事项
- 优先搜索近一年内的面经
- 如果某个 URL 无法访问，跳过并尝试其他链接
- 保持客观，忠实于原文内容，不要编造面试题`

func NewExecutor(cfg *config.Config) (chains.Chain, error) {
	llmOpts := []openai.Option{
		openai.WithModel(cfg.OpenAIModel),
		openai.WithToken(cfg.OpenAIKey),
	}
	if cfg.OpenAIBaseURL != "" {
		llmOpts = append(llmOpts, openai.WithBaseURL(cfg.OpenAIBaseURL))
	}

	llm, err := openai.New(llmOpts...)
	if err != nil {
		return nil, fmt.Errorf("初始化 LLM 失败: %w", err)
	}

	agentTools := []lc_tools.Tool{
		tools.WebSearchTool{ApiKey: cfg.TavilyKey},
		tools.WebScrapeTool{MaxContentLength: cfg.MaxContentLength},
	}

	mem := memory.NewConversationBuffer()

	executor, err := agents.Initialize(
		llm,
		agentTools,
		agents.ZeroShotReactDescription,
		agents.WithMaxIterations(cfg.MaxIterations),
		agents.WithMemory(mem),
		agents.WithPromptPrefix(systemPrompt), // 注入 System Prompt
		agents.WithCallbacksHandler(callbacks.LogHandler{}), // 开启详细日志，方便看 Agent 的思考过程
		agents.WithParserErrorHandler(agents.NewParserErrorHandler(func(s string) string {
			return "解析输出格式出错，请重新组织你的回答。确保严格遵循 Thought/Action/Action Input 格式。"
		})),
	)
	if err != nil {
		return nil, fmt.Errorf("初始化 Agent 失败: %w", err)
	}

	return executor, nil
}

func Run(ctx context.Context, executor chains.Chain, query string) (string, error) {
	result, err := chains.Run(ctx, executor, query)
	if err != nil {
		return "", fmt.Errorf("Agent 执行失败: %w", err)
	}
	return result, nil
}
