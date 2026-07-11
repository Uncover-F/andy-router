package llama

import (
	"os"
	"os/exec"
	"strconv"
)

// * BLOCKING * //
func LlamaServer(modelName string, quantization string, port int) error {
	var cmd *exec.Cmd

	if quantization == "" {
		cmd = exec.Command("llama", "server", "-hf", modelName, "--port", strconv.Itoa(port), "--chat-template", "chatml")
	} else {
		cmd = exec.Command("llama", "server", "-hf", modelName, "--hf-file", quantization, "--port", strconv.Itoa(port), "--chat-template", "chatml")
	}

	// Print output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
