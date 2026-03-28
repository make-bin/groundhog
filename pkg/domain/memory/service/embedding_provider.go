package service

import "context"

// EmbeddingProvider is the domain port for generating text embeddings.
type EmbeddingProvider interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}
