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
	"os"
	"os/exec"
	"strconv"
)

// * BLOCKING * //
func LlamaServer(modelName string, quantization string, contextLength int, port int) error {
	var cmd *exec.Cmd

	if quantization == "" {
		cmd = exec.Command("llama", "server", "-hf", modelName, "--port", strconv.Itoa(port), "--chat-template", "chatml", "-c", strconv.Itoa(contextLength))
	} else {
		cmd = exec.Command("llama", "server", "-hf", modelName, "--hf-file", quantization, "--port", strconv.Itoa(port), "--chat-template", "chatml", "-c", strconv.Itoa(contextLength))
	}

	// Print output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
