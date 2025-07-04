package app

import (
	"context"
	"log"
	"rag_biddingchunks/internal/domain"
	"rag_biddingchunks/internal/util"
)

type CatalogService struct {
	Retriever domain.Retriever
}

func NewCatalogService(r domain.Retriever) *CatalogService {
	return &CatalogService{Retriever: r}
}

func (s *CatalogService) GetCatalogChunks(ctx context.Context, req *domain.CatalogRequest) ([]domain.RetrievalChunk, error) {
	util.LogWithIP(ctx, "调用工具:【get_catalog_chunks】Keywords=%q, TopK=%d, Score=%.2f",
		req.Keywords, *req.TopK, *req.Score)

	// 设置默认值
	defaultTopK := 5
	if req.TopK == nil || *req.TopK <= 0 {
		req.TopK = &defaultTopK
	}

	defaultScore := 0.5
	if req.Score == nil || *req.Score <= 0 {
		req.Score = &defaultScore
	}

	datasetIDs := []string{"01cf583657cf11f0b5690242ac1a0003"}

	chunks, err := s.Retriever.SearchChunks(ctx, datasetIDs, req.Keywords, *req.TopK, *req.Score)
	if err != nil {
		return nil, err
	}
	log.Printf("从 RAGFlow 获取到 %d 条 chunk:", len(chunks))

	for i, c := range chunks {
		summary := util.SummarizeLog(c.Content, 50)
		log.Printf(" [%d] 【%s】相似度=%.2f | 内容: %s", i+1, c.DocumentKeyword, c.Similarity, summary)
	}

	log.Printf("Filtered to %d chunks by score >= %.2f", len(chunks), *req.Score)
	return chunks, nil
}
