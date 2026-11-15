import os
import sys

# Ensure UTF-8 output when piped to Go (prevents Windows cp1252 crashes)
if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

try:
    from docx import Document
    from pypdf import PdfReader
except ImportError as e:
    print(f"Missing dependency: {e}. Run: pip install pypdf python-docx", file=sys.stderr)
    sys.exit(1)


def load_document(path: str) -> str:
    ext = os.path.splitext(path)[1].lower()

    if ext == ".pdf":
        reader = PdfReader(path)
        text = "\n".join(page.extract_text() or "" for page in reader.pages)

    elif ext == ".docx":
        doc = Document(path)
        parts = [p.text for p in doc.paragraphs if p.text.strip()]
        # Also extract text from tables
        for table in doc.tables:
            for row in table.rows:
                row_text = " | ".join(cell.text.strip() for cell in row.cells if cell.text.strip())
                if row_text:
                    parts.append(row_text)
        text = "\n".join(parts)

    elif ext in (".txt", ".md"):
        with open(path, "r", encoding="utf-8", errors="replace") as f:
            text = f.read()

    else:
        raise ValueError(f"Unsupported file type: {ext or 'none'}")

    if not text.strip():
        raise ValueError("Document contains no extractable text (it may be empty or a scanned image).")

    return text


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Error: No file path provided", file=sys.stderr)
        sys.exit(1)

    file_path = sys.argv[1]
    try:
        content = load_document(file_path)
        print(content)
    except Exception as e:
        print(f"Error reading document: {e}", file=sys.stderr)
        sys.exit(1)