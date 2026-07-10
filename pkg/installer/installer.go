package installer

import (
	"errors"
	"os/exec"
	"runtime"
)

const PosixInstallCommand = "curl https://llama.app/install.sh | sh"
const WindowsInstallCommand = "irm https://llama.app/install.ps1 | iex"

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
		err = errors.New("unsupported OS: " + runtime.GOOS)
	}

	return err
}
