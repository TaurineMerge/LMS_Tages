package main

import (
	"flag"
	"log"

	"github.com/TaurineMerge/LMS_Tages/utils/knowledge-base-test-data-generator/internal/config"
	"github.com/TaurineMerge/LMS_Tages/utils/knowledge-base-test-data-generator/internal/generator"
	"github.com/TaurineMerge/LMS_Tages/utils/knowledge-base-test-data-generator/internal/ollama"
	"github.com/TaurineMerge/LMS_Tages/utils/knowledge-base-test-data-generator/internal/writer"
)

func main() {
	// Define and parse command-line flags
	configPath := flag.String("config", "config.yml", "Path to the configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize components
	ollamaClient := ollama.NewClient(cfg.OllamaURL, cfg.OllamaModel)
	sqlWriter, err := writer.NewSQLWriter(cfg.Output)
	if err != nil {
		log.Fatalf("Failed to create SQL writer: %v", err)
	}
	defer sqlWriter.Close()

	gen := generator.New(cfg, ollamaClient, sqlWriter)

	// Run the generator
	if err := gen.Run(); err != nil {
		log.Fatalf("An error occurred during generation: %v", err)
	}

	log.Printf("Successfully generated test data to %s\n", cfg.Output)
}
