package main

import (
	"flag"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/joho/godotenv"

	"github.com/dongsinhho/ai-repos/mcp-server/tools"
)

func main() {
	envFile := flag.String("env", ".env", "Path to environment file")
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		fmt.Printf("Warning: Error loading env file %s: %v\n", *envFile, err)
	}
	mcpServer := server.NewMCPServer(
		"Fetch Kit",
		"1.0.0",
		server.WithLogging(),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	// enableTools := strings.Split(os.Getenv("ENABLE_TOOLS"), ",")
	// allToolsEnabled := len(enableTools) == 1 && enableTools[0] == ""

	// isEnabled := func(toolName string) bool {
	// 	return allToolsEnabled || slices.Contains(enableTools, toolName)
	// }

	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
	// tools.RegisterMailTools(mcpServer)
	// tools.RegisterAwsInsight(mcpServer)
	mcpServer.AddTool(tool, tools.HelloHandler)

	// if err := server.ServeStdio(mcpServer); err != nil {
	// 	panic(fmt.Sprintf("Server error: %v", err))
	// }
	sse := server.NewSSEServer(mcpServer, server.WithBasePath("http://localhost:5000"))
	err := sse.Start(":5000")
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
