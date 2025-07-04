package domain

type ContentRequest struct {
	Project  string  `json:"project"`
	Type     string  `json:"type"`
	Keywords string  `json:"keyword"`
	Score    float64 `json:"score,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PageSize *int    `json:"number,omitempty"`
}

type CatalogRequest struct {
	Keywords string   `json:"keyword"`
	Score    *float64 `json:"score,omitempty"`
	Page     *int     `json:"page,omitempty"`
	PageSize *int     `json:"number,omitempty"`
}

type RetrievalChunk struct {
	ID               string  `json:"id"`
	Content          string  `json:"content"`
	DocumentID       string  `json:"document_id"`
	DocumentKeyword  string  `json:"document_keyword"`
	Similarity       float64 `json:"similarity"`
	TermSimilarity   float64 `json:"term_similarity"`
	VectorSimilarity float64 `json:"vector_similarity"`
	Highlight        string  `json:"highlight"`
}

type RetrievalResponse struct {
	Code int `json:"code"`
	Data struct {
		Chunks []RetrievalChunk `json:"chunks"`
		Total  int              `json:"total"`
	} `json:"data"`
	Message string `json:"message,omitempty"`
}

type RagConf struct {
	RagFlow struct {
		Address    string              `yaml:"address"`
		Port       int                 `yaml:"port"`
		APIKey     string              `yaml:"api_key"`
		DatasetMap map[string][]string `yaml:"dataset_map"`
	} `yaml:"ragflow"`
}
