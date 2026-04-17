package milvus

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/MoScenix/industrial-fault-tree-ai/app/document/conf"
	embedollama "github.com/cloudwego/eino-ext/components/embedding/ollama"
	embedopenai "github.com/cloudwego/eino-ext/components/embedding/openai"
	indexermilvus "github.com/cloudwego/eino-ext/components/indexer/milvus2"
	retrievermilvus "github.com/cloudwego/eino-ext/components/retriever/milvus2"
	retrieversearchmode "github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

var (
	Client    *milvusclient.Client
	Indexer   indexer.Indexer
	Retriever retriever.Retriever
)

var errUnsupportedEmbeddingProvider = errors.New("unsupported embedding provider")

func Init() {
	cfg := conf.GetConf()
	milvusPassword := fmt.Sprintf(cfg.Milvus.Password, os.Getenv("MILVUS_PASSWORD"))
	if cfg.Milvus.Address == "" {
		return
	}

	ctx := context.Background()
	clientConfig := &milvusclient.ClientConfig{
		Address:  cfg.Milvus.Address,
		Username: cfg.Milvus.Username,
		Password: milvusPassword,
		DBName:   cfg.Milvus.Database,
	}

	var err error
	Client, err = milvusclient.New(ctx, clientConfig)
	if err != nil {
		panic(err)
	}

	embedder, err := newEmbedder(ctx)
	if err != nil {
		panic(err)
	}

	Indexer, err = indexermilvus.NewIndexer(ctx, &indexermilvus.IndexerConfig{
		Client:     Client,
		Collection: cfg.Milvus.CollectionName,
		Vector: &indexermilvus.VectorConfig{
			VectorField:  "vector",
			Dimension:    cfg.Embedding.Dimension,
			MetricType:   indexermilvus.COSINE,
			IndexBuilder: indexermilvus.NewHNSWIndexBuilder(),
		},
		Embedding: embedder,
	})
	if err != nil {
		panic(err)
	}

	defaultTopK := int(cfg.Embedding.TopK)
	Retriever, err = retrievermilvus.NewRetriever(ctx, &retrievermilvus.RetrieverConfig{
		Client:       Client,
		Collection:   cfg.Milvus.CollectionName,
		VectorField:  "vector",
		OutputFields: []string{"id", "content", "metadata"},
		SearchMode:   retrieversearchmode.NewApproximate(retrievermilvus.COSINE),
		Embedding:    embedder,
		TopK:         defaultTopK,
	})
	if err != nil {
		panic(err)
	}
}

func newEmbedder(ctx context.Context) (embedding.Embedder, error) {
	cfg := conf.GetConf()
	switch cfg.Embedding.Provider {
	case "", "openai":
		apiKey := resolveEmbeddingAPIKey(cfg.Embedding.APIKey)
		if strings.TrimSpace(apiKey) == "" {
			return nil, errors.New("embedding api key is required")
		}
		var dimensions *int
		if cfg.Embedding.Dimension > 0 {
			d := int(cfg.Embedding.Dimension)
			dimensions = &d
		}

		return embedopenai.NewEmbedder(ctx, &embedopenai.EmbeddingConfig{
			APIKey:     apiKey,
			BaseURL:    cfg.Embedding.BaseURL,
			Model:      cfg.Embedding.Model,
			Dimensions: dimensions,
		})
	case "ollama":
		return embedollama.NewEmbedder(ctx, &embedollama.EmbeddingConfig{
			BaseURL: cfg.Embedding.BaseURL,
			Model:   cfg.Embedding.Model,
		})
	default:
		return nil, errUnsupportedEmbeddingProvider
	}
}

func resolveEmbeddingAPIKey(template string) string {
	if key := strings.TrimSpace(os.Getenv("EMBEDDING_API_KEY")); key != "" {
		return fmt.Sprintf(template, key)
	}
	if key := strings.TrimSpace(os.Getenv("DASHSCOPE_API_KEY")); key != "" {
		return fmt.Sprintf(template, key)
	}
	return strings.TrimSpace(fmt.Sprintf(template, ""))
}
