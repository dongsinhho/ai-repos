package main

import (
	"flag"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"

	"github.com/dongsinhho/ai-repos/mcp-server/tools"
)

func main() {
	envFile := flag.String("env", ".env", "Path to environment file")
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		//fmt.Printf("Warning: Error loading env file %s: %v\n", *envFile, err)
	}
	mcpServer := server.NewMCPServer(
		"Demo",
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

	tools.RegisterMailTools(mcpServer)
	tools.RegisterFilesystemTools(mcpServer)

	// if err := server.ServeStdio(mcpServer); err != nil {
	// 	panic(fmt.Sprintf("Server error: %v", err))
	// }
	sse := server.NewSSEServer(mcpServer,
		server.WithBaseURL("http://localhost:8082"))
	err := sse.Start(":8082")
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
