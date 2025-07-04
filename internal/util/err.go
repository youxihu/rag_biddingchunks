package util

import "github.com/ThinkInAIXYZ/go-mcp/protocol"

// ErrorResult 返回标准格式的错误内容
func ErrorResult(msg string) *protocol.CallToolResult {
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: msg,
			},
		},
	}
}
