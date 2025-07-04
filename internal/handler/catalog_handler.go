// handler/catalog_handler.go
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"rag_biddingchunks/internal/app"
	"rag_biddingchunks/internal/domain"
	"rag_biddingchunks/internal/util"
	"strconv"
)

type CatalogHandler struct {
	Service *app.CatalogService
}

func NewCatalogHandler(service *app.CatalogService) *CatalogHandler {
	return &CatalogHandler{Service: service}
}

func (h *CatalogHandler) GetCatalogChunks(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var rawArgs map[string]interface{}
	if err := json.Unmarshal(req.RawArguments, &rawArgs); err != nil {
		return util.ErrorResult(fmt.Sprintf("参数解析失败: %v", err)), nil
	}

	// 提取关键词
	keywords, _ := rawArgs["keyword"].(string)
	if keywords == "" {
		return util.ErrorResult("关键词不能为空"), nil
	}

	// 处理 top_k
	var topK int
	if rawTopK, ok := rawArgs["top_k"]; ok {
		switch v := rawTopK.(type) {
		case float64:
			topK = int(v)
		case string:
			if v == "" {
				topK = 5 // 使用默认值
			} else if n, err := strconv.Atoi(v); err == nil && n > 0 {
				topK = n
			} else {
				topK = 5
			}
		default:
			topK = 5
		}
	} else {
		topK = 5
	}

	// 处理 score
	var score float64
	if rawScore, ok := rawArgs["score"]; ok {
		switch v := rawScore.(type) {
		case float64:
			score = v
		case string:
			if v == "" {
				score = 0.7 // 使用默认值
			} else if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
				score = f
			} else {
				score = 0.7
			}
		default:
			score = 0.7
		}
	} else {
		score = 0.7
	}

	// 构造请求对象（可选）
	catalogReq := &domain.CatalogRequest{
		Keywords: keywords,
		TopK:     &topK,
		Score:    &score,
	}

	// 真正调用业务逻辑
	chunks, err := h.Service.GetCatalogChunks(ctx, catalogReq)
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

func (h *CatalogHandler) RegisterTools(srv *server.Server) {
	tool, _ := protocol.NewTool(
		"get_catalog_chunks",
		"根据关键词从目录中检索内容片段",
		domain.CatalogRequest{},
	)

	srv.RegisterTool(tool, func(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		return h.GetCatalogChunks(ctx, req)
	})
}
