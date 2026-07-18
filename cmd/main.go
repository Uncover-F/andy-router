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

package main

import (
	"fmt"
	"os/exec"

	"github.com/Uncover-F/andy-router/pkg/api"
	"github.com/Uncover-F/andy-router/pkg/llama"
	"github.com/charmbracelet/log"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/pflag"
)

// Global variables (across all functions)
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

const Version = "1.2.3"

// Global variables (inside main.go)
var port int = 8000
var key string
var model string
var forceAPI bool
var showVersion bool

func main() {
	// Prase CLI flags
	pflag.IntVar(&port, "port", 8000, "local port")
	pflag.StringVar(&key, "key", "", "andyPI key")
	pflag.BoolVar(&forceAPI, "api", false, "force andyAPI")
	pflag.StringVarP(&model, "model", "m", "", "model to use")
	pflag.BoolVarP(&showVersion, "version", "v", false, "show version information")

	pflag.Usage = printHelp
	pflag.Parse()

	// Handle/Validate CLI Inputs
	if port < 1 || port > 65535 {
		log.Fatalf("invalid port %v: must be a number between 1 and 65535", port)
	}

	if showVersion {
		fmt.Println(Green + "andy-router version " + Version + Reset)
		return
	}

	switch model {
	case "":
		// no model specified; auto-detect below
	case "Andy-4.2-Micro":
		model = "Mindcraft-CE/Andy-4.2-Micro-GGUF"
	case "Andy-4.2-Air":
		model = "Mindcraft-CE/Andy-4.2-Air-GGUF"
	case "Andy-4.2":
		model = "Mindcraft-CE/Andy-4.2-GGUF"
	default:
		log.Fatalf("invalid model name %q: must be one of: Andy-4.2-Micro, Andy-4.2-Air, Andy-4.2", model)
	}

	if forceAPI && model != "" {
		log.Fatal("conflict: cannot specify both --api and --model. Selected models (Andy-4.2-Micro, Andy-4.2-Air, Andy-4.2) are only supported for local inference.")
	}

	log.Info("starting andy-router...")

	// Handle --api
	if forceAPI {
		log.Info("force API enabled, bypassing router")
		runAndy()
		return
	}

	// Verify that the client is capable of running llama.cpp
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	if model == "" {
		if v.Total >= llama.MinimumMemory*1024*1024*1024 {
			runLlama()
		} else {
			log.Warn("insufficient memory, falling back to andyAPI instead", "memory", v.Total/1024/1024/1024)
			runAndy()
		}
	} else {
		log.Warn("memory checks bypassed, llama.cpp may run out of memory")
		runLlama()
	}
}

func runLlama() {
	// Ensure llama.cpp is installed
	_, err := exec.LookPath("llama")
	if err != nil {
		log.Info("installing llama.cpp... (this may take a while)")
		err = llama.InstallLlama()
		if err != nil {
			log.Error("failed to install llama.cpp, falling back to andyAPI instead", "error", err)
			runAndy()
			return
		}
	}

	// Handle user-specified models
	if model != "" {
		log.Info("using user-specified model", "model", model)
		err = llama.LlamaServer(model, "", 32000, port)
		if err != nil {
			log.Error("failed to start llama server, falling back to andyAPI instead", "error", err)
			runAndy()
		}
		return
	}

	// Select model using router
	selectedModel, selectedQuant, selectedContextLength, err := llama.SelectModel()
	if err != nil {
		log.Error("failed to select model, falling back to andyAPI instead", "error", err)
		runAndy()
		return
	}
	if selectedModel == "" {
		log.Warn("weak performance detected, falling back to andyAPI instead")
		runAndy()
		return
	}

	// Start llama server
	log.Info("starting llama server...", "model", selectedModel)
	err = llama.LlamaServer(selectedModel, selectedQuant, selectedContextLength, port)
	if err != nil {
		log.Error("failed to start llama server, falling back to andyAPI instead", "error", err)
		runAndy()
		return
	}
}

func runAndy() {
	// Validate API key
	if key == "" {
		log.Warn(Red + "API key not provided, usage limits will apply. get an API key at: https://andy.mindcraft-ce.com/signup" + Reset)
	} else {
		isValid, err := api.VerifyAndyKey(key)
		if err != nil {
			log.Fatal("failed to verify andyAPI key", "err", err)
		}
		if !isValid {
			log.Fatal("invalid andyAPI key")
		}
	}

	// Start andyAPI proxy
	err := api.AndyProxy(key, port)
	if err != nil {
		log.Fatal("failed to start andyAPI", "err", err)
	}
}

func printHelp() {
	fmt.Println(Green + "andy-router-" + Version + " - made by @Uncover-F" + Reset)
	fmt.Println(Green + "discord support: https://discord.gg/mindcraft-ce" + Reset)
	fmt.Println("")
	fmt.Println(Blue + "./andy-router [--port PORT] [--key KEY] [--api] [--model MODEL] [--help]" + Reset)
	fmt.Println("")
	fmt.Println("--port PORT" + Yellow + "         Local port to bind to (default: 8000)" + Reset)
	fmt.Println("--key KEY" + Yellow + "           Optional Andy API key" + Reset)
	fmt.Println("--api" + Yellow + "               Force using the Andy API regardless of compute" + Reset)
	fmt.Println("--model, -m MODEL" + Yellow + "   Specify a model to use (bypasses auto-detection)" + Reset)
	fmt.Println("--version, -v" + Yellow + "          Show version information" + Reset)
	fmt.Println("--help, -h" + Yellow + "          Show help message" + Reset)
}
