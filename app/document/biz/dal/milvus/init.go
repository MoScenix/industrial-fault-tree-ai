package milvus

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/MoScenix/industrial-fault-tree-ai/app/document/conf"
	embedopenai "github.com/cloudwego/eino-ext/components/embedding/openai"
	indexermilvus "github.com/cloudwego/eino-ext/components/indexer/milvus"
	retrievermilvus "github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

var (
	Client    client.Client
	Indexer   indexer.Indexer
	Retriever retriever.Retriever
)

var errUnsupportedEmbeddingProvider = errors.New("unsupported embedding provider")

func Init() {
	cfg := conf.GetConf()
	milvusPassword := fmt.Sprintf(cfg.Milvus.Password, os.Getenv("MILVUS_PASSWORD"))
	embeddingAPIKey := fmt.Sprintf(cfg.Embedding.APIKey, os.Getenv("EMBEDDING_API_KEY"))
	if cfg.Milvus.Address == "" || embeddingAPIKey == "" {
		return
	}

	ctx := context.Background()
	var err error
	Client, err = client.NewClient(ctx, client.Config{
		Address:  cfg.Milvus.Address,
		Username: cfg.Milvus.Username,
		Password: milvusPassword,
	})
	if err != nil {
		panic(err)
	}

	if cfg.Milvus.Database != "" {
		if err := Client.UsingDatabase(ctx, cfg.Milvus.Database); err != nil {
			panic(err)
		}
	}

	embedder, err := newEmbedder(ctx, embeddingAPIKey)
	if err != nil {
		panic(err)
	}

	Indexer, err = indexermilvus.NewIndexer(ctx, &indexermilvus.IndexerConfig{
		Client:     Client,
		Collection: cfg.Milvus.CollectionName,
		Embedding:  embedder,
		MetricType: indexermilvus.IP,
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
		MetricType:   entity.IP,
		Embedding:    embedder,
		TopK:         defaultTopK,
	})
	if err != nil {
		panic(err)
	}
}

func newEmbedder(ctx context.Context, apiKey string) (embedding.Embedder, error) {
	cfg := conf.GetConf()
	switch cfg.Embedding.Provider {
	case "", "openai":
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
	default:
		return nil, errUnsupportedEmbeddingProvider
	}
}
