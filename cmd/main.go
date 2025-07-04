package main

import (
	"log"
	"net/http"
	"rag_biddingchunks/internal/util"

	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"rag_biddingchunks/config"
	"rag_biddingchunks/internal/app"
	"rag_biddingchunks/internal/handler"
	"rag_biddingchunks/internal/infra"
)

func init() {
	if err := config.LoadConfig("/app-acc/configs/online.auth.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
}

func main() {
	retriever := infra.NewRagflowRetriever(&config.Cfg)

	// 初始化服务
	contentService := app.NewContentService(retriever)
	catalogService := app.NewCatalogService(retriever)

	// 初始化 handler
	contentHandler := handler.NewContentHandler(contentService)
	catalogHandler := handler.NewCatalogHandler(catalogService)

	// 创建 MCP Server
	transportImpl, httpHandler, err := transport.NewStreamableHTTPServerTransportAndHandler(
		transport.WithStreamableHTTPServerTransportAndHandlerOptionStateMode(transport.Stateless))
	if err != nil {
		log.Fatalf("创建 MCP Transport 失败: %v", err)
	}

	mcpServer, _ := server.NewServer(transportImpl)

	contentHandler.RegisterTools(mcpServer)
	catalogHandler.RegisterTools(mcpServer)
	finalHandler := util.WrapMCPHandler(httpHandler.HandleMCP())
	localIP := util.GetLocalIP()
	log.Printf("🚀 MCP HTTP Streaming 服务运行于 http://%s:25003/mcp", localIP)
	log.Fatal(http.ListenAndServe(":25003", finalHandler))
}
