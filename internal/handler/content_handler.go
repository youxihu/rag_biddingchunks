package handler

import (
	"context"
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"rag_biddingchunks/internal/app"
	"rag_biddingchunks/internal/domain"
	"rag_biddingchunks/internal/util"
)

type ContentHandler struct {
	Service *app.ContentService
}

func NewContentHandler(service *app.ContentService) *ContentHandler {
	return &ContentHandler{Service: service}
}

func (h *ContentHandler) GetContentChunks(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var params domain.ContentRequest
	if err := protocol.VerifyAndUnmarshal(req.RawArguments, &params); err != nil {
		return util.ErrorResult(fmt.Sprintf("参数解析失败: %v", err)), nil
	}
	chunks, err := h.Service.GetContentChunks(ctx, &params)
	if err != nil {
		return util.ErrorResult(fmt.Sprintf("检索失败: %v", err)), nil
	}
	contents := util.ToTextContent(chunks)
	var result []protocol.Content
	for _, c := range contents {
		result = append(result, &protocol.TextContent{Type: "text", Text: c})
	}
	return &protocol.CallToolResult{Content: result}, nil
}

func (h *ContentHandler) RegisterTools(srv *server.Server) {
	tool, _ := protocol.NewTool(
		"get_content_chunks",
		"根据关键词检索标书内容片段",
		domain.ContentRequest{},
	)

	srv.RegisterTool(tool, func(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		return h.GetContentChunks(ctx, req)
	})
}
