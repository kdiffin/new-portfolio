# new-portfolio

Macwright-style, full-section publishing site scaffold using:



- Go stdlib (`net/http`, `html/template`)
- Tailwind CSS


[ai generated code walkthrough](./CODE_WALKTHROUGH.md) \
[how templating works (inspired by lets go by A. Edwards)](./how-templating-works.md)
## Run

1. Install dependencies:
   - `npm install`
2. Build CSS:
   - `npm run css:build`
3. Start server:
   - `go run ./cmd/web`
4. Open:
   - `http://localhost:4000`

## Template architecture (Let’s Go style)

- `ui/html/base.tmpl` defines the shared page shell.
- `ui/html/partials/*.tmpl` define reusable snippets (`nav`, `footer`, etc.).
- `ui/html/pages/*.tmpl` define page-level `title` and `main` templates.
- The server executes `base` for each response.

## Notes

- Current content is placeholder data from an in-memory store.
- Next implementation step is file-backed content loading for your original writing.
