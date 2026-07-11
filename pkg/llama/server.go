package llama

import (
	"os"
	"os/exec"
	"strconv"
)

// * BLOCKING * //
func LlamaServer(modelName string, port int) error {
	cmd := exec.Command("llama", "server", "-hf", modelName, "--port", strconv.Itoa(port), "--chat-template", "chatml")

	// Print output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
