// mcp-pagination-test is a minimal streamable-HTTP MCP server that exposes a
// configurable number of no-op tools and paginates tools/list. It is useful for
// checking whether an MCP client or gateway follows nextCursor.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type noArgs struct{}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			log.Fatalf("invalid %s=%q: must be a positive integer", key, v)
		}
		return n
	}
	return def
}

func main() {
	var (
		toolCount = envInt("TOOL_COUNT", 40)
		pageSize  = envInt("PAGE_SIZE", 15)
		port      = envInt("PORT", 9123)
	)

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-pagination-test", Version: "1.0.0"}, &mcp.ServerOptions{PageSize: pageSize})
	for i := 1; i <= toolCount; i++ {
		name := fmt.Sprintf("tool_%02d", i)
		mcp.AddTool(server, &mcp.Tool{Name: name, Description: fmt.Sprintf("Test tool number %d", i)},
			func(context.Context, *mcp.CallToolRequest, noArgs) (*mcp.CallToolResult, any, error) {
				return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "called " + name}}}, nil, nil
			})
	}

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, nil)
	log.Printf("serving %d tools, %d per page, on :%d/mcp", toolCount, pageSize, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s session=%s", r.Method, r.URL.Path, r.Header.Get("Mcp-Session-Id"))
		handler.ServeHTTP(w, r)
	})))
}
