# mcp-pagination-test

A tiny MCP server (streamable HTTP, built on the official
[Go SDK](https://github.com/modelcontextprotocol/go-sdk)) that exposes a set of
no-op tools and **paginates `tools/list`**. Use it to check whether an MCP
client, gateway, or UI follows `nextCursor` or silently stops at the first page.

By default it serves 40 tools (`tool_01` … `tool_40`), 15 per page, so a client
that ignores pagination shows exactly 15 tools.

## Run

```bash
docker build -t mcp-pagination-test .
docker run --rm -p 9123:9123 mcp-pagination-test
```

Or without Docker:

```bash
go run .
```

The MCP endpoint is `http://localhost:9123/mcp`.

## Configuration

| Env var      | Default | Meaning                         |
|--------------|---------|---------------------------------|
| `TOOL_COUNT` | `40`    | Number of tools to expose       |
| `PAGE_SIZE`  | `15`    | Tools per `tools/list` page     |
| `PORT`       | `9123`  | Listen port                     |

## Checking a client

Talking to the server directly, `tools/list` returns three pages:

| Page | Tools                   | `nextCursor` |
|------|-------------------------|--------------|
| 1    | `tool_01` – `tool_15`   | set          |
| 2    | `tool_16` – `tool_30`   | set          |
| 3    | `tool_31` – `tool_40`   | absent       |

A client that handles pagination correctly shows all 40 tools. One that does not
shows only `tool_01` – `tool_15`.
