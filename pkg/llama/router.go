// MIT License

// Copyright (c) 2026 Uncover-F

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package llama

import (
	"fmt"

	"github.com/charmbracelet/log"
)

// Global constants (across all functions)
const MinimumMemory uint64 = 6 // minimum memory required to use llama.cpp

func SelectModel() (selectedModel string, selectedQuant string, selectedContextLength int, err error) {
	// selectedModel "" (zero value) = use andyAPI
	// selectedQuant "" (zero value) = let llama.cpp decide (do not specify)

	// Perform benchmark
	tps, err := Benchmark()
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to benchmark performance: %w", err)
	}

	// Route between models
	// TODO: Add logic for routing between different models/quantizations based on benchmark results
	if tps < 300 {
		log.Info("router: model selected", "model", "none", "tps", tps)
		return "", "", 0, nil
	} else if tps < 700 {
		log.Info("router: model selected", "model", "Mindcraft-CE/Andy-4.2-Micro-GGUF", "tps", tps)
		return "Mindcraft-CE/Andy-4.2-Micro-GGUF", "", 32000, nil
	} else if tps < 1200 {
		log.Info("router: model selected", "model", "Mindcraft-CE/Andy-4.2-Air-GGUF", "tps", tps)
		return "Mindcraft-CE/Andy-4.2-Air-GGUF", "", 32000, nil
	} else {
		log.Info("router: model selected", "model", "Mindcraft-CE/Andy-4.2-GGUF", "tps", tps)
		return "Mindcraft-CE/Andy-4.2-GGUF", "", 32000, nil
	}
}
