package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIKey        string
	OpenAIBaseURL    string
	OpenAIModel      string
	TavilyKey        string
	MaxIterations    int
	MaxContentLength int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到 .env 文件，将从系统环境变量读取")
	}

	maxIter, _ := strconv.Atoi(getEnvOrDefault("AGENT_MAX_ITERATIONS", "5"))
	maxLen, _ := strconv.Atoi(getEnvOrDefault("AGENT_MAX_CONTENT_LENGTH", "4000"))

	return &Config{
		OpenAIKey:        os.Getenv("OPENAI_API_KEY"),
		OpenAIBaseURL:    getEnvOrDefault("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		OpenAIModel:      getEnvOrDefault("OPENAI_MODEL", "gpt-4o"),
		TavilyKey:        os.Getenv("TAVILY_API_KEY"),
		MaxIterations:    maxIter,
		MaxContentLength: maxLen,
	}
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
