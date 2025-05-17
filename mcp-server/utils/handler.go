package utils

import (
	"context"
	"fmt"
	"runtime"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func HandleError(handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return ErrorGuard(handler)
}

func ErrorGuard(handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (result *mcp.CallToolResult, err error) {
		defer func() {
			if r := recover(); r != nil {
				// Get stack trace
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, true)
				stackTrace := string(buf[:n])

				result = mcp.NewToolResultError(fmt.Sprintf("Panic: %v\nStack trace:\n%s", r, stackTrace))
			}
		}()
		result, err = handler(ctx, request)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}
		return result, nil
	}
}
