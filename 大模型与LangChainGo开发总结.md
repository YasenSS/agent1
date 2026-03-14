# 大模型应用与 LangChainGo 框架实战总结

本项目（面经Scout Agent）是一个完整的大模型 Agent 落地案例。本文档总结了在开发过程中涉及的**大模型应用开发核心理论**，以及 **LangChainGo 框架的具体用法与功能**。

---

## 第一部分：大模型应用开发核心知识点

在本项目中，我们不仅是简单地调用大模型 API，而是构建了一个具备自主行动能力的 Agent。涉及的核心理论包括：

### 1. Agent 架构与 ReAct 范式
* **概念**：Agent = LLM（大脑） + Memory（记忆） + Tools（工具）。
* **ReAct (Reason + Act)**：大模型在解决复杂问题时，交替进行“思考（Thought）”和“行动（Action）”。
  * *Thought*：分析当前情况，决定下一步需要什么信息。
  * *Action*：决定调用哪个工具（如 `web_search`）。
  * *Observation*：观察工具返回的结果，进入下一轮 Thought，直到得出 *Final Answer*。
* **本项目应用**：Agent 接收到“搜索字节跳动面经”的指令后，自主决定先调用搜索工具，再根据搜索结果的 URL 调用抓取工具，最后进行总结。

### 2. Prompt Engineering (提示词工程)
* **角色设定**：通过 System Prompt 赋予 LLM 明确的身份（“你是一个资深的后端技术专家和面试顾问”），这能显著提升回答的专业度。
* **SOP (标准作业程序) 约束**：在 Prompt 中明确规定 Agent 的工作流程（搜索 -> 挑选 -> 读取 -> 总结），防止大模型“偷懒”或产生幻觉。
* **输出格式化**：通过提供 Markdown 模板，强制大模型按照“面试概况”、“高频考点”、“详细面试题”的结构输出，保证系统输出的稳定性。

### 3. 上下文窗口管理 (Context Window Management)
* **概念**：每个大模型都有最大 Token 限制（如 8K, 128K）。如果输入的文本过长，会导致 API 报错或模型遗忘。
* **本项目应用**：在 `WebScrapeTool` 中实现了**文本截断机制**（`MaxContentLength`）。当抓取的网页正文超过 4000 字符时，自动截断并追加 `...[内容已截断]`，保护大模型不被超长网页撑爆。

### 4. 幻觉抑制与兜底策略 (Hallucination Mitigation)
* **概念**：大模型在缺乏信息时容易“一本正经地胡说八道”。
* **本项目应用**：在 Prompt 中明确要求“保持客观，忠实于原文内容，不要编造面试题”；在工具层面，如果网页 404 或 403，工具会返回明确的错误字符串给大模型，让大模型知道“此路不通”，从而尝试其他链接，而不是瞎编内容。

---

## 第二部分：LangChainGo 框架核心用法总结

`github.com/tmc/langchaingo` 是 Go 语言生态中最主流的 LLM 开发框架。本项目深度使用了其核心模块：

### 1. LLM 客户端接入 (`llms/openai`)
* **功能**：标准化了与各类大模型的交互接口。
* **用法**：虽然包名叫 `openai`，但通过修改 `WithBaseURL` 和 `WithModel`，可以完美兼容 DeepSeek、通义千问等所有支持 OpenAI 格式的国产大模型。
```go
llm, err := openai.New(
    openai.WithModel("deepseek-chat"),
    openai.WithToken("sk-..."),
    openai.WithBaseURL("https://api.deepseek.com"),
)
```

### 2. 自定义工具开发 (`tools.Tool` 接口)
* **功能**：让大模型拥有“手和脚”，可以与外部世界交互。
* **用法**：实现 `Name()`, `Description()`, `Call()` 三个方法。
  * `Name`：工具的唯一标识（如 `web_search`）。
  * `Description`：**极其重要**，这是给大模型看的说明书，告诉大模型什么时候用这个工具、输入什么参数、输出什么结果。
  * `Call`：Go 语言编写的实际业务逻辑（如发起 HTTP 请求）。

### 3. Agent 组装与执行 (`agents` 模块)
* **功能**：将 LLM、Tools、Memory 绑定在一起，驱动 ReAct 循环。
* **用法**：使用 `agents.Initialize` 初始化，采用 `ZeroShotReactDescription` 模式。
```go
executor, err := agents.Initialize(
    llm,
    agentTools,
    agents.ZeroShotReactDescription,
    agents.WithMaxIterations(15), // 防止死循环，限制最大思考步数
    agents.WithPromptPrefix(systemPrompt), // 注入系统提示词
)
```

### 4. 记忆模块 (`memory.NewConversationBuffer`)
* **功能**：保存多轮对话历史，让 Agent 具备上下文理解能力。
* **用法**：初始化 `memory.NewConversationBuffer()` 并通过 `agents.WithMemory(mem)` 注入给 Agent。用户追问“上面提到的第一道题怎么解”时，Agent 能知道“上面”指的是什么。

### 5. 调试与回调机制 (`callbacks.LogHandler`)
* **功能**：透视 Agent 的“内心世界”，用于开发期调试。
* **用法**：通过 `agents.WithCallbacksHandler(callbacks.LogHandler{})` 注入。开启后，终端会打印出 Agent 每一步的 Thought、Action 和 Observation，非常有助于排查 Agent 卡在哪一步。

### 6. 输出解析错误处理 (`ParserErrorHandler`)
* **功能**：解决大模型输出格式不符合 LangChain 预期导致的崩溃问题。
* **用法**：通过 `agents.WithParserErrorHandler` 注入。当大模型没有严格按照 `Thought/Action/Action Input` 格式输出时，拦截错误，并将错误信息作为 Observation 返回给大模型，促使其“自我纠正”。

---

## 第三部分：项目实战避坑经验

1. **工具返回值的艺术**：在 `Call` 方法中，如果是**业务错误**（如网页 404、无搜索结果），**不要返回 Go 的 `error`**，而是返回一段描述错误的 `string`。这样大模型才能“看到”错误并改变策略。只有系统级崩溃（如网络断开）才返回 Go `error`。
2. **编码问题**：抓取国内网页（如 CSDN）时经常遇到 GBK 编码导致乱码，必须使用 `golang.org/x/net/html/charset` 动态探测并转换为 UTF-8，否则大模型读到的全是乱码，无法总结。
3. **迭代次数限制**：复杂的 Agent 任务（搜索+多次抓取+总结）通常需要 8-12 步，默认的 5 步极易触发 `agent not finished before max iterations` 错误，需要适当调大 `MaxIterations`。