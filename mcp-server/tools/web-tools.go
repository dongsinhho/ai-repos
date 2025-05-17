package tools

import (
	"context"

	"github.com/dongsinhho/ai-repos/mcp-server/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Browser tool request/response structs

type OpenPageRequest struct {
	URL string `json:"url"`
}
type OpenPageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type ClickRequest struct {
	Selector string `json:"selector"`
}
type ClickResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type TypeRequest struct {
	Selector string `json:"selector"`
	Text     string `json:"text"`
}
type TypeResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type ScreenshotRequest struct {
	Selector string `json:"selector,omitempty"`
}
type ScreenshotResponse struct {
	ImageBase64 string `json:"image_base64"`
	Error       string `json:"error,omitempty"`
}

type EvaluateRequest struct {
	Script string `json:"script"`
}
type EvaluateResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

// Tool handler skeletons

func browserOpenPageHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, _ := request.Params.Arguments["url"].(string)
	// TODO: Implement browser automation logic
	return mcp.NewToolResultText("Opened page: " + url), nil
}

func browserClickHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector, _ := request.Params.Arguments["selector"].(string)
	// TODO: Implement browser automation logic
	return mcp.NewToolResultText("Clicked: " + selector), nil
}

func browserTypeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector, _ := request.Params.Arguments["selector"].(string)
	text, _ := request.Params.Arguments["text"].(string)
	// TODO: Implement browser automation logic
	return mcp.NewToolResultText("Typed '" + text + "' into: " + selector), nil
}

func browserScreenshotHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector, _ := request.Params.Arguments["selector"].(string)
	// TODO: Implement browser automation logic
	msg := "Screenshot taken (not implemented)"
	if selector != "" {
		msg += " for selector: " + selector
	}
	return mcp.NewToolResultText(msg), nil
}

func browserEvaluateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	script, _ := request.Params.Arguments["script"].(string)
	// TODO: Implement browser automation logic
	return mcp.NewToolResultText("Evaluated script: " + script), nil
}

func RegisterWebTools(s *server.MCPServer) {
	// Open Page tool
	openPageTool := mcp.NewTool("browser_open_page",
		mcp.WithDescription("Open a web page in the browser"),
		mcp.WithString("url", mcp.Required(), mcp.Description("URL to open")),
	)
	s.AddTool(openPageTool, utils.ErrorGuard(browserOpenPageHandler))

	// Click tool
	clickTool := mcp.NewTool("browser_click",
		mcp.WithDescription("Click an element in the browser by selector"),
		mcp.WithString("selector", mcp.Required(), mcp.Description("CSS selector to click")),
	)
	s.AddTool(clickTool, utils.ErrorGuard(browserClickHandler))

	// Type tool
	typeTool := mcp.NewTool("browser_type",
		mcp.WithDescription("Type text into an element in the browser by selector"),
		mcp.WithString("selector", mcp.Required(), mcp.Description("CSS selector to type into")),
		mcp.WithString("text", mcp.Required(), mcp.Description("Text to type")),
	)
	s.AddTool(typeTool, utils.ErrorGuard(browserTypeHandler))

	// Screenshot tool
	screenshotTool := mcp.NewTool("browser_screenshot",
		mcp.WithDescription("Take a screenshot of the page or an element"),
		mcp.WithString("selector", mcp.Description("CSS selector to screenshot (optional)")),
	)
	s.AddTool(screenshotTool, utils.ErrorGuard(browserScreenshotHandler))

	// Evaluate tool
	evaluateTool := mcp.NewTool("browser_evaluate",
		mcp.WithDescription("Evaluate JavaScript in the browser context"),
		mcp.WithString("script", mcp.Required(), mcp.Description("JavaScript code to evaluate")),
	)
	s.AddTool(evaluateTool, utils.ErrorGuard(browserEvaluateHandler))
}
