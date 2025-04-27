package main

import (
	"flag"
	"fmt"

	"github.com/mark3labs/mcp-go/server"

	"github.com/joho/godotenv"
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

	// tools.RegisterMailTools(mcpServer)
	// tools.RegisterAwsInsight(mcpServer)

	if err := server.ServeStdio(mcpServer); err != nil {
		panic(fmt.Sprintf("Server error: %v", err))
	}
}
