package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	"github.com/Uncover-F/andy-router/pkg/installer"
	"github.com/Uncover-F/andy-router/pkg/utils"
	"github.com/charmbracelet/log"
	"github.com/shirou/gopsutil/v4/mem"
)

var port int = 8000
var key string = ""
var model string = ""
var forceAPI bool = false

func main() {
	// Accept flags
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
		}
	}

	if forceAPI && model != "" {
		log.Fatal("conflict: cannot specify both --api and --model. Selected models (micro, air, ye) are only supported for local inference.")
	}

	log.Info("starting andy-router...")

	if forceAPI {
		log.Info("force API enabled, bypassing compute checks")
		andyAPI(port, key)
		return
	}

	// Check if client is capable of using llama.cpp
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	if v.Total/1024/1024/1024 > 6 {
		llamaRuntime(port)
	} else {
		log.Warn("insufficient memory, falling back to andyAPI", "memory", v.Total/1024/1024/1024)
		andyAPI(port, key)
	}

}

func llamaRuntime(port int) {
	err := exec.Command("llama").Run()
	if err != nil {
		log.Info("installing llama.cpp... (this may take a while)")
		err = installer.InstallLlama()
		if err != nil {
			log.Error("failed to install llama.cpp, falling back to andyAPI", "error", err)
			andyAPI(port, key)
			return
		}
	}

	if model != "" {
		log.Info("using user-specified model", "model", model)
		llamaServer(model, port)
		return
	}

	tps, err := utils.Benchmark()
	if err != nil {
		log.Error("failed to benchmark performance, falling back to andyAPI", "error", err)
		andyAPI(port, key)
		return
	} else {
		log.Info("benchmark results", "tps", tps)
	}

	if tps < 300 {
		log.Warn("weak performance detected, falling back to andyAPI")
		andyAPI(port, key)
		return
	} else if tps < 700 {
		log.Info("moderate performance detected, running Andy-4.2-Micro")
		llamaServer("Mindcraft-CE/Andy-4.2-Micro-GGUF", port)
	} else if tps < 1200 {
		log.Info("good performance detected, running Andy-4.2-Air")
		llamaServer("Mindcraft-CE/Andy-4.2-Air-GGUF", port)
	} else {
		log.Info("excellent performance detected, running Andy-4.2")
		llamaServer("Mindcraft-CE/Andy-4.2-GGUF", port)
	}

}

func andyAPI(port int, key string) {
	modelLog := "auto"
	if model != "" {
		modelLog = model
	}
	log.Info("starting andyAPI...", "model", modelLog, "port", port)
	if key == "" {
		log.Warn("API key not provided, daily limits will apply. get an API key at: https://andy.mindcraft-ce.com/signup")
	} else {
		isValid, err := verifyAndyKey(key)
		if err != nil {
			log.Fatal("failed to verify Andy API key: ", "err", err)
		}
		if !isValid {
			log.Fatal("invalid Andy API key")
		}
	}

	target, err := url.Parse("https://andy.mindcraft-ce.com/api/")
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path + req.URL.Path
		req.Host = target.Host

		req.Header.Set("Content-Type", "application/json")

		if key != "" {
			req.Header.Set("Authorization", "Bearer "+key)
		}
	}
	server := &http.Server{
		Addr:    "127.0.0.1:" + strconv.Itoa(port),
		Handler: proxy,
	}

	log.Info("API server running on http://127.0.0.1:" + strconv.Itoa(port))
	log.Fatal(server.ListenAndServe())
}

// HELPER FUNCTIONS

func llamaServer(modelName string, port int) {
	log.Info("starting llama server... (this may take a while)", "model", modelName, "port", port)

	cmd := exec.Command("llama", "server", "-hf", modelName, "--port", strconv.Itoa(port), "--chat-template", "chatml")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Error("failed to start llama server, falling back to andyAPI", "error", err)
		andyAPI(port, key)
		return
	}
}

func verifyAndyKey(key string) (bool, error) {
	reqBody := []byte(`{
		"model": "auto",
		"messages": [
			{
				"role": "user",
				"content": "ping"
			}
		],
		"max_tokens": 1
	}`)

	req, err := http.NewRequest(
		"POST",
		"https://andy.mindcraft-ce.com/api/v1/chat/completions",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusForbidden {
		return false, errors.New("invalid Andy API key: " + string(body))
	}

	// Some APIs return 200 but include an error message
	if bytes.Contains(bytes.ToLower(body), []byte("invalid api key")) ||
		bytes.Contains(bytes.ToLower(body), []byte("invalid key")) ||
		bytes.Contains(bytes.ToLower(body), []byte("unauthorized")) {
		return false, errors.New("invalid Andy API key: " + string(body))
	}

	return true, nil
}

func printHelp() {
	fmt.Println(`
andy-router-v1.1.0 - made by @Uncover-F
discord support: https://discord.gg/mindcraft-ce

example CURL request: curl http://127.0.0.1:8000/v1/chat/completions -H "Content-Type: application/json" -d '{"model":"auto","messages":[{"role":"user","content":"Hello, who are you?"}]}'

Usage:
andy-router [--port PORT] [--key KEY] [--api] [--model MODEL] [--help]

Options:
--port PORT         Local port to bind to (default: 8000)
--key KEY           Optional Andy API key
--api               Force using the Andy API regardless of compute
--model, -m MODEL   Specify a model to use (bypasses auto-detection)
--help, -h          Show this help message`)
}
