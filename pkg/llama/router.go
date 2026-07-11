package llama

import (
	"fmt"

	"github.com/charmbracelet/log"
)

func SelectModel() (selectedModel string, selectedQuant string, err error) {
	// selectedModel "" (zero value) = use to andyAPI
	// selectedQuant "" (zero value) = let llama.cpp decide (do not specify)

	// Perform benchmark
	tps, err := Benchmark()
	if err != nil {
		return "", "", fmt.Errorf("failed to benchmark performance: %w", err)
	}

	// Route between models
	// TODO: Add logic for routing between different models/quantizations based on benchmark results
	if tps < 300 {
		log.Info("router: model selected", "model", "none", "tps", tps)
		return "", "", nil
	} else if tps < 700 {
		log.Info("router: model selected", "model", "Mindcraft-CE/Andy-4.2-Micro-GGUF", "tps", tps)
		return "Mindcraft-CE/Andy-4.2-Micro-GGUF", "", nil
	} else if tps < 1200 {
		log.Info("router: model selected", "model", "Mindcraft-CE/Andy-4.2-Air-GGUF", "tps", tps)
		return "Mindcraft-CE/Andy-4.2-Air-GGUF", "", nil
	} else {
		log.Info("router: model selected", "model", "Mindcraft-CE/Andy-4.2-GGUF", "tps", tps)
		return "Mindcraft-CE/Andy-4.2-GGUF", "", nil
	}
}
