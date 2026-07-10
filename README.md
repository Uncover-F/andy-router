# andy-router

**andy-router** is a lightweight OpenAI-compatible router that simplifies using **Andy models** and other models from the **Mindcraft ecosystem**.

Instead of manually choosing models, configuring local runtimes, managing downloads, or setting up API endpoints, andy-router automatically decides the best way to run Andy based on your system capabilities.

## Features

* 🚀 **Automatic model selection**

  * Detects your hardware capabilities and selects the best available Andy model.
  * Runs larger local models when your machine can handle them.
  * Falls back to the Andy API when local inference is not practical.

* 🧠 **OpenAI-compatible API**

  * Works with tools and applications that support the OpenAI API format.
  * Provides a local endpoint:

  ```
  http://127.0.0.1:8000/v1/chat/completions
  ```

* 🔀 **Automatic routing**

  * Clients can simply use:

    ```json
    {
      "model": "auto"
    }
    ```

  * andy-router handles the actual model selection.

* 💻 **Local inference support**

  * Uses llama.cpp for local Andy models.
  * Automatically installs required components when needed.
  * Benchmarks performance before selecting a model.

* ☁️ **Andy API fallback**

  * Systems that cannot efficiently run local models automatically use:

    ```
    https://andy.mindcraft-ce.com/api/
    ```

  * Optional API key support for authenticated requests.

* 🔑 **API key validation**

  * Validates provided Andy API keys before starting.
  * Prevents running with invalid credentials.

---

## Installation

## Downloads

Pre-built binaries are available for Windows, macOS, and Linux.

Download the latest release for your platform:

| Platform | Architecture  | File                            |
| -------- | ------------- | ------------------------------- |
| Windows  | x64           | `andy-router-windows-amd64.exe` |
| Windows  | ARM64         | `andy-router-windows-arm64.exe` |
| Linux    | x64           | `andy-router-linux-amd64`       |
| Linux    | ARM64         | `andy-router-linux-arm64`       |
| macOS    | Intel         | `andy-router-darwin-amd64`      |
| macOS    | Apple Silicon | `andy-router-darwin-arm64`      |

After downloading:

### Linux / macOS

Make the binary executable:

```bash
chmod +x andy-router
```

Run:

```bash
./andy-router
```

### Windows

Run:

```powershell
andy-router-windows-amd64.exe
```

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

```
http://127.0.0.1:8000
```

Example request:

```bash
curl http://127.0.0.1:8000/v1/chat/completions \
-H "Content-Type: application/json" \
-d '{"model":"auto","messages":[{"role":"user","content":"Hello, who are you?"}]}'
```

---

## Options

```
andy-router [options]
```

### `--port`

Change the local API port.

Example:

```bash
andy-router --port 3344
```

Runs:

```
http://127.0.0.1:3344
```

---

### `--key`

Provide an Andy API key.

Example:

```bash
andy-router --key YOUR_API_KEY
```

Without a key, the router uses the public Andy API limits.

Get an API key:

```
https://andy.mindcraft-ce.com/signup
```

---

## Model Selection

When running locally, andy-router automatically chooses a model based on benchmark results.

Example routing:

| Performance     | Model          |
| --------------- | -------------- |
| Low performance | Andy API       |
| Moderate        | Andy-4.2-Micro |
| Good            | Andy-4.2-Air   |
| High            | Andy-4.2       |

Applications do not need to know the exact model name.

Simply use:

```json
{
  "model": "auto"
}
```

---

## OpenAI Client Configuration

andy-router works with OpenAI-compatible clients.

Example:

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://127.0.0.1:8000/v1",
    api_key="unused"
)

response = client.chat.completions.create(
    model="auto",
    messages=[
        {
            "role": "user",
            "content": "Hello!"
        }
    ]
)

print(response.choices[0].message.content)
```

---

## Mindcraft-CE Integration

To run **andy-router** as the LLM backend for a **Mindcraft-CE** bot, you can configure your bot's profile to route request completions to the router's local endpoint.

### Example profile.json
```json
{
    "name": "andy-router",

    "model": {
        "api": "vllm",
        "model": "auto",
        "url": "http://127.0.0.1:8000/v1"
    }
}
```

### Steps to Use

1. **Start the Router**:
   Ensure `andy-router` is running:
   ```bash
   andy-router
   ```
2. **Configure Bot Profile**:
   Copy the json config (or merge its settings) into the `profiles/` directory of your Mindcraft-CE project (e.g., `profiles/profile.json`).
   
   > The `"name"` field in your profile JSON must exactly match the Minecraft username of the account your bot is using. Change `"andy-router"` in the template to your bot's actual in-game name to avoid communication issues.

3. **Launch the Bot**:
   Run the bot from the Mindcraft-CE project directory using the `--profiles` flag:
   ```bash
   node main.js --profiles ./profiles/profile.json
   ```
---

## Why use andy-router?

Using AI models often requires knowing:

* Which model should I run?
* Can my hardware handle it?
* Should I use local inference or an API?
* Which endpoint format does my application expect?

andy-router removes those decisions.

Applications connect to one simple OpenAI-compatible endpoint, and the router handles the rest.

---

## Supported Models

Currently focused on Andy models:

* `Mindcraft-CE/Andy-4.2-Micro-GGUF`
* `Mindcraft-CE/Andy-4.2-Air-GGUF`
* `Mindcraft-CE/Andy-4.2-GGUF`

More models may be supported in the future.

---

## License

See the repository license for details.

---

Made by **@Uncover-F**

Mindcraft CE community support:

https://discord.gg/mindcraft-ce
