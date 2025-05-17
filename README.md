# MCP Server for AI Integration

This repository contains the `mcp-server` source code, a modular server for integrating AI tools, automation, and external APIs using the Model Context Protocol (MCP).

## Features
- Modular tool registration (search, mail, filesystem, web automation, etc.)
- Supports user/session-based API key management for secure multi-user scenarios
- Easily extendable with new tools and services
- Written in Go for performance and portability

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/dongsinhho/ai-repos.git
   cd ai-repos/mcp-server
   ```
2. **Install dependencies:**
   ```bash
   go mod tidy
   ```
3. **Build the server:**
   ```bash
   go build -o server.exe main.go
   ```
4. **Run the server:**
   ```bash
   ./server.exe
   ```

## Tool System Overview
- Tools are registered in `tools/` (e.g., `mail-tools.go`, `search-tools.go`, `web-tools.go`)
- API keys for each user/service can be set and managed via the `set_api_key` tool
- Example tools: Brave Search, Gmail, Filesystem, Web Automation (browser)

## API Key Management
- API keys are never committed to the repository
- Use the `set_api_key` tool to set keys per user and service at runtime
- See `.gitignore` for sensitive file exclusions

## Development
- Add new tools in the `tools/` directory and register them in `main.go`
- Utility and session management code is in `utils/` and `session/`

## License
This project is licensed under the MIT License. See `LICENSE` for details.