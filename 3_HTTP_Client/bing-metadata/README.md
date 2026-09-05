## Bing Document Metadata Scraper

Searches a target domain for public Office files (`.docx`, `.xlsx`, `.pptx`) 
via Bing, downloads each match, and extracts author/editor/software metadata 
from the embedded XML — a classic OSINT technique (see: FOCA) for gathering 
real employee names and software versions with zero direct interaction with 
the target.

### Usage
```bash
go run main.go <domain>
```
Example: `go run main.go mit.edu` — checks all three file types automatically.

### Note

Only Open XML formats are supported (`.docx`/`.xlsx`/`.pptx`) — they're ZIP 
archives internally. PDF and legacy `.doc`/`.xls`/`.ppt` use different binary 
formats and would need a separate parser.