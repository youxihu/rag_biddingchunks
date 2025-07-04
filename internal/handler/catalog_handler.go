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

	// 不再解析 top_k，固定 1024
	const fixedTopK = 1024

	// 处理 score
	var score float64
	if rawScore, ok := rawArgs["score"]; ok {
		switch v := rawScore.(type) {
		case float64:
			score = v
		case string:
			if v == "" {
				score = 0.7 // 默认值
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

	// 解析 page 参数，默认1
	var page int = 1
	if rawPage, ok := rawArgs["page"]; ok {
		switch v := rawPage.(type) {
		case float64:
			page = int(v)
		case string:
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				page = n
			}
		}
		if page <= 0 {
			page = 1
		}
	}

	// 解析 page_size 参数，默认5
	var pageSize int = 5
	if rawPageSize, ok := rawArgs["page_size"]; ok {
		switch v := rawPageSize.(type) {
		case float64:
			pageSize = int(v)
		case string:
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				pageSize = n
			}
		}
		if pageSize <= 0 {
			pageSize = 5
		}
	}

	// 构造请求对象，去掉TopK，新增page和pageSize
	catalogReq := &domain.CatalogRequest{
		Keywords: keywords,
		Score:    &score,
		Page:     &page,
		PageSize: &pageSize,
	}

	// 调用业务逻辑
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
