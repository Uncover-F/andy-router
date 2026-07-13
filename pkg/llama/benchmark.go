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
	"encoding/json"
	"fmt"
	"os/exec"
)

type BenchResult struct {
	NPrompt int     `json:"n_prompt"`
	NGen    int     `json:"n_gen"`
	AvgTS   float64 `json:"avg_ts"`
	Model   string  `json:"model_filename"`
}

func Benchmark() (int, error) {
	cmd := exec.Command(
		"llama", "bench",
		"-hf", "SupraLabs/Supra-Router-51M-gguf",
		"--hf-file", "Supra-Router-51M-Q1_0.gguf",
		"--output", "json",
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Parse results
	var results []BenchResult
	if err := json.Unmarshal(output, &results); err != nil {
		return 0, fmt.Errorf("failed to unmarshal benchmark results: %w", err)
	}
	for _, r := range results {
		if r.NGen > 0 {
			return int(r.AvgTS), err
		}
	}

	return 0, fmt.Errorf("no generation benchmark results found")
}
