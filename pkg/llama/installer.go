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
