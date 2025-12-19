package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// PromptConfig holds the templates for Ollama prompts.
type PromptConfig struct {
	CategoryTitle      string `yaml:"category_title"`
	CourseContent      string `yaml:"course_content"`
	LessonTitle        string `yaml:"lesson_title"`
	LessonContentBlock string `yaml:"lesson_content_block"`
}

// Config holds the application configuration.
type Config struct {
	Output      string       `yaml:"output"`
	OllamaURL   string       `yaml:"ollama_url"`
	OllamaModel string       `yaml:"ollama_model"`
	Counts      CountConfig  `yaml:"counts"`
	Prompts     PromptConfig `yaml:"prompts"`
}

// CountConfig holds the configuration for the number of items to generate.
type CountConfig struct {
	Categories    string `yaml:"categories"`
	Courses       string `yaml:"courses"`
	Lessons       string `yaml:"lessons"`
	ContentBlocks string `yaml:"content_blocks"`
}

// Load reads a configuration file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config yaml: %w", err)
	}

	return &cfg, nil
}

// ParseRange parses a string like "5" or "2:6" into a min and max integer.
func ParseRange(rangeStr string) (int, int, error) {
	if strings.Contains(rangeStr, ":") {
		parts := strings.Split(rangeStr, ":")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid range format: %s", rangeStr)
		}
		min, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid min value in range: %s", parts[0])
		}
		max, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid max value in range: %s", parts[1])
		}
		if min > max {
			return 0, 0, fmt.Errorf("min value cannot be greater than max value in range")
		}
		return min, max, nil
	}

	val, err := strconv.Atoi(rangeStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid number format: %s", rangeStr)
	}
	return val, val, nil
}