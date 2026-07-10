package main

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/Uncover-F/andy-router/pkg/installer"
	"github.com/charmbracelet/log"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {
	log.Info("starting andy-router...")

	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	if v.Available/1024/1024/1024 > 6 {
		llamaRuntime()
	} else {
		andyAPI()
	}

	// Accept Ports
	if os.Args[1] == "--port" {
		if len(os.Args) > 2 {
			port, err := strconv.Atoi(os.Args[2])
			if err != nil || port < 1 || port > 65535 {
				log.Fatalf("invalid port %q: must be a number between 1 and 65535", os.Args[2])
			}
		} else {
			log.Fatal("missing value for --port: expected a port number (1-65535)")
		}
	}

}

func llamaRuntime() {
	err := exec.Command("llama").Run()
	if err != nil {
		log.Info("installing llama.cpp...")
		err = installer.InstallLlama()
		if err != nil {
			log.Error("failed to install llama.cpp, falling back to andyAPI", "error", err)
			andyAPI()
		}
	}
}

func andyAPI() {

}
