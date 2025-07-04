package infra

import (
	"context"
	"fmt"
	"strings"
	"time"

	"rag_biddingchunks/internal/domain"
	"rag_biddingchunks/internal/util"
)

type RagflowRetriever struct {
	cfg *domain.RagConf
}

func NewRagflowRetriever(cfg *domain.RagConf) *RagflowRetriever {
	return &RagflowRetriever{cfg: cfg}
}

func (r *RagflowRetriever) ResolveDatasetIDs(docType string) ([]string, error) {
	if docType == "" {
		var all []string
		for _, ids := range r.cfg.RagFlow.DatasetMap {
			all = append(all, ids...)
		}
		if len(all) == 0 {
			return nil, fmt.Errorf("no datasets found in dataset_map")
		}
		return all, nil
	}

	types := strings.Split(docType, ",")
	idSet := make(map[string]struct{})

	for _, t := range types {
		t = strings.TrimSpace(t)
		var key string
		switch t {
		case "采购", "货物", "3":
			key = "3"
		case "工程", "建筑", "1":
			key = "1"
		case "服务", "2":
			key = "2"
		default:
			return nil, fmt.Errorf("unknown document type: %s", t)
		}

		ids, ok := r.cfg.RagFlow.DatasetMap[key]
		if !ok || len(ids) == 0 {
			return nil, fmt.Errorf("no dataset IDs found for type: %s", key)
		}

		for _, id := range ids {
			idSet[id] = struct{}{}
		}
	}

	// map -> slice 去重后返回
	var result []string
	for id := range idSet {
		result = append(result, id)
	}

	return result, nil
}

func (r *RagflowRetriever) SearchChunks(ctx context.Context, datasetIDs []string, keywords string, topK int, score float64, page int, pageSize int) ([]domain.RetrievalChunk, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	url := fmt.Sprintf("http://%s:%d/api/v1/retrieval", r.cfg.RagFlow.Address, r.cfg.RagFlow.Port)
	body := map[string]interface{}{
		"question":             keywords,
		"dataset_ids":          datasetIDs,
		"top_k":                topK,
		"similarity_threshold": score,
		"keyword":              true,
		"highlight":            false,
		"page":                 page,
		"page_size":            pageSize,
	}

	respBody, err := util.PostJSON(ctx, url, r.cfg.RagFlow.APIKey, body)
	if err != nil {
		return nil, fmt.Errorf("多次请求失败: %w", err)
	}

	var result domain.RetrievalResponse
	if err := util.UnmarshalJSON(respBody, &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("检索失败: %s (code=%d)", result.Message, result.Code)
	}

	return result.Data.Chunks, nil
}
