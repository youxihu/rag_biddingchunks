package util

import (
	"fmt"
	"rag_biddingchunks/internal/domain"
)

func ToTextContent(chunks []domain.RetrievalChunk) []string {
	var contents []string
	for _, chunk := range chunks {
		if chunk.Similarity >= 0 {
			contents = append(contents, fmt.Sprintf("【%s】%.2f\n%s", chunk.DocumentKeyword, chunk.Similarity, chunk.Content))
		}
	}
	if len(contents) == 0 {
		contents = append(contents, "没有找到符合条件的内容。")
	}
	return contents
}
