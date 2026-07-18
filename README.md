# andy-router 🧠

![Version Badge](https://img.shields.io/badge/version-1.2.3-blue) ![License Badge](https://img.shields.io/badge/license-MIT-yellow) ![Language Badge](https://img.shields.io/badge/language-Go-green)

**andy-router** is a lightweight, OpenAI-compatible model router for **Andy models** and the **Mindcraft ecosystem**.

Instead of manually selecting models, configuring runtimes, downloading files, or managing API endpoints, andy-router automatically determines the best way to run Andy based on your system capabilities.

It provides a single OpenAI-compatible endpoint that applications can connect to while handling model selection, local inference, and API fallback automatically.

---

## Features

### 🚀 Automatic model selection

* Detects your hardware capabilities.
* Selects the best available Andy model for your system.
* Uses larger local models when your machine can handle them.
* Falls back to the Andy API when local inference is not practical.

### 🧠 OpenAI-compatible API

Works with applications that support the OpenAI API format.

Local endpoint:

```text
http://127.0.0.1:8000/v1/chat/completions
```

---

## Installation ⬇️

### One-line installer

Install automatically:

```bash
curl -fsSL https://raw.githubusercontent.com/Uncover-F/andy-router/main/cdn/rinstall.sh | sh
```

The installer will:

* Detect your operating system and architecture.
* Download the correct binary.
* Install it to `~/.local/bin`.
* Configure instructions for adding it to your `PATH` if needed.

> Windows PowerShell installer support is planned.

---

## Manual Downloads

Pre-built binaries are available for:

| Platform | Architecture  | File                            |
| -------- | ------------- | ------------------------------- |
| Windows  | x64           | `andy-router-windows-amd64.exe` |
| Windows  | ARM64         | `andy-router-windows-arm64.exe` |
| Linux    | x64           | `andy-router-linux-amd64`       |
| Linux    | ARM64         | `andy-router-linux-arm64`       |
| macOS    | Intel         | `andy-router-darwin-amd64`      |
| macOS    | Apple Silicon | `andy-router-darwin-arm64`      |

### Linux / macOS

Make the binary executable:

```bash
chmod +x andy-router-your-platform-your-architecture
```

Run:

```bash
./andy-router-your-platform-your-architecture
```

### Windows

Run:

```powershell
andy-router-windows-your-architecture.exe
```

---

## Requirements

* Linux, macOS, or Windows
* `curl` (only required for the installer)
* A supported local inference backend for local models
* An Andy API key is optional

---

## Building From Source

Clone the repository:

```bash
git clone https://github.com/Uncover-F/andy-router
cd andy-router
```

Build:

```bash
go build -o andy-router ./cmd
```

Run:

```bash
./andy-router
```

---

## Usage

Start the router:

```bash
andy-router
```

The API will be available at:

```text
http://127.0.0.1:8000
```

Example request:

```bash
curl http://127.0.0.1:8000/v1/chat/completions \
-H "Content-Type: application/json" \
-d '{"model":"auto","messages":[{"role":"user","content":"Hello, who are you?"}]}'
```

---

# Command Options

```text
andy-router [options]
```

## `--port`

Change the local API port.

Example:

```bash
andy-router --port 3344
```

Runs:

```text
http://127.0.0.1:3344
```

---

## `--key`

Provide an Andy API key.

Example:

```bash
andy-router --key YOUR_API_KEY
```

Without a key, the router uses the public Andy API limits.

Get an API key:

```text
https://andy.mindcraft-ce.com/signup
```

---

## `--api`

Force API mode.

This disables local model detection and always routes requests through the Andy API.

Example:

```bash
andy-router --api
```

---

## `--model` / `-m`

Select a specific local model instead of automatic selection.

Supported models:

| Name             | Resolved Model                     |
| ---------------- | ---------------------------------- |
| `Andy-4.2-Micro` | `Mindcraft-CE/Andy-4.2-Micro-GGUF` |
| `Andy-4.2-Air`   | `Mindcraft-CE/Andy-4.2-Air-GGUF`   |
| `Andy-4.2`       | `Mindcraft-CE/Andy-4.2-GGUF`       |

Example:

```bash
andy-router --model Andy-4.2-Air
```

or:

```bash
andy-router -m Andy-4.2-Air
```

> [!IMPORTANT]
> `--model` cannot be combined with `--api`. Local models are only available through local inference.

---

# Model Selection

When running locally, andy-router benchmarks your system and automatically chooses the most suitable model.

Example:

| System Performance | Selected Model |
| ------------------ | -------------- |
| Low performance    | Andy API       |
| Moderate           | Andy-4.2-Micro |
| Good               | Andy-4.2-Air   |
| High               | Andy-4.2       |

Applications should simply request:

```json
{
  "model": "auto"
}
```

Currently supported models:

* `Mindcraft-CE/Andy-4.2-Micro-GGUF`
* `Mindcraft-CE/Andy-4.2-Air-GGUF`
* `Mindcraft-CE/Andy-4.2-GGUF`

More models may be added in future releases.

---

# Mindcraft-CE Integration

andy-router can be used as the LLM backend for a **Mindcraft-CE** bot.

Example `profile.json`:

```json
{
    "name": "andy",
    "model": {
        "api": "vllm",
        "model": "auto",
        "url": "http://127.0.0.1:8000/v1"
    }
}
```

## Setup

### 1. Start the router

```bash
andy-router
```

### 2. Configure your bot profile

Place the profile inside your Mindcraft-CE `profiles/` directory.

The `"name"` field must exactly match your Minecraft username.

Example:

```json
"name": "YourMinecraftUsername"
```

### 3. Launch Mindcraft-CE

```bash
node main.js --profiles ./profiles/profile.json
```

---

# Why use andy-router?

Using AI models usually requires answering:

* Which model should I run?
* Can my hardware handle it?
* Should I use local inference or an API?
* Which endpoint format does my application need?

andy-router removes those decisions.

Applications connect to one OpenAI-compatible endpoint, and the router handles the rest.

---

# License

andy-router is released under the MIT License.

See the repository license for details.

---

Made by **@Uncover-F**

Mindcraft CE community support:

https://discord.gg/mindcraft-ce
