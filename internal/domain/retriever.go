package domain

import (
	"context"
)

type Retriever interface {
	ResolveDatasetIDs(docType string) ([]string, error)
	SearchChunks(ctx context.Context, datasetIDs []string, keywords string, topK int, score float64, page int, pageSize int) ([]RetrievalChunk, error)
}
