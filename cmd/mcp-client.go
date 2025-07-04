package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/client"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"log"
)

func main() {
	//serverURL := "http://47.97.157.44:25007/mcp"
	serverURL := "http://192.168.0.43:25003/mcp"
	transportClient, err := transport.NewStreamableHTTPClientTransport(serverURL)
	if err != nil {
		log.Fatalf("连接 MCP 服务失败: %v", err)
	}

	mcpClient, err := client.NewClient(transportClient)
	if err != nil {
		log.Fatalf("创建 MCP 客户端失败: %v", err)
	}
	defer mcpClient.Close()

	// 获取可用工具列表
	toolsResult, err := mcpClient.ListTools(context.Background())
	if err != nil {
		log.Fatalf("获取工具失败: %v", err)
	}

	fmt.Println("🔍 可用工具:")
	for _, tool := range toolsResult.Tools {
		fmt.Printf(" - %s: %s\n", tool.Name, tool.Description)
	}

	// 调用工具
	callTool(mcpClient, "get_content_chunks", map[string]interface{}{
		"project":   "项目A",
		"type":      "1",
		"keyword":   "钻孔灌注桩施工方案",
		"score":     0.5,
		"page_size": 1,
	})

	callTool(mcpClient, "get_catalog_chunks", map[string]interface{}{
		"keyword":   "大型土石方项目",
		"score":     0.1,
		"page_size": 1,
	})

}

func callTool(c *client.Client, toolName string, args map[string]interface{}) {
	fmt.Printf("\n🛠️ 正在调用工具: %s\n", toolName)

	body, _ := json.Marshal(args)
	req := &protocol.CallToolRequest{
		Name:         toolName,
		RawArguments: body,
	}

	result, err := c.CallTool(context.Background(), req)
	if err != nil {
		fmt.Printf("❌ 调用工具失败: %v\n", err)
		return
	}

	for _, content := range result.Content {
		if textContent, ok := content.(*protocol.TextContent); ok {
			fmt.Println(" =>", textContent.Text)
		} else {
			fmt.Printf(" => 不支持的内容类型: %T\n", content)
		}
	}
}
