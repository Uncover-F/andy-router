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
	"os/exec"
	"runtime"
)

const PosixInstallCommand = "curl https://raw.githubusercontent.com/Uncover-F/andy-router/refs/heads/main/cdn/install.sh | sh"
const WindowsInstallCommand = "irm https://raw.githubusercontent.com/Uncover-F/andy-router/refs/heads/main/cdn/install.ps1 | iex"

func installPosix() error {
	cmd := exec.Command("sh", "-c", PosixInstallCommand)
	return cmd.Run()
}

func installWindows() error {
	cmd := exec.Command("powershell", "-Command", WindowsInstallCommand)
	return cmd.Run()
}

func InstallLlama() error {
	var err error
	switch runtime.GOOS {
	case "darwin", "linux", "freebsd", "openbsd":
		err = installPosix()
	case "windows":
		err = installWindows()
	default:
		err = fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return err
}
