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
	util.LogWithIP(ctx, "调用工具:【get_catalog_chunks】Keywords=%q, Score=%.2f, Page=%d, PageSize=%d",
		req.Keywords, *req.Score, 1, *req.PageSize)

	datasetIDs := []string{"01cf583657cf11f0b5690242ac1a0003"}

	pageSize := 5
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = *req.PageSize
	}

	chunks, err := s.Retriever.SearchChunks(ctx, datasetIDs, req.Keywords, 1024, *req.Score, 1, pageSize)

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
