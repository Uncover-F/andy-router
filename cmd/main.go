package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/Uncover-F/andy-router/pkg/api"
	"github.com/Uncover-F/andy-router/pkg/llama"
	"github.com/charmbracelet/log"
	"github.com/shirou/gopsutil/v4/mem"
)

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

const Version = "1.1.0"

var port int = 8000
var key string = ""
var model string = ""
var forceAPI bool = false

func main() {
	// Accept cmd input flags
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--help", "-h":
			printHelp()
			return
		case "--port":
			if i+1 >= len(os.Args) {
				log.Fatal("missing value for --port: expected a port number (1-65535)")
			}
			p, err := strconv.Atoi(os.Args[i+1])
			if err != nil || p < 1 || p > 65535 {
				log.Fatalf("invalid port %q: must be a number between 1 and 65535", os.Args[i+1])
			}
			port = p
			i++
		case "--key":
			if i+1 >= len(os.Args) {
				log.Fatal("missing value for --key: expected a key string")
			}
			if len(os.Args[i+1]) > 0 && os.Args[i+1][0] == '-' {
				log.Fatalf("invalid value for --key: %q looks like an option", os.Args[i+1])
			}
			key = os.Args[i+1]
			i++
		case "--api":
			forceAPI = true
		case "--model", "-m":
			if i+1 >= len(os.Args) {
				log.Fatal("missing value for --model: expected a model name")
			}
			val := os.Args[i+1]
			switch val {
			case "Andy-4.2-Micro":
				model = "Mindcraft-CE/Andy-4.2-Micro-GGUF"
			case "Andy-4.2-Air":
				model = "Mindcraft-CE/Andy-4.2-Air-GGUF"
			case "Andy-4.2":
				model = "Mindcraft-CE/Andy-4.2-GGUF"
			default:
				log.Fatalf("invalid model name %q: must be one of: Andy-4.2-Micro, Andy-4.2-Air, Andy-4.2", val)
			}
			i++
		default:
			log.Error(Red + "unknown option " + os.Args[i] + Reset)
			printTruncatedHelp()
			os.Exit(1)
		}
	}
	if forceAPI && model != "" {
		log.Fatal("conflict: cannot specify both --api and --model. Selected models (micro, air, ye) are only supported for local inference.")
	}

	log.Info("starting andy-router...")

	// Handle --api
	if forceAPI {
		log.Info("force API enabled, bypassing router")
		runAndy()
		return
	}

	// Verify the client is capable of running llama.cpp
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	if model == "" {
		if v.Total/1024/1024/1024 > 6 {
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
			log.Error("failed to install llama.cpp, falling back to andyAPI", "error", err)
			runAndy()
			return
		}
	}

	// Handle user-specified models
	if model != "" {
		log.Info("using user-specified model", "model", model)
		err = llama.LlamaServer(model, port)
		if err != nil {
			log.Error("failed to start llama server, falling back to andyAPI instead", "error", err)
			runAndy()
		}
		return
	}

	// Benchmark system performance & select model
	tps, err := llama.Benchmark()
	if err != nil {
		log.Error("failed to benchmark performance, falling back to andyAPI", "error", err)
		runAndy()
		return
	} else {
		log.Info("benchmark finished", "tps", tps)
	}
	selectedModel := llama.SelectModel(tps)
	if selectedModel == "" {
		log.Warn("weak performance detected, falling back to andyAPI")
		runAndy()
		return
	}

	// Start llama server
	err = llama.LlamaServer(selectedModel, port)
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
	// andy-router-v1.1.0 - made by @Uncover-F
	// discord support: https://discord.gg/mindcraft-ce

	// Usage:
	// andy-router [--port PORT] [--key KEY] [--api] [--model MODEL] [--help]

	// Options:
	// --port PORT         Local port to bind to (default: 8000)
	// --key KEY           Optional Andy API key
	// --api               Force using the Andy API regardless of compute
	// --model, -m MODEL   Specify a model to use (bypasses auto-detection)
	// --help, -h          Show this help message

	fmt.Println(Green + "andy-router-" + Version + " - made by @Uncover-F" + Reset)
	fmt.Println(Green + "discord support: https://discord.gg/mindcraft-ce" + Reset)
	fmt.Println("")
	fmt.Println(Blue + "./andy-router [--port PORT] [--key KEY] [--api] [--model MODEL] [--help]" + Reset)
	fmt.Println("")
	fmt.Println("--port PORT" + Yellow + "         Local port to bind to (default: 8000)" + Reset)
	fmt.Println("--key KEY" + Yellow + "           Optional Andy API key" + Reset)
	fmt.Println("--api" + Yellow + "               Force using the Andy API regardless of compute" + Reset)
	fmt.Println("--model, -m MODEL" + Yellow + "   Specify a model to use (bypasses auto-detection)" + Reset)
	fmt.Println("--help, -h" + Yellow + "          Show this help message" + Reset)
}

func printTruncatedHelp() {
	// andy-router-v1.1.0 - made by @Uncover-F
	// discord support: https://discord.gg/mindcraft-ce

	// Usage:
	// andy-router [--port PORT] [--key KEY] [--api] [--model MODEL] [--help]

	// Options:
	// --help, -h          Show help message

	fmt.Println("")
	fmt.Println(Green + "andy-router-" + Version + " - made by @Uncover-F" + Reset)
	fmt.Println(Green + "discord support: https://discord.gg/mindcraft-ce" + Reset)
	fmt.Println("")
	fmt.Println(Blue + "./andy-router [--port PORT] [--key KEY] [--api] [--model MODEL] [--help]" + Reset)
	fmt.Println("")
	fmt.Println("--help, -h" + Yellow + "          Show help message" + Reset)
}
