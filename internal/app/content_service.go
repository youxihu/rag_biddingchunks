package app

import (
	"context"
	"log"
	"rag_biddingchunks/internal/domain"
	"rag_biddingchunks/internal/util"
)

type ContentService struct {
	Retriever domain.Retriever
}

func NewContentService(r domain.Retriever) *ContentService {
	return &ContentService{Retriever: r}
}

func (s *ContentService) GetContentChunks(ctx context.Context, req *domain.ContentRequest) ([]domain.RetrievalChunk, error) {
	log.Printf("【Content】 Received request: %+v", req)

	datasetIDs, err := s.Retriever.ResolveDatasetIDs(req.Type)
	if err != nil {
		log.Printf("Failed to resolve dataset IDs: %v", err)
		return nil, err
	}
	log.Printf("Resolved dataset IDs: %v", datasetIDs)

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 5
	}

	chunks, err := s.Retriever.SearchChunks(ctx, datasetIDs, req.Keywords, 1024, req.Score, 1, pageSize)
	if err != nil {
		log.Printf("Search chunks failed: %v", err)
		return nil, err
	}
	log.Printf("📤 从 RAGFlow 获取到 %d 条 chunk:", len(chunks))

	for i, c := range chunks {
		summary := util.SummarizeLog(c.Content, 40)
		log.Printf(" [%d] 【%s】相似度=%.2f | 内容: %s", i+1, c.DocumentKeyword, c.Similarity, summary)
	}

	log.Printf("Filtered to %d chunks by score >= %.2f", len(chunks), req.Score)

	return chunks, nil
}
