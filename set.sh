#!/usr/bin/env bash
set -e

echo "=========================================="
echo "       CLILLM - Project Setup             "
echo "=========================================="

# 1. Check Go installation
if ! command -v go &> /dev/null; then
    echo "[!] Go is not installed. Please install Go (https://go.dev/dl/) first."
    exit 1
fi
echo "[+] Go found: $(go version)"

# 2. Check Python3 & pip installation
if ! command -v python3 &> /dev/null; then
    echo "[!] Python3 is not installed. Installing python3 and pip..."
    sudo apt-get update && sudo apt-get install -y python3 python3-pip
fi
echo "[+] Python3 found: $(python3 --version)"

# 3. Install Python packages directly for default 'python3' (no venv activation needed!)
echo "[*] Installing Python dependencies (pypdf, python-docx) for default python3..."
python3 -m pip install --user --break-system-packages pypdf python-docx 2>/dev/null || \
python3 -m pip install --user pypdf python-docx

# 4. Check Ollama installation
if ! command -v ollama &> /dev/null; then
    echo "[*] Ollama not found. Installing Ollama..."
    curl -fsSL https://ollama.com/install.sh | sh
else
    echo "[+] Ollama found: $(ollama --version)"
fi

# 5. Start Ollama if it isn't running yet
if ! pgrep -x "ollama" > /dev/null; then
    echo "[*] Starting Ollama service..."
    sudo systemctl start ollama 2>/dev/null || (ollama serve >/dev/null 2>&1 & sleep 2)
fi

# 6. Pull a default model if none are installed
if [ -z "$(ollama list 2>/dev/null | tail -n +2)" ]; then
    echo "[*] Pulling default lightweight model (qwen2.5:0.5b)..."
    ollama pull qwen2.5:0.5b
fi

# 7. Initialize Go module if go.mod is missing
if [ ! -f "go.mod" ]; then
    echo "[*] Initializing Go module..."
    go mod init github.com/tracebyte8/CLILLM
fi

echo "=========================================="
echo " Setup Complete!"
echo " Now you can simply run:"
echo "   go run ."
echo "=========================================="