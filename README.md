<div align="center">

```text
   ██████╗██╗     ██╗██╗     ██╗     ███╗   ███╗
  ██╔════╝██║     ██║██║     ██║     ████╗ ████║
  ██║     ██║     ██║██║     ██║     ██╔████╔██║
  ██║     ██║     ██║██║     ██║     ██║╚██╔╝██║
  ╚██████╗███████╗██║███████╗███████╗██║ ╚═╝ ██║
   ╚═════╝╚══════╝╚═╝╚══════╝╚══════╝╚═╝     ╚═╝
             LOCAL LLM TERMINAL
```

# CLILLM

**A fast, lightweight terminal-based AI assistant built with Go and powered by Ollama.**

[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Python](https://img.shields.io/badge/Python-3.8+-3776AB?style=flat-square&logo=python&logoColor=white)](https://www.python.org/)
[![Ollama](https://img.shields.io/badge/Ollama-Local_LLM-000000?style=flat-square)](https://ollama.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

</div>

---

## Features

- **100% Local & Private** — Runs entirely on your machine via the local [Ollama](https://ollama.com/) HTTP API (`localhost:11434`).
- **Interactive Model Selector** — Automatically detects all installed Ollama models and lets you pick one at startup.
- **Document Upload & Q&A (`/upload`)** — Parse and query local `.pdf`, `.docx`, `.txt`, and `.md` files directly from the prompt.
- **Session Context Memory** — Maintains a sliding conversation history window so the model remembers previous questions and uploaded documents.
- **Graceful Signal Handling** — Clean exit on `Ctrl+C` (`SIGINT`).

---

## Requirements

- **[Go](https://go.dev/dl/)** (1.20 or higher)
- **[Python 3](https://www.python.org/)** (`pypdf` and `python-docx` for document parsing)
- **[Ollama](https://ollama.com/)** with at least one pulled model (e.g., `llama3.2` or `qwen2.5:0.5b`)

---

## Quick Start

### 1. Clone the repository
```bash
git clone [https://github.com/tracebyte8/CLILLM.git](https://github.com/tracebyte8/CLILLM.git)
cd CLILLM
```

### 2. Automated Setup (Recommended)
Run the setup script to install Python document parsers, verify Go & Ollama, and pull a default model if needed:
```bash
chmod +x setup.sh
./setup.sh
```

### 3. Manual Setup
If you prefer setting up manually:
```bash
# Install Python dependencies for document upload
python3 -m pip install --user --break-system-packages pypdf python-docx

# Start Ollama and pull a model
ollama serve &
ollama pull llama3.2
```

### 4. Run CLILLM
```bash
go run .
```

---

## Usage & Commands

Once inside the terminal prompt (`You >`), you can chat naturally or use built-in slash commands:

| Command | Description | Example |
| :--- | :--- | :--- |
| `/upload <path> [prompt]` | Extracts text from a `.pdf`, `.docx`, `.txt`, or `.md` file and sends it to the LLM with an optional question. | `/upload ./report.pdf Summarize the key findings` |
| `/upload <path>` | Uploads and automatically summarizes the document. | `/upload ./notes.txt` |
| `/clear` | Clears the current conversation history memory. | `/clear` |
| `Ctrl + C` | Gracefully exits the terminal session. | `^C` |

---

## Project Structure

```text
CLILLM/
├── main.go          # Core CLI loop, Ollama HTTP client, model selector, and history state
├── readoc.py        # Document extraction helper (.pdf, .docx, .txt, .md)
├── setup.sh         # Automated dependency and environment installer
├── go.mod           # Go module definition
└── README.md        # Project documentation
```

---

## License

Distributed under the MIT License.