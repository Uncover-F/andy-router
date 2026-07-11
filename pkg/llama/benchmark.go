package llama

import (
	"encoding/json"
	"errors"
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
		return 0, err
	}
	for _, r := range results {
		if r.NGen > 0 {
			return int(r.AvgTS), err
		}
	}

	return 0, errors.New("no generation benchmark results found")
}
