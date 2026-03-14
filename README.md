# 🚀 Interview-Insight-Agent | 大厂后端面经收割机

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)
![LLM](https://img.shields.io/badge/LLM-DeepSeek--V3-blue)
![Framework](https://img.shields.io/badge/Framework-langchaingo-orange)
![License](https://img.shields.io/badge/License-MIT-green)

**Interview-Insight-Agent** 是一个基于 Go 语言和 `langchaingo` 框架开发的自主智能体 (Autonomous Agent)。它旨在通过 AI 技术解决后端开发者在找工作时，“面经搜索碎片化、考点总结低效”的痛点。

---

## 💡 核心功能

- **🔎 自主多源搜索**：集成 **Tavily AI Search**，能够绕过传统搜索广告，精准定位牛客网、知乎、掘金等平台的深度面经。
- **📄 零爬虫网页解析**：利用 **Jina Reader** 协议，将复杂的富文本网页瞬间转化为 LLM 友好的 Markdown 纯文本。
- **🧠 结构化考点提炼**：基于 **DeepSeek-V3** 模型，自动对海量面经进行语义聚合，输出高频考点。
- **💬 交互式追问**：支持对话记忆，你可以追问：“针对刚才提到的 Redis 部分，字节跳动更喜欢问哪些底层原理？”

---

## 🏗️ 系统架构 (Architecture)

项目采用经典的 **ReAct (Reasoning and Acting)** 决策循环：

1. **用户指令** (User Input) -> "帮我找腾讯 Go 后端面经"
2. **推理思考** (Thought) -> "我需要先在网上搜索相关链接"
3. **执行行动** (Action) -> 调用 `WebSearchTool`
4. **观察结果** (Observation) -> 获取网页列表
5. **再次迭代** -> 调用 `JinaReaderTool` 读取正文 -> 总结 -> 输出结果

---

## 🛠️ 技术栈

| 组件 | 技术选型 | 备注 |
| :--- | :--- | :--- |
| **编程语言** | Go (Golang) | 高性能、工程化友好 |
| **Agent 框架** | [langchaingo](https://github.com/tmc/langchaingo) | 核心驱动框架 |
| **大脑 (LLM)** | DeepSeek-V3 / R1 | 性价比之王，逻辑推理极强 |
| **搜索工具** | Tavily API | 专为 AI 设计的搜索引擎 |
| **抓取协议** | Jina Reader | 免去维护复杂爬虫逻辑的烦恼 |

---

## 🚀 快速开始

### 1. 环境准备
- 确保本地已安装 Go 1.20 或更高版本。
- 获取 DeepSeek API Key 和 Tavily API Key。

### 2. 克隆与安装
```bash
git clone [https://github.com/YasenSS/agent1.git](https://github.com/YasenSS/agent1.git)
cd agent1
go mod tidy

