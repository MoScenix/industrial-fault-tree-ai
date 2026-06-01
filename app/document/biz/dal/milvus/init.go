package milvus

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/MoScenix/industrial-fault-tree-ai/app/document/conf"
	"github.com/bytedance/sonic"
	embedollama "github.com/cloudwego/eino-ext/components/embedding/ollama"
	embedopenai "github.com/cloudwego/eino-ext/components/embedding/openai"
	indexermilvus "github.com/cloudwego/eino-ext/components/indexer/milvus2"
	retrievermilvus "github.com/cloudwego/eino-ext/components/retriever/milvus2"
	retrieversearchmode "github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

var (
	Client    *milvusclient.Client
	Indexer   indexer.Indexer
	Retriever retriever.Retriever
	Embedder  embedding.Embedder
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
	Embedder = embedder
	if err = ensureCollectionSchema(ctx, cfg); err != nil {
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
		DocumentConverter: func(ctx context.Context, docs []*schema.Document, vectors [][]float64) ([]column.Column, error) {
			return convertDocsToColumns(docs, vectors)
		},
	})
	if err != nil {
		panic(err)
	}

	defaultTopK := int(cfg.Embedding.TopK)
	Retriever, err = retrievermilvus.NewRetriever(ctx, &retrievermilvus.RetrieverConfig{
		Client:           Client,
		Collection:       cfg.Milvus.CollectionName,
		VectorField:      "vector",
		OutputFields:     []string{"id", "content", "metadata"},
		SearchMode:       retrieversearchmode.NewApproximate(retrievermilvus.COSINE),
		Embedding:        embedder,
		TopK:             defaultTopK,
		ConsistencyLevel: retrievermilvus.ConsistencyLevelStrong,
	})
	if err != nil {
		panic(err)
	}
}

func FlushCollection(ctx context.Context) error {
	cfg := conf.GetConf()
	if Client == nil || cfg.Milvus.CollectionName == "" {
		return nil
	}

	task, err := Client.Flush(ctx, milvusclient.NewFlushOption(cfg.Milvus.CollectionName))
	if err != nil {
		return err
	}
	return task.Await(ctx)
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

func ensureCollectionSchema(ctx context.Context, cfg *conf.Config) error {
	hasCollection, err := Client.HasCollection(ctx, milvusclient.NewHasCollectionOption(cfg.Milvus.CollectionName))
	if err != nil {
		return err
	}
	if hasCollection {
		return nil
	}

	idField := entity.NewField().
		WithName("id").
		WithDataType(entity.FieldTypeVarChar).
		WithMaxLength(512).
		WithIsPrimaryKey(true)
	contentField := entity.NewField().
		WithName("content").
		WithDataType(entity.FieldTypeVarChar).
		WithMaxLength(65535)
	metadataField := entity.NewField().
		WithName("metadata").
		WithDataType(entity.FieldTypeJSON)
	vectorField := entity.NewField().
		WithName("vector").
		WithDataType(entity.FieldTypeFloatVector).
		WithDim(cfg.Embedding.Dimension)
	ownerTypeField := entity.NewField().
		WithName("owner_type").
		WithDataType(entity.FieldTypeVarChar).
		WithMaxLength(32)
	ownerIDField := entity.NewField().
		WithName("owner_id").
		WithDataType(entity.FieldTypeVarChar).
		WithMaxLength(256)

	sch := entity.NewSchema().
		WithField(idField).
		WithField(contentField).
		WithField(metadataField).
		WithField(vectorField).
		WithField(ownerTypeField).
		WithField(ownerIDField)

	if err = Client.CreateCollection(ctx, milvusclient.NewCreateCollectionOption(cfg.Milvus.CollectionName, sch)); err != nil {
		return err
	}

	idx := index.NewHNSWIndex(entity.COSINE, 16, 100)
	task, err := Client.CreateIndex(ctx, milvusclient.NewCreateIndexOption(cfg.Milvus.CollectionName, "vector", idx))
	if err != nil {
		return err
	}
	if err = task.Await(ctx); err != nil {
		return err
	}

	loadTask, err := Client.LoadCollection(ctx, milvusclient.NewLoadCollectionOption(cfg.Milvus.CollectionName))
	if err != nil {
		return err
	}
	return loadTask.Await(ctx)
}

func convertDocsToColumns(docs []*schema.Document, vectors [][]float64) ([]column.Column, error) {
	ids := make([]string, 0, len(docs))
	contents := make([]string, 0, len(docs))
	metadatas := make([][]byte, 0, len(docs))
	vecs := make([][]float32, 0, len(docs))
	ownerTypes := make([]string, 0, len(docs))
	ownerIDs := make([]string, 0, len(docs))

	for i, doc := range docs {
		ids = append(ids, doc.ID)
		contents = append(contents, doc.Content)

		ownerType := valueFromMeta(doc.MetaData, "owner_type")
		ownerID := valueFromMeta(doc.MetaData, "owner_id")
		ownerTypes = append(ownerTypes, ownerType)
		ownerIDs = append(ownerIDs, ownerID)

		meta, err := sonic.Marshal(doc.MetaData)
		if err != nil {
			return nil, err
		}
		metadatas = append(metadatas, meta)

		var sourceVec []float64
		if len(vectors) == len(docs) {
			sourceVec = vectors[i]
		} else {
			sourceVec = doc.DenseVector()
		}
		vec := make([]float32, len(sourceVec))
		for j, v := range sourceVec {
			vec[j] = float32(v)
		}
		vecs = append(vecs, vec)
	}

	dim := 0
	if len(vecs) > 0 {
		dim = len(vecs[0])
	}
	return []column.Column{
		column.NewColumnVarChar("id", ids),
		column.NewColumnVarChar("content", contents),
		column.NewColumnJSONBytes("metadata", metadatas),
		column.NewColumnFloatVector("vector", dim, vecs),
		column.NewColumnVarChar("owner_type", ownerTypes),
		column.NewColumnVarChar("owner_id", ownerIDs),
	}, nil
}

func valueFromMeta(meta map[string]any, key string) string {
	if meta == nil {
		return ""
	}
	if v, ok := meta[key].(string); ok {
		return v
	}
	return ""
}
