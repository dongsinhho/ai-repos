package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/dongsinhho/ai-repos/mcp-server/services"
	"github.com/dongsinhho/ai-repos/mcp-server/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func RegisterMailTools(s *server.MCPServer) {
	// Search tool
	searchTool := mcp.NewTool("gmail_search",
		mcp.WithDescription("Search emails in Gmail using Gmail's search syntax"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Gmail search query. Follow Gmail's search syntax")),
	)
	s.AddTool(searchTool, utils.ErrorGuard(gmailSearchHandler))

	// Read email tool
	readEmailTool := mcp.NewTool("gmail_read_email",
		mcp.WithDescription("Read a specific email's full content including headers and body"),
		mcp.WithString("message_id", mcp.Required(), mcp.Description("ID of the email message to read")),
		mcp.WithBoolean("include_attachments", mcp.Description("Whether to include attachment information")),
	)
	s.AddTool(readEmailTool, utils.ErrorGuard(gmailReadEmailHandler))

	// Mark as read tool
	markReadTool := mcp.NewTool("gmail_mark_read",
		mcp.WithDescription("Mark a specific email as read (remove UNREAD label)"),
		mcp.WithString("message_id", mcp.Required(), mcp.Description("ID of the email message to mark as read")),
	)
	s.AddTool(markReadTool, utils.ErrorGuard(gmailMarkReadHandler))

	// Send email tool
	sendMailTool := mcp.NewTool("gmail_send_email",
		mcp.WithDescription("Send an email using Gmail API"),
		mcp.WithString("to", mcp.Required(), mcp.Description("Recipient email address(es), comma separated")),
		mcp.WithString("subject", mcp.Required(), mcp.Description("Email subject")),
		mcp.WithString("body", mcp.Required(), mcp.Description("Email body (plain text)")),
	)
	s.AddTool(sendMailTool, utils.ErrorGuard(gmailSendEmailHandler))
}

var gmailService = sync.OnceValue(func() *gmail.Service {
	ctx := context.Background()
	pwd, _ := os.Getwd()
	log.Println("pwd: ", pwd)
	credentialsFile := pwd + "/google-credential/credentials.json"
	tokenFile := pwd + "/google-credential/token.json"

	// tokenFile := os.Getenv("GOOGLE_TOKEN_FILE")
	// if tokenFile == "" {
	// 	panic("GOOGLE_TOKEN_FILE environment variable must be set")
	// }

	// credentialsFile := os.Getenv("GOOGLE_CREDENTIALS_FILE")
	// if credentialsFile == "" {
	// 	panic("GOOGLE_CREDENTIALS_FILE environment variable must be set")
	// }

	client := services.GoogleHttpClient(tokenFile, credentialsFile)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		panic(fmt.Sprintf("failed to create Gmail service: %v", err))
	}

	return srv
})

func gmailSearchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return mcp.NewToolResultError("query must be a string"), nil
	}

	user := "me"
	var allMessages []*gmail.Message
	pageToken := ""
	for {
		listCall := gmailService().Users.Messages.List(user).Q(query).MaxResults(100)
		if pageToken != "" {
			listCall.PageToken(pageToken)
		}
		resp, err := listCall.Do()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to search emails: %v", err)), nil
		}
		if len(resp.Messages) == 0 {
			break
		}
		for _, msg := range resp.Messages {
			message, err := gmailService().Users.Messages.Get(user, msg.Id).Do()
			if err != nil {
				log.Printf("Failed to get message %s: %v", msg.Id, err)
				continue
			}
			allMessages = append(allMessages, message)
		}
		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d emails:\n\n", len(allMessages)))
	for _, message := range allMessages {
		details := make(map[string]string)
		for _, header := range message.Payload.Headers {
			switch header.Name {
			case "From":
				details["from"] = header.Value
			case "Subject":
				details["subject"] = header.Value
			case "Date":
				details["date"] = header.Value
			}
		}
		result.WriteString(fmt.Sprintf("Message ID: %s\n", message.Id))
		result.WriteString(fmt.Sprintf("From: %s\n", details["from"]))
		result.WriteString(fmt.Sprintf("Subject: %s\n", details["subject"]))
		result.WriteString(fmt.Sprintf("Date: %s\n", details["date"]))
		result.WriteString(fmt.Sprintf("Snippet: %s\n", message.Snippet))
		result.WriteString("-------------------\n")
	}
	return mcp.NewToolResultText(result.String()), nil
}

func gmailReadEmailHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	messageID, ok := request.Params.Arguments["message_id"].(string)
	if !ok {
		return mcp.NewToolResultError("message_id must be a string"), nil
	}

	includeAttachments, _ := request.Params.Arguments["include_attachments"].(bool)

	// Get the full email message
	message, err := gmailService().Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get email: %v", err)), nil
	}

	var result strings.Builder

	// Extract headers
	headers := make(map[string]string)
	for _, header := range message.Payload.Headers {
		switch header.Name {
		case "From", "To", "Cc", "Subject", "Date":
			headers[header.Name] = header.Value
			result.WriteString(fmt.Sprintf("%s: %s\n", header.Name, header.Value))
		}
	}
	result.WriteString("\n")

	// Extract body
	body := extractMessageBody(message.Payload)
	result.WriteString("Body:\n")
	result.WriteString(body)
	result.WriteString("\n")

	// Handle attachments if requested
	if includeAttachments && len(message.Payload.Parts) > 0 {
		result.WriteString("\nAttachments:\n")
		for _, part := range message.Payload.Parts {
			if part.Filename != "" {
				result.WriteString(fmt.Sprintf("- %s (Size: %d bytes)\n",
					part.Filename, part.Body.Size))
			}
		}
	}

	return mcp.NewToolResultText(result.String()), nil
}

func gmailMarkReadHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	messageID, ok := request.Params.Arguments["message_id"].(string)
	if !ok {
		return mcp.NewToolResultError("message_id must be a string"), nil
	}

	user := "me"
	modifyReq := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}
	_, err := gmailService().Users.Messages.Modify(user, messageID, modifyReq).Do()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to mark email as read: %v", err)), nil
	}
	return mcp.NewToolResultText("Email marked as read."), nil
}

func gmailSendEmailHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	to, ok := request.Params.Arguments["to"].(string)
	if !ok || to == "" {
		return mcp.NewToolResultError("to must be a non-empty string"), nil
	}
	subject, ok := request.Params.Arguments["subject"].(string)
	if !ok {
		return mcp.NewToolResultError("subject must be a string"), nil
	}
	body, ok := request.Params.Arguments["body"].(string)
	if !ok {
		return mcp.NewToolResultError("body must be a string"), nil
	}

	// Encode subject as UTF-8 base64 for proper header formatting
	subjectHeader := fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
	msgStr := fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s", to, subjectHeader, body)
	msg := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(msgStr)),
	}
	_, err := gmailService().Users.Messages.Send("me", msg).Do()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to send email: %v", err)), nil
	}
	return mcp.NewToolResultText("Email sent successfully."), nil
}

func extractMessageBody(payload *gmail.MessagePart) string {
	if payload.MimeType == "text/plain" && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err != nil {
			return fmt.Sprintf("Error decoding body: %v", err)
		}
		return string(data)
	}

	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if part.MimeType == "text/plain" {
				data, err := base64.URLEncoding.DecodeString(part.Body.Data)
				if err != nil {
					continue
				}
				return string(data)
			}
		}
	}

	return "No readable text body found"
}
